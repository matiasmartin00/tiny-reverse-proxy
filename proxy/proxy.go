package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/matiasmartin00/tiny-reverse-proxy/loadbalancer"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	target := loadbalancer.GetLoadBalancer().GetNextBackend(r)

	if target == "" {
		logger.GetLogger().Error("Not available backends")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	logger.GetLogger().Debug("Proxying request to: ", target)

	loadbalancer.GetLoadBalancer().IncrementConnection(target)
	defer loadbalancer.GetLoadBalancer().DecrementConnection(target)

	targetURL, err := url.Parse(target)
	if err != nil {
		logger.GetLogger().Error("Error parsing target URL: ", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)

	logger.GetLogger().Debug("Request proxied successfully. Target: ", target)
}
