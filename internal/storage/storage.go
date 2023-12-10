package storage

import "chat-app/pkg/models"

type Storage interface {
	CreateUser(user models.UserDTO) (*models.User, error)
	GetUserById(int) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	RemoveUserById(int) (*models.User, error)
	UpdateUser(*models.User) (*models.User, error)
	CreateRoom(models.Room) (*models.Room, error)
	GetAllRooms() ([]*models.Room, error)
	GetRoomById(int) (*models.Room, error)
	RemoveRoomById(int) (*models.Room, error)
}
