package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"iforgetgo/services"
)

type SmsChannel struct {
	Name     string
	provider services.ServiceProvider
	i        int64
}

func NewSmsChannel() *SmsChannel {
	runmode := beego.AppConfig.String("runmode") != "dev"
	provider := services.NewSmsServiceProvider(
		beego.AppConfig.String("sms_gate_url"),
		beego.AppConfig.String("sms_gate_user"),
		beego.AppConfig.String("sms_gate_pass"),
		beego.AppConfig.String("sms_gate_from"),
		runmode)

	return &SmsChannel{"phone", provider, 0}
}

func (this *SmsChannel) Send(message *ChannelMessage) {
	msg := this.prepareMessage(message.Message)
	err := this.provider.Send(message.UserName, message.Address, msg)
	this.i++
	if err == nil {
		fmt.Println("Отправлено sms: ", this.i)
	} else {
		fmt.Println("Ошибка отправки sms: ", err)
	}

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
