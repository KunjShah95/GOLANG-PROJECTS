package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"api-gateway/config"      // Corrected import path
	"api-gateway/middlewares" // Corrected import path
)

// Helper function to parse duration from string with parameter name
func parseDuration(paramName, d string) time.Duration {
	duration, err := time.ParseDuration(d)
	if err != nil {
		log.Fatalf("Failed to parse duration for %s (%s): %v", paramName, d, err)
	}
	return duration
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Set default values if necessary
	if cfg.Port == "" {
		cfg.Port = "8080"
		log.Println("Port not set. Using default port 8080.")
	}
	if cfg.WriteTimeout == "" {
		cfg.WriteTimeout = "15s"
		log.Println("WriteTimeout not set. Using default WriteTimeout of 15s.")
	}
	if cfg.ReadTimeout == "" {
		cfg.ReadTimeout = "15s"
		log.Println("ReadTimeout not set. Using default ReadTimeout of 15s.")
	}
	if cfg.IdleTimeout == "" {
		cfg.IdleTimeout = "60s"
		log.Println("IdleTimeout not set. Using default IdleTimeout of 60s.")
	}
	if cfg.ShutdownTimeout == "" {
		cfg.ShutdownTimeout = "30s"
		log.Println("ShutdownTimeout not set. Using default ShutdownTimeout of 30s.")
	}

	// Validate BackendServers configuration
	if len(cfg.BackendServers) == 0 {
		log.Fatal("No backend servers configured. Please set BACKEND_SERVERS in the configuration.")
	}
	log.Printf("Configured Backend Servers: %v\n", cfg.BackendServers)

	// Initialize LoadBalancer with backend servers from config
	loadBalancer := middlewares.NewLoadBalancer(cfg.BackendServers)
	log.Println("LoadBalancer initialized with backend servers.")

	// Start LoadBalancer health checks
	go func() {
		log.Println("Starting LoadBalancer health checks.")
		loadBalancer.HealthCheck(30 * time.Second)
	}()

	// Initialize a new router
	r := mux.NewRouter()

	// Create a subrouter for health check endpoints to bypass middleware
	healthCheckRouter := r.PathPrefix("/health").Subrouter()
	healthCheckRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Gateway is running"))
	}).Methods("GET")
	healthCheckRouter.HandleFunc("/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health check for v1"))
	}).Methods("GET")
	healthCheckRouter.HandleFunc("/v2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health check for v2"))
	}).Methods("GET")

	log.Println("Health check routes registered.")

	// Remove duplicate health check handlers to avoid routing conflicts
	// (These have been omitted as they were causing conflicts.)

	// Apply middlewares after defining healthCheckRouter to ensure health endpoints bypass them
	r.Use(middlewares.Logging)
	r.Use(middlewares.SlidingWindowRateLimit)
	r.Use(middlewares.CircuitBreaker)
	r.Use(loadBalancer.RoundRobinLoadBalancer) // Load balancer before cache and compression
	r.Use(middlewares.Cache)
	r.Use(middlewares.GzipCompression)
	r.Use(middlewares.TimeoutMiddleware)
	r.Use(middlewares.RetryMiddleware) // Retry after timeout

	log.Println("Middlewares applied to the router.")

	// Define additional routes if necessary
	// Example:
	// r.HandleFunc("/api/v1/resource", ResourceHandler).Methods("GET")

	// Create the HTTP server with configuration parameters
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + cfg.Port,
		WriteTimeout: parseDuration("WriteTimeout", cfg.WriteTimeout),
		ReadTimeout:  parseDuration("ReadTimeout", cfg.ReadTimeout),
		IdleTimeout:  parseDuration("IdleTimeout", cfg.IdleTimeout),
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting API Gateway on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("API Gateway is up and running.")

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown signal received. Shutting down server...")

	// Create a deadline to wait for ongoing operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), parseDuration("ShutdownTimeout", cfg.ShutdownTimeout))
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
