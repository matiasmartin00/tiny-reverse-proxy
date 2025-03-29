package main

import (
	"log"

	"github.com/matiasmartin00/tiny-reverse-proxy/server"
)

func main() {
	log.Println("Starting server on :8080")
	server.Server()
}
