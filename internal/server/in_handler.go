package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/avnovoselov/live_debugger/internal/util"
	"github.com/avnovoselov/live_debugger/pkg/live_debugger"
)

// inHandlerEB - inHandler error builder contains base error
var inHandlerEB = util.ErrorBuilder(errors.New("in handler error"))

var (
	InHandlerUnexpectedIncomingMessageTypeError = inHandlerEB(errors.New("unexpected incoming message type"))
	InHandlerUnmarshallRequestError             = inHandlerEB(UnmarshallRequestError)
	InHandlerMarshallResponseError              = inHandlerEB(MarshallResponseError)
	InHandlerSendResponseError                  = inHandlerEB(errors.New("send response error"))
	InHandlerWsUpgradeError                     = inHandlerEB(errors.New("ws upgrade error"))
)

// InHandler - handle incoming logging streams
type InHandler struct {
	queue queue[live_debugger.LogDTO]
	upg   upgrader
}

// NewInHandler - InHandler constructor
func NewInHandler(queue queue[live_debugger.LogDTO], upg upgrader) *InHandler {
	return &InHandler{
		queue: queue,
		upg:   upg,
	}
}

// ServeHTTP - http.Handler interface implementation
func (h InHandler) ServeHTTP(w responseWriter, r *http.Request) {
	var (
		conn   *websocket.Conn
		req    live_debugger.LogDTO
		err    error
		offset uint64

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

	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()

		fields["message"] = string(message)
		fields["messageType"] = messageType

		log.Debug().Fields(fields).Msg("Receive message")

		if err != nil {
			log.Error().Fields(fields).Err(err)
			break
		}
		if req, err = h.decodeRequest(messageType, message); err != nil {
			log.Error().Fields(fields).Err(err)
			continue
		}
		if offset, err = h.handleMessage(conn, req); err != nil {
			log.Error().Fields(fields).Err(err)
			break
		}

		fields["offset"] = offset
		log.Debug().Fields(fields).Msg("Message handled")
	}
}

// wsUpgrade - upgrades the HTTP server connection to the WebSocket protocol.
// Using https://github.com/gorilla/websocket
func (h InHandler) wsUpgrade(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, err error) {
	if conn, err = h.upg.Upgrade(w, r, nil); err != nil {
		err = errors.Join(InHandlerWsUpgradeError, err)
	}

	return
}

// decodeRequest - checks messageType and try to parse request body
func (h InHandler) decodeRequest(messageType int, message []byte) (req live_debugger.LogDTO, err error) {
	if websocket.TextMessage != messageType {
		err = InHandlerUnexpectedIncomingMessageTypeError

		return
	}

	if req, err = live_debugger.ParseJSON[live_debugger.LogDTO](message); err != nil {
		err = errors.Join(InHandlerUnmarshallRequestError, err)
	}

	return
}

// handleMessage - handle request message and send json encoded response
func (h InHandler) handleMessage(connection *websocket.Conn, req live_debugger.LogDTO) (offset uint64, err error) {
	var message []byte

	offset = h.queue.Append(req)
	res := live_debugger.LogCreatedDTO{Offset: &offset}

	if message, err = live_debugger.EncodeJSON(res); err != nil {
		err = errors.Join(InHandlerMarshallResponseError, err)

		return
	}

	if err = connection.WriteMessage(websocket.TextMessage, message); err != nil {
		err = errors.Join(InHandlerSendResponseError)

		return
	}

	return
}
