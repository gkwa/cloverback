package cloverback

import (
	"bytes"
	"log/slog"
	"text/template"
)

func renderTmpl(pushes []Push, tmplStr string) bytes.Buffer {
	tmpl, err := template.New("pushTemplate").Parse(tmplStr)
	if err != nil {
		slog.Error("parsing template", "error", err)
		return bytes.Buffer{}
	}

	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, pushes)
	if err != nil {
		slog.Error("executing template", "error", err)
		return bytes.Buffer{}
	}

	return outputBuffer
}

func genOrgMode(pushes []Push, renderer func(pushes []Push, tmplStr string) bytes.Buffer) bytes.Buffer {
	tmplStr := `{{range .}}
*** {{.Title}}

{{.URL}}
{{end}}
`

	return renderer(pushes, tmplStr)
}
