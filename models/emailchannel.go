package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"notifier/services"
)

type EmailChannel struct {
	Name     string
	provider services.ServiceProvider
	i        int64
}

func NewEmailChannel() *EmailChannel {
	runmode := beego.AppConfig.String("runmode") != "dev"
	provider := services.NewSmtpMailServiceProvider(
		beego.AppConfig.String("gmailUsername"),
		beego.AppConfig.String("gmailPassword"),
		beego.AppConfig.String("gmailHost"),
		beego.AppConfig.String("gmailPort"),
		runmode)
	return &EmailChannel{"Email", provider, 0}
}

func (this *EmailChannel) Send(message *ChannelMessage) {
	msg := this.prepareMessage(message.Message)
	err := this.provider.Send(message.UserName, message.Address, msg)

	if err == nil {
		this.i++
		fmt.Println("Отправлено email: ", this.i)
	} else {
		fmt.Println("Ошибка отправки email", err)
	}
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
