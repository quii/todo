package views

import (
	"embed"
	"html/template"
)

var (
	//go:embed "templates/*"
	todoTemplates embed.FS
)

func NewTemplates() (*template.Template, error) {
	return template.ParseFS(todoTemplates, "templates/*/*.gohtml", "templates/*.gohtml")
}
