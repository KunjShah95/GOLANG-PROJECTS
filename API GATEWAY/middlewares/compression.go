package middlewares

import (
	"net/http"
)

// Other middleware functions can be defined here if needed.
// Ensure that GzipCompression is not redeclared and is used from middlewares/gzip.go.

// Example of another middleware function
func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example middleware logic
		next.ServeHTTP(w, r)
	})
}
