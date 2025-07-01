package notification

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/mail"
	"strings"

	"github.com/machilan1/cruise/internal/business/sdk/mailer"
)

//go:embed templates
var templates embed.FS

type Core struct {
	mailer         mailer.Mailer
	frontendOrigin string
}

func NewCore(mailer mailer.Mailer, frontendOrigin string) *Core {
	return &Core{
		mailer:         mailer,
		frontendOrigin: frontendOrigin,
	}
}

func (c *Core) SendPasswordResetEmail(ctx context.Context, email mail.Address, token string) error {
	tmpl, err := template.ParseFS(templates, "templates/password-reset.tmpl")
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	data := struct {
		Link string
	}{
		// Note: the frontend developer should ensure that the link is valid.
		Link: fmt.Sprintf("%s/auth/reset-password?token=%s", c.frontendOrigin, token),
	}
	w := new(strings.Builder)
	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	msg := mailer.Message{
		To:      []string{email.Address},
		Subject: "重設您的密碼",
		Body:    w.String(),
	}
	if _, err := c.mailer.Send(ctx, msg); err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
