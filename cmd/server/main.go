package main

import (
	"net/http"

	"github.com/quii/todo/adapters/todohttp"
	"github.com/quii/todo/domain/todo"
)

func main() {
	list := todo.List{}
	list.Add("Bake a cake")
	list.Add("Feed the cat")
	list.Add("Take out the rubbish")

	handler := todohttp.NewTodoHandler(&list)

	http.ListenAndServe(":8000", handler)
}
