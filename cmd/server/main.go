package main

import (
	"fmt"
	"github.com/gorilla/websocket"

	"github.com/avnovoselov/live_debugger/internal/queue"
	"github.com/avnovoselov/live_debugger/internal/request"
	"github.com/avnovoselov/live_debugger/internal/server"
)

func main() {
	upg := &websocket.Upgrader{}
	q := queue.NewQueue[request.LogRequest](1000)

	inHandler := server.NewInHandler(q, upg)
	outHandler := server.NewInHandler(q, upg)

	s := server.NewServer("1.0.0", "/ws", "/out", "127.0.0.1:8080", inHandler, outHandler)
	err := s.Run()
	fmt.Println(err)
}
