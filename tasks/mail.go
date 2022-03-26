package tasks

import (
	"github.com/jordan-wright/email"
	"github.com/myOmikron/RustymonBackend/configs"
	"github.com/myOmikron/echotools/worker"
	"net"
	"net/smtp"
	"strconv"
	"time"
)

func NewMailTask(mailConfig *configs.Mail, receiver []string, subject string, body string) worker.Task {
	return worker.NewTask(func() error {
		// Authentication.
		authentication := smtp.PlainAuth("", mailConfig.User, mailConfig.Password, mailConfig.Host)

		// Sending email.
		e := email.NewEmail()
		e.From = mailConfig.From
		e.To = receiver
		e.Subject = subject
		e.Text = []byte(body)
		e.Headers = map[string][]string{
			"Date": {time.Now().Format(time.RFC1123Z)},
		}

		return e.Send(net.JoinHostPort(mailConfig.Host, strconv.Itoa(int(mailConfig.Port))), authentication)
	})
}
