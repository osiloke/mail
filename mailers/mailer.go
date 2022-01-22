package mailers

import "context"

//Config required
type Config struct {
	Mailer  string `json:"mailer"`
	Mailgun struct {
		Domain string `json:"domain"`
		Key    string `json:"key"`
	} `json:"mailgun"`
	Mailjet struct {
		ApiKey    string `json:"apiKey"`
		SecretKey string `json:"secretKey"`
	} `json:"mailjet"`
	Postmark struct {
		APIToken    string `json:"apiToken"`
		ServerToken string `json:"serverToken"`
	} `json:"postmark"`
	SMTP struct {
		Server   string `json:"server"`
		Username string `json:"username"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		SSL      bool   `json:"ssl"`
	} `json:"smtp"`
	Sendgrid struct {
		Domain string `json:"domain"`
		Key    string `json:"key"`
	} `json:"sendgrid"`
}

//Params required
type Params struct {
	BodyTemplate      string `json:"bodyTemplate"`
	SubjectTemplate   string `json:"subjectTemplate"`
	RecipientTemplate string `json:"recipientTemplate"`
	Sender            string `json:"sender"`
	Attachment        string `json:"attachment"`
}

// Mailer defines an object that can send a mail
type Mailer interface {
	Send(ctx context.Context, p *MailParams) error
}

type MailParams struct {
	Sender, Subject, Text, Recipient, Html, Attachment string
}
