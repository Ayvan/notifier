package models

type Message struct {
	Id string
	Sender string
	Receiver *User
	Message string
}


func NewMessage(id string, sender string, receiver *User, message string) *Message {
	return &Message{id , sender , receiver , message }
}
