package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/avnovoselov/live_debugger/internal/server/internal"
)

// Server - WebSocket server having two location
type Server struct {
	// inLocation - incoming logging stream. Debugging system writes logs this location
	inLocation string
	// outLocation - outgoing logging stream. Log readers get logs this location
	outLocation string

	// version - server version
	version string

	// address - host:port formatted TCP address for the server to listen on
	address string

	// inHandler - handler processes incoming logging stream
	inHandler httpHandler
	// outHandler - handler processes outgoing logging stream
	outHandler httpHandler

	// server - http.Server instance
	server *http.Server
}

// NewServer - Server instance constructor
func NewServer(version string, inLocation string, outLocation string, address string, inHandler httpHandler, outHandler httpHandler) *Server {
	inLocation = internal.NormalizeLocation(inLocation)
	outLocation = internal.NormalizeLocation(outLocation)

	return &Server{
		address:     address,
		version:     version,
		inLocation:  inLocation,
		outLocation: outLocation,
		inHandler:   inHandler,
		outHandler:  outHandler,
		server:      &http.Server{Addr: address},
	}
}

// Run - configure server and start listening http connections
func (s *Server) Run() {
	http.Handle(s.inLocation, s.inHandler)
	http.Handle(s.outLocation, s.outHandler)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go s.serve(wg)

	wg.Wait()
}

// Stop - stop server listening http connection
func (s *Server) Stop(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Println("shutdown server err: ", err)
	}
}

// serve - start listening http connection
func (s *Server) serve(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("ListenAndServe():", err)
	}
}
