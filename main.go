package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var backends = []string{
	"http://localhost:5001",
	"http://localhost:5002",
	"http://localhost:5003",
}

var current uint64

func nextBackend() string {
	current = (current + 1) % uint64(len(backends))
	return backends[current]
}

func handler(w http.ResponseWriter, r *http.Request) {
	target := nextBackend()
	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
