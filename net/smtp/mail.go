package smtp

import (
	"bytes"
	"net/http"
)

type Sender interface {
	SendMail(raw []byte) error
}

type Mail interface {
	// Subject returns mail subject
	Subject() string

	// Body returns mail body
	Body() []byte
}

func Send(address string, mail Mail, sender Sender) error {
	buffer := bytes.NewBuffer(nil)

	// Writing mail headers
	if err := defaultHeaders(mail.Subject(), address).Write(buffer); err != nil {
		return err
	}

	// Writing mail body
	if _, err := buffer.Write(mail.Body()); err != nil {
		return err
	}

	// Sending
	return sender.SendMail(buffer.Bytes())
}

func defaultHeaders(subject, to string) http.Header {
	headers := http.Header{}

	headers.Set("MIME-version", "1.0")
	headers.Set("Content-Type", "text/html; charset=\"UTF-8\"")
	headers.Set("Subject", subject)
	headers.Set("To", to)

	return headers
}
