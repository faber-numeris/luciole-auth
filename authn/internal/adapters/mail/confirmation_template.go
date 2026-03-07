package mail

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
)

type ConfirmationTemplate struct {
	ConfirmationURL string
	Code            string
}

func (t ConfirmationTemplate) Subject() string {
	return "Confirm your email address"
}

func (t ConfirmationTemplate) BodyHTML() string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button { 
            display: inline-block; 
            padding: 12px 24px; 
            background-color: #4F46E5; 
            color: white; 
            text-decoration: none; 
            border-radius: 4px; 
            margin: 20px 0;
        }
        .footer { font-size: 12px; color: #666; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Confirm your email address</h2>
        <p>Thank you for registering. Please confirm your email address by clicking the button below:</p>
        <p><a href="%s" class="button">Confirm Email</a></p>
        <p>Or use this code: <strong>%s</strong></p>
        <p>If you didn't create an account, please ignore this email.</p>
        <div class="footer">
            <p>This email was sent by Luciole Auth.</p>
        </div>
    </div>
</body>
</html>`, t.ConfirmationURL, t.Code)
}

func (t ConfirmationTemplate) BodyPlainText() string {
	return fmt.Sprintf(`Confirm your email address

Thank you for registering. Please confirm your email address by visiting the following link:

%s

Or use this code: %s

If you didn't create an account, please ignore this email.

This email was sent by Luciole Auth.`, t.ConfirmationURL, t.Code)
}

func (t ConfirmationTemplate) BuildMail(from string, to []string) Mail {
	return Mail{
		From:        from,
		To:          to,
		Subject:     t.Subject(),
		Body:        t.BodyPlainText(),
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
	partText.Write([]byte(tpl.BodyPlainText()))

	// 2. Create the HTML part
	htmlHeader := make(textproto.MIMEHeader)
	htmlHeader.Set("Content-Type", "text/html; charset=utf-8")
	partHtml, err := writer.CreatePart(htmlHeader)
	if err != nil {
		return "", "", err
	}
	partHtml.Write([]byte(tpl.BodyHTML()))

	// 3. Close to add the final --boundary--
	writer.Close()

	// writer.FormDataContentType() returns "multipart/form-data; boundary=..."
	// We want "multipart/alternative", so we swap the prefix:
	contentType := fmt.Sprintf("multipart/alternative; boundary=%s", writer.Boundary())

	return buf.String(), contentType, nil
}
