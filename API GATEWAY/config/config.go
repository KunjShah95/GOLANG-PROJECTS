package config

import (
	"log"
	"os"
	"strings" // Import strings package

	"github.com/joho/godotenv"
)

// LoadConfig loads the configuration from environment variables
func LoadConfig() Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Load configuration from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	writeTimeout := os.Getenv("WRITE_TIMEOUT")
	if writeTimeout == "" {
		writeTimeout = "15s"
	}

	readTimeout := os.Getenv("READ_TIMEOUT")
	if readTimeout == "" {
		readTimeout = "15s"
	}

	idleTimeout := os.Getenv("IDLE_TIMEOUT")
	if idleTimeout == "" {
		idleTimeout = "60s"
	}

	shutdownTimeout := os.Getenv("SHUTDOWN_TIMEOUT")
	if shutdownTimeout == "" {
		shutdownTimeout = "5s"
	}

	backendServers := os.Getenv("BACKEND_SERVERS")
	if backendServers == "" {
		backendServers = "http://localhost:8081,http://localhost:8082"
	}

	return Config{
		Port:            port,
		WriteTimeout:    writeTimeout,
		ReadTimeout:     readTimeout,
		IdleTimeout:     idleTimeout,
		ShutdownTimeout: shutdownTimeout,
		BackendServers:  strings.Split(backendServers, ","),
	}
}

// Config holds the configuration for the API Gateway
type Config struct {
	Port            string
	WriteTimeout    string
	ReadTimeout     string
	IdleTimeout     string
	ShutdownTimeout string
	BackendServers  []string
}
