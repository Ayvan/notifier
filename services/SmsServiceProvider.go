package services

import (
	"log"
	"fmt"
	"net/http"
	"net/url"
	"iforgetgo/models"
)

type SmsServiceProvider struct {
	host       string // https://www.stramedia.ru/modules/send_sms.php
	user       string // kreddy
	pass       string // Ht2411s
	from       string //89269116791

}

func NewSmsServiceProvider(host string, user string, pass  string) *SmsServiceProvider {

	return &SmsServiceProvider{host , user , pass, "89269116791"}
}

func (this *SmsServiceProvider) Send(message *models.ChannelMessage) {

	// Подготовка данных для POST
	values := make(url.Values)
	values.Set("username", this.user)
	values.Add("password", this.pass)
	values.Add("to", message.Address)
	values.Add("from", this.pass)
	values.Add("coding", "2")
	values.Add("text", message.Message)
	values.Add("priority", "0")
	values.Add("mclass", "1")
	values.Add("dlrmask", "31")

	// Submit form
	resp, err := http.PostForm(this.host, values)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.ContentLength == 0 {
		log.Fatal("SMS Gate Error")
	}


}
