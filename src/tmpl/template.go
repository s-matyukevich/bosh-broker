package tmpl

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

func NewTemplate(text string) (*Template, error) {
	t, err := template.New("bosh-template").Parse(text)
	if err != nil {
		return nil, err
	}
	return &Template{t}, nil
}

type Template struct {
	t *template.Template
}

func (t *Template) ExecuteAndSave(data map[string]interface{}, path string) error {
	b := &bytes.Buffer{}
	err := t.t.Execute(b, data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b.Bytes(), 600)
}

func (t *Template) Execute(data map[string]interface{}) (string, error) {
	b := &bytes.Buffer{}
	err := t.t.Execute(b, data)
	return b.String(), err
}
