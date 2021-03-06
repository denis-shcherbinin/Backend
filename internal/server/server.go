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

func (s *Server) Run(cfg *config.HTTPConfig, handler http.Handler) error {

	s.httpServer = &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: cfg.MaxHeaderMegabytes << shift,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
