package todohttp

import (
	"embed"
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

	service *todo.Service
	templ   *template.Template
}

func NewTodoHandler(service *todo.Service) *TodoHandler {
	router := mux.NewRouter()
	handler := &TodoHandler{
		Handler: router,
		service: service,
	}

	router.HandleFunc("/", handler.index).Methods(http.MethodGet)
	router.HandleFunc("/", handler.sortIndex).Methods(http.MethodPost)
	router.HandleFunc("/add", handler.add).Methods(http.MethodPost)
	router.HandleFunc("/search", handler.search).Methods(http.MethodPost)
	router.HandleFunc("/toggle/{ID}", handler.toggle).Methods(http.MethodPost)
	router.HandleFunc("/{ID}", handler.delete).Methods(http.MethodDelete)

	staticHandler, err := newStaticHandler()
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	templ, err := template.ParseFS(todoTemplates, "templates/*.gohtml")
	if err != nil {
		panic(err)
	}

	handler.templ = templ

	return handler
}

func (t *TodoHandler) index(w http.ResponseWriter, r *http.Request) {
	if err := t.templ.ExecuteTemplate(w, "index.gohtml", t.service.Todos()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TodoHandler) add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t.service.Add(r.FormValue("description"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (t *TodoHandler) toggle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.templ.ExecuteTemplate(w, "item.gohtml", t.service.Toggle(id))
}

func (t *TodoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["ID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.service.Delete(id)
}

func (t *TodoHandler) sortIndex(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t.service.ReOrder(r.Form["id"])
	if err := t.templ.ExecuteTemplate(w, "items.gohtml", t.service.Todos()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TodoHandler) search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if err := t.templ.ExecuteTemplate(w, "items.gohtml", t.service.Search(r.Form.Get("search"))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newStaticHandler() (http.Handler, error) {
	lol, err := fs.Sub(static, "static")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(lol)), nil
}
