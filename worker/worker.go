package worker

import (
	// "context"
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/Masterminds/sprig"
	strip "github.com/grokify/html-strip-tags-go"

	// "github.com/microcosm-cc/bluemonday"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/cavaliercoder/grab"
	"github.com/osiloke/mail/mailers"
	"github.com/osiloke/mail/queues/machinery"
	"github.com/tidwall/gjson"
)

func do(addonConfig, addonParams, data, traceID string) error {
	config := mailers.Config{}
	// ctx := context.Background()
	err := json.Unmarshal([]byte(addonConfig), &config)
	if err != nil {
		return err
	}

	params := mailers.Params{}
	err = json.Unmarshal([]byte(addonParams), &params)
	if err != nil {
		return err
	}
	var d map[string]interface{}
	err = json.Unmarshal([]byte(data), &d)
	if err != nil {
		return err
	}
	var mc mailers.Mailer
	log.DEBUG.Printf("%s Mailer - New Email", config.Mailer)
	switch config.Mailer {
	case "smtp":
		port := 465
		if config.SMTP.Port > 0 {
			port = config.SMTP.Port
		}
		mc = mailers.NewSMTP(config.SMTP.Server, config.SMTP.Username, config.SMTP.Password, port)
	case "mailjet":
		mc = mailers.NewMailjetMailer(config.Mailjet.ApiKey, config.Mailjet.SecretKey)
	case "postmark":
		mc = mailers.NewPostmarkMailer(config.Postmark.ServerToken, config.Postmark.APIToken)
	case "mailgun":
		mc = mailers.NewMailgunMailer(config.Mailgun.Domain, config.Mailgun.Key)
	case "sendgrid":
		mc = mailers.NewSendgridMailer(config.Sendgrid.Domain, config.Sendgrid.Key)
	default:
		return errors.New("no mailer specified")
	}
	// Create an instance of the Mailgun Client
	if len(params.SubjectTemplate) == 0 {
		return errors.New("missing subject template")
	}
	if len(params.RecipientTemplate) == 0 {
		return errors.New("missing recipient template")
	}
	if len(params.BodyTemplate) == 0 {
		return errors.New("missing body template")
	}
	// p := bluemonday.UGCPolicy()
	// p.AllowStandardURLs()

	// // We only allow <p> and <a href="">
	// p.AllowAttrs("href").OnElements("a")

	sender := params.Sender
	subject := params.SubjectTemplate
	subjectTpl, err := template.New("").Funcs(sprig.FuncMap()).Delims("[[", "]]").Parse(params.SubjectTemplate)
	if err == nil {
		var tplBuffer bytes.Buffer
		if err := subjectTpl.Execute(&tplBuffer, d); err == nil {
			subject = tplBuffer.String()
		}
	}
	bodyData, _ := b64.StdEncoding.DecodeString(params.BodyTemplate)
	bd := string(bodyData)
	bodyTpl, err := template.New("").Funcs(sprig.FuncMap()).Delims("[[", "]]").Parse(bd)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	if err := bodyTpl.Execute(&tpl, d); err != nil {
		return err
	}

	body := tpl.String()
	recipient := params.RecipientTemplate
	r := gjson.Get(data, params.RecipientTemplate)
	if r.Exists() {
		recipient = r.String()
	}

	attachment := params.Attachment
	attachmentTpl, err := template.New("").Funcs(sprig.FuncMap()).Delims("[[", "]]").Parse(params.Attachment)
	if err == nil {
		var tplBuffer bytes.Buffer
		if err := attachmentTpl.Execute(&tplBuffer, d); err == nil {
			attachment = tplBuffer.String()
		}
	}

	if len(attachment) > 0 {
		resp, err := grab.Get(".", attachment)
		if err != nil {
			return err
		} else {
			attachment = resp.Filename
			defer os.Remove(resp.Filename)
		}
	}

	// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
	html := body // p.Sanitize(body)
	text := strip.StripTags(html)
	// The message object allows you to add attachments and Bcc recipients

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	log.DEBUG.Println("sending email to :", recipient)
	// Send the message	with a 10 second timeout
	err = mc.Send(ctx, &mailers.MailParams{Sender: sender, Subject: subject, Text: text,
		Recipient: recipient, Html: html, Attachment: attachment})
	if err != nil {
		err2 := fmt.Errorf("%s - sender: %s, subject: %s, recipient: %s", err, sender, subject, recipient)
		log.ERROR.Println(err, err2)
		return err2
	}
	return nil
}

// Worker a fcm worker that sends messages to centrifuge
type Worker struct {
	ID    string `help:"worker id"`
	Build string `help:"build"`
}

// Run run the worker
func (w *Worker) Run() error {
	return machinery.Worker(w.ID, map[string]interface{}{
		"email": do,
	})
}

// NewWorker new worker
func NewWorker(build string) *Worker {
	return &Worker{Build: build}
}
