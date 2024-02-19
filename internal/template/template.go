package template

import (
	"bytes"
	"html/template"
)

var (
	registrationTmpl *template.Template
)

func PrepareTemplates() error {
	t, err := template.ParseFiles("./html/registration.html")
	if err != nil {
		return err
	}
	registrationTmpl = t
	return nil
}

func Register(link string) ([]byte, error) {
	var buff bytes.Buffer
	err := registrationTmpl.Execute(&buff, link)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}