package mail

import (
	"bytes"
	"html/template"
	"net/http"
	"reflect"
	"sync"
)

type htmlMail struct {
	sender Sender

	model    Model
	template *template.Template

	buffer   *bytes.Buffer
	bufferMu sync.Mutex
}

func NewHtmlMail(sender Sender, model Model) (Mail, error) {
	htmlTemplate, err := template.ParseFiles(model.Template())
	if err != nil {
		return nil, err
	}

	// Trying to execute test data
	buffer := bytes.NewBuffer(nil)

	// Writing mail headers
	_ = model.Headers().Write(buffer)

	// Writing mail body
	if err := htmlTemplate.Execute(buffer, model); err != nil {
		return nil, err
	}

	buffer.Reset()
	return &htmlMail{
		sender:   sender,
		model:    model,
		template: htmlTemplate,
		buffer:   buffer,
	}, nil
}

func (m *htmlMail) Send(to string, data interface{}) error {
	if reflect.TypeOf(m.model) != reflect.TypeOf(data) {
		panic("mail: wrong data type")
	}

	m.bufferMu.Lock()

	// Writing mail headers
	headers := m.model.Headers()
	addHtmlHeaders(headers, to)

	_ = headers.Write(m.buffer)

	// Writing mail body
	// No error can occur here as we have tested
	// template execution inside a constructor
	_ = m.template.Execute(m.buffer, data)

	mail := m.buffer.Bytes()
	m.buffer.Reset()

	// Unlock before making request to mailing service
	// to avoid long mutex lock time
	m.bufferMu.Unlock()

	return m.sender.SendMail(mail)
}

func addHtmlHeaders(headers http.Header, to string) {
	headers.Set("MIME-version", "1.0")
	headers.Set("Content-Type", "text/html; charset=\"UTF-8\"")
	headers.Set("To", to)
}
