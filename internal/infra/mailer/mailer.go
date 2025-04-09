package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"html/template"

	"github.com/bernardinorafael/go-boilerplate/pkg/fault"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"
	"github.com/resend/resend-go/v2"
)

//go:embed "templates"
var templateFS embed.FS

type SendParams struct {
	From    string
	To      string
	Subject string
	File    string
	Data    any
}

type Config struct {
	APIKey     string
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration
}

type Mailer struct {
	ctx    context.Context
	log    logging.Logger
	client *resend.Client
	config Config
}

func New(ctx context.Context, log logging.Logger, config Config) *Mailer {
	return &Mailer{
		ctx:    ctx,
		log:    log,
		client: resend.NewClient(config.APIKey),
		config: config,
	}
}

func (m *Mailer) Send(p SendParams) error {
	tmplLocation := fmt.Sprintf("templates/%s", p.File)

	tmpl, err := template.New("email").ParseFS(templateFS, tmplLocation)
	if err != nil {
		m.log.Errorw(m.ctx, "error on parse template", logging.Err(err))
		return fault.New(
			"failed to parse template",
			fault.WithHTTPCode(http.StatusInternalServerError),
			fault.WithTag(fault.MAILER_ERROR),
			fault.WithError(err),
		)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, p.Data); err != nil {
		m.log.Errorw(m.ctx, "error on execute template", logging.Err(err))
		return fault.New(
			"failed to execute template",
			fault.WithHTTPCode(http.StatusInternalServerError),
			fault.WithTag(fault.MAILER_ERROR),
			fault.WithError(err),
		)
	}

	params := &resend.SendEmailRequest{
		From:    p.From,
		To:      []string{p.To},
		Html:    body.String(),
		Subject: p.Subject,
	}

	return m.send(params, m.config.MaxRetries)
}

func (m *Mailer) send(p *resend.SendEmailRequest, maxRetries int) error {
	var mailerErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(m.ctx, m.config.Timeout)
		defer cancel()

		_, err := m.client.Emails.SendWithContext(ctx, p)
		if err == nil {
			m.log.Infow(ctx, "email sent",
				logging.Int("attempt", attempt),
				logging.String("to", p.To[0]),
			)
			return nil
		}

		m.log.Errorw(ctx, "error on send email",
			logging.Err(err),
			logging.Int("attempt", attempt),
			logging.String("to", p.To[0]),
		)

		mailerErr = err
		time.Sleep(m.config.RetryDelay)
	}

	return fault.New(
		fmt.Sprintf("error on send email after %d attemps", maxRetries),
		fault.WithHTTPCode(http.StatusInternalServerError),
		fault.WithTag(fault.MAILER_ERROR),
		fault.WithError(mailerErr),
	)
}
