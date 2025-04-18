package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Config interface {
	GetAllBackends() []Backend
	GetBackendsForPath(path string) ([]Backend, error)
	GetLoggingLevel() string
	GetLoadBalancerStrategy() string
	GetServerPort() int
	GetServerReadTimeout() time.Duration
	GetServerWriteTimeout() time.Duration
	GetServerIdleTimeout() time.Duration
	GetServerMaxHeaderBytes() int
	loadConfig()
	watchConfig()
}

type Backend interface {
	GetURL() string
	GetHealthPath() string
	GetHealthTimeout() time.Duration
	GetWeight() int
}

type server struct {
	Port           int           `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"read-timeout"`
	WriteTimeout   time.Duration `yaml:"write-timeout"`
	IdleTimeout    time.Duration `yaml:"idle-timeout"`
	MaxHeaderBytes int           `yaml:"max-header-bytes"`
}

type logging struct {
	Level string `yaml:"level"`
}

type backendHealth struct {
	Path   string        `yaml:"path"`
	Timout time.Duration `yaml:"timeout"`
}

type backend struct {
	URL    string        `yaml:"url"`
	Health backendHealth `yaml:"health"`
	Weight int           `yaml:"weight"`
}

func (b *backend) GetURL() string {
	return b.URL
}

func (b *backend) GetHealthPath() string {
	return b.Health.Path
}

func (b *backend) GetHealthTimeout() time.Duration {
	return b.Health.Timout
}

func (b *backend) GetWeight() int {
	return b.Weight
}

type loadBalancer struct {
	Strategy string `yaml:"strategy"`
}

type configurationFile struct {
	Logging logging              `yaml:"logging"`
	Routes  map[string][]backend `yaml:"routes"`
	LB      loadBalancer         `yaml:"loadbalancer"`
	Server  server               `yaml:"server"`
}

type configuration struct {
	file     string
	cf       configurationFile
	mutex    *sync.Mutex
	backends []Backend
}

var instance Config
var once sync.Once

func (c *configuration) GetAllBackends() []Backend {
	return c.backends
}

func (c *configuration) GetLoggingLevel() string {
	return c.cf.Logging.Level
}

func (c *configuration) GetLoadBalancerStrategy() string {
	return c.cf.LB.Strategy
}

func (c *configuration) GetBackendsForPath(path string) ([]Backend, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for route, backends := range c.cf.Routes {
		if len(path) >= len(route) && path[:len(route)] == route {
			var result []Backend
			for i := range backends {
				result = append(result, &backends[i])
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("no backends found for path: %s", path)
}

func (c *configuration) GetServerPort() int {
	return c.cf.Server.Port
}

func (c *configuration) GetServerReadTimeout() time.Duration {
	return c.cf.Server.ReadTimeout
}

func (c *configuration) GetServerWriteTimeout() time.Duration {
	return c.cf.Server.WriteTimeout
}

func (c *configuration) GetServerIdleTimeout() time.Duration {
	return c.cf.Server.IdleTimeout
}

func (c *configuration) GetServerMaxHeaderBytes() int {
	return c.cf.Server.MaxHeaderBytes
}

func (c *configuration) loadConfig() {
	// Load configuration from file
	log.Println("Loading configuration")

	file, err := os.Open(c.file)

	if err != nil {
		log.Fatalf("Error opening configuration file: %v", err)
	}

	defer file.Close()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	decoder := yaml.NewDecoder(file)
	c.cf = configurationFile{}

	err = decoder.Decode(&c.cf)
	if err != nil {
		log.Fatalf("Error decoding configuration file: %v", err)
	}

	c.loadBackends()
	log.Println("Configuration loaded")
}

func (c *configuration) loadBackends() {
	c.backends = []Backend{}
	for _, backends := range c.cf.Routes {
		for i := range backends {
			c.backends = append(c.backends, &backends[i])
		}
	}
}

func (c *configuration) watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}

	defer watcher.Close()

	err = watcher.Add(c.file)
	if err != nil {
		log.Fatalf("Error adding watcher: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Println("Configuration file modified")
				time.Sleep(1 * time.Second) // Wait for the file to be written
				c.loadConfig()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

func New() Config {
	once.Do(func() {
		instance = &configuration{
			file:  "config.yaml",
			mutex: &sync.Mutex{},
		}
		instance.loadConfig()
		go instance.watchConfig()
	})
	return instance
}

func GetConfig() Config {
	return instance
}
