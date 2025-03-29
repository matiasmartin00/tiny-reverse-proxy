package loadbalancer

import (
	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
)

var backendConnections = make(map[string]int)

func InitConnections() {
	for _, backend := range config.Config.Backends {
		backendConnections[backend.URL] = 0
	}
}

func IncrementConnection(backend string) {
	backendConnections[backend]++
}

func DecrementConnection(backend string) {
	backendConnections[backend]--
}

func GetLeastConnectionsBackend() string {
	minConnections := int(^uint(0) >> 1)
	var minConnectionsBackend string

	for _, backend := range config.Config.Backends {
		if healthcheck.IsNotBackendHealthy(backend.URL) {
			continue
		}

		if backendConnections[backend.URL] < minConnections {
			minConnections = backendConnections[backend.URL]
			minConnectionsBackend = backend.URL
		}
	}

	return minConnectionsBackend
}
