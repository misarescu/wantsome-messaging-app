package models

import "fmt"

type NotFoundError struct {
	Id int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Item not found with id: %d", e.Id)
}

type BroadcastError struct {
	Users []*User
}

func (e *BroadcastError) Error() string {
	var idSlice []int
	for _, user := range e.Users {
		idSlice = append(idSlice, user.Id)
	}
	return fmt.Sprintf("Error broadcasting to connections with ids:\n %d\n", idSlice)
}

type ConnectionError struct {}

func (e *ConnectionError) Error() string {
	return "Connection error"
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad request: %s", e.Message)
}