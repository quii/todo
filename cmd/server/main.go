package main

import (
	"log"
	"net/http"

	"github.com/quii/todo/adapters/todohttp"
	"github.com/quii/todo/domain/todo"
)

const addr = ":8000"

func main() {
	list := todo.List{}
	list.Add("Bake a cake")
	list.Add("Feed the cat")
	list.Add("Take out the rubbish")

	handler, err := todohttp.NewTodoHandler(&list)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
