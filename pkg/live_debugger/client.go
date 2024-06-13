package live_debugger

import (
	"github.com/gorilla/websocket"
	netURL "net/url"
)

type Client struct {
	scheme     string
	addr       string
	inLocation string
}

func NewClient(scheme, addr, inLocation string) *Client {
	return &Client{
		scheme:     scheme,
		addr:       addr,
		inLocation: inLocation,
	}
}

func (c Client) Send(dto LogDTO) (err error) {
	var (
		url        netURL.URL
		connection *websocket.Conn
		message    []byte
	)

	url = netURL.URL{
		Scheme: c.scheme,
		Host:   c.addr,
		Path:   c.inLocation,
	}

	if connection, _, err = websocket.DefaultDialer.Dial(url.String(), nil); err != nil {
		return err
	}
	if message, err = EncodeJSON(dto); err != nil {
		return err
	}

	//nolint:errcheck
	//goland:noinspection GoUnhandledErrorResult
	defer connection.Close()

	return connection.WriteMessage(websocket.TextMessage, message)
}
