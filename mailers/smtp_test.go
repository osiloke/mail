package mailers

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-gomail/gomail"
)

func TestNewSMTP(t *testing.T) {
	type args struct {
		server   string
		username string
		password string
		port     int
	}
	tests := []struct {
		name string
		args args
		want *SMTP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSMTP(tt.args.server, tt.args.username, tt.args.password, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSMTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSMTP_Send(t *testing.T) {
	type fields struct {
		mailer *gomail.Dialer
	}
	type args struct {
		ctx       context.Context
		sender    string
		subject   string
		text      string
		recipient string
		html      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"test",
			fields{gomail.NewDialer("smtp.mailtrap.io", 587, "632af48614f62d", "3314c577d16a37")},
			args{context.Background(), "me@osiloke.com", "Hello", "hello", "me@osiloke.com", "<b>hello</b>"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SMTP{
				mailer: tt.fields.mailer,
			}
			if err := s.Send(tt.args.ctx, tt.args.sender, tt.args.subject, tt.args.text, tt.args.recipient, tt.args.html); (err != nil) != tt.wantErr {
				t.Errorf("SMTP.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
