package storage

import (
	"sync"

	"github.com/gorilla/websocket"
)

type UserConnectionsMap struct {
	m  sync.Mutex;
	userConnections  map[*websocket.Conn]string;
}

func (ucm *UserConnectionsMap) SetConnection(){ }

func (ucm *UserConnectionsMap) GetConnection(){ }

func (ucm *UserConnectionsMap) DeleteConnection(){ }

