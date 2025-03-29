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
	config.LoadConfig()
	logger.InitLogger()
	go config.WatchConfig()
	healthcheck.AsyncVerifier()
	loadbalancer.InitConnections()
	log.Println("Starting server on :8080")
	server.Server()
}
