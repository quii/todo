package todo

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	ID          uuid.UUID
	Description string
	CreatedAt   time.Time
	Complete    bool
}

func (t *Todo) ToggleDone() {
	t.Complete = !t.Complete
}
