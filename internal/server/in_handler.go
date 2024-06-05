package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"

	"github.com/avnovoselov/live_debugger/internal/request"
)

type InHandler struct {
	queue queue[request.LogRequest]
	upg   upgrader
}

func NewInHandler(queue queue[request.LogRequest], upg upgrader) *InHandler {
	return &InHandler{
		queue: queue,
		upg:   upg,
	}
}

func (h *InHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection, err := h.upg.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err.Error())
		return
	}

	//goland:noinspection GoUnhandledErrorResult
	defer connection.Close()

	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			fmt.Println("read:", err.Error())
			break
		}
		if mt != websocket.TextMessage {
			fmt.Println("skip message:", mt)
			continue
		}
		req, err := request.ParseJSON(message)
		if err != nil {
			fmt.Println("parse error:", err.Error())
			continue
		}

		offset := h.queue.Append(req)
		res := request.LogResponse{Offset: &offset}

		b, err := request.EncodeJSON(res)
		if err != nil {
			fmt.Println("encode error: ", err)
		}

		err = connection.WriteMessage(mt, b)
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}
