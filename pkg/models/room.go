package models

import (
	"chat-app/pkg/loggers"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	m               sync.Mutex `json:"-"`
	UserConnections map[*websocket.Conn]*User `json:"-"`
	Id              int `json:"id,omitempty"` // id 0 is forbidden
	Name						string `json:"name"`
}

func (r *Room) CreateConnection(connection *websocket.Conn, user *User) error {
	// create connection only if it doesn't exist
	r.m.Lock()
	if _, ok := r.UserConnections[connection]; !ok {
		
		r.UserConnections[connection] = user
		r.m.Unlock()
		return nil
	} else {
		r.m.Unlock()
		return fmt.Errorf("user %d already has a connection", user.Id)
	}
}

func (r *Room) UpdateConnection(connection *websocket.Conn, user *User) error {
	// update connection only if it does exist
	r.m.Lock()
	if _, ok := r.UserConnections[connection]; ok {
		r.UserConnections[connection] = user
		r.m.Unlock()
		return nil
	} else {
		r.m.Unlock()
		return &NotFoundError{Id: user.Id}
	}
}

func (r *Room) GetUserByConnection(connection *websocket.Conn) (*User, error) {
	// update connection only if it does exist
	r.m.Lock()
	if user, ok := r.UserConnections[connection]; ok {
		r.m.Unlock()
		return user, nil
	} else {
		r.m.Unlock()
		return nil, &NotFoundError{Id: user.Id}
	}
}

func (r *Room) DeleteConnection(connection *websocket.Conn) {
	// delete connection
	r.m.Lock()
	delete(r.UserConnections, connection)
	r.m.Unlock()

}

func (r *Room) DeleteAllConnection() {
	// delete connections
	r.m.Lock()
	for conn := range r.UserConnections {
		conn.Close()
		delete(r.UserConnections, conn)
	}
	r.m.Unlock()
}

func (r *Room) BroadcastMessage(msg ResponseMessage, ignoreUser *User) error {
	berr := &BroadcastError{}
	r.m.Lock()
	loggers.InfoLogger.Printf("connections present: %+v\n", r.UserConnections)
	for conn, user := range r.UserConnections {
		if user != ignoreUser{
			err := conn.WriteJSON(msg);
			if err != nil {
				user := r.UserConnections[conn]
				berr.Users = append(berr.Users, user)
				conn.Close()
				delete(r.UserConnections, conn)
				loggers.WarningLogger.Printf("deleted conn %+v", conn)
			}
		}
	}
	r.m.Unlock()
	if len(berr.Users) == 0 {
		return nil
	} else {
		return berr
	}
}
