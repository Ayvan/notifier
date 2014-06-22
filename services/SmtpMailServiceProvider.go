package services

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"log"
	"encoding/base64"
)

type SmtpMailServiceProvider struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Auth        smtp.Auth
}

func NewSmtpMailServiceProvider(username, password, host, port string) *SmtpMailServiceProvider {
	return &SmtpMailServiceProvider{username, password, host, port, nil}
}

func (this *SmtpMailServiceProvider) Authenticate() {
	this.Auth = smtp.PlainAuth(
		"",
		this.Username,
		this.Password,
		this.Host,
	)
}

func (this *SmtpMailServiceProvider) Send(userName, userEmail, message string) {
	from := mail.Address{"iForget", this.Username}
	to := mail.Address{userName, userEmail}
	body := message
	subject := "Напоминание"

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

	smtpServer := this.Host + ":" + this.Port

	connection, error := smtp.Dial(smtpServer)
	if error != nil {
		log.Panic(error)
	}

	host, _, _ := net.SplitHostPort(smtpServer)
	tlc := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	if error = connection.StartTLS(tlc); error != nil {
		log.Panic(error)
	}

	if error = connection.Auth(this.Auth); error != nil {
		log.Panic(error)
	}

	if error = connection.Mail(from.Address); error != nil {
		log.Panic(error)
	}

	if error = connection.Rcpt(to.Address); error != nil {
		log.Panic(error)
	}

	dataCloser, error := connection.Data()
	if error != nil {
		log.Panic(error)
	}

	_, error = dataCloser.Write([]byte(fullMessage))
	if error != nil {
		log.Panic(error)
	}

	error = dataCloser.Close()
	if error != nil {
		log.Panic(error)
	}

	connection.Quit()
}
