package server

import (
	"context"
	"fmt"
	"github.com/avnovoselov/live_debugger/internal/server/internal"
	"net/http"
	"sync"
)

type Server struct {
	version     string
	inLocation  string
	outLocation string
	address     string
	inHandler   httpHandler
	outHandler  httpHandler
	server      *http.Server
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
		server:      &http.Server{Addr: address},
	}
}

func (s *Server) Run() {
	http.Handle(s.inLocation, s.inHandler)
	http.Handle(s.outLocation, s.outHandler)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go s.serve(wg)

	wg.Wait()
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Println("shutdown server err: ", err)
	}
}

func (s *Server) serve(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("ListenAndServe():", err)
	}
}
