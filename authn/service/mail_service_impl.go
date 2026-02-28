package service

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/configuration"
)

type MailService struct {
	cfg configuration.IMailConfig
}

func NewMailService(cfg configuration.IMailConfig) IMailService {
	return &MailService{
		cfg: cfg,
	}
}

func (s *MailService) SendConfirmationEmail(ctx context.Context, email, code string) error {
	from := s.cfg.MailFrom()
	to := []string{email}
	subject := "Confirm your registration"
	body := fmt.Sprintf("Your confirmation code is: %s\n\nThis code will expire in 10 minutes.", code)

	msg := buildMessage(from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.cfg.MailHost(), s.cfg.MailPort())

	var auth smtp.Auth
	if s.cfg.MailUsername() != "" {
		auth = smtp.PlainAuth("", s.cfg.MailUsername(), s.cfg.MailPassword(), s.cfg.MailHost())
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to mail server: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set mail from: %w", err)
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %w", err)
	}
	defer wc.Close()

	_, err = bytes.NewReader([]byte(msg)).WriteTo(wc)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func buildMessage(from string, to []string, subject, body string) string {
	var buf bytes.Buffer
	buf.WriteString("From: " + from + "\r\n")
	buf.WriteString("To: " + to[0] + "\r\n")
	buf.WriteString("Subject: " + subject + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buf.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(body)
	return buf.String()
}

var _ IMailService = (*MailService)(nil)
