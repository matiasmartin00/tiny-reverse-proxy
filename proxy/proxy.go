package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/matiasmartin00/tiny-reverse-proxy/loadbalancer"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

type ReverseProxy interface {
	ReverseProxyHandler(w http.ResponseWriter, r *http.Request)
}

type reverseProxyImpl struct {
	lb loadbalancer.LoadBalancer
}

func (rp *reverseProxyImpl) ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	target := rp.lb.GetNextBackend(r)

	if target == "" {
		logger.GetLogger().Error("Not available backends")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	logger.GetLogger().Debug("Proxying request to: ", target)

	rp.lb.IncrementConnection(target)
	defer rp.lb.DecrementConnection(target)

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

func NewReverseProxy(lb loadbalancer.LoadBalancer) ReverseProxy {
	return &reverseProxyImpl{
		lb: lb,
	}
}
