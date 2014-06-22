package models

import (
	"fmt"
)

type SmsChannel struct {
	Name string
	Params []string
}

func NewSmsChannel() *SmsChannel {
	params := []string{"124"}
	return &SmsChannel{"Phone",params}
}

func (this *SmsChannel) Send(message *ChannelMessage) {
	msg := this.prepareMessage(message.Message)
	fmt.Println("SmsChannel.Send: ","Отправляем SMS с текстом '",msg,"'")
}

func (this *SmsChannel) GetName() string {
	return this.Name
}

/**
	подготовка сообщения к отправке, например рендер темплейта
 */
func (this *SmsChannel) prepareMessage(message string) string {
	return message
}
