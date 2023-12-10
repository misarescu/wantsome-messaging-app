package main

import (
	"chat-app/internal/client"
	"flag"
)

func main() {
	roomId := flag.Int("room", 0, "room id to be joined")
	userId := flag.Int("userid", 0, "user id to join")

	flag.Parse()

	c := client.NewClient("ws://localhost:8080/chatroom",*roomId,*userId,)
	
	c.RunClient()
}
