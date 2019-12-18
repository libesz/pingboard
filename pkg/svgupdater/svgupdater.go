package svgupdater

import (
	"errors"
	"log"
	"time"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/scheduler"
	"github.com/libesz/pingboard/pkg/svgmanip"
)

func Run(requestChan <-chan chan *etree.Document, resultChan chan scheduler.ResultChange, origSvg *etree.Document, allUpdateRules []config.Target) {
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
				log.Println("Error during SVG update: ", err)
			}
		case clientSvgChan := <-requestChan:
			//fmt.Println("sent data:", svg)
			clientSvgChan <- svg
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
	return nil, errors.New("Timeout happened when tried to gather SVG content")
}
