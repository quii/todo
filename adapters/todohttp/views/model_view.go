package views

import (
	"html/template"
	"net/http"
)

/*
FAO snarks

Yes, this is overly abstract, i'm just playing around, chill.
*/

type ModelView[T any] struct {
	templ     *template.Template
	modelName string
}

func NewModelView[T any](templ *template.Template, modelName string) *ModelView[T] {
	return &ModelView[T]{templ: templ, modelName: modelName}
}

func (t *ModelView[T]) List(w http.ResponseWriter, results []T) {
	t.renderOr500(w, t.modelName+"s", results)
}

func (t *ModelView[T]) View(w http.ResponseWriter, item T) {
	t.renderOr500(w, "view_"+t.modelName, item)
}

func (t *ModelView[T]) Edit(w http.ResponseWriter, item T) {
	t.renderOr500(w, "edit_"+t.modelName, item)
}

func (t *ModelView[T]) renderOr500(w http.ResponseWriter, templateName string, viewModel any) {
	if err := t.templ.ExecuteTemplate(w, templateName, viewModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
