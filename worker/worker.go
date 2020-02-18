package worker

import (
	// "context"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"github.com/grokify/html-strip-tags-go"
	mailgun "github.com/mailgun/mailgun-go/v3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/osiloke/mail/queues/machinery"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasttemplate"
)

//Config required
type Config struct {
	Mailgun struct{ 
		Domain string `json:"domain"`
		Key    string `json:"key"`
	} `json:"mailgun"`
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
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(config.Mailgun.Domain, config.Mailgun.Key)
	if len(params.SubjectTemplate) == 0 {
		return errors.New("missing subject template")
	}
	if len(params.RecipientTemplate) == 0 {
		return errors.New("missing recipient template")
	}
	if len(params.BodyTemplate) == 0 {
		return errors.New("missing body template")
	}
	sender := params.Sender
	subjectTpl := fasttemplate.New(params.SubjectTemplate, "[[", "]]")
	subject := subjectTpl.ExecuteString(d)
	bodyData, _ := b64.StdEncoding.DecodeString(params.BodyTemplate)
	bd := string(bodyData)
	bodyTpl := fasttemplate.New(bd, "[[", "]]")
	body := bodyTpl.ExecuteString(d)
	recipient := gjson.Get(data, params.RecipientTemplate)

	p := bluemonday.UGCPolicy()
	p.AllowStandardURLs()

	// We only allow <p> and <a href="">
	p.AllowAttrs("href").OnElements("a")

	// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
	html := body// p.Sanitize(body)
	text := strip.StripTags(html)
	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, text, recipient.String())
	message.SetTracking(true)
	message.SetHtml(html)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
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
		"mail": do,
	})
}

// NewWorker new worker
func NewWorker(build string) *Worker {
	return &Worker{Build: build}
}
