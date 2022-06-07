package models

import (
	"database/sql"

	"github.com/lib/pq"
)

// Article represent an article in the DB
type Article struct {
	ID         int
	Title      string
	Body       string
	Categories []string
}

// ArticleModel wraps the DB connection pool
type ArticleModel struct {
	DB *sql.DB
}

// Create a article
func (m ArticleModel) Create(title, body string, categories []string) error {
	stmt := "INSERT INTO articles (title, body, categories) VALUES($1, $2, $3)"
	args := []interface{}{title, body, pq.Array(categories)}

	_, err := m.DB.Exec(stmt, args...)

	if err != nil {
		return err
	}

	return nil
}
