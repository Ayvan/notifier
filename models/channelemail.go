package models

import (
	"fmt"
)

type ChannelEmail struct {
	Id int
	Name string
	Params []string
}

func NewChannelEmail() ChannelEmail {
	params := []string{"124"}
	channel := ChannelEmail{1,"email",params}

	return channel
}

func (this ChannelEmail) Send(message *ChannelMessage) {
	fmt.Println(message.Message)
}

func (this ChannelEmail) GetId() int {
	return this.Id
}
