package loadbalancer

import (
	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

var backendConnections = make(map[string]int)

func IncrementConnection(backend string) {
	backendConnections[backend]++
}

func DecrementConnection(backend string) {
	backendConnections[backend]--
}

func getLeastConnectionsBackend(backends []config.Backend) string {
	logger.GetLogger().Debug("Least Connections Load Balancer")
	minConnections := int(^uint(0) >> 1)
	var minConnectionsBackend string

	for _, backend := range backends {
		if healthcheck.GetVerifier().IsNotBackendHealthy(backend.GetURL()) {
			continue
		}

		if _, ok := backendConnections[backend.GetURL()]; !ok {
			backendConnections[backend.GetURL()] = 0
		}

		if backendConnections[backend.GetURL()] < minConnections {
			minConnections = backendConnections[backend.GetURL()]
			minConnectionsBackend = backend.GetURL()
		}
	}

	return minConnectionsBackend
}
