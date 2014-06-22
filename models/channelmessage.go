package models

type ChannelMessage struct {
	Id      string
	Channel string
	Message string
	Address string
	UserName string
	//Channel *Channel
	//Message *Message

}

func NewChannelMessage(id string, channel string, message string, address string, username string) *ChannelMessage{
	return &ChannelMessage{id, channel, message, address, username}
}
