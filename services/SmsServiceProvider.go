package services

import (
	"net/http"
	"net/url"
	"errors"
)

type SmsServiceProvider struct {
	host       string
	user       string
	pass       string
	from       string
	runmode    bool

}

func NewSmsServiceProvider(host string, user string, pass string, from string, runmode bool) *SmsServiceProvider {

	return &SmsServiceProvider{host , user , pass, from, runmode}
}

func (this *SmsServiceProvider) Send(userName, address, message string) {
	_ = userName
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
	if this.runmode {
		resp, error := http.PostForm(this.host, values)
		if error != nil {
			return error
		}
		defer resp.Body.Close()

		if resp.ContentLength == 0 {
			return errors.New("SMS Gate Error")
		}
		return nil
	}
	return nil
}
