package storage

import "chat-app/pkg/models"

type MongoStorage struct{}

func (s *MongoStorage) Get(id int) *models.User {
	return &models.User{
		Id:   1,
		Name: "Foo",
	}
}
