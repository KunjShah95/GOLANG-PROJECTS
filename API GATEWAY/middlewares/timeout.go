package middlewares

import (
	"context"
	"net/http"
	"time"
)

var timeout = 5 * time.Second // Timeout duration for backend requests

// TimeoutMiddleware sets a timeout for backend requests
func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
