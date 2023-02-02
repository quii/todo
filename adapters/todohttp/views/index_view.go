package views

import (
	"html/template"
	"net/http"

	"github.com/quii/todo/domain/todo"
)

type IndexView struct {
	templ *template.Template
}

func NewIndexView(templ *template.Template) *IndexView {
	return &IndexView{templ: templ}
}

func (t *IndexView) Index(w http.ResponseWriter, todos []todo.Todo) {
	var viewModel any = todos
	if err := t.templ.ExecuteTemplate(w, "index", viewModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
