package services

import (
	"errors"
	"net/http"
	"net/url"
)

type SmsServiceProvider struct {
	host    string
	user    string
	pass    string
	from    string
	runmode bool
}

func NewSmsServiceProvider(host string, user string, pass string, from string, runmode bool) *SmsServiceProvider {

	return &SmsServiceProvider{host, user, pass, from, runmode}
}

func (this *SmsServiceProvider) Send(userName, address, message string) error {
	_ = userName
	// Подготовка данных для POST
	values := make(url.Values)
	values.Set("username", this.user)
	values.Add("password", this.pass)
	values.Add("to", address)
	values.Add("from", this.from)
	values.Add("coding", "2")
	values.Add("text", message)
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
