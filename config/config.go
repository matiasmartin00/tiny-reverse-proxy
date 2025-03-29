package config

import "log"

var Backends = []string{
	"http://localhost:5001",
	"http://localhost:5002",
	"http://localhost:5003",
}

func LoadConfig() {
	// Load configuration from file
	log.Println("Configuration loaded")
}
