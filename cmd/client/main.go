package main

import (
	"fmt"
	netURL "net/url"

	"github.com/gorilla/websocket"
	"github.com/manifoldco/promptui"

	"github.com/avnovoselov/live_debugger/pkg/live_debugger"
)

const (
	defaultSchema      = "ws"
	defaultAddress     = "127.0.0.1"
	defaultPort        = "8080"
	defaultLocation    = "/in"
	defaultType        = "message"
	defaultSource      = "live_debugger_client"
	defaultFingerprint = "00000000-0000-0000-0000-000000000000"
)

func main() {
	var (
		address    string
		inLocation string
		port       string

		l int
		t string
		m string
		s string
		f string

		connection *websocket.Conn
		err        error

		message []byte
	)

	fmt.Println("\nServer")

	addressPrompt, portPrompt, inLocationPrompt := provideServerPrompts()

	address, _ = addressPrompt.Run()
	port, _ = portPrompt.Run()
	inLocation, _ = inLocationPrompt.Run()

	url := netURL.URL{Scheme: defaultSchema, Host: fmt.Sprintf("%s:%s", address, port), Path: inLocation}
	if connection, _, err = websocket.DefaultDialer.Dial(url.String(), nil); err != nil {
		fmt.Println(err.Error())
		return
	}

	messageNumber := 0
	for {
		fmt.Printf("\nMessage %d\n", messageNumber)
		messageNumber += 1
		levelSelect, typePrompt, messagePrompt, sourcePrompt, fingerprintPrompt := provideMessagePrompts()

		l, _, _ = levelSelect.Run()
		t, _ = typePrompt.Run()
		m, _ = messagePrompt.Run()
		s, _ = sourcePrompt.Run()
		f, _ = fingerprintPrompt.Run()

		dto := live_debugger.LogDTO{
			Level:       live_debugger.Level(l),
			Type:        live_debugger.Type(t),
			Message:     live_debugger.Message(m),
			Source:      live_debugger.Source(s),
			Fingerprint: live_debugger.Fingerprint(f),
		}

		if message, err = live_debugger.EncodeJSON(dto); err != nil {
			fmt.Println(err.Error())
			return
		}

		if err = connection.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func provideServerPrompts() (addressPrompt promptui.Prompt, portPrompt promptui.Prompt, inLocationPrompt promptui.Prompt) {
	addressPrompt = promptui.Prompt{
		Label:     "Address",
		Default:   defaultAddress,
		AllowEdit: true,
	}
	portPrompt = promptui.Prompt{
		Label:     "Port",
		Default:   defaultPort,
		AllowEdit: true,
	}
	inLocationPrompt = promptui.Prompt{
		Label:     "Location",
		Default:   defaultLocation,
		AllowEdit: true,
	}

	return
}

func provideMessagePrompts() (levelSelect promptui.Select, typePrompt promptui.Prompt, messagePrompt promptui.Prompt, sourcePrompt promptui.Prompt, fingerprintPrompt promptui.Prompt) {
	levelSelect = promptui.Select{
		Label: "Level",
		Items: []string{
			"Debug",
			"Info",
			"Warning",
			"Error",
		},
		CursorPos: 1,
	}
	typePrompt = promptui.Prompt{
		Label:     "Type",
		Default:   defaultType,
		AllowEdit: true,
	}
	messagePrompt = promptui.Prompt{
		Label:     "Message",
		AllowEdit: true,
	}
	sourcePrompt = promptui.Prompt{
		Label:     "Source",
		Default:   defaultSource,
		AllowEdit: true,
	}
	fingerprintPrompt = promptui.Prompt{
		Label:     "Fingerprint",
		Default:   defaultFingerprint,
		AllowEdit: true,
	}

	return
}
