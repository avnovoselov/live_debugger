package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/avnovoselov/live_debugger/internal/server"
	mocks "github.com/avnovoselov/live_debugger/mocks/server"
)

func TestServer_Run(t *testing.T) {
	ctx := context.Background()
	version := "1.0.0"
	inLocation := "/in"
	outLocation := "/out"
	address := "127.0.0.1:36768"
	inHandler := mocks.NewHttpHandler(t)
	outHandler := mocks.NewHttpHandler(t)

	srv := server.NewServer(version, inLocation, outLocation, address, inHandler, outHandler)
	//goland:noinspection GoUnhandledErrorResult
	go srv.Run()

	time.Sleep(100 * time.Millisecond)

	srv.Stop(ctx)
}
