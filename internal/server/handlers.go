package server

import (
	"chat-app/pkg/loggers"
	"chat-app/pkg/models"
	"encoding/json"
	"errors"
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

func (s *Server) handleReadRoomChat(w http.ResponseWriter, r *http.Request) error {
	strRoomId := mux.Vars(r)["id"]
	roomId, err := strconv.Atoi(strRoomId) 
	if err != nil {
		return fmt.Errorf("expected room id to be integer, but recieved: %s", strRoomId)
	} 

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("got error upgrading connection %s", err.Error())
	}
	defer conn.Close()

	room, err := s.store.GetRoomById(roomId)

	if err != nil {
		return err
	}

	result := func() error {
		for {
			var msg models.UserMessage = models.UserMessage{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				loggers.WarningLogger.Printf("got error reading message %s\n", err.Error())
				room.DeleteConnection(conn)
				return &models.ConnectionError{}
			} 

			user, err := s.store.GetUserById(msg.UserId)

			if err != nil {
				return &models.NotFoundError{Id: msg.UserId}
			}
			
			room.CreateConnection(conn, user)
			loggers.InfoLogger.Printf("message recieved: %+v\n", msg)
			s.broadcast <- models.RoomMessage{RoomId: roomId, UserMessage: msg}
		}
	}()

	return result
}

func (s *Server) handleBroadcastRoomChat() {
	for {
		select {
		case msg:= <-s.broadcast :
			room, err := s.store.GetRoomById(msg.RoomId)
			if err != nil {
				loggers.ErrorLogger.Printf("%s\n", err.Error())
				continue
			}

			if user, err := s.store.GetUserById(msg.UserMessage.UserId); err == nil{
				room.BroadcastMessage(models.ResponseMessage{Message:msg.UserMessage.Message, FromUser:user.Name, ErrorStatus: false})
			} else {
				loggers.ErrorLogger.Printf("err: %s", err.Error())
			}
		}
	}
}

func (s *Server) handleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	users, _ := s.store.GetAllUsers()
	writeJSON(w, http.StatusOK, users)

	return nil
}

func (s *Server) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	user, err := s.store.GetUserById(id)
	if err != nil {
		return err
	}

	writeJSON(w,http.StatusOK, user)
	
	return nil
}

func (s *Server) handleRemoveUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}
	user, err := s.store.RemoveUserById(id)
	if err != nil {
		return err 
	}
	loggers.WarningLogger.Printf("user %+v was deleted\n", user)
	writeJSON(w,http.StatusOK, user)

	return nil
}

func (s *Server) handleUpdateUserById(w http.ResponseWriter, r *http.Request) error {
	if headerContentType := r.Header.Get("Content-Type"); headerContentType != "application/json" {
		return &models.BadRequestError{Message: "Content type needs to be json"}
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil{
		return &models.BadRequestError{Message: "id type needs to be integer"}
	}
	user, err := s.store.GetUserById(id)

	if err != nil {
		return &models.BadRequestError{Message: err.Error()}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(user); err != nil{
		var unmarshallErr *json.UnmarshalTypeError

		if errors.As(err, &unmarshallErr) {
			// return models.NewRequestError("Bad Request: Wrong type provided")
			return &models.BadRequestError{Message: "Wrong type provided"}
		} else {
			return &models.BadRequestError{Message: err.Error()}
		}
	}

	retUser, err := s.store.UpdateUser(user)

	if err != nil {
		return &models.BadRequestError{Message: err.Error()}
	}

	writeJSON(w,http.StatusOK, retUser)

	return nil

}