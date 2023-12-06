package storage

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type NotFoundError struct {
	Id string
}

func (c *NotFoundError) Error() string {
	return fmt.Sprintf("Connection not found with id: %s", c.Id)
}

type BroadcastError struct {
	Ids []string
}

func (b *BroadcastError) Error() string {
	return fmt.Sprintf("Error broadcasting to connections with ids:\n %s\n", b.Ids)
}

type ConnectionsMap struct {
	m  sync.Mutex;
	userConnections  map[*websocket.Conn]string;
}

func (cm *ConnectionsMap) CreateConnection(connection *websocket.Conn, id string) *NotFoundError { 
	// create connection only if it doesn't exist
	cm.m.Lock()
	if _, ok := cm.userConnections[connection]; !ok{
		cm.userConnections[connection] = id
		cm.m.Unlock()
		return nil
	} else {
		cm.m.Unlock()
		return &NotFoundError{Id: id}
	} 
}

func (cm *ConnectionsMap) UpdateConnection(connection *websocket.Conn, id string) *NotFoundError { 
	// update connection only if it does exist
	cm.m.Lock()
	if _, ok := cm.userConnections[connection]; ok{
		cm.userConnections[connection] = id
		cm.m.Unlock()
		return nil
	} else {
		cm.m.Unlock()
		return &NotFoundError{Id: id}
	} 
}

func (cm *ConnectionsMap) GetConnectionId(connection *websocket.Conn) (string, *NotFoundError) { 
	// update connection only if it does exist
	cm.m.Lock()
	if id, ok := cm.userConnections[connection]; ok{
		cm.userConnections[connection] = id
		cm.m.Unlock()
		return id, nil
	} else {
		cm.m.Unlock()
		return "", &NotFoundError{Id: id}
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
			id := cm.userConnections[conn]
			berr.Ids = append(berr.Ids, id)
			conn.Close()
			delete(cm.userConnections,conn)
		}
	}
	if len(berr.Ids) == 0 {
		return nil
	} else {
		return berr
	}
}