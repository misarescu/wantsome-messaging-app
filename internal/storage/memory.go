package storage

import (
	"chat-app/pkg/models"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	rm sync.Mutex
	um sync.Mutex
)

type MemoryStorage struct {
	users map[int]*models.User // all users in the app
	rooms map[int]*models.Room // all rooms in the app
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users: map[int]*models.User{
			0: {
				Id:   0,
				Name: "Alice",
			},
			1: {
				Id:   1,
				Name: "Bob",
			},
			2: {
				Id:   2,
				Name: "Charlie",
			},
			3: {
				Id:   3,
				Name: "Dan",
			},
			4: {
				Id:   4,
				Name: "Elly",
			},
			5: {
				Id:   5,
				Name: "Finn",
			},
			6: {
				Id:   6,
				Name: "Gin",
			},
			7: {
				Id:   7,
				Name: "Henry",
			},
			8: {
				Id:   8,
				Name: "Irina",
			},
			9: {
				Id:   9,
				Name: "Jake",
			},
		},
		rooms: map[int]*models.Room{
			0: {
				Id: 0,
				UserConnections: make(map[*websocket.Conn]*models.User),
			},
			1: {
				Id: 1,
				UserConnections: make(map[*websocket.Conn]*models.User),
			},
		},
	}
}

func (s *MemoryStorage) GetUserById(id int) (*models.User, error) {
	um.Lock()
	if result, ok := s.users[id]; ok {
		um.Unlock()
		return result, nil
	} else {
		um.Unlock()
		return nil, &models.NotFoundError{Id: id}
	}
}

func (s *MemoryStorage) GetAllUsers() ([]*models.User, error) {
	var result []*models.User
	um.Lock()
	for _, user := range s.users {
		result = append(result, user)
	}
	um.Unlock()
	return result, nil
}

func (s *MemoryStorage) RemoveUserById(id int) (*models.User, error) {
	um.Lock()
	if user, err := s.GetUserById(id); err == nil {
		delete(s.users, id)
		um.Unlock()
		return user, nil
	} else {
		um.Unlock()
		return nil, &models.NotFoundError{Id: id}
	}
}

func (s *MemoryStorage) UpdateUser(u *models.User) (*models.User, error) {
	id := u.Id
	um.Lock()
	if _, ok := s.users[id]; ok {
		s.users[id] = u
		um.Unlock()
		return s.users[id], nil
	} else {
		um.Unlock()
		return nil, &models.NotFoundError{Id: id}
	}
}

func (s *MemoryStorage) GetAllRooms() ([]*models.Room, error){
	var result []*models.Room
	rm.Lock()
	for _, room := range s.rooms {
		result = append(result, room)
	}
	rm.Unlock()
	return result, nil
}

func (s *MemoryStorage) GetRoomById(id int) (*models.Room, error) {
	rm.Lock()
	if room, ok := s.rooms[id]; ok{
		rm.Unlock()
		return room, nil
	} else {
		rm.Unlock()
		return nil, &models.NotFoundError{Id:id}
	}
}

func (s *MemoryStorage) RemoveRoomById(id int) (*models.Room, error){
	rm.Lock()
	if room, err := s.GetRoomById(id); err == nil {
		delete(s.rooms, room.Id)
		rm.Unlock()
		return room, nil
	} else {
		rm.Unlock()
		return nil, &models.NotFoundError{Id:id}
	}
}
