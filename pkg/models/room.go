package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	m               sync.Mutex
	userConnections map[*websocket.Conn]*User
	Id              int
}

func (r *Room) CreateConnection(connection *websocket.Conn, user *User) *NotFoundError {
	// create connection only if it doesn't exist
	r.m.Lock()
	if _, ok := r.userConnections[connection]; !ok {
		r.userConnections[connection] = user
		r.m.Unlock()
		return nil
	} else {
		r.m.Unlock()
		return &NotFoundError{Id: user.Id}
	}
}

func (r *Room) UpdateConnection(connection *websocket.Conn, user *User) *NotFoundError {
	// update connection only if it does exist
	r.m.Lock()
	if _, ok := r.userConnections[connection]; ok {
		r.userConnections[connection] = user
		r.m.Unlock()
		return nil
	} else {
		r.m.Unlock()
		return &NotFoundError{Id: user.Id}
	}
}

func (r *Room) GetConnectionId(connection *websocket.Conn) (*User, *NotFoundError) {
	// update connection only if it does exist
	r.m.Lock()
	if user, ok := r.userConnections[connection]; ok {
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
	delete(r.userConnections, connection)
	r.m.Unlock()

}

func (r *Room) DeleteAllConnection() {
	// delete connections
	r.m.Lock()
	for conn := range r.userConnections {
		conn.Close()
		delete(r.userConnections, conn)
	}
	r.m.Unlock()
}

func (r *Room) BroadcastMessage(msg string) *BroadcastError {
	berr := &BroadcastError{}
	r.m.Lock()
	for conn := range r.userConnections {
		if err := conn.WriteJSON(msg); err != nil {
			user := r.userConnections[conn]
			berr.Users = append(berr.Users, user)
			conn.Close()
			delete(r.userConnections, conn)
		}
	}
	if len(berr.Users) == 0 {
		return nil
	} else {
		return berr
	}
}
