package views

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/quii/todo/domain/todo"
)

var (
	//go:embed "templates/*"
	todoTemplates embed.FS
)

type TodoView struct {
	templ *template.Template
}

func NewTodoView(templ *template.Template) *TodoView {
	return &TodoView{templ: templ}
}

func (t *TodoView) List(w http.ResponseWriter, results []todo.Todo) {
	t.renderOr500(w, "todos", results)
}

func (t *TodoView) View(w http.ResponseWriter, todo todo.Todo) {
	t.renderOr500(w, "view_todo", todo)
}

func (t *TodoView) Edit(w http.ResponseWriter, item todo.Todo) {
	t.renderOr500(w, "edit_todo", item)
}

func (t *TodoView) renderOr500(w http.ResponseWriter, templateName string, viewModel any) {
	if err := t.templ.ExecuteTemplate(w, templateName, viewModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
