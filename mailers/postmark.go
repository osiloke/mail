package mailers

import (
	"context"

	"github.com/apex/log"
	"github.com/keighl/postmark"
)

// NewPostmarkMailer create new mailjet mailer
func NewPostmarkMailer(serverToken, accountToken string) *PostmarkMailer {
	pm := postmark.NewClient(serverToken, accountToken)
	return &PostmarkMailer{pm, "email"}
}

// PostmarkMailer a mailer that sends emails with mailjet
type PostmarkMailer struct {
	pm       *postmark.Client
	customID string
}

// Send send an email
func (m *PostmarkMailer) Send(ctx context.Context, sender, subject, text, recipient, html string) error {
	email := postmark.Email{
		From:       sender,
		To:         recipient,
		Subject:    subject,
		HtmlBody:   html,
		TextBody:   text,
		Tag:        m.customID,
		TrackOpens: true,
	}

	res, err := m.pm.SendEmail(email)
	if err == nil {
		log.Debugf("Message Sent - %s", res)
	}
	return err
}
