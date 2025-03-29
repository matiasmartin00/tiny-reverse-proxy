package healthcheck

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
)

type backendHealthy struct {
	url     string
	healthy bool
}

var statusBackends = make(map[string]bool)

func AsyncVerifier() {
	go startVerifier()
}

func IsNotBackendHealthy(url string) bool {
	return !IsBackendHealthy(url)
}

func IsBackendHealthy(url string) bool {
	return statusBackends[url]
}

func startVerifier() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		verifyBackends()
	}
}

func verifyBackends() {
	logger.Debug("Verifying backends")

	var wg sync.WaitGroup
	results := make(chan backendHealthy, len(config.Config.Backends))

	for _, backend := range config.Config.Backends {
		wg.Add(1)
		go func(url string, healthPath string) {
			defer wg.Done()
			results <- backendHealthy{url, isBackendHealthy(url, healthPath)}
		}(backend.URL, backend.HealthPath)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		statusBackends[res.url] = res.healthy
	}

	logger.Debug("Backends verified")
}

func isBackendHealthy(url string, healthPath string) bool {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	requestURL := fmt.Sprintf("%s%s", url, healthPath)

	logger.Debug("Checking health of ", requestURL)

	resp, err := client.Get(requestURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Debug("Backend is not healthy: ", url)
		return false
	}

	logger.Debug("Backend is healthy: ", url)
	return true
}
