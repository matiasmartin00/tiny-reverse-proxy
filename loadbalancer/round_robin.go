package loadbalancer

import (
	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

var current uint64

// RoundRobin algorithm will return the next available backend in the list
func getNextRoundRobinBackend(backends []config.Backend) string {
	logger.GetLogger().Debug("Round Robin Load Balancer")
	activeServers := make([]string, 0, len(backends))

	for _, backend := range backends {
		if healthcheck.GetVerifier().IsNotBackendHealthy(backend.GetURL()) {
			continue
		}
		activeServers = append(activeServers, backend.GetURL())
	}

	if len(activeServers) == 0 {
		return ""
	}

	current = (current + 1) % uint64(len(activeServers))
	return activeServers[current]
}

// Weighted Round Robin will return the next available backend in the list based on the weight
func getWeightedBackend(backends []config.Backend) string {
	logger.GetLogger().Debug("Weighted Load Balancer")
	weightedList := []string{}

	for _, backend := range backends {
		if healthcheck.GetVerifier().IsNotBackendHealthy(backend.GetURL()) {
			continue
		}

		for i := 0; i < backend.GetWeight(); i++ {
			weightedList = append(weightedList, backend.GetURL())
		}
	}

	if len(weightedList) == 0 {
		return ""
	}

	current = (current + 1) % uint64(len(weightedList))
	return weightedList[current]
}
