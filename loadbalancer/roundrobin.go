package loadbalancer

import (
	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

var current uint64

// RoundRobin algorithm will return the next available backend in the list
func getNextRoundRobinBackend() string {
	logger.Debug("Round Robin Load Balancer")
	activeServers := make([]string, 0, len(config.Config.Backends))

	for _, backend := range config.Config.Backends {
		if healthcheck.IsNotBackendHealthy(backend.URL) {
			continue
		}
		activeServers = append(activeServers, backend.URL)
	}

	if len(activeServers) == 0 {
		return ""
	}

	current = (current + 1) % uint64(len(activeServers))
	return activeServers[current]
}

// Weighted Round Robin will return the next available backend in the list based on the weight
func getWeightedBackend() string {
	logger.Debug("Weighted Load Balancer")
	weightedList := []string{}

	for _, backend := range config.Config.Backends {
		if healthcheck.IsNotBackendHealthy(backend.URL) {
			continue
		}

		for i := 0; i < backend.Weight; i++ {
			weightedList = append(weightedList, backend.URL)
		}
	}

	if len(weightedList) == 0 {
		return ""
	}

	current = (current + 1) % uint64(len(weightedList))
	return weightedList[current]
}
