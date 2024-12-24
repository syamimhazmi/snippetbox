package models

import (
	"context"
	"errors"
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
	stmt := `insert into snippets(title, content, expired_at, created_at, updated_at)
	values(
		$1,
		$2, 
		current_timestamp at time zone 'utc' + $3 * interval '1 day',
		current_timestamp at time zone 'utc',
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

func (sm *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `select id, title, content, expired_at, created_at
	from snippets
	where expired_at > (current_timestamp at time zone 'utc')
	and id = $1
	`

	// the reason why we need to initialise this into a variable
	// Go doesn't let us to directly accessing the attribute without
	// declare it into variable.
	var s Snippet
	err := sm.DB.QueryRow(context.Background(), stmt, id).Scan(
		&s.ID,
		&s.Title,
		&s.Content,
		&s.Expires,
		&s.Created,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}

		return Snippet{}, err
	}

	return s, nil
}

func Latest() ([]Snippet, error) {
	return nil, nil
}
