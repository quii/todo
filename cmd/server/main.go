package main

import (
	"log"
	"net/http"

	"github.com/quii/todo/adapters/todohttp"
	"github.com/quii/todo/adapters/todohttp/views"
	"github.com/quii/todo/domain/todo"
)

const addr = ":8000"

func main() {
	list := todo.List{}
	list.Add("Bake a cake")
	list.Add("Feed the cat")
	list.Add("Take out the rubbish")

	templates, err := views.NewTemplates()

	if err != nil {
		log.Fatal(err)
	}

	handler, err := todohttp.NewTodoHandler(&list, views.NewTodoView(templates), views.NewIndexView(templates))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
