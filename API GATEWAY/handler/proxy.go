package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ProxyHandler returns an HTTP handler that proxies requests to the specified backend service
func ProxyHandler(target string) http.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Could not parse target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		// Additional logic for request modification can go here
		proxy.ServeHTTP(w, r)
	}
}
