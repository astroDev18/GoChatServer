package gochat

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

// Global Variables used by the rest of the app

// map where the key is a pointer to web socket
var clients = make(map[*websocket.Conn]bool) // connected clients
// Queue for messages sent by clients
var broadcast = make(chan Message) // Broadcast Channel
// Eventually have goroutine to read new messages from the channel and send them to other clients connected

// Create instance for upgrader
// Object with methods for taking a normal HTTP connection and upgrading it to a Web Socket
var upgrader = websocket.Upgrader{}

// Define object to hold messages
// Struct with string attributes for email address, username and message

type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

// Create static fileserver and tie that to the "/" route
// User access and view index.html and any assets
func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Define "/ws" that hanndle requests for initiating a Web Socket
	http.HandleFunc("/ws", handleConnections)
	// Start listening for incoming chat messages
	go handleMessages()

	// Print helpful messages to the console and start the webserver logging any errors
	log.Println("http serveer started on :8000")

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

// Create functions to handle incoming WebSocket connections
// Use upgrader's "Upgrade()" method to change our initial GET request to a full on WebSocket
// If error is found log but dont exit

// Defer statement lets Go know to close out our WebSocket connection when the function returns
// Saving us from writing multiple "Close()" statements depending on how the function returns
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Printf("Closing failed")
		}
	}(ws)
	// Register new client by adding it to global "clients" map we created

	clients[ws] = true

	// Infinite loop that continuously waits for a new message to be written to the WebSocket
	// deserializes it from JSON to a Message object
	// throws it into the broadcast channel

	// If error detected we assume the client has disconnected
	// Log the error and remove client from the global "clients" map, so we don't contact the client again

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

// handleMessages Function
// A loop that continuously reads from "broadcast" channel
// relays the message to all the clients with their respective WebSocket connection
// Again, if an error is detected writing to the WebSocket, we close the connection and remove from "Clients" map

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				_ = client.Close()
				delete(clients, client)
			}
		}
	}
}




