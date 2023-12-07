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
	users map[int]*models.User	// all users in the app
	rooms map[int]*models.Room		// all rooms in the app
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users: []*models.User{
			{
				Id:   0,
				Name: "Alice",
			},
			{
				Id:   1,
				Name: "Bob",
			},
			{
				Id:   2,
				Name: "Charlie",
			},
			{
				Id:   3,
				Name: "Dan",
			},
			{
				Id:   4,
				Name: "Elly",
			},
			{
				Id:   5,
				Name: "Finn",
			},
			{
				Id:   6,
				Name: "Gin",
			},
			{
				Id:   7,
				Name: "Henry",
			},
			{
				Id:   8,
				Name: "Irina",
			},
			{
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
	for _, user := range s.users {
		result = append(result, user)
	}
	return result, nil
}

func (s *MemoryStorage) RemoveUserById(id int) (*models.User, *models.NotFoundError) {
	if user, err := s.GetUserById(id); err != nil{
		delete(s.users, id)
		return user, nil
	} else {
		return nil, &models.NotFoundError{Id:id}
	}
}

func (s *MemoryStorage) UpdateUser(u *models.User) (*models.User, *models.NotFoundError) {
	id := u.Id
	if _, ok := s.users[id]; ok {
		s.users[id] = u
		return s.users[id], nil
	} else {
		return nil, &models.NotFoundError{Id:id}
	}
}
