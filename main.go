package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/scheduler"
	"github.com/libesz/pingboard/pkg/svgmanip"
	"github.com/libesz/pingboard/pkg/svgupdater"
)

func main() {
	configData, err := config.Get(os.Args[1])
	if err != nil {
		panic(err)
	}
	svg := etree.NewDocument()
	if err = svg.ReadFromFile(configData.SvgPath); err != nil {
		panic(err)
	}
	if err = svgmanip.CheckDoc(svg, configData); err != nil {
		panic(err)
	}

	resultChan := make(chan scheduler.ResultChange)
	go scheduler.Run(context.Background(), configData.Targets, resultChan)
	requestChan := make(chan chan *etree.Document)
	go svgupdater.Run(requestChan, resultChan, svg, configData.Targets)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handleSvg(requestChan, w, req)
	}))

	err = http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleSvg(requestChan chan chan *etree.Document, w http.ResponseWriter, req *http.Request) {
	log.Println("Got request from client: " + req.RemoteAddr)
	svg, err := svgupdater.Get(requestChan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.Println("Error 500 happened with error: " + err.Error())
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	svg.WriteTo(w)
	log.Println("Sent response to client: " + req.RemoteAddr)
	return
}
