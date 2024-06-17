package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/avnovoselov/live_debugger/internal/configuration"
	"github.com/avnovoselov/live_debugger/internal/util"
	"github.com/avnovoselov/live_debugger/pkg/live_debugger"
)

// OutHandler - handle outgoing logging streams
type OutHandler struct {
	upg   upgrader
	queue queue[live_debugger.LogDTO]

	amountOfSequentiallyErrorToBreak int

	throttleDurationMs        time.Duration
	sleepAfterErrorDurationMs time.Duration
}

// NewOutHandler - OutHandler constructor
func NewOutHandler(
	queue queue[live_debugger.LogDTO],
	upg upgrader,
	server configuration.Server,
) *OutHandler {
	return &OutHandler{
		upg:                              upg,
		queue:                            queue,
		amountOfSequentiallyErrorToBreak: server.AmountOfSequentiallyErrorToBreak,
		throttleDurationMs:               util.IntToMillisecond(server.ThrottleDurationMs),
		sleepAfterErrorDurationMs:        util.IntToMillisecond(server.SleepAfterErrorDurationMs),
	}
}

// ServeHTTP - http.Handler interface implementation
func (h OutHandler) ServeHTTP(w responseWriter, r *http.Request) {
	var (
		logDTO live_debugger.LogDTO
		conn   *websocket.Conn

		err                       error
		amountOfSequentiallyError int
		previousOffset            uint64
		currentOffset             uint64

		fields = make(map[string]any)
	)

	if conn, err = h.wsUpgrade(w, r); err != nil {
		log.Error().Err(err)
		return
	}

	if conn.LocalAddr() != nil {
		fields["localAddress"] = conn.LocalAddr().String()
	}
	if conn.RemoteAddr() != nil {
		fields["remoteAddress"] = conn.RemoteAddr().String()
	}

	log.Info().Fields(fields).Msg("New connection")

	//nolint:errcheck
	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()

	for {
		logDTO, currentOffset, err = h.queue.GetLast()
		if err != nil {
			time.Sleep(h.sleepAfterErrorDurationMs)
			continue
		}
		if currentOffset == previousOffset {
			time.Sleep(h.sleepAfterErrorDurationMs)
			continue
		}

		if err = h.sendMessage(conn, logDTO); err != nil {
			log.Error().Err(err)
			amountOfSequentiallyError += 1
		} else {
			previousOffset = currentOffset
			amountOfSequentiallyError = 0
		}

		if amountOfSequentiallyError > h.amountOfSequentiallyErrorToBreak {
			break
		}

		time.Sleep(h.throttleDurationMs)
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

// sendMessage - send message to client
func (h OutHandler) sendMessage(conn *websocket.Conn, log live_debugger.LogDTO) (err error) {
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
