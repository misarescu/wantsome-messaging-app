package server

import (
	"chat-app/pkg/models"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	broadcast = make(chan models.Message)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// func home(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello world from my server!")
// }

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("got error upgrading connection %s\n", err)
		return
	}
	defer conn.Close()

	m.Lock()
	userConnections[conn] = ""
	m.Unlock()
	fmt.Printf("Connected client")

	for {
		var msg models.Message = models.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("got error reading message %s\n", err)
			m.Lock()
			delete(userConnections, conn)
			m.Unlock()
			return
		}
		m.Lock()
		userConnections[conn] = msg.UserName
		m.Unlock()
		broadcast <- msg
	}
}

func handleMsg() {
	for {
		msg := <-broadcast

		m.Lock()
		for client := range userConnections {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("got error broadcating message to client %s", err)
				client.Close()
				delete(userConnections, client)
			}
		}
		m.Unlock()
	}
}
