package storage

import "chat-app/pkg/models"

type Storage interface {
	GetUserById(int) *models.User
	GetAllUsers() []*models.User
	RemoveUserById(int) *models.User
	UpdateUser(*models.User) *models.User
}
