package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/matiasmartin00/tiny-reverse-proxy/loadbalancer"
)

func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	target := loadbalancer.NextBackend()
	if target == "" {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}
