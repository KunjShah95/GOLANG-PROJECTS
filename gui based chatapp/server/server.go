package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	_ "modernc.org/sqlite" // SQLite driver
)

var (
	clients     = make(map[*websocket.Conn]bool)
	broadcast   = make(chan Message)
	upgrader    = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	db          *sql.DB
	rateLimiter = make(map[*websocket.Conn]time.Time)
)

// Message represents a chat message.
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

// Initialize the database.
func init() {
	var err error
	db, err = sql.Open("sqlite", "./chat.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		content TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

// Sanitize input to prevent malicious content.
func sanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// Save message to the database.
func saveMessage(msg Message) {
	_, err := db.Exec("INSERT INTO messages (username, content) VALUES (?, ?)", msg.Username, msg.Content)
	if err != nil {
		fmt.Println("Error saving message:", err)
	}
}

// Handle WebSocket connections.
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	// Authentication
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err = ws.ReadJSON(&creds)
	if err != nil || creds.Password != "securepass" { // Simplified authentication
		ws.WriteJSON("Authentication failed")
		ws.Close()
		return
	}

	// Listen for messages
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(clients, ws)
			break
		}

		msg.Content = sanitizeInput(msg.Content)

		// Rate limiting
		if lastTime, ok := rateLimiter[ws]; ok {
			if time.Since(lastTime) < 1*time.Second {
				ws.WriteJSON("You are sending messages too quickly.")
				continue
			}
		}
		rateLimiter[ws] = time.Now()

		broadcast <- msg
	}
}

// Broadcast messages to all clients.
func handleMessages() {
	for {
		msg := <-broadcast
		saveMessage(msg)
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static")) // Optional static file server
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Server running on https://localhost:8080")
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil) // TLS-enabled server
	if err != nil {
		log.Fatal("Server startup error:", err)
	}
}
