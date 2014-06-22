package models

type Message struct {
	Id string
	Sender string
	Receiver string
	Message string
}


func NewMessage(id string, sender string, receiver string, message string) *Message {
	return &Message{id , sender , receiver , message }
}
