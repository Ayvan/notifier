package models

type Channel interface {
	Send(message ChannelMessage)
}
