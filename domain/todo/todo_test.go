package todo_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/todo/domain/todo"
)

func TestService(t *testing.T) {
	t.Run("can add todo and toggle completion", func(t *testing.T) {
		service := todo.Service{}
		assert.Equal(t, 0, len(service.Todos()))

		id := uuid.New()
		service.Add(todo.Todo{
			ID: id,
		})

		service.Add(todo.Todo{ID: uuid.New()})
		service.Add(todo.Todo{ID: uuid.New()})
		service.Add(todo.Todo{ID: uuid.New()})

		todos := service.Todos()
		assert.False(t, todos[0].Complete)

		service.Toggle(id)

		todos = service.Todos()
		assert.True(t, todos[0].Complete)

		service.Toggle(id)

		todos = service.Todos()
		assert.False(t, todos[0].Complete)

	})
}
