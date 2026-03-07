package mail

import (
	"context"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/faber-numeris/luciole-auth/authn/configuration"
	"github.com/faber-numeris/luciole-auth/authn/model"
)

type Mailpit struct {
	configuration configuration.IMailConfig
}

func NewMailpit(cfg configuration.IMailConfig) IMailService {
	return &Mailpit{configuration: cfg}
}

func (m *Mailpit) SendConfirmationEmail(ctx context.Context, confirmation model.UserConfirmation) error {
	from := m.configuration.MailFrom()
	fromName := m.configuration.MailFromName()
	smtpHost := m.configuration.SMTPHost()
	smtpPort := m.configuration.SMTPPort()

	toAddr, err := mail.ParseAddress(strings.TrimSpace(confirmation.UserEmail))
	if err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}

	fromAddr := fmt.Sprintf("%s <%s>", fromName, from)
	to := []string{toAddr.Address}

	built := BuildConfirmationMail(fromAddr, to, confirmation.Token, confirmation.Token)
	msg := built.BuildMessage()

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	// In the local development setup with Mailpit, there is usually no
	// authentication or TLS required.
	// The smtp.PlainAuth method can be used, but with empty username and password.
	//auth := smtp.PlainAuth("", "", "", addr)

	// Send the email
	return smtp.SendMail(addr, nil, from, to, []byte(msg))

}
