package main

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/avnovoselov/live_debugger/internal/queue"
	"github.com/avnovoselov/live_debugger/internal/server"
	"github.com/avnovoselov/live_debugger/pkg/live_debugger"
)

func main() {
	var (
		logger *zap.Logger
		err    error
	)

	upg := &websocket.Upgrader{}
	q := queue.NewQueue[live_debugger.LogDTO](1000)

	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}

	//nolint:errcheck
	//goland:noinspection GoUnhandledErrorResult
	defer logger.Sync()

	inHandler := server.NewInHandler(q, upg, logger)
	outHandler := server.NewOutHandler(q, upg, logger)

	s := server.NewServer("1.0.0", "/in", "/out", "127.0.0.1:8080", inHandler, outHandler)
	s.Run()
}
