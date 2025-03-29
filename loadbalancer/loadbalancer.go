package loadbalancer

import (
	"net/http"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

func GetNextBackend(r *http.Request) string {
	logger.Debug("Load Balancer Strategy: ", config.Config.LB.Strategy)
	switch config.Config.LB.Strategy {
	case "round_robin":
		return getNextRoundRobinBackend()
	case "weighted":
		return getWeightedBackend()
	case "least_connections":
		return getLeastConnectionsBackend()
	case "ip_hash":
		return getIPHashBackend(r.RemoteAddr)
	default:
		logger.Debug("Load Balancer Strategy not found, using round robin")
		return getNextRoundRobinBackend()
	}
}
