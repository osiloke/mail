package mailers

import (
	"context"

	"github.com/apex/log"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// NewMailgunMailer create new mailgun mailer
func NewSendgridMailer(domain, key string) *SendgridMailer {
	mg := sendgrid.NewSendClient(key)
	return &SendgridMailer{mg}
}

// SendgridMailer a mailer that sends emails with sendgrid
type SendgridMailer struct {
	mg *sendgrid.Client
}

// Send send an email
func (m *SendgridMailer) Send(ctx context.Context, p *MailParams) error {
	from := mail.NewEmail(p.Sender, p.Sender)
	to := mail.NewEmail(p.Recipient, p.Recipient)
	message := mail.NewSingleEmail(from, p.Subject, to, p.Text, p.Html)
	resp, err := m.mg.Send(message)
	log.Debugf("Message Sent - %s - %s", resp.StatusCode, resp.Body)
	return err
}
