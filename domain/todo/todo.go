package todo

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type Todo struct {
	ID          uuid.UUID
	Description string
	CreatedAt   time.Time
	Complete    bool
}

type Service struct {
	todos []Todo
}

func (s *Service) Add(item Todo) {
	s.todos = append(s.todos, item)
}

func (s *Service) Toggle(id uuid.UUID) Todo {
	i := slices.IndexFunc(s.todos, func(todo Todo) bool {
		return todo.ID == id
	})
	s.todos[i].Complete = !s.todos[i].Complete
	return s.todos[i]
}

func (s *Service) Todos() []Todo {
	return s.todos
}

func (s *Service) Delete(id uuid.UUID) {
	i := slices.IndexFunc(s.todos, func(todo Todo) bool {
		return todo.ID == id
	})
	s.todos = append(s.todos[:i], s.todos[i+1:]...)
}
