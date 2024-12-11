package middlewares

import (
	"net/http"
	"sync"
	"time"
)

var (
	circuitOpen     bool
	circuitMutex    sync.Mutex
	lastFailureTime time.Time
)

const circuitTimeout = 30 * time.Second // Circuit timeout duration

// CircuitBreaker middleware prevents request forwarding if the circuit is open
func CircuitBreaker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		circuitMutex.Lock()
		defer circuitMutex.Unlock()

		// Check if the circuit is open and within the timeout period
		if circuitOpen && time.Since(lastFailureTime) < circuitTimeout {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		// Allow the request if the circuit is closed
		next.ServeHTTP(w, r)
	})
}

// SetCircuitOpen sets the circuit to an open state
func SetCircuitOpen() {
	circuitMutex.Lock()
	defer circuitMutex.Unlock()
	circuitOpen = true
	lastFailureTime = time.Now()
}
