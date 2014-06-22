package models

import (
	"fmt"
)

type EmailChannel struct {
	Name string
	Params []string
}

func NewEmailChannel() *EmailChannel {
	params := []string{"124"}
	return &EmailChannel{"Mail",params}
}

func (this *EmailChannel) Send(message *ChannelMessage) {
	msg := this.prepareMessage(message.Message)
	fmt.Println("EmailChannel.Send: ","Отправляем сообщение с текстом '",msg,"'")
}

func (this *EmailChannel) GetName() string {
	return this.Name
}

/**
	подготовка сообщения к отправке, например рендер темплейта
 */
func (this *EmailChannel) prepareMessage(message string) string {
	return message
}
