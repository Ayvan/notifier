package models

type Message struct {
	NoticeId string
	Message  string
	Receiver User
}

func NewMessage(noticeId string, message string, receiver User) *Message {
	return &Message{noticeId, message, receiver}
}
