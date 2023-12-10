package models

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ValidateUser(u *User) bool {
	return true
}
