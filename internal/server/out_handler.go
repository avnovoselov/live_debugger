package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/avnovoselov/live_debugger/internal/dto"
)

type OutHandler struct {
	upg    upgrader
	logger *zap.Logger
	queue  queue[dto.LogDTO]
}

func NewOutHandler(queue queue[dto.LogDTO], upg upgrader, logger *zap.Logger) *OutHandler {
	return &OutHandler{
		upg:    upg,
		logger: logger,
		queue:  queue,
	}
}

// ServeHTTP - http.Handler interface implementation
func (h OutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		log     dto.LogDTO
		conn    *websocket.Conn
		err     error
		offset  uint64
		message []byte
	)

	if conn, err = h.wsUpgrade(w, r); err != nil {
		h.logger.Error(err.Error())
		return
	}

	fields := []zap.Field{
		zap.String("localAddress", conn.LocalAddr().String()),
		zap.String("remoteAddress", conn.RemoteAddr().String()),
	}

	h.logger.Info("New connection", fields...)

	//nolint:errcheck
	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()

	for {
		log, offset, err = h.queue.GetLast()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	if message, err = dto.EncodeJSON[dto.LogDTO](log); err != nil {
		err = errors.Join(InHandlerMarshallResponseError, err)

		return
	}

	if err = conn.WriteMessage(websocket.TextMessage, message); err != nil {
		err = errors.Join(InHandlerSendResponseError)

		return
	}
	fmt.Println(offset)

}

// wsUpgrade - upgrades the HTTP server connection to the WebSocket protocol.
// Using https://github.com/gorilla/websocket
func (h OutHandler) wsUpgrade(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, err error) {
	if conn, err = h.upg.Upgrade(w, r, nil); err != nil {
		err = errors.Join(InHandlerWsUpgradeError, err)
	}

	return
}
