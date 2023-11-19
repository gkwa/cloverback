package cloverback

import (
	"bytes"
	"log/slog"
	"text/template"
)

func genOrgMode(pushes []Push) bytes.Buffer {
	tmplStr := `{{range .}}
*** {{.Title}}

{{.URL}}
{{end}}
`

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
