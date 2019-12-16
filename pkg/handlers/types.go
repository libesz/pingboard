package handlers

import (
	"context"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/handlers/ping"
)

type TargetHandler interface {
	Run(context.Context) bool
}

func NewHandler(handlerType config.Handler, target string) TargetHandler {
	//return ping.New(target)
	return &ping.PingConfig{Target: target}
	//switch handlerType {
	//case "ping":
	//}
}
