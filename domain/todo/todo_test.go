package todo_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/todo/domain/todo"
)

func TestService(t *testing.T) {
	t.Run("can add todo and toggle completion", func(t *testing.T) {
		service := todo.Service{}
		assert.Equal(t, 0, len(service.Todos()))

		someTask := "kill react"
		service.Add(someTask)
		service.Add("blah blah 1")
		service.Add("blah blah 2")
		service.Add("blah blah 3")

		todos := service.Todos()
		assert.False(t, todos[0].Complete)
		id := todos[0].ID

		service.Toggle(id)

		todos = service.Todos()
		assert.True(t, todos[0].Complete)

		service.Toggle(id)

		todos = service.Todos()
		assert.False(t, todos[0].Complete)
	})
}
