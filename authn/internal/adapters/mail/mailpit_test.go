package mail

import (
	"context"
	"errors"
	"mime/multipart"
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

	t.Run("invalid recipient", func(t *testing.T) {
		err := m.SendConfirmationEmail(ctx, domain.UserConfirmation{UserEmail: "invalid"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid recipient address")
	})
}

func TestMail_Methods(t *testing.T) {
	m := Mail{
		From:    "from@example.com",
		To:      []string{"to@example.com"},
		Subject: "Subject",
		Body:    "Body",
	}

	assert.True(t, m.Validate())
	msg := m.BuildMessage()
	assert.Contains(t, msg, "From: from@example.com")
	assert.Contains(t, msg, "Subject: Subject")

	m2 := Mail{
		From:    "from@example.com",
		To:      []string{"to@example.com"},
		Subject: "Subject",
		Body:    "Body",
		ContentType: "text/html",
	}
	msg2 := m2.BuildMessage()
	assert.Contains(t, msg2, "Content-Type: text/html")

	mEmpty := Mail{}
	assert.False(t, mEmpty.Validate())
	assert.False(t, (&Mail{From: "f"}).Validate())
	assert.False(t, (&Mail{From: "f", To: []string{"t"}}).Validate())
	assert.False(t, (&Mail{From: "f", To: []string{"t"}, Subject: "s"}).Validate())
}

func TestConfirmationTemplate(t *testing.T) {
	tpl := ConfirmationTemplate{
		ConfirmationURL: "http://confirm",
		Code:            "123",
	}

	assert.Equal(t, "Confirm your email address", tpl.Subject())
	
	html, err := tpl.BodyHTML()
	assert.NoError(t, err)
	assert.Contains(t, html, "http://confirm")

	text, err := tpl.BodyPlainText()
	assert.NoError(t, err)
	assert.Contains(t, text, "123")

	mail := tpl.BuildMail("from", []string{"to"})
	assert.Equal(t, "from", mail.From)
}

func TestBuildConfirmationMail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := BuildConfirmationMail("from", []string{"to"}, "url", "code")
		assert.Equal(t, "from", m.From)
		assert.Contains(t, m.ContentType, "multipart/alternative")
	})

	t.Run("multipart failure fallback", func(t *testing.T) {
		// Mock multipart.NewWriter to fail by patching its return if possible, 
		// but since it's a factory function returning a struct pointer, 
		// let's try to trigger an error in buildMultipartMail by patching CreatePart.
		
		dummyWriter := &multipart.Writer{}
		patches := gomonkey.ApplyFunc(multipart.NewWriter, func(any) *multipart.Writer {
			return dummyWriter
		})
		defer patches.Reset()
		
		patches.ApplyMethod(dummyWriter, "CreatePart", func(_ *multipart.Writer, _ any) (any, error) {
			return nil, errors.New("multipart error")
		})

		m := BuildConfirmationMail("from", []string{"to"}, "url", "code")
		assert.Equal(t, "from", m.From)
		assert.Equal(t, "text/plain; charset=utf-8", m.ContentType)
	})
}

func TestNewService(t *testing.T) {
	mockCfg := mocks.NewMockIMailConfig(t)
	svc := NewService(mockCfg)
	assert.NotNil(t, svc)
}
