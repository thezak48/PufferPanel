package email

import (
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/gofrs/uuid/v5"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"strings"
)

type smtpProvider struct {
	Provider
}

func init() {
	providers["smtp"] = smtpProvider{}
}

func (smtpProvider) Send(to, subject, body string) error {
	from := config.EmailFrom.Value()
	if from == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailFrom.Key())
	}

	host := config.EmailHost.Value()
	if host == "" {
		return pufferpanel.ErrSettingNotConfigured(config.EmailHost.Key())
	}

	var auth sasl.Client
	if username := config.EmailUsername.Value(); username != "" {
		auth = sasl.NewPlainClient("", username, config.EmailPassword.Value())
	} else {
		auth = sasl.NewAnonymousClient("")
	}

	refId, _ := uuid.NewV4()
	refIdStr := strings.ReplaceAll(refId.String(), "-", "")

	str := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Message-ID: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"+
		from, to, refIdStr, subject, body,
	)

	data := strings.NewReader(str)
	return smtp.SendMail(host, auth, from, []string{to}, data)
}
