package middlewares

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

var cache = make(map[string]string)
var cacheMu sync.RWMutex
var cacheTTL = 5 * time.Minute

// Cache middleware caches responses for a certain time
func Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cacheMu.RLock()
		cachedResponse, exists := cache[r.URL.Path]
		cacheMu.RUnlock()

		if exists {
			w.Write([]byte(cachedResponse))
			return
		}

		// Capture response to cache it using NewResponseWriter
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		// Store response in cache if successful
		if rw.statusCode == http.StatusOK {
			cacheMu.Lock()
			cache[r.URL.Path] = rw.body.String()
			cacheMu.Unlock()

			// Expire cache after TTL
			time.AfterFunc(cacheTTL, func() {
				cacheMu.Lock()
				delete(cache, r.URL.Path)
				cacheMu.Unlock()
			})
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	body       *strings.Builder
	statusCode int
	flusher    http.Flusher
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	if rw.body == nil {
		rw.body = &strings.Builder{}
	}
	rw.body.Write(p)
	return rw.ResponseWriter.Write(p)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Flush() {
	if rw.flusher != nil {
		rw.flusher.Flush()
	}
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	if flusher, ok := w.(http.Flusher); ok {
		return &responseWriter{w, nil, 0, flusher}
	}
	return &responseWriter{w, nil, 0, nil}
}
