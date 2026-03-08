package config

type IMailConfig interface {
	SMTPHost() string
	SMTPPort() int
	SMTPUsername() string
	SMTPPassword() string
	MailFrom() string
	MailFromName() string
	ConfirmationURLFormat() string
}

var _ IMailConfig = (*MailConfig)(nil)

type MailConfig struct {
	SMTPHost_              string `env:"SMTP_HOST,required"`
	SMTPPort_              int    `env:"SMTP_PORT" envDefault:"1025"`
	SMTPUsername_          string `env:"SMTP_USERNAME"`
	SMTPPassword_          string `env:"SMTP_PASSWORD"`
	MailFrom_              string `env:"MAIL_FROM,required"`
	MailFromName_          string `env:"MAIL_FROM_NAME,required"`
	ConfirmationURLFormat_ string `env:"CONFIRMATION_URL_FORMAT" envDefault:"http://localhost:8080/v1/authn/confirm?token=%s"`
}

func (c MailConfig) SMTPHost() string {
	return c.SMTPHost_
}

func (c MailConfig) SMTPPort() int {
	return c.SMTPPort_
}

func (c MailConfig) SMTPUsername() string {
	return c.SMTPUsername_
}

func (c MailConfig) SMTPPassword() string {
	return c.SMTPPassword_
}

func (c MailConfig) MailFrom() string {
	return c.MailFrom_
}

func (c MailConfig) MailFromName() string {
	return c.MailFromName_
}

func (c MailConfig) ConfirmationURLFormat() string {
	return c.ConfirmationURLFormat_
}
