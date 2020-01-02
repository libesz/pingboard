package svgupdater

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/scheduler"
	"github.com/libesz/pingboard/pkg/svgmanip"
)

type SvgUpdater struct {
	resultChan     chan scheduler.ResultChange
	origSvg        *etree.Document
	allUpdateRules []config.Target
	adhocRequest   chan chan *etree.Document
	observers      []chan *etree.Document
	register       chan chan *etree.Document
	deregister     chan chan *etree.Document
}

func New(resultChan chan scheduler.ResultChange, origSvg *etree.Document, allUpdateRules []config.Target) SvgUpdater {
	adhocRequest := make(chan chan *etree.Document)
	register := make(chan chan *etree.Document)
	deregister := make(chan chan *etree.Document)
	return SvgUpdater{
		resultChan:     resultChan,
		origSvg:        origSvg,
		allUpdateRules: allUpdateRules,
		adhocRequest:   adhocRequest,
		register:       register,
		deregister:     deregister,
	}
}

func (s *SvgUpdater) Register() chan *etree.Document {
	observer := make(chan *etree.Document)
	s.register <- observer
	return observer
}

func (s *SvgUpdater) DeRegister(observer chan *etree.Document) {
	s.deregister <- observer
	return
}

func (s *SvgUpdater) Run(ctx context.Context) {
	log.Println("[svgupdater] Started up")
	var actualUpdateRules []svgmanip.Target
	svg := s.origSvg.Copy()
	for {
		select {
		case result := <-s.resultChan:
			if result.Value {
				for i, v := range s.allUpdateRules {
					if v.ID == result.ID {
						actualUpdateRules = append(actualUpdateRules, svgmanip.Target{ID: result.ID, Fill: s.allUpdateRules[i].Fill, LastChange: "Last change at: " + time.Now().Format(time.RFC1123)})
					}
				}
			} else {
				for i, v := range actualUpdateRules {
					if v.ID == result.ID {
						actualUpdateRules = append(actualUpdateRules[:i], actualUpdateRules[i+1:]...)
					}
				}
			}
			svg = s.origSvg.Copy()
			if err := svgmanip.UpdateDoc(svg, actualUpdateRules); err != nil {
				log.Println("[svgupdater] Error during SVG update:", err)
			}
			for _, client := range s.observers {
				log.Println("[svgupdater] Sending update to observer:", client)
				client <- svg
			}
		case adhocChan := <-s.adhocRequest:
			fmt.Println("[svgupdater] Ad-hoc request served:", adhocChan)
			adhocChan <- svg
		case newObserver := <-s.register:
			s.observers = append(s.observers, newObserver)
			log.Println("[svgupdater] Registered observer", newObserver, ", total:", len(s.observers))
		case oldObserver := <-s.deregister:
			for i, v := range s.observers {
				if v == oldObserver {
					s.observers = append(s.observers[:i], s.observers[i+1:]...)
				}
			}
			close(oldObserver)
			log.Println("[svgupdater] Deregistered observer", oldObserver, ", total:", len(s.observers))
		case <-ctx.Done():
			log.Println("[svgupdater] Exiting")
			return
		}
	}
}

func (s *SvgUpdater) Get() (*etree.Document, error) {
	updateChan := make(chan *etree.Document)
	s.adhocRequest <- updateChan
	timeout := time.After(5 * time.Second)
	select {
	case <-timeout:
		break
	case svg := <-updateChan:
		return svg, nil
	}
	return nil, errors.New("[svgupdater] Timeout happened when tried to gather SVG content from updater")
}
