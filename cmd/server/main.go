package main

import (
	"net/http"

	"github.com/quii/todo/adapters/todohttp"
	"github.com/quii/todo/domain/todo"
)

func main() {
	service := todo.Service{}
	service.Add("Bake a cake")
	service.Add("Feed the cat")
	service.Add("Take out the rubbish")

	handler := todohttp.NewTodoHandler(&service)

	http.ListenAndServe(":8000", handler)
}
