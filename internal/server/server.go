package server

import (
	"net/http"

	"github.com/avnovoselov/live_debugger/internal/server/internal"
)

type Server struct {
	version     string
	inLocation  string
	outLocation string
	address     string
	inHandler   httpHandler
	outHandler  httpHandler
}

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
	}
}

func (s *Server) Run() error {
	http.Handle(s.inLocation, s.inHandler)
	http.Handle(s.outLocation, s.outHandler)

	return http.ListenAndServe(s.address, nil)
}
