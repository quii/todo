package todohttp_test

import (
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

func TestNewTodoHandler(t *testing.T) {
	todoList := &todo.List{}
	handler, err := todohttp.NewTodoHandler(todoList)
	assert.NoError(t, err)

	server := httptest.NewServer(handler)
	defer server.Close()

	launcher := launcher.New().Headless(true).MustLaunch()
	rod := rod.New().Timeout(20 * time.Second).SlowMotion(1 * time.Millisecond).ControlURL(launcher).MustConnect()
	page := rod.MustPage(server.URL)

	el := page.MustElement(`[name="description"]`)
	el.MustInput("Eat cheese")
	el.MustType(input.Enter)
	page = rod.MustPage(server.URL)

	el = page.MustElement(`[name="description"]`)
	el.MustInput("Drink port")
	el.MustType(input.Enter)

	assert.Equal(t, 2, len(todoList.Todos()))
	assert.Equal(t, "Eat cheese", todoList.Todos()[0].Description)
	assert.Equal(t, "Drink port", todoList.Todos()[1].Description)

	t.Run("todo: attempts at testing drag and drog", func(t *testing.T) {
		t.Skip()
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
		page.MustScreenshot("b.png")
	})

}
