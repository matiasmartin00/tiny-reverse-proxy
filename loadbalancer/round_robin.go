package loadbalancer

import (
	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

type RoundRobin interface {
	GetNextRoundRobinBackend(backends []config.Backend) string
	GetWeightedBackend(backends []config.Backend) string
}

type roundRobinImpl struct {
	verifier healthcheck.Verifier
	current  uint64
}

// RoundRobin algorithm will return the next available backend in the list
func (rb *roundRobinImpl) GetNextRoundRobinBackend(backends []config.Backend) string {
	logger.GetLogger().Debug("Round Robin Load Balancer")
	activeServers := make([]string, 0, len(backends))

	for _, backend := range backends {
		if rb.verifier.IsNotBackendHealthy(backend.GetURL()) {
			continue
		}
		activeServers = append(activeServers, backend.GetURL())
	}

	if len(activeServers) == 0 {
		return ""
	}

	rb.current = (rb.current + 1) % uint64(len(activeServers))
	return activeServers[rb.current]
}

// Weighted Round Robin will return the next available backend in the list based on the weight
func (rb *roundRobinImpl) GetWeightedBackend(backends []config.Backend) string {
	logger.GetLogger().Debug("Weighted Load Balancer")
	weightedList := []string{}

	for _, backend := range backends {
		if rb.verifier.IsNotBackendHealthy(backend.GetURL()) {
			continue
		}

		for i := 0; i < backend.GetWeight(); i++ {
			weightedList = append(weightedList, backend.GetURL())
		}
	}

	if len(weightedList) == 0 {
		return ""
	}

	rb.current = (rb.current + 1) % uint64(len(weightedList))
	return weightedList[rb.current]
}

func newRoundRobin(verifier healthcheck.Verifier) RoundRobin {
	return &roundRobinImpl{
		verifier: verifier,
		current:  0,
	}
}
