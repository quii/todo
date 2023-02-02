package views

import "html/template"

func NewTemplates() (*template.Template, error) {
	return template.ParseFS(todoTemplates, "templates/*/*.gohtml", "templates/*.gohtml")
}
