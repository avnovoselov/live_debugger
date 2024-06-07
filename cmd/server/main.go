package main

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/avnovoselov/live_debugger/internal/queue"
	"github.com/avnovoselov/live_debugger/internal/request"
	"github.com/avnovoselov/live_debugger/internal/server"
)

func main() {
	var (
		logger *zap.Logger
		err    error
	)

	upg := &websocket.Upgrader{}
	q := queue.NewQueue[request.LogRequest](1000)

	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer logger.Sync()

	inHandler := server.NewInHandler(q, upg, logger)
	outHandler := server.NewInHandler(q, upg, logger)

	s := server.NewServer("1.0.0", "/ws", "/out", "127.0.0.1:8080", inHandler, outHandler)
	s.Run()
}
