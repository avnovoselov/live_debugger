package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/avnovoselov/live_debugger/internal/dto"
	"github.com/avnovoselov/live_debugger/internal/util"
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

// InHandler - handler processes incoming logging stream
type InHandler struct {
	queue  queue[dto.LogDTO]
	upg    upgrader
	logger *zap.Logger
}

// NewInHandler - InHandler constructor
func NewInHandler(queue queue[dto.LogDTO], upg upgrader, logger *zap.Logger) *InHandler {
	return &InHandler{
		queue:  queue,
		upg:    upg,
		logger: logger,
	}
}

// ServeHTTP - http.Handler interface implementation
func (h InHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		conn   *websocket.Conn
		req    dto.LogDTO
		err    error
		offset uint64
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

	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()

		fields = append(
			fields,
			zap.String("message", string(message)),
			zap.Int("messageType", messageType),
		)
		h.logger.Info("Receive message", fields...)

		if err != nil {
			h.logger.Error(err.Error(), fields...)
			break
		}
		if req, err = h.decodeRequest(messageType, message); err != nil {
			h.logger.Error(err.Error(), fields...)
			continue
		}
		if offset, err = h.handleMessage(conn, req); err != nil {
			h.logger.Error(err.Error(), fields...)
			break
		}

		fields = append(fields, zap.Uint64("offset", offset))
		h.logger.Info("Send message", fields...)
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
func (h InHandler) decodeRequest(messageType int, message []byte) (req dto.LogDTO, err error) {
	if websocket.TextMessage != messageType {
		err = InHandlerUnexpectedIncomingMessageTypeError

		return
	}

	if req, err = dto.ParseJSON[dto.LogDTO](message); err != nil {
		err = errors.Join(InHandlerUnmarshallRequestError, err)
	}

	return
}

// handleMessage - handle request message and send json encoded response
func (h InHandler) handleMessage(connection *websocket.Conn, req dto.LogDTO) (offset uint64, err error) {
	var message []byte

	offset = h.queue.Append(req)
	res := dto.LogCreatedDTO{Offset: &offset}

	if message, err = dto.EncodeJSON(res); err != nil {
		err = errors.Join(InHandlerMarshallResponseError, err)

		return
	}

	if err = connection.WriteMessage(websocket.TextMessage, message); err != nil {
		err = errors.Join(InHandlerSendResponseError)

		return
	}

	return
}
