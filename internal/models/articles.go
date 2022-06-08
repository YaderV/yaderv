package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// Article represent an article in the DB
type Article struct {
	ID         int
	Title      string
	Body       string
	Categories []string
	CreatedAt  time.Time
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

// List return all the articles
func (m ArticleModel) List() ([]Article, error) {
	stmt := "SELECT id, title, body, categories, created_at FROM articles ORDER BY id DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	articles := []Article{}

	for rows.Next() {
		a := Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Body, pq.Array(&a.Categories), &a.CreatedAt)

		if err != nil {
			return nil, err
		}

		articles = append(articles, a)
	}

	return articles, nil
}
