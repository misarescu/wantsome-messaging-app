package storage

import "chat-app/pkg/models"

type Storage interface {
	GetUserById(int) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	RemoveUserById(int) (*models.User, error)
	UpdateUser(*models.User) (*models.User, error)
	GetAllRooms() ([]*models.Room, error)
	GetRoomById(int) (*models.Room, error)
	RemoveRoomById(int) (*models.Room, error)
}
