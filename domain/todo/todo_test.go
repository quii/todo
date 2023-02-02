package todo_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/todo/domain/todo"
)

func TestService(t *testing.T) {
	t.Run("can add todoview and toggle completion", func(t *testing.T) {
		service := todo.List{}
		assert.Equal(t, 0, len(service.Todos()))

		someTask := "kill react"
		service.Add(someTask)
		service.Add("blah blah 1")
		service.Add("blah blah 2")
		service.Add("blah blah 3")

		todos := service.Todos()
		assert.False(t, todos[0].Complete)
		id := todos[0].ID

		service.ToggleDone(id)

		todos = service.Todos()
		assert.True(t, todos[0].Complete)

		service.ToggleDone(id)

		todos = service.Todos()
		assert.False(t, todos[0].Complete)
	})

	t.Run("rename", func(t *testing.T) {
		service := todo.List{}
		service.Add("kill react")
		service.Add("blah blah 1")

		todos := service.Todos()
		assert.Equal(t, "kill react", todos[0].Description)
		id := todos[0].ID

		service.Rename(id, "kill react and redux")

		todos = service.Todos()
		assert.Equal(t, "kill react and redux", todos[0].Description)
		assert.Equal(t, "blah blah 1", todos[1].Description)
	})

	t.Run("delete", func(t *testing.T) {
		service := todo.List{}
		service.Add("kill react")
		service.Add("blah blah 1")
		service.Add("blah blah 2")
		service.Add("blah blah 3")

		todos := service.Todos()
		assert.Equal(t, 4, len(todos))
		id := todos[0].ID

		service.Delete(id)

		todos = service.Todos()
		assert.Equal(t, 3, len(todos))
		assert.Equal(t, "blah blah 1", todos[0].Description)
		assert.Equal(t, "blah blah 2", todos[1].Description)
		assert.Equal(t, "blah blah 3", todos[2].Description)
	})

	t.Run("reorder", func(t *testing.T) {
		service := todo.List{}
		service.Add("kill react")
		service.Add("blah blah 1")
		service.Add("blah blah 2")
		service.Add("blah blah 3")

		todos := service.Todos()
		assert.Equal(t, 4, len(todos))
		assert.Equal(t, "kill react", todos[0].Description)

		// reorder
		service.ReOrder([]string{
			todos[3].ID.String(),
			todos[2].ID.String(),
			todos[1].ID.String(),
			todos[0].ID.String(),
		})

		todos = service.Todos()
		assert.Equal(t, 4, len(todos))
		assert.Equal(t, "blah blah 3", todos[0].Description)
		assert.Equal(t, "blah blah 2", todos[1].Description)
		assert.Equal(t, "blah blah 1", todos[2].Description)
		assert.Equal(t, "kill react", todos[3].Description)
	})

	t.Run("search", func(t *testing.T) {
		service := todo.List{}
		service.Add("kill react")
		service.Add("blah blah 1")
		service.Add("blah blah 2")
		service.Add("blah blah 3")

		todos := service.Search("blah")
		assert.Equal(t, 3, len(todos))
		assert.Equal(t, "blah blah 1", todos[0].Description)
		assert.Equal(t, "blah blah 2", todos[1].Description)
		assert.Equal(t, "blah blah 3", todos[2].Description)
	})

	t.Run("empty the list", func(t *testing.T) {
		service := todo.List{}
		service.Add("kill react")
		service.Add("blah blah 1")
		service.Add("blah blah 2")
		service.Add("blah blah 3")

		todos := service.Todos()
		assert.Equal(t, 4, len(todos))

		service.Empty()

		todos = service.Todos()
		assert.Equal(t, 0, len(todos))
	})

}
