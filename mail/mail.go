package mail

import (
	"net/http"
)

type Mail interface {
	Send(to string, data interface{}) error
}

type Sender interface {
	SendMail(raw []byte) error
}

type Model interface {
	Template() string
	Headers() http.Header
}
