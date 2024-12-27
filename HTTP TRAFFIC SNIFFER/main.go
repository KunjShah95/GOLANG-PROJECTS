package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Function to check if the request should be logged based on filters
func shouldLogRequest(req *http.Request, ipFilter, urlFilter, methodFilter string) bool {
	if ipFilter != "" && req.RemoteAddr != ipFilter {
		return false
	}
	if urlFilter != "" && !strings.Contains(req.URL.String(), urlFilter) {
		return false
	}
	if methodFilter != "" && req.Method != methodFilter {
		return false
	}
	return true
}

// Function to log the request to a file
func logRequest(req *http.Request, logFile *os.File, status string) {
	// Read request body
	body, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Log entry format
	logEntry := fmt.Sprintf("%s - %s\nTimestamp: %s\nMethod: %s\nURL: %s\nRemoteAddr: %s\nHeaders: %v\nBody: %s\nStatus: %s\n\n",
		req.Header.Get("User-Agent"),
		status,
		time.Now().Format(time.RFC3339),
		req.Method,
		req.URL.String(),
		req.RemoteAddr,
		req.Header,
		string(body),
		status)
	_, _ = fmt.Fprintln(logFile, logEntry)
}

// Rate limiter to control the logging rate
var rateLimiter = &RateLimiter{
	limit: 10, // Max 10 logs per second
	burst: 10,
	tick:  1 * time.Second,
}

type RateLimiter struct {
	limit  int
	burst  int
	tick   time.Duration
	mu     sync.Mutex
	cond   *sync.Cond
	tokens int
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// If there are tokens available, allow logging
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

func (r *RateLimiter) Refill() {
	ticker := time.NewTicker(r.tick)
	for range ticker.C {
		r.mu.Lock()
		// Refill tokens based on rate limiter burst capacity
		if r.tokens < r.burst {
			r.tokens++
		}
		r.mu.Unlock()
	}
}

// Function to rotate log file based on size
func rotateLogFile(logFile *os.File) *os.File {
	// Get file info
	fileInfo, err := logFile.Stat()
	if err != nil {
		log.Fatal("Error checking log file size:", err)
	}

	// If the file exceeds 5MB, rotate the log file
	if fileInfo.Size() > 5*1024*1024 {
		// Close the current log file
		logFile.Close()

		// Rename the old log file to include a timestamp
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		err := os.Rename("http_traffic.log", fmt.Sprintf("http_traffic_%s.log", timestamp))
		if err != nil {
			log.Fatal("Error rotating log file:", err)
		}

		// Create a new log file
		newLogFile, err := os.OpenFile("http_traffic.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return newLogFile
	}
	return logFile
}

func main() {
	// Define filters
	ipFilter := ""     // Allow all IPs
	urlFilter := "/"   // Allow all URLs
	methodFilter := "" // Allow all methods

	// Open or create the log file
	logFile, err := os.OpenFile("http_traffic.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Start rate limiter in a separate goroutine
	go rateLimiter.Refill()

	// Serve HTTP requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request should be logged based on filters and rate limiter
		if shouldLogRequest(r, ipFilter, urlFilter, methodFilter) && rateLimiter.Allow() {
			logRequest(r, logFile, "Logged")
			fmt.Println("Logged request:", r.URL) // Debug print to confirm logging
		} else {
			// Log all requests that are not logged due to filters or rate limiter
			logRequest(r, logFile, "Not Logged - Filter/Rate Limiter")
			fmt.Println("Request not logged due to filters or rate limiter.")
		}

		// Respond with a simple message
		fmt.Fprintf(w, "Hello, World!")
	})

	// Periodically check log file size and rotate if necessary
	go func() {
		for {
			time.Sleep(10 * time.Second) // Check every 10 seconds
			logFile = rotateLogFile(logFile)
		}
	}()

	// Start HTTP server
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
