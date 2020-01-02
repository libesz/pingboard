//go:generate go run generate_static.go
package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/net/websocket"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/scheduler"
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

	configSource, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	configData, err := config.Get(configSource)
	if err != nil {
		panic(err)
	}

	svg := etree.NewDocument()
	if err := svg.ReadFromFile(configData.SvgPath); err != nil {
		panic(err)
	}

	if err = config.Validate(configData, svg); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	resultChan := make(chan scheduler.ResultChange)
	wg.Add(1)
	go func() {
		scheduler.Run(ctx, configData.Targets, resultChan)
		wg.Done()
	}()

	updater := svgupdater.New(resultChan, svg, configData.Targets)
	wg.Add(1)
	go func() {
		updater.Run(ctx)
		wg.Done()
	}()

	//svgHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {handleSvg(requestChan, w, req) })
	svgHandler := websocket.Handler(func(ws *websocket.Conn) {
		handleSvg(ctx, &updater, ws)
	})
	mux := http.NewServeMux()
	mux.Handle("/svg", svgHandler)
	mux.Handle("/", http.FileServer(static))

	server := &http.Server{Addr: ":2003", Handler: mux}
	wg.Add(1)
	go func() {
		log.Println("[main] HTTP server is listening on :2003")
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

func handleSvg(ctx context.Context, updater *svgupdater.SvgUpdater, ws *websocket.Conn) {
	log.Println("[main] Got connection from client: " + ws.RemoteAddr().String())
	clientClosed := make(chan struct{})
	go func() {
		var data []byte
		websocket.Message.Receive(ws, data)
		close(clientClosed)
	}()

	svgChan := updater.Register()
	svg, err := updater.Get()
	if err != nil {
		log.Println("[main] SVG updater error: " + err.Error())
	}
	pushSvgToWs(svg, ws)
	for {
		select {
		case svg := <-svgChan:
			pushSvgToWs(svg, ws)
		case <-ctx.Done():
			log.Println("[main] WebSocket handler returns")
			return
		case <-clientClosed:
			log.Println("[main] Client disconnected:", ws.RemoteAddr().String())
			updater.DeRegister(svgChan)
			return
		}
	}
}

func pushSvgToWs(svg *etree.Document, ws *websocket.Conn) {
	str, err := svg.WriteToBytes()
	if err != nil {
		log.Println("[main] SVG error: " + err.Error())
	}
	str2 := base64.StdEncoding.EncodeToString(str)
	websocket.Message.Send(ws, []byte(str2))
	log.Println("[main] Update sent to client:", ws.RemoteAddr().String())
}
