package scheduler

import (
	"context"
	"log"
	"sync"
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
	log.Println("[scheduler] Started up")
	scheduledHandlers := []scheduledHandler{}
	for _, target := range targets {
		log.Printf("[scheduler] Initialized scheduledHandler: %+v\n", target)
		scheduledHandlers = append(scheduledHandlers, scheduledHandler{ID: target.ID, handler: handlers.NewHandler(target.Method, target.EndPoint)})
	}
	var wg sync.WaitGroup
	for _, handler := range scheduledHandlers {
		wg.Add(1)
		go func(handler scheduledHandler) {
			handler.run(ctx, resultChange)
			wg.Done()
		}(handler)
	}
	<-ctx.Done()
	log.Println("[scheduler] Waiting handlers to return...")
	wg.Wait()
	log.Println("[scheduler] Exiting")
}

func (s *scheduledHandler) run(ctx context.Context, resultChange chan ResultChange) {
	for {
		select {
		case <-ctx.Done():
			log.Println("[scheduler] Handler for " + s.ID + " exiting")
			return
		case <-time.After(5 * time.Second):
			result := s.handler.Run(context.Background())
			if s.lastResult != result {
				log.Println("[scheduler] Update for ID: "+s.ID+", value:", result)
				resultChange <- ResultChange{ID: s.ID, Value: result}
				s.lastResult = result
			}
		}
	}
}
