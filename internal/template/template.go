package template

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
)

var (
	registrationTmpl *template.Template
)

func PrepareTemplates() error {
	path, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	if path == "/" {
		path = ""
	}
	t, err := template.ParseFiles(fmt.Sprintf("%s/html/registration.html", path))
	if err != nil {
		return err
	}
	registrationTmpl = t
	return nil
}

func Register(link string) ([]byte, error) {
	var buff bytes.Buffer
	err := registrationTmpl.Execute(&buff, template.URL(link))
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}