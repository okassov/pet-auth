package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func New(handler http.Handler, serverPort string) *Server {

	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         serverPort,
	}

	s := &Server{
		server:          httpServer,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	s.server.ListenAndServe()

	return s
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
