package loadbalancer

import "github.com/matiasmartin00/tiny-reverse-proxy/config"

var current uint64

func NextBackend() string {
	current = (current + 1) % uint64(len(config.Backends))
	return config.Backends[current]
}
