package todohttp_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/quii/todo/adapters/todohttp"
	"github.com/quii/todo/domain/todo"
)

import "github.com/go-rod/rod"

type TodoPage struct {
	Rod  *rod.Browser
	Page *rod.Page
	URL  string
}

func (t *TodoPage) Add(description string) {
	t.Page.MustElement(`[name="description"]`).MustInput(description).MustType(input.Enter)
	t.Page = t.Rod.MustPage(t.URL)
}

func (t *TodoPage) Edit(from, to string) {
	el := fmt.Sprintf(`[data-description="%s"]`, from)
	t.Page.MustElement(el + ` .edit`).MustClick()
	t.Page.MustElement(el + ` input[type="text"]`).MustInput(to)
	t.Page.MustElement(el + ` input[type="text"]`).MustType(input.Enter)
	t.Page = t.Rod.MustPage(t.URL)
}

func (t *TodoPage) Delete(description string) {
	el := fmt.Sprintf(`[data-description="%s"]`, description)
	t.Page.MustElement(el + ` .delete`).MustClick()
	t.Page = t.Rod.MustPage(t.URL)
}

func (t *TodoPage) Toggle(description string) {
	el := fmt.Sprintf(`[data-description="%s"]`, description)
	t.Page.MustElement(el + ` span`).MustClick()
	t.Page = t.Rod.MustPage(t.URL)
}

func TestNewTodoHandler(t *testing.T) {
	todoList := &todo.List{}
	handler, err := todohttp.NewTodoHandler(todoList)
	assert.NoError(t, err)

	server := httptest.NewServer(handler)
	defer server.Close()

	launcher := launcher.New().Headless(true).MustLaunch()
	rod := rod.New().Timeout(20 * time.Second).ControlURL(launcher).MustConnect()
	page := rod.MustPage(server.URL)

	todoListPage := &TodoPage{
		Rod:  rod,
		Page: page,
		URL:  server.URL,
	}

	t.Run("add some todos", func(t *testing.T) {
		todoListPage.Add("Eat cheese")
		todoListPage.Add("Drink port")

		assert.Equal(t, 2, len(todoList.Todos()))
		assert.Equal(t, "Eat cheese", todoList.Todos()[0].Description)
		assert.Equal(t, "Drink port", todoList.Todos()[1].Description)
	})

	t.Run("edit a todo", func(t *testing.T) {
		todoListPage.Edit("Eat cheese", "Eat cheese and crackers")
		assert.Equal(t, "Eat cheese and crackers", todoList.Todos()[0].Description)
	})

	t.Run("delete a todo", func(t *testing.T) {
		todoListPage.Add("Delete react")
		assert.Equal(t, 3, len(todoList.Todos()))
		assert.Equal(t, "Eat cheese and crackers", todoList.Todos()[0].Description)
		assert.Equal(t, "Drink port", todoList.Todos()[1].Description)
		assert.Equal(t, "Delete react", todoList.Todos()[2].Description)
		todoListPage.Delete("Delete react")
		assert.Equal(t, 2, len(todoList.Todos()))
		assert.Equal(t, "Eat cheese and crackers", todoList.Todos()[0].Description)
		assert.Equal(t, "Drink port", todoList.Todos()[1].Description)
	})

	t.Run("mark a todo as done", func(t *testing.T) {
		todoListPage.Add("Mark this as done")
		assert.Equal(t, 3, len(todoList.Todos()))
		assert.Equal(t, "Mark this as done", todoList.Todos()[2].Description)
		todoListPage.Toggle("Mark this as done")
		assert.True(t, todoList.Todos()[2].Complete)
		todoListPage.Toggle("Mark this as done")
		assert.False(t, todoList.Todos()[2].Complete)
	})

	t.Run("todo: attempts at testing drag and drog", func(t *testing.T) {
		t.Skip("pft")
		portBox, _ := page.MustElement(`[data-description="Drink port"]`).Shape()
		log.Println(portBox.OnePointInside())
		cheeseBox, _ := page.MustElement(`[data-description="Eat cheese"]`).Shape()
		log.Println(cheeseBox.OnePointInside())

		mouse := page.Mouse

		assert.NoError(t, mouse.MoveTo(*portBox.OnePointInside()))
		mouse.MustDown(proto.InputMouseButtonLeft)
		mouse.MoveLinear(*cheeseBox.OnePointInside(), 3)
		mouse.MustUp(proto.InputMouseButtonLeft)

		portBox, _ = page.MustElement(`[data-description="Drink port"]`).Shape()
		log.Println(portBox.OnePointInside())

		page = rod.MustPage(server.URL)
		assert.Equal(t, "Drink port", todoList.Todos()[0].Description)
		assert.Equal(t, "Eat cheese", todoList.Todos()[1].Description)
	})

}
