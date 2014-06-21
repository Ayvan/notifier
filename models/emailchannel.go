package models

import (
	"fmt"
)

type EmailChannel struct {
	Name string
	Params []string
}

func NewChannelEmail() EmailChannel {
	params := []string{"124"}
	return EmailChannel{"Email",params}
}

func (this EmailChannel) Send(message *ChannelMessage) {
	fmt.Println(message.Message)
}

func (this EmailChannel) GetName() string {
	return this.Name
}
