package mail

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/textproto"
	texttemplate "text/template"
)

//go:embed templates/confirmation_template.*
var templatesFS embed.FS

var (
	htmlTpl = template.Must(template.ParseFS(templatesFS, "templates/confirmation_template.html"))
	textTpl = texttemplate.Must(texttemplate.ParseFS(templatesFS, "templates/confirmation_template.txt"))
)

type ConfirmationTemplate struct {
	ConfirmationURL string
	Code            string
}

func (t ConfirmationTemplate) Subject() string {
	return "Confirm your email address"
}

func (t ConfirmationTemplate) BodyHTML() (string, error) {
	var buf bytes.Buffer
	if err := htmlTpl.Execute(&buf, t); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t ConfirmationTemplate) BodyPlainText() (string, error) {
	var buf bytes.Buffer
	if err := textTpl.Execute(&buf, t); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t ConfirmationTemplate) BuildMail(from string, to []string) Mail {
	body, err := t.BodyPlainText()
	if err != nil {
		body = "Error generating email body"
	}

	return Mail{
		From:        from,
		To:          to,
		Subject:     t.Subject(),
		Body:        body,
		ContentType: "text/plain; charset=utf-8",
	}
}

func BuildConfirmationMail(from string, to []string, confirmationURL string, code string) Mail {
	tpl := ConfirmationTemplate{
		ConfirmationURL: confirmationURL,
		Code:            code,
	}

	// We generate the body and the dynamic Content-Type header together
	body, contentType, err := buildMultipartMail(tpl)
	if err != nil {
		// Fallback to plain text if multipart fails
		return tpl.BuildMail(from, to)
	}

	return Mail{
		From:        from,
		To:          to,
		Subject:     tpl.Subject(),
		Body:        body,
		ContentType: contentType, // This will now include the boundary
	}
}

func buildMultipartMail(tpl ConfirmationTemplate) (string, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 1. Create the Plain Text part
	textHeader := make(textproto.MIMEHeader)
	textHeader.Set("Content-Type", "text/plain; charset=utf-8")
	partText, err := writer.CreatePart(textHeader)
	if err != nil {
		return "", "", err
	}
	textBody, err := tpl.BodyPlainText()
	if err != nil {
		return "", "", err
	}
	partText.Write([]byte(textBody))

	// 2. Create the HTML part
	htmlHeader := make(textproto.MIMEHeader)
	htmlHeader.Set("Content-Type", "text/html; charset=utf-8")
	partHtml, err := writer.CreatePart(htmlHeader)
	if err != nil {
		return "", "", err
	}
	htmlBody, err := tpl.BodyHTML()
	if err != nil {
		return "", "", err
	}
	partHtml.Write([]byte(htmlBody))

	// 3. Close to add the final --boundary--
	writer.Close()

	// writer.FormDataContentType() returns "multipart/form-data; boundary=..."
	// We want "multipart/alternative", so we swap the prefix:
	contentType := fmt.Sprintf("multipart/alternative; boundary=%s", writer.Boundary())

	return buf.String(), contentType, nil
}
