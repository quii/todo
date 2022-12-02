package main

import (
	"net/http"

	"github.com/quii/todo/adapters/todohttp"
	"github.com/quii/todo/domain/todo"
)

func main() {
	service := todo.Service{}
	handler := todohttp.NewTodoHandler(&service)

	http.ListenAndServe(":8000", handler)
}
