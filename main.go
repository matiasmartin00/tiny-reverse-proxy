package main

import (
	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/loadbalancer"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
	"github.com/matiasmartin00/tiny-reverse-proxy/proxy"
	"github.com/matiasmartin00/tiny-reverse-proxy/server"
)

func main() {
	config := config.New()
	logger.NewLogger(config)
	verifier := healthcheck.NewVerifier(config)
	lb := loadbalancer.NewLoadBalancer(config, verifier)
	proxy := proxy.NewReverseProxy(lb)
	server := server.NewServer(proxy)
	server.StartServer()
}
