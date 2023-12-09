package server

import (
	"chat-app/internal/storage"
	"chat-app/pkg/loggers"
	"chat-app/pkg/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	// broadcast = make(chan models.Message)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)


func (s *Server) handleRoomChat(w http.ResponseWriter, r *http.Request) error {
	strRoomId := mux.Vars(r)["id"]
	roomId, err := strconv.Atoi(strRoomId) 
	if err != nil {
		return fmt.Errorf("expected room id to be intiger, but recieved: %s", strRoomId)
	} 

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("got error upgrading connection %s\n", err)
	}
	defer conn.Close()

	room, err := s.store.GetRoomById(roomId)

	if err != nil {
		return err
	}

	loggers.InfoLogger.Println("Client connected")

	for {
		var msg models.Message = models.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("got error reading message %s\n", err)
			room.DeleteConnection(conn)
			return &models.ConnectionError{}
		} 

		user, err := s.store.GetUserById(msg.UserId)

		if err != nil {
			return &models.NotFoundError{Id: msg.UserId}
		}
		
		room.CreateConnection(conn, user)
		
		room.Broadcast <- msg
	}
}

func handleMsg(store *storage.Storage) {
	for {

		

		// msg := <-room.broadcast

		// m.Lock()
		// for client := range userConnections {
		// 	err := client.WriteJSON(msg)
		// 	if err != nil {
		// 		fmt.Printf("got error broadcating message to client %s", err)
		// 		client.Close()
		// 		delete(userConnections, client)
		// 	}
		// }
		// m.Unlock()
	}
}

func (s *Server) handleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	users, _ := s.store.GetAllUsers()
	writeJSON(w, http.StatusOK, users)

	return nil
}