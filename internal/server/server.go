package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/avnovoselov/live_debugger/internal"
	"github.com/avnovoselov/live_debugger/internal/request"
)

var (
	serverError = errors.New("server error")
)

const (
	tempDirectoryTemplate = "live_debugger_%s"
)

type Server struct {
	version internal.Version
	queue   []request.Request
}

func NewServer(version internal.Version) *Server {
	return &Server{
		version: version,
		queue:   make([]request.Request, 100),
	}
}

func (s *Server) Run(ctx context.Context) {
	http.HandleFunc("/log", func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			_, _ = res.Write([]byte("read all error"))
		}

		dto := request.Request{}
		err = json.Unmarshal(body, &dto)
		if err != nil {
			fmt.Println(err, string(body))
			_, _ = res.Write([]byte("json unmarshall error"))
		}
	})

	err := http.ListenAndServe(":8088", nil)
	fmt.Println(err)
}

func (s *Server) tempDirectory() string {
	return fmt.Sprintf(tempDirectoryTemplate, s.version)
}
