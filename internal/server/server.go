package server

import (
	"context"
	"github.com/PolyProjectOPD/Backend/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

const (
	// Necessary shift to translate from megabytes to bytes
	shift = 20
)

func NewServer(cfg *config.HTTPConfig, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        handler,
			MaxHeaderBytes: cfg.MaxHeaderMegabytes << shift,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
