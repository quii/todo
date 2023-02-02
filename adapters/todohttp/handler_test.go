package todohttp_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/quii/todo/adapters/todohttp"
	"github.com/quii/todo/adapters/todohttp/views"
	"github.com/quii/todo/domain/todo"
)

import "github.com/go-rod/rod"

func TestNewTodoHandler(t *testing.T) {
	todoList := &todo.List{}
	templates, err := views.NewTemplates()
	assert.NoError(t, err)
	handler, err := todohttp.NewTodoHandler(todoList, views.NewTodoView(templates), views.NewIndexView(templates))
	assert.NoError(t, err)

	server := httptest.NewServer(handler)
	defer server.Close()

	launcher := launcher.New().Headless(true).MustLaunch()
	rod := rod.New().Timeout(20 * time.Second).ControlURL(launcher).MustConnect()

	todoListPage := &TodoPage{
		Rod:  rod,
		Page: rod.MustPage(server.URL),
		URL:  server.URL,
	}

	it := func(description string, f func(*testing.T)) {
		t.Helper()
		todoList.Empty()
		todoListPage.Home()
		t.Run(description, f)
	}

	it("can add some todos", func(t *testing.T) {
		todoListPage.Add("Eat cheese")
		todoListPage.Add("Drink port")

		assert.Equal(t, 2, len(todoList.Todos()))
		assert.Equal(t, "Eat cheese", todoList.Todos()[0].Description)
		assert.Equal(t, "Drink port", todoList.Todos()[1].Description)
	})

	it("can edit a todoview", func(t *testing.T) {
		todoListPage.Add("Eat cheese")
		todoListPage.Edit("Eat cheese", "Eat cheese and crackers")
		assert.Equal(t, "Eat cheese and crackers", todoList.Todos()[0].Description)
	})

	it("can delete a todoview", func(t *testing.T) {
		todoListPage.Add("Eat cheese")
		todoListPage.Add("Drink port")
		assert.Equal(t, 2, len(todoList.Todos()))

		todoListPage.Delete("Drink port")
		assert.Equal(t, 1, len(todoList.Todos()))
		assert.Equal(t, "Eat cheese", todoList.Todos()[0].Description)
	})

	it("can mark a todoview as done", func(t *testing.T) {
		todoListPage.Add("Mark this as done")
		assert.False(t, todoList.Todos()[0].Complete)

		todoListPage.Toggle("Mark this as done")
		assert.True(t, todoList.Todos()[0].Complete)
		todoListPage.Toggle("Mark this as done")
		assert.False(t, todoList.Todos()[0].Complete)
	})

	//t.Run("todoview: attempts at testing drag and drog", func(t *testing.T) {
	//	t.Skip("pft")
	//	portBox, _ := page.MustElement(`[data-description="Drink port"]`).Shape()
	//	log.Println(portBox.OnePointInside())
	//	cheeseBox, _ := page.MustElement(`[data-description="Eat cheese"]`).Shape()
	//	log.Println(cheeseBox.OnePointInside())
	//
	//	mouse := page.Mouse
	//
	//	assert.NoError(t, mouse.MoveTo(*portBox.OnePointInside()))
	//	mouse.MustDown(proto.InputMouseButtonLeft)
	//	mouse.MoveLinear(*cheeseBox.OnePointInside(), 3)
	//	mouse.MustUp(proto.InputMouseButtonLeft)
	//
	//	portBox, _ = page.MustElement(`[data-description="Drink port"]`).Shape()
	//	log.Println(portBox.OnePointInside())
	//
	//	page = rod.MustPage(server.URL)
	//	assert.Equal(t, "Drink port", todoList.Todos()[0].Description)
	//	assert.Equal(t, "Eat cheese", todoList.Todos()[1].Description)
	//})

}
