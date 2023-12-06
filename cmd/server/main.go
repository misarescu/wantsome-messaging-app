package main

import (
	"chat-app/internal/server"
	"chat-app/internal/storage"
	"flag"
	"fmt"
)

func main() {
	listenAddr := flag.String("addr", "localhost", "the server address")
	listenPort := flag.String("port", "8080", "the server port")
	flag.Parse()

	store := storage.NewMemoryStorage()

	s := server.NewServer(fmt.Sprintf("%s:%s", *listenAddr, *listenPort), store)
	s.RunServer()
}
