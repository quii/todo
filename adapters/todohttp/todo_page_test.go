package todohttp_test

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type TodoPage struct {
	Rod  *rod.Browser
	Page *rod.Page
	URL  string
}

func (t *TodoPage) Home() {
	t.Page = t.Rod.MustPage(t.URL)
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
