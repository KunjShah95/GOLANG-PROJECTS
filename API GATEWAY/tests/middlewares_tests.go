package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/middlewares"
)

// Example unit test for Logging middleware
func TestLoggingMiddleware(t *testing.T) {
	// Create a sample handler to test
	handler := middlewares.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Additional checks can be added to verify log output if needed
}
