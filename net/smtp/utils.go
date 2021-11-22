package smtp

import (
	"bytes"
	"html/template"

	"github.com/intale-llc/foundation/errors"
)

func ParseHtmlTemplate(path string, data interface{}) ([]byte, error) {
	bodyTemplate, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	if err := bodyTemplate.Execute(buffer, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ParseHtmlTemplateMust(path string, data interface{}) []byte {
	result, err := ParseHtmlTemplate(path, data)
	if err != nil {
		errors.Panicf("smtp: unexpected html parse error: %v", err)
	}

	return result
}
