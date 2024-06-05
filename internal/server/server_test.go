package server_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/server"
	mocks "github.com/avnovoselov/live_debugger/mocks/server"
)

func TestServer_Run(t *testing.T) {
	version := "1.0.0"
	inLocation := "/in"
	outLocation := "/out"
	address := "127.0.0.1:11000"
	inHandler := mocks.NewHttpHandler(t)
	outHandler := mocks.NewHttpHandler(t)

	inHandler.EXPECT().ServeHTTP(mock.Anything, mock.Anything)
	outHandler.EXPECT().ServeHTTP(mock.Anything, mock.Anything)

	go func() {
		srv := server.NewServer(version, inLocation, outLocation, address, inHandler, outHandler)
		_ = srv.Run()
	}()

	//goland:noinspection HttpUrlsUsage
	resp, err := http.Get("http://" + address + inLocation)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	//goland:noinspection HttpUrlsUsage
	resp, err = http.Get("http://" + address + outLocation)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
}
