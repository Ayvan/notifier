package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"iforgetgo/services"
)

type EmailChannel struct {
	Name string
	provider services.ServiceProvider

}

func NewEmailChannel() *EmailChannel {
	provider := services.NewSmtpMailServiceProvider(beego.AppConfig.String("gmailUsername"), beego.AppConfig.String("gmailPassword"), beego.AppConfig.String("gmailHost"), beego.AppConfig.String("gmailPort"))
	return &EmailChannel{"Mail",provider}
}

func (this *EmailChannel) Send(message *ChannelMessage) {
	msg := this.prepareMessage(message.Message)
	this.provider.Send(message.UserName, message.Address, msg)
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
