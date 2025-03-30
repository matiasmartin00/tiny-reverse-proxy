package main

import (
	"log"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
	"github.com/matiasmartin00/tiny-reverse-proxy/server"
)

func main() {
	config := config.New()
	logger.NewLogger(config)
	healthcheck.NewVerifier(config)
	log.Println("Starting server on :8080")
	server.Server()
}
