package mail

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"strings"
	"time"

	"html/template"

	"github.com/charmbracelet/log"
	"github.com/resend/resend-go/v2"
)

//go:embed "templates"
var templateFS embed.FS

type SendParams struct {
	From    EmailSender
	To      string
	Subject string
	File    Template
	Data    any
}

type Config struct {
	APIKey  string
	Timeout time.Duration
}

type Mail struct {
	ctx     context.Context
	log     *log.Logger
	client  *resend.Client
	timeout time.Duration
}

func New(ctx context.Context, log *log.Logger, APIKey string, timeout time.Duration) *Mail {
	return &Mail{
		ctx:     ctx,
		log:     log,
		timeout: timeout,
		client:  resend.NewClient(APIKey),
	}
}

func (m *Mail) Send(p SendParams) error {
	tmplLocation := fmt.Sprintf("templates/%s", p.File)

	tmpl, err := template.New("email").ParseFS(templateFS, tmplLocation)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, p.Data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    string(p.From),
		To:      []string{p.To},
		Html:    body.String(),
		Subject: p.Subject,
	}

	_, err = m.client.Emails.SendWithContext(m.ctx, params)
	if err == nil {
		m.log.Debug(
			"email sent successfully",
			"details", strings.Join(
				[]string{
					fmt.Sprintf("from: %s", p.From),
					fmt.Sprintf("to: %s", p.To),
					fmt.Sprintf("subject: %s", p.Subject),
				},
				"\n",
			),
		)

		return nil
	}

	return fmt.Errorf("failed to send email: %w", err)
}
