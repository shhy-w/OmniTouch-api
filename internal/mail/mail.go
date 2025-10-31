package mail

import (
	"errors"

	"git.uozi.org/uozi/cosy-example-api/settings"
	"github.com/spf13/cast"
	"gopkg.in/gomail.v2"
)

// Dialer interface to allow mocking
type Dialer interface {
	DialAndSend(...*gomail.Message) error
}

func NewDialer() *gomail.Dialer {
	return gomail.NewDialer(settings.MailSettings.Host, cast.ToInt(settings.MailSettings.Port),
		settings.MailSettings.Email, settings.MailSettings.Password)
}

func SendNotification(subject, content string) (err error) {
	err = Send(NewDialer(), settings.MailSettings.To, subject, content)
	return
}

func SendToUser(to, subject, content string) (err error) {
	err = Send(NewDialer(), to, subject, content)
	return
}

func Send(d Dialer, to, subject, content string) error {
	if to == "" {
		return errors.New("cannot send to a empty email")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", settings.MailSettings.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	return d.DialAndSend(m)
}
