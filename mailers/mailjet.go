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
func (m *MailjetMailer) Send(ctx context.Context, sender, subject, text, recipient, html string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: sender,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: recipient,
				},
			},
			Subject:  subject,
			TextPart: text,
			HTMLPart: html,
			CustomID: m.customID,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := m.mj.SendMailV31(&messages)
	log.Debugf("Message Sent - %s", res)
	return err
}
