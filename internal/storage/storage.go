package storage

import "chat-app/pkg/models"

type Storage interface {
	Get(int) *models.User
	GetAll() []*models.User
	Remove(int) *models.User
	Update(int, *models.User) *models.User
}
