package server

import (
	"context"
	"github.com/PolyProjectOPD/Backend/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Config, handler http.Handler) error {
	httpConfig := cfg.HTTP

	s.httpServer = &http.Server{
		Addr:           ":" + httpConfig.Port,
		Handler:        handler,
		MaxHeaderBytes: httpConfig.MaxHeaderMegabytes << 20,
		ReadTimeout:    httpConfig.ReadTimeout,
		WriteTimeout:   httpConfig.WriteTimeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
