package mailgun

import (
	"context"
	"time"

	"github.com/fahmifan/mailmerger"
	"github.com/mailgun/mailgun-go/v4"
)

var _ mailmerger.Transporter = (*MailgunTransporter)(nil)

type MailgunClient interface {
	Send(ctx context.Context, m *mailgun.Message) (string, string, error)
	NewMessage(from, subject, text string, to ...string) *mailgun.Message
}

type MailgunTransporter struct {
	client MailgunClient
}

func (m *MailgunTransporter) Send(ctx context.Context, subject, from, to string, body []byte) (err error) {
	message := m.client.NewMessage(from, subject, "", to)
	message.SetHtml(string(body))
	ctx2, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err = m.client.Send(ctx2, message)
	return err
}
