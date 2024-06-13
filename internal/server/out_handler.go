package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/avnovoselov/live_debugger/pkg/live_debugger"
)

type OutHandler struct {
	upg    upgrader
	logger *zap.Logger
	queue  queue[live_debugger.LogDTO]
}

func NewOutHandler(queue queue[live_debugger.LogDTO], upg upgrader, logger *zap.Logger) *OutHandler {
	return &OutHandler{
		upg:    upg,
		logger: logger,
		queue:  queue,
	}
}

// ServeHTTP - http.Handler interface implementation
func (h OutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		log                live_debugger.LogDTO
		conn               *websocket.Conn
		err                error
		errorSequenceCount int
		previousOffset     uint64
		currentOffset      uint64
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
		time.Sleep(time.Second / 2)
		log, currentOffset, err = h.queue.GetLast()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if currentOffset == previousOffset {
			time.Sleep(time.Second)
			continue
		}

		if err = h.handleMessage(conn, log); err != nil {
			h.logger.Error(err.Error())
			errorSequenceCount += 1
		} else {
			previousOffset = currentOffset
			errorSequenceCount = 0
		}

		// todo configuration and logging
		if errorSequenceCount > 10 {
			break
		}
	}
}

// wsUpgrade - upgrades the HTTP server connection to the WebSocket protocol.
// Using https://github.com/gorilla/websocket
func (h OutHandler) wsUpgrade(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, err error) {
	if conn, err = h.upg.Upgrade(w, r, nil); err != nil {
		err = errors.Join(InHandlerWsUpgradeError, err)
	}

	return
}

func (h OutHandler) handleMessage(conn *websocket.Conn, log live_debugger.LogDTO) (err error) {
	var message []byte

	if message, err = live_debugger.EncodeJSON[live_debugger.LogDTO](log); err != nil {
		err = errors.Join(InHandlerMarshallResponseError, err)

		return
	}

	if err = conn.WriteMessage(websocket.TextMessage, message); err != nil {
		err = errors.Join(InHandlerSendResponseError)

		return
	}

	return
}
