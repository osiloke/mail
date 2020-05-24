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
	"time"

	"github.com/Masterminds/sprig"
	"github.com/apex/log"
	strip "github.com/grokify/html-strip-tags-go"

	// "github.com/microcosm-cc/bluemonday"
	"github.com/osiloke/mail/mailers"
	"github.com/osiloke/mail/queues/machinery"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasttemplate"
)

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
}

//Params required
type Params struct {
	BodyTemplate      string `json:"bodyTemplate"`
	SubjectTemplate   string `json:"subjectTemplate"`
	RecipientTemplate string `json:"recipientTemplate"`
	Sender            string `json:"sender"`
}

func do(addonConfig, addonParams, data, traceID string) error {
	config := Config{}
	// ctx := context.Background()
	err := json.Unmarshal([]byte(addonConfig), &config)
	if err != nil {
		return err
	}

	params := Params{}
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
	log.Debugf("%s Mailer - New Email", config.Mailer)
	switch config.Mailer {
	case "smtp":
		port := 587
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
	subjectTpl := fasttemplate.New(params.SubjectTemplate, "[[", "]]")
	subject := subjectTpl.ExecuteString(d)
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

	// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
	html := body // p.Sanitize(body)
	text := strip.StripTags(html)
	// The message object allows you to add attachments and Bcc recipients

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	err = mc.Send(ctx, sender, subject, text, recipient, html)
	if err != nil {
		return fmt.Errorf("%s - sender: %s, subject: %s, recipient: %s", err, sender, subject, recipient)
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
