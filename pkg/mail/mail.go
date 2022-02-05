package mail

import (
	"gohub/pkg/config"
	"sync"
)

// From 发送对象
type From struct {
	Address string
	Name    string
}

// Email
type Email struct {
	From    From
	To      []string
	Bcc     []string
	Cc      []string
	Subject string
	Text    []byte // Plaintext message (optional)
	HTML    []byte // Html message (optional)
}

// Mailer
type Mailer struct {
	Driver Driver
}

var once sync.Once

var internalMailer *Mailer

// NewMailer 单例模式获取
func NewMailer() *Mailer {
	once.Do(func() {
		internalMailer = &Mailer{
			Driver: &SMTP{},
		}
	})
	return internalMailer
}

// Send 发送邮件
func (mailer *Mailer) Send(email Email) bool {
	return mailer.Driver.Send(email, config.GetStringMapString("mail.smtp"))
}
