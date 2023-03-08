package main

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

type mailServer struct {
	host        string
	port        int
	username    string
	password    string
	ConnTimeOut time.Duration
	SendTimeOut time.Duration
}

type message struct {
	from    string
	to      string
	message string
}

func (config *Config) sendMail(msg message) {
	config.wg.Add(1)
	config.mailSendChan <- msg
}

func (config *Config) concurrencyMailing() {
	for {
		select {
		case msg := <-config.mailSendChan:
			go config.mailing(msg)
		case <-config.mailDoneChan:
			return
		}
	}
}

func (config *Config) mailing(msg message) {
	defer config.wg.Done()
	client := config.setUpMailClient()
	email := config.newMsg(msg)
	err := email.Send(client)
	config.err(err)
}

func (config *Config) setUpMailClient() *mail.SMTPClient {
	server := mail.NewSMTPClient()
	server.Host = config.MailServer.host
	server.Port = config.MailServer.port
	server.Username = config.MailServer.username
	server.Password = config.MailServer.password
	server.Encryption = mail.EncryptionNone
	server.KeepAlive = false
	server.ConnectTimeout = config.MailServer.ConnTimeOut
	server.SendTimeout = config.MailServer.SendTimeOut
	client, err := server.Connect()
	config.err(err)
	return client
}
func (config *Config) newMsg(msg message) *mail.Email {
	email := mail.NewMSG()
	email.SetFrom(msg.from).AddTo(msg.to)
	email.SetBody(mail.TextPlain, msg.message)
	config.infoLog.Println(email.GetFrom(), email.GetMessage())
	config.errorLog.Println(email.GetError())
	return email
}
