package server

import (
	"log"
	"net/http"

	"github.com/matiasmartin00/tiny-reverse-proxy/logger"
	"github.com/matiasmartin00/tiny-reverse-proxy/proxy"
)

func Server() {
	http.HandleFunc("/", proxy.ReverseProxyHandler)
	logger.GetLogger().Info("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
