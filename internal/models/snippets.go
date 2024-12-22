package models

import (
	"time"

	"github.com/jackc/pgx/v5"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgx.Conn
}

func Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func Latest() ([]Snippet, error) {
	return nil, nil
}
