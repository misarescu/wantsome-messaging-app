package storage

import "chat-app/pkg/models"

type MemoryStorage struct {
	users []*models.User
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

func (s *MemoryStorage) GetUserById(id int) *models.User {
	return s.users[id]
}

func (s *MemoryStorage) GetAllUsers() []*models.User {
	return s.users
}

func (s *MemoryStorage) RemoveUserById(id int) *models.User {
	user := s.GetUserById(id)

	s.users = append(s.users[:id], s.users[id+1:]...)

	return user
}

func (s *MemoryStorage) UpdateUser(u *models.User) *models.User {
	id := u.Id
	s.users[id] = u

	return s.users[id]
}
