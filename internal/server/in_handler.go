package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/avnovoselov/live_debugger/internal/request"
	"github.com/avnovoselov/live_debugger/internal/response"
	"github.com/avnovoselov/live_debugger/internal/util"
)

var inHandlerEB = util.ErrorBuilder(errors.New("in handler error"))

var (
	InHandlerUnexpectedMessageTypeError = inHandlerEB(errors.New("unexpected message type"))
	InHandlerUnmarshallRequestError     = inHandlerEB(UnmarshallRequestError)
	InHandlerMarshallResponseError      = inHandlerEB(MarshallResponseError)
	InHandlerSendResponseError          = inHandlerEB(errors.New("send response error"))
	InHandlerWsUpgradeError             = inHandlerEB(errors.New("ws upgrade error"))
)

type InHandler struct {
	queue  queue[request.LogRequest]
	upg    upgrader
	logger *zap.Logger
}

func NewInHandler(queue queue[request.LogRequest], upg upgrader, logger *zap.Logger) *InHandler {
	return &InHandler{
		queue:  queue,
		upg:    upg,
		logger: logger,
	}
}

func (h InHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		conn   *websocket.Conn
		req    request.LogRequest
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

		fields = append(fields, zap.String("message", string(message)))
		fields = append(fields, zap.Int("messageType", messageType))
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

func (h InHandler) wsUpgrade(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, err error) {
	if conn, err = h.upg.Upgrade(w, r, nil); err != nil {
		err = errors.Join(InHandlerWsUpgradeError, err)
	}

	return
}

func (h InHandler) decodeRequest(messageType int, message []byte) (req request.LogRequest, err error) {
	if websocket.TextMessage != messageType {
		err = InHandlerUnexpectedMessageTypeError

		return
	}

	if req, err = request.ParseJSON(message); err != nil {
		err = errors.Join(InHandlerUnmarshallRequestError, err)
	}

	return
}

func (h InHandler) handleMessage(connection *websocket.Conn, req request.LogRequest) (offset uint64, err error) {
	var message []byte

	offset = h.queue.Append(req)
	res := response.LogResponse{Offset: &offset}

	if message, err = response.EncodeJSON(res); err != nil {
		err = errors.Join(InHandlerMarshallResponseError, err)

		return
	}

	if err = connection.WriteMessage(websocket.TextMessage, message); err != nil {
		err = errors.Join(InHandlerSendResponseError)

		return
	}

	return
}
