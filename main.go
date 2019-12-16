package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/scheduler"
	"github.com/libesz/pingboard/pkg/svgmanip"
	"gopkg.in/yaml.v2"
)

func main() {
	filename := os.Args[1]
	configSource, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	config := config.Config{}

	err = yaml.Unmarshal([]byte(configSource), &config)
	if err != nil {
		panic(err)
	}

	svg := etree.NewDocument()
	if err := svg.ReadFromFile(config.SvgPath); err != nil {
		panic(err)
	}
	if err = svgmanip.CheckDoc(svg, config); err != nil {
		panic(err)
	}
	resultChan := make(chan scheduler.ResultChange)
	go scheduler.Run(context.Background(), config.Targets, resultChan)
	requestChan := make(chan chan *etree.Document)
	go updater(requestChan, resultChan, svg, config.Targets)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handleSvg(requestChan, w, req)
	}))
	err = http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleSvg(requestChan chan chan *etree.Document, w http.ResponseWriter, req *http.Request) {
	fmt.Println("Got request")
	svg := updatee(requestChan)
	if svg == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		fmt.Println("500")
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	svg.WriteTo(w)
	fmt.Println("Sent response")
	return
}

func updater(requestChan <-chan chan *etree.Document, resultChan chan scheduler.ResultChange, origSvg *etree.Document, allUpdateRules []config.Target) {
	var actualUpdateRules []config.Target
	svg := origSvg.Copy()
	for {
		select {
		case result := <-resultChan:
			if result.Value {
				for i, v := range allUpdateRules {
					if v.SvgID == result.ID {
						actualUpdateRules = append(actualUpdateRules, allUpdateRules[i])
					}
				}
			} else {
				for i, v := range actualUpdateRules {
					if v.SvgID == result.ID {
						actualUpdateRules = append(actualUpdateRules[:i], actualUpdateRules[i+1:]...)
					}
				}
			}
			svg = origSvg.Copy()
			if err := svgmanip.UpdateDoc(svg, actualUpdateRules); err != nil {
				fmt.Println("Error during update: ", err)
			}
		case clientSvgChan := <-requestChan:
			fmt.Println("sent data:", svg)
			clientSvgChan <- svg
		}
	}
}

func updatee(requestChan chan<- chan *etree.Document) *etree.Document {
	updateChan := make(chan *etree.Document)
	requestChan <- updateChan
	timeout := time.After(5 * time.Second)
	select {
	case <-timeout:
		fmt.Println("timeout")
	case svg := <-updateChan:
		fmt.Println("rec data:", svg)
		return svg
	}
	return nil
}
