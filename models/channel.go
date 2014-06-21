package models

type Channel interface {
	Send(message *ChannelMessage)
	GetId() int
}
