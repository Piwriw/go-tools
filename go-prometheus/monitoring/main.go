package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	// create a new mux server
	server := http.NewServeMux()
	// register a new handler for the /metrics endpoint
	server.Handle("/metrics", promhttp.Handler())
	// start an http server using the mux server
	http.ListenAndServe(":9001", server)

}
