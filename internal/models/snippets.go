package models

import (
	"context"
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

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `insert into snippets(title, content, expired_at, created_at)
	values(
		$1,
		$2, 
		current_timestamp at time zone 'utc' + $3 * interval '1 day',
		current_timestamp at time zone 'utc' 
	)
	returning id
	`

	var id int
	err := sm.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func Latest() ([]Snippet, error) {
	return nil, nil
}
