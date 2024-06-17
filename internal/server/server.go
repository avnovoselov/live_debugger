package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/avnovoselov/live_debugger/internal/configuration"
	"github.com/avnovoselov/live_debugger/internal/server/internal"
)

// Server - WebSocket server having two location
type Server struct {
	// inLocation - incoming logging stream. Debugging system writes logs this location
	inLocation string
	// outLocation - outgoing logging stream. Log readers get logs this location
	outLocation string

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
func NewServer(configuration configuration.Server, inHandler httpHandler, outHandler httpHandler) *Server {
	inLocation := internal.NormalizeLocation(configuration.InLocation)
	outLocation := internal.NormalizeLocation(configuration.OutLocation)

	address := internal.BuildAddress(configuration.IP, configuration.Port)

	return &Server{
		address:     address,
		inLocation:  inLocation,
		outLocation: outLocation,
		inHandler:   inHandler,
		outHandler:  outHandler,
		server:      &http.Server{Addr: address},
	}
}

// Run - configure server and start listening http connections
func (s *Server) Run() {
	log.Debug().Msg("Run server")

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
		log.Error().Err(err).Msg("shutdown server error")
	}
}

// serve - start listening http connection
func (s *Server) serve(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Err(err).Msg("server.ListenAndServe error")
	}
}
