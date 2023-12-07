package storage

import (
	"chat-app/pkg/models"
	"sync"
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
	}
}

func (s *MemoryStorage) GetUserById(id int) (*models.User, *models.NotFoundError) {
	um.Lock()
	if result, ok := s.users[id]; ok {
		um.Unlock()
		return result, nil
	} else {
		um.Unlock()
		return nil, &models.NotFoundError{Id: id}
	}
}

func (s *MemoryStorage) GetAllUsers() ([]*models.User, *models.NotFoundError) {
	var result []*models.User
	um.Lock()
	for _, user := range s.users {
		result = append(result, user)
	}
	um.Unlock()
	return result, nil
}

func (s *MemoryStorage) RemoveUserById(id int) (*models.User, *models.NotFoundError) {
	um.Lock()
	if user, err := s.GetUserById(id); err != nil {
		delete(s.users, id)
		um.Unlock()
		return user, nil
	} else {
		um.Unlock()
		return nil, &models.NotFoundError{Id: id}
	}
}

func (s *MemoryStorage) UpdateUser(u *models.User) (*models.User, *models.NotFoundError) {
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
