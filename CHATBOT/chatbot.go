package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Respond function generates responses based on user input.
func Respond(input string) string {
	loweredInput := strings.ToLower(strings.TrimSpace(input))

	// Expanded responses
	switch loweredInput {
	case "hello", "hi":
		return "Hello! How can I assist you today?"
	case "how are you":
		return "I'm just a program, but I'm doing great! How about you?"
	case "bye", "goodbye":
		return "Goodbye! Have a fantastic day!"
	case "what's your name?", "what is your name?":
		return "I'm GoBot, your friendly chatbot!"
	case "what can you do?":
		return "I can chat with you and provide information. More features are coming soon!"
	case "tell me a joke":
		return "Why do programmers prefer dark mode? Because light attracts bugs!"
	default:
		return "I'm not sure I understand that. Can you ask something else?"
	}
}

func main() {
	// Initialize the Gin router
	r := gin.Default()

	// Serve a simple HTML page for the chatbot UI
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Chatbot</title>
				<style>
					body { font-family: Arial, sans-serif; margin: 20px; }
					#chat { margin-bottom: 20px; }
					.message { margin: 5px 0; }
					.user { color: blue; }
					.bot { color: green; }
				</style>
			</head>
			<body>
				<h1>Chat with GoBot</h1>
				<div id="chat"></div>
				<input type="text" id="message" placeholder="Type a message..." />
				<button onclick="sendMessage()">Send</button>

				<script>
					async function sendMessage() {
						const input = document.getElementById('message');
						const chat = document.getElementById('chat');
						const userMessage = input.value;
						if (!userMessage) return;

						// Display user message
						const userDiv = document.createElement('div');
						userDiv.className = 'message user';
						userDiv.innerText = 'You: ' + userMessage;
						chat.appendChild(userDiv);

						// Send to server and get response
						const response = await fetch('/chat', {
							method: 'POST',
							headers: { 'Content-Type': 'application/json' },
							body: JSON.stringify({ message: userMessage })
						});
						const data = await response.json();

						// Display bot response
						const botDiv = document.createElement('div');
						botDiv.className = 'message bot';
						botDiv.innerText = 'GoBot: ' + data.response;
						chat.appendChild(botDiv);

						// Clear input
						input.value = '';
					}
				</script>
			</body>
			</html>
		`)
	})

	// Endpoint for processing chatbot messages
	r.POST("/chat", func(c *gin.Context) {
		var request struct {
			Message string `json:"message"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Generate a response
		response := Respond(request.Message)
		c.JSON(http.StatusOK, gin.H{"response": response})
	})

	// Start the server
	r.Run(":8080") // Run on http://localhost:8080
}
