package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/todo/adapters/todohttp"
	"github.com/quii/todo/domain/todo"
)

func main() {
	service := todo.Service{}
	service.Add(todo.Todo{ID: uuid.New(), Description: "Bake a cake"})
	service.Add(todo.Todo{ID: uuid.New(), Description: "Feed the cat"})
	service.Add(todo.Todo{ID: uuid.New(), Description: "Take out the rubbish"})

	handler := todohttp.NewTodoHandler(&service)

	http.ListenAndServe(":8000", handler)
}
