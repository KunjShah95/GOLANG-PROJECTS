package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const maxConnections = 100

var activeConnections = make(chan struct{}, maxConnections)

func main() {
	// Configurable host and port
	port := flag.String("port", "8080", "Port to run the TCP server on")
	useTLS := flag.Bool("tls", false, "Enable TLS (requires server.crt and server.key)")
	flag.Parse()

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	address := fmt.Sprintf("%s:%s", host, *port)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var listener net.Listener
	var err error

	if *useTLS {
		// TLS configuration
		cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
		if err != nil {
			log.Fatalf("Failed to load TLS certificates: %s", err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err = tls.Listen("tcp", address, config)
		log.Println("TLS enabled")
	} else {
		// Plain TCP listener
		listener, err = net.Listen("tcp", address)
	}

	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	defer listener.Close()

	log.Printf("TCP server listening on %s\n", address)

	// Goroutine to handle graceful shutdown
	go func() {
		<-stop
		log.Println("\nShutting down server...")
		listener.Close()
		close(activeConnections)
		os.Exit(0)
	}()

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Limit active connections
		activeConnections <- struct{}{} // Reserve a slot
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConnection(conn)
			<-activeConnections // Release the slot
		}()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Connection details
	log.Printf("New connection from %s\n", conn.RemoteAddr().String())

	// Set read/write timeout
	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	// Simple handshake authentication
	conn.Write([]byte("Enter password: "))
	reader := bufio.NewReader(conn)
	password, _ := reader.ReadString('\n')
	if strings.TrimSpace(password) != "secret" {
		conn.Write([]byte("Authentication failed. Disconnecting.\n"))
		log.Printf("Authentication failed for %s\n", conn.RemoteAddr().String())
		return
	}
	conn.Write([]byte("Authentication successful. Welcome!\n"))

	// Handle client communication
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection closed by %s\n", conn.RemoteAddr().String())
			return
		}

		// Process and respond to the client
		message = strings.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		log.Printf("Received from %s: %s\n", conn.RemoteAddr().String(), message)

		// Check if the message is "Hello"
		if strings.ToLower(message) == "hello" {
			response := "Hello! How can I help you?\n"
			conn.Write([]byte(response))
			log.Printf("Responded with: %s\n", response)
		} else {
			// Echo any other message
			response := fmt.Sprintf("Echo: %s\n", message)
			conn.Write([]byte(response))
		}
	}
}
