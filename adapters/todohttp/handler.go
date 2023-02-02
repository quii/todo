package todohttp

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quii/todo/domain/todo"
)

var (
	//go:embed "templates/*"
	todoTemplates embed.FS

	//go:embed static
	static embed.FS
)

type TodoHandler struct {
	http.Handler

	list  *todo.List
	templ *template.Template
}

func NewTodoHandler(service *todo.List) (*TodoHandler, error) {
	router := mux.NewRouter()
	handler := &TodoHandler{
		Handler: router,
		list:    service,
	}

	router.HandleFunc("/", handler.index).Methods(http.MethodGet)

	router.HandleFunc("/todos", handler.add).Methods(http.MethodPost)
	router.HandleFunc("/todos", handler.search).Methods(http.MethodGet)
	router.HandleFunc("/todos/sort", handler.reOrder).Methods(http.MethodPost)
	router.HandleFunc("/todos/{ID}/edit", handler.edit).Methods(http.MethodGet)
	router.HandleFunc("/todos/{ID}/toggle", handler.toggle).Methods(http.MethodPost)
	router.HandleFunc("/todos/{ID}", handler.delete).Methods(http.MethodDelete)
	router.HandleFunc("/todos/{ID}", handler.rename).Methods(http.MethodPatch)

	staticHandler, err := newStaticHandler()
	if err != nil {
		return nil, fmt.Errorf("problem making static resources handler: %w", err)
	}
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	templ, err := template.ParseFS(todoTemplates, "templates/*/*.gohtml", "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	handler.templ = templ

	return handler, nil
}

func (t *TodoHandler) index(w http.ResponseWriter, r *http.Request) {
	t.renderOr500(w, "index", t.list.Todos())
}

func (t *TodoHandler) add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t.list.Add(r.FormValue("description"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (t *TodoHandler) toggle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.renderTodo(w, t.list.ToggleDone(id))
}

func (t *TodoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.list.Delete(id)
}

func (t *TodoHandler) reOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t.list.ReOrder(r.Form["id"])
	t.renderTodos(w, t.list.Todos())
}

func (t *TodoHandler) search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")
	results := t.list.Search(searchTerm)
	t.renderTodos(w, results)
}

func (t *TodoHandler) rename(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newName := r.Form["name"][0]
	todo := t.list.Rename(id, newName)

	t.renderTodo(w, todo)
}

func (t *TodoHandler) edit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item := t.list.Get(id)
	t.renderOr500(w, "edit_todo", item)
}

func newStaticHandler() (http.Handler, error) {
	lol, err := fs.Sub(static, "static")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(lol)), nil
}

func (t *TodoHandler) renderTodos(w http.ResponseWriter, results []todo.Todo) {
	t.renderOr500(w, "todos", results)
}

func (t *TodoHandler) renderTodo(w http.ResponseWriter, todo todo.Todo) {
	t.renderOr500(w, "view_todo", todo)
}

func (t *TodoHandler) renderOr500(w http.ResponseWriter, templateName string, viewModel any) {
	if err := t.templ.ExecuteTemplate(w, templateName, viewModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
