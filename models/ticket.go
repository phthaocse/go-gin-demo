package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

const (
	TICKET_REPORTER_FKEY_CONSTRAINT = "ticket_assignee_fkey"
	FOREIGN_KEY_VIOLATION           = "foreign_key_violation"
)

type Ticket struct {
	Id        int        `json:"id"`
	Reporter  int        `json:"reporter"`
	Assignee  *int       `json:"assignee"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (t *Ticket) Create(db *sql.DB) error {
	row := db.QueryRow(`INSERT INTO ticket
						(reporter, assignee, title, content) 
						VALUES ($1, $2, $3, $4) RETURNING *`, t.Reporter, t.Assignee, t.Title, t.Content)
	err := row.Scan(&t.Id, &t.Reporter, &t.Assignee, &t.Title, &t.Content, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		err := err.(*pq.Error)
		if err.Code.Name() == FOREIGN_KEY_VIOLATION && err.Constraint == TICKET_REPORTER_FKEY_CONSTRAINT {
			return errors.New("The assignee isn't existed")
		}
		return err
	}
	return nil
}
