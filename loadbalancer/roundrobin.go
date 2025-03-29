package loadbalancer

import (
	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
)

var current uint64

func NextBackend() string {
	activeServers := make([]string, 0, len(config.Config.Backends))

	for _, backend := range config.Config.Backends {
		if !healthcheck.IsBackendHealthy(backend) {
			continue
		}
		activeServers = append(activeServers, backend)
	}

	current = (current + 1) % uint64(len(activeServers))
	return activeServers[current]
}
