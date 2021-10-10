package mail

import (
	"bytes"
	"html/template"
	"net/http"
)

type HtmlModel interface {
	// File returns path to HTML template
	File() string

	// Headers returns additional headers for the mail
	Headers() http.Header

	// Senders returns all possible senders to send mail
	Senders() []Sender
}

type Sender interface {
	SendMail(raw []byte) error
}

func SendHTML(model HtmlModel, to string) error {
	htmlTemplate, err := template.ParseFiles(model.File())
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)

	// Writing mail headers
	headers := model.Headers()
	addDefaultHeaders(headers, to)

	if err := headers.Write(buffer); err != nil {
		return err
	}

	// Writing mail body
	if err := htmlTemplate.Execute(buffer, model); err != nil {
		return err
	}

	var sendErr error

	// Trying to send on all resources
	senders := model.Senders()
	for _, sender := range senders {
		sendErr = sender.SendMail(buffer.Bytes())

		if sendErr == nil {
			break
		}
	}

	return sendErr
}

func addDefaultHeaders(headers http.Header, to string) {
	headers.Set("MIME-version", "1.0")
	headers.Set("Content-Type", "text/html; charset=\"UTF-8\"")
	headers.Set("To", to)
}
