package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

var Config configuration

type backend struct {
	URL        string `yaml:"url"`
	HealthPath string `yaml:"health-path"`
	Weight     int    `yaml:"weight"`
}

type loadbalancer struct {
	Strategy string `yaml:"strategy"`
}

type configuration struct {
	Backends []backend    `yaml:"backends"`
	LB       loadbalancer `yaml:"loadbalancer"`
}

var mutex = &sync.Mutex{}

func LoadConfig() {
	// Load configuration from file
	log.Println("Loading configuration")

	file, err := os.Open("config.yaml")

	if err != nil {
		log.Fatalf("Error opening configuration file: %v", err)
	}

	defer file.Close()

	mutex.Lock()
	defer mutex.Unlock()

	decoder := yaml.NewDecoder(file)
	Config = configuration{}

	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalf("Error decoding configuration file: %v", err)
	}

	log.Println("Configuration loaded")
}

func WatchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}

	defer watcher.Close()

	err = watcher.Add("config.yaml")
	if err != nil {
		log.Fatalf("Error adding watcher: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Println("Configuration file modified")
				time.Sleep(1 * time.Second) // Wait for the file to be written
				LoadConfig()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}
