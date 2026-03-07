package mail

import (
	"bytes"
	"time"
)

type Mail struct {
	From        string
	To          []string
	Subject     string
	Body        string
	ContentType string
}

func (e *Mail) BuildMessage() string {
	var buf bytes.Buffer
	buf.WriteString("From: " + e.From + "\r\n")
	buf.WriteString("To: " + e.To[0] + "\r\n")
	buf.WriteString("Subject: " + e.Subject + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	if e.ContentType == "" {
		e.ContentType = "text/plain; charset=\"utf-8\""
	}
	buf.WriteString("Content-Type: " + e.ContentType + "\r\n")
	buf.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(e.Body)
	return buf.String()
}

func (e *Mail) Validate() bool {
	if e.From == "" || len(e.To) == 0 || e.Subject == "" || e.Body == "" {
		return false
	}
	return true
}
