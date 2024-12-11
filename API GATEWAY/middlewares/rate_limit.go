package middlewares

import (
	"net/http"
	"sync"
	"time"
)

var (
	windowSize  = 60 * time.Second // Duration of the time window
	maxRequests = 5                // Maximum allowed requests per IP within the time window
	mu          sync.Mutex         // Mutex for thread-safe access
)

// requestCount tracks the timestamps of requests for each IP
var requestCount = make(map[string][]time.Time)

// SlidingWindowRateLimit applies rate limiting based on a sliding time window
func SlidingWindowRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		now := time.Now()

		mu.Lock() // Lock to ensure safe concurrent access
		defer mu.Unlock()

		// Retrieve or initialize the list of timestamps for the IP
		requests := requestCount[ip]
		if requests == nil {
			requests = []time.Time{}
		}

		// Filter out timestamps that are outside the sliding window
		windowStart := now.Add(-windowSize)
		validRequests := make([]time.Time, 0, len(requests))
		for _, timestamp := range requests {
			if timestamp.After(windowStart) {
				validRequests = append(validRequests, timestamp)
			}
		}

		// Check if the number of valid requests exceeds the limit
		if len(validRequests) >= maxRequests {
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Append the current request timestamp and update the map
		validRequests = append(validRequests, now)
		requestCount[ip] = validRequests

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
