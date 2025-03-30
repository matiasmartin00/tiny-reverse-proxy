package loadbalancer

import (
	"net/http"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

func GetNextBackend(r *http.Request) string {
	logger.GetLogger().Debug("Load Balancer Strategy: ", config.GetConfig().GetLoadBalancerStrategy())
	backends, err := config.GetConfig().GetBackendsForPath(r.URL.Path)

	if err != nil {
		logger.GetLogger().Error("Error getting backends for path: ", r.URL.Path, " - ", err)
		return ""
	}

	if len(backends) == 0 {
		logger.GetLogger().Error("No backends available for path: ", r.URL.Path)
		return ""
	}

	switch config.GetConfig().GetLoadBalancerStrategy() {
	case "round_robin":
		return getNextRoundRobinBackend(backends)
	case "weighted":
		return getWeightedBackend(backends)
	case "least_connections":
		return getLeastConnectionsBackend(backends)
	case "ip_hash":
		return getIPHashBackend(r.RemoteAddr, backends)
	default:
		logger.GetLogger().Debug("Load Balancer Strategy not found, using round robin")
		return getNextRoundRobinBackend(backends)
	}
}
