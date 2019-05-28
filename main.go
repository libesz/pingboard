package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/svgmanip"
	"github.com/tatsushid/go-fastping"
	"gopkg.in/yaml.v2"
)

var requestChan chan chan *etree.Document

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
	requestChan = make(chan chan *etree.Document)
	go updater(requestChan, svg)
	http.Handle("/", http.HandlerFunc(handleSvg))
	err = http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	} //updatee(requestChan)

	//svg.WriteToFile("new.svg")
	//ping()
}

func handleSvg(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Got request")
	svg := updatee(requestChan)
	if svg == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		fmt.Println("500")
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	svg.WriteTo(w)
	fmt.Println("Sent response")
	return
}

func updater(requestChan <-chan chan *etree.Document, svg *etree.Document) {
	for {
		updateChan := <-requestChan
		fmt.Println("sent data:", svg)
		updateChan <- svg
	}
}

func updatee(requestChan chan<- chan *etree.Document) *etree.Document {
	updateChan := make(chan *etree.Document)
	requestChan <- updateChan
	timeout := time.After(3 * time.Second)
	select {
	case <-timeout:
		fmt.Println("timeout")
	case svg := <-updateChan:
		fmt.Println("rec data:", svg)
		return svg
	}
	return nil
}

func ping() {
	//pinger.SetPrivileged(true)
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	p.MaxRTT = time.Millisecond * 1000
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
	//time.Sleep(10 * time.Second)
}
