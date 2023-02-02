package views

import (
	"html/template"

	"github.com/quii/todo/domain/todo"
)

func NewTodoView(templ *template.Template) *ModelView[todo.Todo] {
	return NewModelView[todo.Todo](templ, "todo")
}
