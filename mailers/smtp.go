package mailers

import (
	"context"

	"github.com/go-gomail/gomail"
)

// NewSMTP new smtp mailer
func NewSMTP(server, username, password string, port int) *SMTP {
	d := gomail.NewDialer(server, port, username, password)
	return &SMTP{
		mailer: d,
	}
}

// SMTP sends mails using cmtp
type SMTP struct {
	mailer *gomail.Dialer
}

// Send send an email
func (s *SMTP) Send(ctx context.Context, sender, subject, text, recipient, html string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetHeader("text/plain", text)
	m.SetBody("text/html", html)

	err := s.mailer.DialAndSend(m)
	return err
}
