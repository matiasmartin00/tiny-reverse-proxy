package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var Config configuration

type configuration struct {
	Backends []string `yaml:"backends"`
}

func LoadConfig() {
	// Load configuration from file
	log.Println("Loading configuration")

	file, err := os.Open("config.yaml")

	if err != nil {
		log.Fatalf("Error opening configuration file: %v", err)
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	Config = configuration{}

	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalf("Error decoding configuration file: %v", err)
	}

	log.Println("Configuration loaded")
}
