package main

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/*.html
var templateFS embed.FS

func render(w io.Writer, p page) error {
	tmpl, err := template.New("layout.html").Funcs(template.FuncMap{
		"tr": func(lang, en, zh string) string {
			if lang == "zh" {
				return zh
			}
			return en
		},
	}).Option("missingkey=error").ParseFS(templateFS, "templates/layout.html", "templates/"+p.ID+".html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, p)
}
