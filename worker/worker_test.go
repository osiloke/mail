package worker

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_do(t *testing.T) {
	data, err := ioutil.ReadFile("./test.html.b64")
	if err != nil {
		panic(err)
	}
	body := string(data)
	params := fmt.Sprintf(`{"bodyTemplate": "%s", "sender": "osi@progwebtech.com", "subjectTemplate": "Welcome", "recipientTemplate": "recipient"}`, body)
	type args struct {
		addonConfig string
		addonParams string
		data        string
		traceID     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{
				`{"mailgun":{"domain":"sandbox5cd5276fe1664858921893702b5cb2f7.mailgun.org", "key": "key-537107024f22430067f5f40e48ffdb76"}}`,
				params,
				`{"UserName": "user", "recipient":"me@osiloke.com","FirstName":"Osiloke","CompanyName":"Test Company","ContactEmail":"test@testcompany.com"}`,
				`trace`,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := do(tt.args.addonConfig, tt.args.addonParams, tt.args.data, tt.args.traceID); (err != nil) != tt.wantErr {
				t.Errorf("do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
