package cloverback

import (
	"bytes"
	"log/slog"
	"text/template"
)

func genOrgMode(reply PushbulletHTTReply) bytes.Buffer {
	tmplStr := `{{range .Pushes}}
*** {{.Title}}
**** summary 

{{.URL}}
{{end}}`

	tmpl, err := template.New("pushTemplate").Parse(tmplStr)
	if err != nil {
		slog.Error("parsing template", "error", err.Error())
		return bytes.Buffer{}
	}

	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, reply)
	if err != nil {
		slog.Error("executing template", "error", err.Error())
		return bytes.Buffer{}
	}

	return outputBuffer
}
