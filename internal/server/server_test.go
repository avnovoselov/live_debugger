package server_test

import (
	"testing"

	"github.com/avnovoselov/live_debugger/internal/server"
	mocks "github.com/avnovoselov/live_debugger/mocks/server"
)

func TestServer_Run(t *testing.T) {
	version := "1.0.0"
	inLocation := "/in"
	outLocation := "/out"
	address := "127.0.0.1:36768"
	inHandler := mocks.NewHttpHandler(t)
	outHandler := mocks.NewHttpHandler(t)

	go func() {
		srv := server.NewServer(version, inLocation, outLocation, address, inHandler, outHandler)
		_ = srv.Run()
	}()
}
