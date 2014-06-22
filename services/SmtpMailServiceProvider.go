package services

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"encoding/base64"
)

type SmtpMailServiceProvider struct {
	username    string
	password    string
	host        string
	port        string
	auth        smtp.Auth
	runmode     bool
}

func NewSmtpMailServiceProvider(username, password, host, port string, runmode bool) *SmtpMailServiceProvider {
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		host,
	)
	return &SmtpMailServiceProvider{username, password, host, port, auth, runmode}
}

func (this *SmtpMailServiceProvider) Send(userName, address, message string) error {
	from := mail.Address{"iForget", this.username}
	to := mail.Address{userName, address}
	body := message
	subject := "Напоминание"

	//если не "боевой" режим, то дальше ничего не делаем
	if !this.runmode {
		return nil
	}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	fullMessage := ""

	for k, v := range header {
		fullMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	fullMessage += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	smtpServer := this.host + ":" + this.port

	connection, error := smtp.Dial(smtpServer)
	if error != nil {
		return error
	}

	host, _, _ := net.SplitHostPort(smtpServer)
	tlc := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	if error = connection.StartTLS(tlc); error != nil {
		return error
	}

	if error = connection.Auth(this.auth); error != nil {
		return error
	}

	if error = connection.Mail(from.Address); error != nil {
		return error
	}

	if error = connection.Rcpt(to.Address); error != nil {
		return error
	}

	dataCloser, error := connection.Data()
	if error != nil {
		return error
	}

	_, error = dataCloser.Write([]byte(fullMessage))
	if error != nil {
		return error
	}

	error = dataCloser.Close()
	if error != nil {
		return error
	}

	connection.Quit()

	return nil
}
