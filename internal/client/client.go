package client

import (
	"bufio"
	"chat-app/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type Client struct {
	Url    string
	UserId int
	RoomId int
}

func NewClient(baseurl string, roomId, userId int) *Client {
	// baseurl := "ws://localhost:8080/chatroom"
	return &Client{
		Url:    fmt.Sprintf("%s/%d", baseurl, roomId),
		UserId: userId,
		RoomId: roomId,
	}
}

func (c *Client) RunClient() {

	conn, _, err := websocket.DefaultDialer.Dial(c.Url, nil)

	if err != nil {
		log.Fatalf("error dialing %s", err)
	}
	defer conn.Close()

	conn.WriteJSON(models.UserMessage{UserId: c.UserId, Message: "Joined the Chat!"})

	done := make(chan bool)

	// reading server messages
	go func() {
		defer close(done)
		for {
			_, recvMessage, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading: %s\n", err.Error())
				return
			}

			var message models.ResponseMessage
			err = json.Unmarshal(recvMessage, &message)
			if err != nil {
				log.Printf("error reading: %s\n", err.Error())
				return
			}

			fmt.Printf("%s: %s\n", message.FromUser, message.Message)
		}
	}()

	// writing messages to server
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			inputMessage, _ := reader.ReadString('\n')
			if inputMessage == "(exit)\n" {
				fmt.Println("quit...")
				done <- true
			}
			err := conn.WriteJSON(models.UserMessage{UserId: c.UserId, Message: inputMessage})
			if err != nil {
				log.Printf("error writing %s\n", err)
				return
			}
			fmt.Printf("->\t%s\n", inputMessage)
		}
	}()

	<-done
	conn.WriteJSON(models.UserMessage{UserId: c.UserId, Message: "Left the Chat!"})
}
