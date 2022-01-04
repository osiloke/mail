package mailers

import (
	"context"

	"github.com/apex/log"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

// NewMailjetMailer create new mailjet mailer
func NewMailjetMailer(apiKey, secretKey string) *MailjetMailer {
	mj := mailjet.NewMailjetClient(apiKey, secretKey)
	return &MailjetMailer{mj, "email"}
}

// MailjetMailer a mailer that sends emails with mailjet
type MailjetMailer struct {
	mj       *mailjet.Client
	customID string
}

// Send send an email
func (m *MailjetMailer) Send(ctx context.Context, p *MailParams) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: p.Sender,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: p.Recipient,
				},
			},
			Subject:  p.Subject,
			TextPart: p.Text,
			HTMLPart: p.Html,
			CustomID: m.customID,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := m.mj.SendMailV31(&messages)
	log.Debugf("Message Sent - %s", res)
	return err
}
