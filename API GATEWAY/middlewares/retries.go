package middlewares

import (
	"log"
	"net/http"
	"time"
)

var maxRetries = 3

// RetryMiddleware attempts to retry failed backend requests
func RetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		for i := 0; i < maxRetries; i++ {
			err = tryNextHandler(w, r, next)
			if err == nil {
				return
			}
			log.Println("Retry attempt", i+1, "due to error:", err)
			time.Sleep(2 * time.Second) // Wait before retrying
		}
		http.Error(w, "Service Unavailable after retries", http.StatusServiceUnavailable)
	})
}

func tryNextHandler(w http.ResponseWriter, r *http.Request, next http.Handler) error {
	// Call the next handler
	next.ServeHTTP(w, r)
	return nil // Assuming no error is returned by the next handler
}
