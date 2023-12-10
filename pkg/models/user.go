package models

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserDTO struct {
	Name string `json:"name"`
}
