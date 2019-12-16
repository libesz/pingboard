package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/handlers"
)

type scheduledHandler struct {
	ID         string
	handler    handlers.TargetHandler
	lastResult bool
}

type ResultChange struct {
	ID    string
	Value bool
}

func Run(ctx context.Context, targets []config.Target, resultChange chan ResultChange) {
	scheduledHandlers := []scheduledHandler{}
	for _, target := range targets {
		log.Printf("Initialized scheduledHandler: %+v\n", target)
		scheduledHandlers = append(scheduledHandlers, scheduledHandler{ID: target.SvgID, handler: handlers.NewHandler(target.Handler, target.EndPoint)})
	}
	for _, handler := range scheduledHandlers {
		go handler.run(ctx, resultChange)
	}
}

func (s *scheduledHandler) run(ctx context.Context, resultChange chan ResultChange) {
	for {
		time.Sleep(5 * time.Second)
		result := s.handler.Run(ctx)
		if s.lastResult != result {
			log.Println("Update for ID: "+s.ID+", value:", result)
			resultChange <- ResultChange{ID: s.ID, Value: result}
			s.lastResult = result
		}
	}
}
