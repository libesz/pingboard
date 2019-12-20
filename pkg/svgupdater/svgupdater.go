package svgupdater

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/scheduler"
	"github.com/libesz/pingboard/pkg/svgmanip"
)

func Run(ctx context.Context, requestChan <-chan chan *etree.Document, resultChan chan scheduler.ResultChange, origSvg *etree.Document, allUpdateRules []config.Target) {
	log.Println("[svgupdater] Started up")
	var actualUpdateRules []svgmanip.Target
	svg := origSvg.Copy()
	for {
		select {
		case result := <-resultChan:
			if result.Value {
				for i, v := range allUpdateRules {
					if v.ID == result.ID {
						actualUpdateRules = append(actualUpdateRules, svgmanip.Target{ID: result.ID, Fill: allUpdateRules[i].Fill, LastChange: "Last change at: " + time.Now().Format(time.RFC1123)})
					}
				}
			} else {
				for i, v := range actualUpdateRules {
					if v.ID == result.ID {
						actualUpdateRules = append(actualUpdateRules[:i], actualUpdateRules[i+1:]...)
					}
				}
			}
			svg = origSvg.Copy()
			if err := svgmanip.UpdateDoc(svg, actualUpdateRules); err != nil {
				log.Println("Error during SVG update: ", err)
			}
		case clientSvgChan := <-requestChan:
			//fmt.Println("sent data:", svg)
			clientSvgChan <- svg
		case <-ctx.Done():
			log.Println("[svgupdater] Exiting")
			return
		}
	}
}

func Get(requestChan chan<- chan *etree.Document) (*etree.Document, error) {
	updateChan := make(chan *etree.Document)
	requestChan <- updateChan
	timeout := time.After(5 * time.Second)
	select {
	case <-timeout:
		break
	case svg := <-updateChan:
		//fmt.Println("rec data:", svg)
		return svg, nil
	}
	return nil, errors.New("[svgupdater] Timeout happened when tried to gather SVG content from updater")
}
