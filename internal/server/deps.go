package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type queue[Element any] interface {
	Append(element Element) uint64
	GetLast() (Element, uint64, error)
}

type httpHandler interface {
	http.Handler
}

type upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type responseWriter interface {
	http.ResponseWriter
}
