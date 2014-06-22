package models

type ChannelProvider interface {
	Send(message *ChannelMessage)
}
type Channel interface {
	ChannelProvider
	GetName() string
}
