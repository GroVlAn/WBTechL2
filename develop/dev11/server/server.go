package server

import (
	"context"
	"dev11/config"
	"net/http"
)

type HTTPServer struct {
	httpServer *http.Server
}

func NewHttpServer(conf *config.Config, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		httpServer: &http.Server{
			Addr:              ":" + conf.ServConf.Port,
			Handler:           handler,
			MaxHeaderBytes:    conf.ServConf.MaxHeaderBytes,
			ReadHeaderTimeout: conf.ServConf.ReadHeaderTimeout,
			WriteTimeout:      conf.ServConf.WriteTimeout,
		},
	}
}

func (s *HTTPServer) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
