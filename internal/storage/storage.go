package storage

import "chat-app/pkg/models"

type Storage interface {
	GetUserById(int) (*models.User, *models.NotFoundError)
	GetAllUsers() ([]*models.User, *models.NotFoundError)
	RemoveUserById(int) (*models.User, *models.NotFoundError)
	UpdateUser(*models.User) (*models.User, *models.NotFoundError)
	GetAllRooms() ([]*models.Room, *models.NotFoundError)
}
