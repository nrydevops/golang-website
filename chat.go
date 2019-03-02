package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)


var clients map[*websocket.Conn]string // connected clients
var broadcast chan MessageSentJSON       // broadcast channel
var mux	sync.Mutex

// The upgrader
var upgrader websocket.Upgrader

func init() {

	clients = make(map[*websocket.Conn]string) // connected clients
	broadcast = make(chan MessageSentJSON)       // broadcast channel

	// Configure the upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	go func() {

		for _ = range time.NewTicker(10 * time.Second).C {

			mux.Lock()
			// Pinging to clients to identify their states
			for client := range clients {

				err := client.WriteJSON("{}")
				if err != nil {

					client.Close()
					delete(clients, client)
				}
			}
			mux.Unlock()

			// Send to everyone list of chat users
			for ws, _ := range clients {

				ws.WriteJSON(MessageActiveUsersJSON{activeClientUsernames()})
			}
		}
	}()
}

func activeClientUsernames() []string {

	mux.Lock()
	defer mux.Unlock()
	// Creating slice of all usernames
	usernames := make([]string, len(clients), len(clients))
	idx := 0
	for _,value := range clients {

		usernames[idx] = value
		idx++
	}

	return usernames
}

func handleConnections(w http.ResponseWriter, req *http.Request) {

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// User object of this client
	user := GetUserFromRequest(req)

	// Register our new client
	clients[ws] = user.Username

	// Send to client available info about users
	ws.WriteJSON(MessageAllUserInfoJSON{user.Username, activeClientUsernames()})

	// Get recent messages from database
	db_messages := GetRecentMessages(100)
	messages := MessageArraySentJSON{make([]MessageSentJSON, len(db_messages), len(db_messages))}
	for i, message := range db_messages {

		messages.Messages[i] = MessageSentJSON{
			message.Username,
			message.Message,
			strings.Split(message.Date," ")[1],
		}
	}

	// Send to client recent messages ( 100 recent messages )
	ws.WriteJSON(messages)

	for {
		var msg_received MessageReceivedJSON
		// Read in a new message as JSON and map it to a MessageJSON object
		err := ws.ReadJSON(&msg_received)
		if err != nil {
			delete(clients, ws)
			break
		}

		_ = WriteMessageToDB(
			Message{
				user.Username,
				msg_received.Message,
				time.Now().Format("2006-01-02 15:04:05"),
			})

		// Send the newly received message to the broadcast channel
		broadcast <- MessageSentJSON{
			user.Username,
			msg_received.Message,
			time.Now().Format("15:04:05"),
		}
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for client := range clients {

			client.WriteJSON(msg)
		}
	}
}