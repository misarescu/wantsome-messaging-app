package client

import (
	"chat-app/pkg/models"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func RunClient() {
	url := "ws://localhost:8080/ws"
	randId := rand.Intn(10)
	message := models.UserMessage{Message: fmt.Sprintf("Hello world from my client %d!", randId), UserId: randId}

	c, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		log.Fatalf("error dialing %s", err)
	}
	defer c.Close()

	done := make(chan bool)

	// reading server messages
	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("error reading: %s\n", err.Error())
				return
			}
			fmt.Printf("got message: %s\n", msg)
		}
	}()

	// writing messages to server
	go func() {
		for {
			err := c.WriteJSON(message)
			if err != nil {
				log.Printf("error writing %s\n", err)
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	<-done
}
