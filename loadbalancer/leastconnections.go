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
	logger.Debug("Least Connections Load Balancer")
	minConnections := int(^uint(0) >> 1)
	var minConnectionsBackend string

	for _, backend := range backends {
		if healthcheck.IsNotBackendHealthy(backend.URL) {
			continue
		}

		if _, ok := backendConnections[backend.URL]; !ok {
			backendConnections[backend.URL] = 0
		}

		if backendConnections[backend.URL] < minConnections {
			minConnections = backendConnections[backend.URL]
			minConnectionsBackend = backend.URL
		}
	}

	return minConnectionsBackend
}
