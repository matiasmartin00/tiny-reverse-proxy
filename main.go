package main

import (
	"log"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/loadbalancer"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
	"github.com/matiasmartin00/tiny-reverse-proxy/server"
)

func main() {
	config := config.New()
	logger.NewLogger(config)
	verifier := healthcheck.NewVerifier(config)
	loadbalancer.NewLoadBalancer(config, verifier)
	log.Println("Starting server on :8080")
	server.Server()
}
