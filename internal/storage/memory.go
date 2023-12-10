package storage

import (
	"chat-app/pkg/loggers"
	"chat-app/pkg/models"
	"math/rand"
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
			1: {
				Id: 1,
				Name: "Boys",
				UserConnections: make(map[*websocket.Conn]*models.User),
			},
			2: {
				Id: 2,
				Name: "Girls",
				UserConnections: make(map[*websocket.Conn]*models.User),
			},
		},
	}
}

func (s *MemoryStorage) CreateUser(user models.UserDTO) (*models.User, error){
	um.Lock()
	// find if username is unique
	for _, usr := range s.users {
		if usr.Name == user.Name {
			um.Unlock()
			return nil, &models.BadRequestError{Message: "username already exists!"}
		}
	}

	// generate new ids until it finds a unique one
	id := rand.Int()
	for _, ok := s.users[id]; ok ;{
		id = rand.Int()
	}
	loggers.WarningLogger.Printf("found unique id: %d", id)

	newUser := &models.User{Id: id, Name: user.Name}
	s.users[id] = newUser
	um.Unlock()

	return newUser, nil
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
	if user, ok := s.users[id]; ok {
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

func (s *MemoryStorage) CreateRoom(room models.Room) (*models.Room, error){
	rm.Lock()
	// find if room name is unique
	for _, roomIt := range s.rooms {
		if roomIt.Name == room.Name {
			rm.Unlock()
			return nil, &models.BadRequestError{Message: "room already exists!"}
		}
	}

	// generate new ids until it finds a unique one
	id := rand.Int()
	for _, ok := s.users[id]; ok ;{
		id = rand.Int()
	}

	loggers.WarningLogger.Printf("found unique id: %d", id)
	newRoom := &models.Room{UserConnections: make(map[*websocket.Conn]*models.User),Name: room.Name,Id:id,}

	s.rooms[id] = newRoom
	rm.Unlock()
	
	return newRoom, nil
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
	if room, ok := s.rooms[id]; ok {
		delete(s.rooms, room.Id)
		rm.Unlock()
		return room, nil
	} else {
		rm.Unlock()
		return nil, &models.NotFoundError{Id:id}
	}
}
