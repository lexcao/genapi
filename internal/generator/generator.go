package generator

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/lexcao/genapi/internal/model"
)

func GenerateFile(filename string, interfaces []model.Interface) {

}

func generate(tmpl TemplateName, data any) string {
	var buf bytes.Buffer
	template.Must(templates, templates.ExecuteTemplate(&buf, tmpl, data)) // TODO: error handling
	fomatted := must(format.Source(buf.Bytes()))
	return string(fomatted)
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
