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
	return EmailChannel{"Mail",params}
}

func (this EmailChannel) Send(message *ChannelMessage) {
	fmt.Println("EmailChannel.Send: ","Отправляем сообщение с текстом '",message.Message,"'")
}

func (this EmailChannel) GetName() string {
	return this.Name
}
