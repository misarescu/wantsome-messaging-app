package storage

import (
	"chat-app/pkg/models"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type NotFoundError struct {
	Id int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Connection not found with id: %d", e.Id)
}

type BroadcastError struct {
	Users []*models.User
}

func (e *BroadcastError) Error() string {
	var idSlice []int
	for _, user := range e.Users {
		idSlice = append(idSlice, user.Id)
	}
	return fmt.Sprintf("Error broadcasting to connections with ids:\n %d\n", idSlice)
}

type ConnectionsMap struct {
	m  sync.Mutex;
	userConnections  map[*websocket.Conn]*models.User;
}

func (cm *ConnectionsMap) CreateConnection(connection *websocket.Conn, user *models.User) *NotFoundError { 
	// create connection only if it doesn't exist
	cm.m.Lock()
	if _, ok := cm.userConnections[connection]; !ok{
		cm.userConnections[connection] = user
		cm.m.Unlock()
		return nil
	} else {
		cm.m.Unlock()
		return &NotFoundError{Id: user.Id}
	} 
}

func (cm *ConnectionsMap) UpdateConnection(connection *websocket.Conn, user *models.User) *NotFoundError { 
	// update connection only if it does exist
	cm.m.Lock()
	if _, ok := cm.userConnections[connection]; ok{
		cm.userConnections[connection] = user
		cm.m.Unlock()
		return nil
	} else {
		cm.m.Unlock()
		return &NotFoundError{Id: user.Id}
	} 
}

func (cm *ConnectionsMap) GetConnectionId(connection *websocket.Conn) (*models.User, *NotFoundError) { 
	// update connection only if it does exist
	cm.m.Lock()
	if user, ok := cm.userConnections[connection]; ok{
		cm.m.Unlock()
		return user, nil
	} else {
		cm.m.Unlock()
		return nil, &NotFoundError{Id: user.Id}
	} 
}

func (cm *ConnectionsMap) DeleteConnection(connection *websocket.Conn) { 
	// delete connection
	cm.m.Lock()
	delete(cm.userConnections, connection)
	cm.m.Unlock()
	
}

func (cm *ConnectionsMap) DeleteAllConnection() {
	// delete connections
	cm.m.Lock()
	for conn := range cm.userConnections {
		conn.Close()
		delete(cm.userConnections,conn)
	}
	cm.m.Unlock()
}

func (cm *ConnectionsMap) BroadcastMessage(msg string) *BroadcastError {
	berr := &BroadcastError{}
	cm.m.Lock()
	for conn := range cm.userConnections {
		if err := conn.WriteJSON(msg); err != nil{
			user := cm.userConnections[conn]
			berr.Users = append(berr.Users, user)
			conn.Close()
			delete(cm.userConnections,conn)
		}
	}
	if len(berr.Users) == 0 {
		return nil
	} else {
		return berr
	}
}