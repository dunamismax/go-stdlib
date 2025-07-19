package components

import (
	"html/template"
	"io"
)

// BaseTemplate defines the base HTML template
const BaseTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{.Title}}</title>
    <style>{{.CSS}}</style>
    <script src="/static/htmx.min.js"></script>
</head>
<body>
    {{.Content}}
</body>
</html>`

// LayoutData holds data for the base template
type LayoutData struct {
	Title   string
	CSS     string
	Content template.HTML
}

// RenderLayout renders the base layout with the given data
func RenderLayout(w io.Writer, data LayoutData) error {
	tmpl, err := template.New("layout").Parse(BaseTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// Common HTML snippets
func Container(content string) string {
	return `<div class="container">` + content + `</div>`
}

func Card(content string) string {
	return `<div class="card">` + content + `</div>`
}

func Button(text, attrs string) string {
	return `<button ` + attrs + `>` + text + `</button>`
}

func Input(inputType, name, placeholder, attrs string) string {
	return `<input type="` + inputType + `" name="` + name + `" placeholder="` + placeholder + `" ` + attrs + `>`
}

func TextArea(name, placeholder, rows, attrs string) string {
	return `<textarea name="` + name + `" placeholder="` + placeholder + `" rows="` + rows + `" ` + attrs + `></textarea>`
}

func Form(method, action, content string) string {
	return `<form method="` + method + `" action="` + action + `">` + content + `</form>`
}

func Grid(content string) string {
	return `<div class="grid">` + content + `</div>`
}

func Result(id, content string) string {
	return `<div id="` + id + `" class="result">` + content + `</div>`
}
