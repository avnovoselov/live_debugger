package main

import (
	"io/fs"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/avnovoselov/live_debugger/internal/configuration"
	"github.com/avnovoselov/live_debugger/internal/queue"
	"github.com/avnovoselov/live_debugger/internal/server"
	"github.com/avnovoselov/live_debugger/pkg/live_debugger"
)

func main() {
	var (
		fileSystem fs.FS
		cfg        configuration.Common

		err error
	)

	fileSystem = os.DirFS("./")
	parser := configuration.NewParser(fileSystem)

	if cfg, err = parser.Parse("configuration.toml"); err != nil {
		panic(err)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	upg := &websocket.Upgrader{}
	q := queue.NewQueue[live_debugger.LogDTO](cfg.Queue)

	inHandler := server.NewInHandler(q, upg)
	outHandler := server.NewOutHandler(q, upg, cfg.Server)

	s := server.NewServer(cfg.Server, inHandler, outHandler)
	s.Run()
}
