package loadbalancer

import (
	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

type LeastConnections interface {
	GetLeastConnectionsBackend(backends []config.Backend) string
	IncrementConnection(backend string)
	DecrementConnection(backend string)
}

type leastConnectionsImpl struct {
	backendConnections map[string]int
	verifier           healthcheck.Verifier
}

func (lc *leastConnectionsImpl) IncrementConnection(backend string) {
	lc.backendConnections[backend]++
}

func (lc *leastConnectionsImpl) DecrementConnection(backend string) {
	lc.backendConnections[backend]--
}

func (lc *leastConnectionsImpl) GetLeastConnectionsBackend(backends []config.Backend) string {
	logger.GetLogger().Debug("Least Connections Load Balancer")
	minConnections := int(^uint(0) >> 1)
	var minConnectionsBackend string

	for _, backend := range backends {
		if lc.verifier.IsNotBackendHealthy(backend.GetURL()) {
			continue
		}

		if _, ok := lc.backendConnections[backend.GetURL()]; !ok {
			lc.backendConnections[backend.GetURL()] = 0
		}

		if lc.backendConnections[backend.GetURL()] < minConnections {
			minConnections = lc.backendConnections[backend.GetURL()]
			minConnectionsBackend = backend.GetURL()
		}
	}

	return minConnectionsBackend
}

func newLeastConnections(verifier healthcheck.Verifier) LeastConnections {
	return &leastConnectionsImpl{
		backendConnections: make(map[string]int),
		verifier:           verifier,
	}
}
