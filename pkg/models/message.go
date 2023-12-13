package models

type UserMessage struct {
	Message string `json:"message"`
	UserId  int    `json:"user_id"`
}

type RoomMessage struct {
	UserMessage UserMessage `json:"user_message"`
	RoomId      int         `json:"room_id"`
}

type ResponseMessage struct {
	Message     string `json:"message"`
	FromUser    string `json:"from_user,omitempty"`
	ErrorStatus bool   `json:"error"`
}
