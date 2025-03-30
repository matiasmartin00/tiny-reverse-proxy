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
	isBackendHealthy(backend config.Backend) bool
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
		go func(backend config.Backend) {
			defer wg.Done()
			results <- backendHealthy{backend.GetURL(), v.isBackendHealthy(backend)}
		}(backend)
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

func (v *verifier) isBackendHealthy(backend config.Backend) bool {
	client := http.Client{
		Timeout: backend.GetHealthTimeout(),
	}
	requestURL := fmt.Sprintf("%s%s", backend.GetURL(), backend.GetHealthPath())

	logger.GetLogger().Debug("Checking health of ", requestURL)

	resp, err := client.Get(requestURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.GetLogger().Debug("Backend is not healthy: ", backend.GetURL())
		return false
	}

	logger.GetLogger().Debug("Backend is healthy: ", backend.GetURL())
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
