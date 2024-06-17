package server_test

import (
	"context"
	"github.com/avnovoselov/live_debugger/internal/configuration"
	"testing"
	"time"

	"github.com/avnovoselov/live_debugger/internal/server"
	mocks "github.com/avnovoselov/live_debugger/mocks/server"
)

func TestServer_Run(t *testing.T) {
	ctx := context.Background()

	inHandler := mocks.NewHttpHandler(t)
	outHandler := mocks.NewHttpHandler(t)

	cfg := configuration.Server{
		IP:          "127.0.0.1",
		Port:        "36768",
		InLocation:  "/in",
		OutLocation: "/out",
	}

	srv := server.NewServer(cfg, inHandler, outHandler)

	//nolint:errcheck
	//goland:noinspection GoUnhandledErrorResult
	go srv.Run()

	time.Sleep(100 * time.Millisecond)

	srv.Stop(ctx)
}
