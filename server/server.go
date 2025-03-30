package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
	"github.com/matiasmartin00/tiny-reverse-proxy/proxy"
)

type Server interface {
	StartServer() error
}

type server struct {
	proxy    proxy.ReverseProxy
	server   *http.Server
	serveMux *http.ServeMux
}

var instance Server
var once sync.Once

func (s *server) StartServer() error {
	s.serveMux.HandleFunc("/", s.proxy.ReverseProxyHandler)
	logger.GetLogger().Info("Listening on :8080")
	return s.server.ListenAndServe()
}

func NewServer(cfg config.Config, p proxy.ReverseProxy) Server {
	once.Do(func() {
		mux := http.NewServeMux()
		s := &server{
			proxy:    p,
			serveMux: mux,
			server: &http.Server{
				Addr:           fmt.Sprintf(":%d", cfg.GetServerPort()),
				Handler:        mux,
				ReadTimeout:    cfg.GetServerIdleTimeout(),
				WriteTimeout:   cfg.GetServerIdleTimeout(),
				IdleTimeout:    cfg.GetServerIdleTimeout(),
				MaxHeaderBytes: cfg.GetServerMaxHeaderBytes(),
			},
		}
		instance = s
		instance.StartServer()
	})
	return instance
}
