package models

import (
	"fmt"
//	"iforgetgo/services"
)

type SmsChannel struct {
	Name     string
	provider services.ServiceProvider
}

func NewSmsChannel() *SmsChannel {
	runmode := beego.AppConfig.String("runmode") != "dev"
	provider := services.NewSmsServiceProvider(
		beego.AppConfig.String("sms_gate_url"),
		beego.AppConfig.String("sms_gate_user"),
		beego.AppConfig.String("sms_gate_pass"),
		beego.AppConfig.String("sms_gate_from"),
		runmode)

	return &SmsChannel{"Phone", provider}
}

func (this *SmsChannel) Send(message *ChannelMessage) {
	msg := this.prepareMessage(message.Message)
	this.provider.Send(message.UserName, message.Address, msg)
	fmt.Println("SmsChannel.Send: ", "Отправляем SMS с текстом \"", msg, "\"")
}

func (this *SmsChannel) GetName() string {
	return this.Name
}

/**
подготовка сообщения к отправке, например рендер темплейта
*/
func (this *SmsChannel) prepareMessage(message string) string {
	if len(message) > 70 {
		return message[:70]
	}
	return message
}
