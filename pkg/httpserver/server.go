package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":8080"
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server *http.Server
}

func New(handler http.Handler) *Server {

	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &Server{
		server: httpServer,
	}

	s.server.ListenAndServe()

}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
