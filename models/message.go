package models

type Message struct {
	Id int
	Sender int
	Receiver *User
	Message string
}


func NewMessage(id int, sender int, receiver *User, message string) *Message {
	return &Message{id , sender , receiver , message }
}
