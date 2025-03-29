package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/matiasmartin00/tiny-reverse-proxy/loadbalancer"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	target := loadbalancer.GetNextBackend(r)

	if target == "" {
		logger.Error("Not available backends")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	logger.Debug("Proxying request to: ", target)

	loadbalancer.IncrementConnection(target)
	defer loadbalancer.DecrementConnection(target)

	targetURL, err := url.Parse(target)
	if err != nil {
		logger.Error("Error parsing target URL: ", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)

	logger.Debug("Request proxied successfully. Target: ", target)
}
