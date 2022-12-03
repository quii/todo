package todo

import (
	"strings"
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

func (s *Service) Add(description string) {
	s.todos = append(s.todos, Todo{
		ID:          uuid.New(),
		Description: description,
		CreatedAt:   time.Now(),
		Complete:    false,
	})
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
	i := s.indexOf(id)
	s.todos = append(s.todos[:i], s.todos[i+1:]...)
}

func (s *Service) ReOrder(ids []string) {

	var uuids []uuid.UUID
	for _, id := range ids {
		uuids = append(uuids, uuid.MustParse(id))
	}

	var newList []Todo
	for _, id := range uuids {
		newList = append(newList, s.todos[s.indexOf(id)])
	}

	s.todos = newList
}

func (s *Service) indexOf(id uuid.UUID) int {
	return slices.IndexFunc(s.todos, func(todo Todo) bool {
		return todo.ID == id
	})
}

func (s *Service) Search(search string) []Todo {
	search = strings.ToLower(search)
	var results []Todo
	for _, todo := range s.todos {
		if strings.Contains(strings.ToLower(todo.Description), search) {
			results = append(results, todo)
		}
	}
	return results
}
