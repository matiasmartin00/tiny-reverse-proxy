package loadbalancer

import (
	"net/http"
	"sync"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/healthcheck"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

type LoadBalancer interface {
	GetNextBackend(r *http.Request) string
	IncrementConnection(backend string)
	DecrementConnection(backend string)
}

type loadBalancerImpl struct {
	config           config.Config
	roundRobin       RoundRobin
	leastConnections LeastConnections
	ipHash           IPHash
}

var once sync.Once
var instance LoadBalancer

func (lb *loadBalancerImpl) GetNextBackend(r *http.Request) string {
	logger.GetLogger().Debug("Load Balancer Strategy: ", lb.config.GetLoadBalancerStrategy())
	backends, err := lb.config.GetBackendsForPath(r.URL.Path)

	if err != nil {
		logger.GetLogger().Error("Error getting backends for path: ", r.URL.Path, " - ", err)
		return ""
	}

	if len(backends) == 0 {
		logger.GetLogger().Error("No backends available for path: ", r.URL.Path)
		return ""
	}

	switch lb.config.GetLoadBalancerStrategy() {
	case "round_robin":
		return lb.roundRobin.GetNextRoundRobinBackend(backends)
	case "weighted":
		return lb.roundRobin.GetWeightedBackend(backends)
	case "least_connections":
		return lb.leastConnections.GetLeastConnectionsBackend(backends)
	case "ip_hash":
		return lb.ipHash.GetIPHashBackend(r.RemoteAddr, backends)
	default:
		logger.GetLogger().Debug("Load Balancer Strategy not found, using round robin")
		return lb.roundRobin.GetNextRoundRobinBackend(backends)
	}
}

func (lb *loadBalancerImpl) IncrementConnection(backend string) {
	lb.leastConnections.IncrementConnection(backend)
}

func (lb *loadBalancerImpl) DecrementConnection(backend string) {
	lb.leastConnections.DecrementConnection(backend)
}

func NewLoadBalancer(config config.Config, verifier healthcheck.Verifier) LoadBalancer {
	once.Do(func() {
		instance = &loadBalancerImpl{
			config:           config,
			roundRobin:       newRoundRobin(verifier),
			leastConnections: newLeastConnections(verifier),
			ipHash:           newIPHash(verifier),
		}
	})

	return instance
}

func GetLoadBalancer() LoadBalancer {
	return instance
}
