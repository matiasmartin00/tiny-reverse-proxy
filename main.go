package main

import (
	"log"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/server"
)

func main() {
	config.LoadConfig()
	healthcheck.AsyncVerifier()
	log.Println("Starting server on :8080")
	server.Server()
}
