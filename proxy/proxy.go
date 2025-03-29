package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/matiasmartin00/tiny-reverse-proxy/loadbalancer"
)

var strategy = "round_robin"

func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	target := getTarget(r)

	if target == "" {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	loadbalancer.IncrementConnection(target)
	defer loadbalancer.DecrementConnection(target)

	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

func getTarget(r *http.Request) string {
	switch strategy {
	case "round_robin":
		return loadbalancer.GetNextRoundRobinBackend()
	case "weighted":
		return loadbalancer.GetWeightedBackend()
	case "least_connections":
		return loadbalancer.GetLeastConnectionsBackend()
	case "ip_hash":
		return loadbalancer.GetIPHashBackend(r.RemoteAddr)
	default:
		return loadbalancer.GetNextRoundRobinBackend()
	}
}
