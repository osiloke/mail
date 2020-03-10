package mailers

import "context"

// Mailer defines an object that can send a mail
type Mailer interface {
	Send(ctx context.Context, sender, subject, text, recipient, html string) error
}
