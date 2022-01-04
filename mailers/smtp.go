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
func (s *SMTP) Send(ctx context.Context, p *MailParams) error {
	m := gomail.NewMessage()
	m.SetHeader("From", p.Sender)
	m.SetHeader("To", p.Recipient)
	m.SetHeader("Subject", p.Subject)
	m.SetHeader("text/plain", p.Text)
	m.SetBody("text/html", p.Html)
	if len(p.Attachment) > 0 {
		m.Attach(p.Attachment)
	}

	err := s.mailer.DialAndSend(m)
	return err
}
