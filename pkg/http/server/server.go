package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/s-gurman/user-segmentation/config"
)

const (
	_defaultIdleTimeout     = 30 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type HTTPServer interface {
	Start() error
	Stop() error
}

type httpserver struct {
	server *http.Server
}

func New(cfg config.HTTPServerConfig, handler http.Handler) HTTPServer {
	return &httpserver{
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  _defaultIdleTimeout,
		},
	}
}

func (s *httpserver) Start() error {
	return s.server.ListenAndServe()
}

func (s *httpserver) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
