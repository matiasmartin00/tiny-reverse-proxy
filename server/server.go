package server

import (
	"net/http"

	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
	"github.com/matiasmartin00/tiny-reverse-proxy/proxy"
)

type Server interface {
	StartServer() error
}

type server struct {
	proxy proxy.ReverseProxy
}

func (s *server) StartServer() error {
	http.HandleFunc("/", s.proxy.ReverseProxyHandler)
	logger.GetLogger().Info("Listening on :8080")
	return http.ListenAndServe(":8080", nil)
}

func NewServer(p proxy.ReverseProxy) Server {
	return &server{
		proxy: p,
	}
}
