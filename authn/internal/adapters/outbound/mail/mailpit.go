package mail

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/ports/outbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
)

// Verify at compile time that Mailpit implements outboundport.Mailer
var _ outboundport.Mailer = (*Mailpit)(nil)

type Mailpit struct {
	configuration config.IMailConfig
}

func NewMailpit(cfg config.IMailConfig) *Mailpit {
	return &Mailpit{configuration: cfg}
}

func (m *Mailpit) SendConfirmationEmail(ctx context.Context, confirmation domain.UserConfirmation) error {
	from := m.configuration.MailFrom()
	to := []string{confirmation.UserEmail}
	subject := "Confirm your registration"
	body := fmt.Sprintf("Please confirm your registration by clicking on the following link: "+m.configuration.ConfirmationURLFormat(), confirmation.Token)

	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%d", m.configuration.SMTPHost(), m.configuration.SMTPPort())

	return smtp.SendMail(addr, nil, from, to, msg)
}

func (m *Mailpit) SendPasswordResetEmail(ctx context.Context, email string, token string) error {
	from := m.configuration.MailFrom()
	to := []string{email}
	subject := "Reset your password"
	body := fmt.Sprintf("Please reset your password by clicking on the following link: http://localhost:8080/reset-password?token=%s", token)

	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%d", m.configuration.SMTPHost(), m.configuration.SMTPPort())

	return smtp.SendMail(addr, nil, from, to, msg)
}
