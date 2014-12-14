package models

type ChannelMessage struct {
	NoticeId       string
	Channel  string
	Message  string
	Address  string
	UserName string
	//Channel *Channel
	//Message *Message

}

func NewChannelMessage(noticeId string, channel string, message string, address string, username string) *ChannelMessage {
	return &ChannelMessage{noticeId, channel, message, address, username}
}
