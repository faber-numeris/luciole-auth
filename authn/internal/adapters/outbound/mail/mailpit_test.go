package mail

import (
	"context"
	"net/smtp"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestMailpit_SendConfirmationEmail(t *testing.T) {
	mockCfg := mocks.NewMockIMailConfig(t)
	mockCfg.EXPECT().MailFrom().Return("from@example.com").Maybe()
	mockCfg.EXPECT().MailFromName().Return("From Name").Maybe()
	mockCfg.EXPECT().SMTPHost().Return("localhost").Maybe()
	mockCfg.EXPECT().SMTPPort().Return(1025).Maybe()
	mockCfg.EXPECT().SMTPUsername().Return("user").Maybe()
	mockCfg.EXPECT().SMTPPassword().Return("pass").Maybe()
	mockCfg.EXPECT().ConfirmationURLFormat().Return("http://localhost/confirm/%s").Maybe()

	m := NewMailpit(mockCfg)
	ctx := context.Background()
	conf := domain.UserConfirmation{
		UserEmail: "test@example.com",
		Token:     "token-123",
	}

	t.Run("success", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(smtp.SendMail, func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			assert.Equal(t, "localhost:1025", addr)
			return nil
		})
		defer patches.Reset()

		err := m.SendConfirmationEmail(ctx, conf)
		assert.NoError(t, err)
	})
}

func TestNewService(t *testing.T) {
	svc := NewService()
	assert.NotNil(t, svc)
}
