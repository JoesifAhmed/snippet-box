package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: No Matching Record Found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Ceated  time.Time
	Expires time.Time
}
