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

var once sync.Once
var instance Verifier

type Verifier interface {
	IsNotBackendHealthy(url string) bool
	IsBackendHealthy(url string) bool
	startVerifier()
	verifyBackends()
	isBackendHealthy(url string, healthPath string) bool
}

type verifier struct {
	config         config.Config
	statusBackends map[string]bool
}

func (v *verifier) IsNotBackendHealthy(url string) bool {
	return !v.IsBackendHealthy(url)
}

func (v *verifier) IsBackendHealthy(url string) bool {
	return v.statusBackends[url]
}

func (v *verifier) startVerifier() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		v.verifyBackends()
	}
}

func (v *verifier) verifyBackends() {
	logger.GetLogger().Debug("Verifying backends")

	var wg sync.WaitGroup
	results := make(chan backendHealthy, len(config.GetConfig().GetAllBackends()))

	for _, backend := range config.GetConfig().GetAllBackends() {
		wg.Add(1)
		go func(url string, healthPath string) {
			defer wg.Done()
			results <- backendHealthy{url, v.isBackendHealthy(url, healthPath)}
		}(backend.GetURL(), backend.GetHealthPath())
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		v.statusBackends[res.url] = res.healthy
	}

	logger.GetLogger().Debug("Backends verified")
}

func (v *verifier) isBackendHealthy(url string, healthPath string) bool {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	requestURL := fmt.Sprintf("%s%s", url, healthPath)

	logger.GetLogger().Debug("Checking health of ", requestURL)

	resp, err := client.Get(requestURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.GetLogger().Debug("Backend is not healthy: ", url)
		return false
	}

	logger.GetLogger().Debug("Backend is healthy: ", url)
	return true
}

func NewVerifier(cfg config.Config) Verifier {
	once.Do(func() {
		instance = &verifier{
			config:         cfg,
			statusBackends: make(map[string]bool),
		}
		go instance.startVerifier()
		logger.GetLogger().Debug("Verifier initialized")
	})
	return instance
}

func GetVerifier() Verifier {
	return instance
}
