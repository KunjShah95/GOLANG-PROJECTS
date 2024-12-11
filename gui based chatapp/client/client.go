package main

import (
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
)

var ws *websocket.Conn

// connect to the server
func connectToServer(username, password string) error {

	u := url.URL{Scheme: "wss", Host: "localhost:8080", Path: "/ws"}

	var err error
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	//send authentication details
	auth := map[string]string{"username": username, "password": password}
	err = ws.WriteJSON(auth)
	if err != nil {
		return err
	}
	return nil
}

// start the chat GUI
func startChatGUI(username string) {
	a := app.New()
	w := a.NewWindow("Chat App")

	messageList := widget.NewMultiLineEntry()
	messageList.SetPlaceHolder("Chat Messages")
	messageList.Disable()

	messageInput := widget.NewEntry()
	messageInput.SetPlaceHolder("Type a message...")

	sendButton := widget.NewButton("Send", func() {
		msg := map[string]string{"username": username, "message": messageInput.Text}
		err := ws.WriteJSON(msg)
		if err != nil {
			log.Println("Error sending message:", err)
		}
		messageInput.SetText("")
	})

	go func() {
		for {
			var msg map[string]string
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println("Error receiving message:", err)
				return
			}
			messageList.SetText(messageList.Text + "\n" + msg["username"] + ": " + msg["content"])
		}
	}()

	w.SetContent(container.NewVBox(
		messageList,
		messageInput,
		sendButton,
	))
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Login")

	username := widget.NewEntry()
	username.SetPlaceHolder("Username")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	loginButton := widget.NewButton("Login", func() {
		err := connectToServer(username.Text, password.Text)
		if err != nil {
			log.Println("Connection error:", err)
		} else {
			myWindow.Close()
			startChatGUI(username.Text)
		}
	})

	myWindow.SetContent(container.NewVBox(
		username,
		password,
		loginButton,
	))
	myWindow.ShowAndRun()
}
