package generator

import (
	"text/template"
)

type TemplateName = string

const tmplInterface TemplateName = "tmplInterface"
const tmplMethod TemplateName = "tmplMethod"
const tmplMethodBody TemplateName = "tmplMethodBody"

var templates = template.Must(
	template.New("Templates").Parse(`
// CODE GENERATED BY genapi. DO NOT EDIT.
package {{.Package}}

{{range .Interfaces}}
{{template "tmplInterface" .}}
{{end}}
`))

func init() {
	template.Must(templates.New(tmplInterface).Parse(`
{{- $impl := printf "impl%s" .Name}}
type {{ $impl }} struct {
	client genapi.HttpClient
}

// setHttpClient implments genapi.Interface
func (i *{{ $impl }}) setHttpClient(client genapi.HttpClient) {
	i.client = client
}

{{range .Methods}}
{{template "tmplMethod" .}}
{{end}}
`))

	template.Must(templates.New(tmplMethod).Parse(`
func (i *impl{{ .Interface }}) {{ .Name }}(
{{- range $i, $p := .Params }}
	{{- if $i }}, {{ end }}{{ $p.Name }} {{ $p.Type }}
{{- end -}}
)
{{- if gt (len .Results) 1 }}({{- end -}}
	{{- range $i, $r := .Results }}
	{{- if $i }}, {{ end }}{{ $r.Name }} {{ $r.Type }}
	{{- end -}}
{{- if gt (len .Results) 1 }}){{- end -}}
{{- " " -}}
{
{{- template "tmplMethodBody" . -}}
}
`))

	template.Must(templates.New(tmplMethodBody).Parse(`
	{{- with .Bindings.Results -}}
	{{- .Assignment -}} :=
	{{- end -}}i.client.Do(&genapi.Request{
		{{- with .Bindings.Method }}
		Method: "{{.}}",
		{{- end }}
		{{- with .Bindings.Path }}
		Path: "{{.}}",
		{{- end }}
		{{- with .Bindings.PathParams }}
		PathParams: {{.}},
		{{- end }}
		{{- with .Bindings.Queries }}
		Queries: {{.}},
		{{- end }}
		{{- with .Bindings.Headers }}
		Headers: {{.}},
		{{- end }}
		{{- with .Bindings.Context }}
		Context: {{.}},
		{{- end }}
		{{- with .Bindings.Body }}
		Body: {{.}},
		{{- end }}
	})
	{{ with .Bindings.Results -}}
	return {{ .Statement -}}
	{{- end -}}
	`))
}

// TODO: 5. Handle Result & Error
