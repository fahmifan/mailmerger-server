package smtp

import (
	"context"
	_ "embed"

	"github.com/fahmifan/mailmerger"
	"github.com/flosch/pongo2"
	"gopkg.in/gomail.v2"
)

var _ mailmerger.Transporter = (*SMTP)(nil)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SMTP struct {
	mail     *gomail.Dialer
	cfg      *Config
	template *pongo2.Template
}

func NewSmtpClient(cfg *Config) (smtp *SMTP, err error) {
	smtp = &SMTP{
		cfg: cfg,
		mail: gomail.NewDialer(
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Password,
		),
	}

	smtp.template, err = pongo2.FromString(template)
	if err != nil {
		return
	}

	closer, err := smtp.mail.Dial()
	if err != nil {
		return nil, err
	}
	closer.Close()

	return smtp, nil
}

//go:embed template.html
var template string

func (m *SMTP) Send(ctx context.Context, subject, from, to string, body []byte) (err error) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	newBody, err := m.template.Execute(pongo2.Context{"body": string(body)})
	if err != nil {
		return
	}
	msg.SetBody("text/html", newBody)

	return m.mail.DialAndSend(msg)
}
