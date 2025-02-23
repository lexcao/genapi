package generator

import (
	"bytes"
	"go/format"

	"github.com/lexcao/genapi/internal/binder"
	"github.com/lexcao/genapi/internal/model"
)

func GenerateFile(filename string, interfaces []model.Interface) {

}

func generateInterface(tmpl TemplateName, data model.Interface) (string, error) {
	if err := binder.Bind(&data); err != nil {
		return "", err
	}

	return generate(tmpl, data)
}

func generateMethod(tmpl TemplateName, data model.Method) (string, error) {
	if err := binder.BindMethod(&data); err != nil {
		return "", err
	}

	return generate(tmpl, data)
}

func generate(tmpl TemplateName, data any) (string, error) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, tmpl, data); err != nil {
		return "", err
	}

	fomatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(fomatted), nil
}
