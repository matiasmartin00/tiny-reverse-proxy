package healthcheck

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
)

type backendHealthy struct {
	url     string
	healthy bool
}

var statusBackends = make(map[string]bool)

func AsyncVerifier() {
	go startVerifier()
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
	log.Println("Verifying backends")

	var wg sync.WaitGroup
	results := make(chan backendHealthy, len(config.Config.Backends))

	for _, backend := range config.Config.Backends {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			results <- backendHealthy{url, isBackendHealthy(url)}
		}(backend)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		statusBackends[res.url] = res.healthy
	}

	log.Println("Backends verified")
}

func isBackendHealthy(url string) bool {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("%s/ping", url))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
