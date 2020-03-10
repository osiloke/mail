package mailers

import (
	"context"

	"github.com/apex/log"
	mailgun "github.com/mailgun/mailgun-go/v3"
)

// NewMailgunMailer create new mailgun mailer
func NewMailgunMailer(domain, key string) *MailgunMailer {
	mg := mailgun.NewMailgun(domain, key)
	return &MailgunMailer{mg}
}

// MailgunMailer a mailer that sends emails with mailgun
type MailgunMailer struct {
	mg *mailgun.MailgunImpl
}

// Send send an email
func (m *MailgunMailer) Send(ctx context.Context, sender, subject, text, recipient, html string) error {
	message := m.mg.NewMessage(sender, subject, text, recipient)
	message.SetTracking(true)
	message.SetHtml(html)
	resp, id, err := m.mg.Send(ctx, message)
	log.Debugf("Message Sent - %s - %s", id, resp)
	return err
}
