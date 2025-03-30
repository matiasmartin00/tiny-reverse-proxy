package loadbalancer

import (
	"hash/fnv"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

func getIPHashBackend(clientIP string, backends []config.Backend) string {
	logger.GetLogger().Debug("IP Hash Load Balancer")
	activeServers := []string{}

	for _, backend := range backends {
		if healthcheck.GetVerifier().IsNotBackendHealthy(backend.GetURL()) {
			continue
		}
		activeServers = append(activeServers, backend.GetURL())
	}

	if len(activeServers) == 0 {
		return ""
	}

	hash := fnv.New32a()
	hash.Write([]byte(clientIP))
	index := int(hash.Sum32()) % len(activeServers)
	return activeServers[index]
}
