package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/scheduler"
	"github.com/libesz/pingboard/pkg/svgmanip"
	"github.com/libesz/pingboard/pkg/svgupdater"
)

func signalHandler(cancel context.CancelFunc, sigs chan os.Signal) {
	sig := <-sigs
	log.Println("[main] Signal received: " + sig.String())
	cancel()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go signalHandler(cancel, sigs)

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

	var wg sync.WaitGroup

	resultChan := make(chan scheduler.ResultChange)
	wg.Add(1)
	go func() {
		scheduler.Run(ctx, configData.Targets, resultChan)
		wg.Done()
	}()

	requestChan := make(chan chan *etree.Document)
	wg.Add(1)
	go func() {
		svgupdater.Run(ctx, requestChan, resultChan, svg, configData.Targets)
		wg.Done()
	}()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { handleSvg(requestChan, w, req) })
	server := &http.Server{Addr: ":2003", Handler: handler}
	wg.Add(1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("[main] ListenAndServe: " + err.Error())
		}
		wg.Done()
	}()
	<-ctx.Done()
	log.Println("[main] Exiting, waiting everybody to return...")
	server.Close()
	wg.Wait()
	log.Println("[main] Exiting, done")
}

func handleSvg(requestChan chan chan *etree.Document, w http.ResponseWriter, req *http.Request) {
	log.Println("[main] Got request from client: " + req.RemoteAddr)
	svg, err := svgupdater.Get(requestChan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.Println("[main] Error 500 happened with error: " + err.Error())
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	svg.WriteTo(w)
	log.Println("[main] Sent response to client: " + req.RemoteAddr)
	return
}
