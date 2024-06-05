package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type queue[T any] interface {
	Append(element T) uint64
}

type httpHandler interface {
	http.Handler
}

type upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}
