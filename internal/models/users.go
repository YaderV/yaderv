package models

import (
	"database/sql"
	"time"
)

// User represents a user in the database
type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash []byte
	Created      time.Time
}

// UserModel wraps the DB connection pool
type UserModel struct {
	DB *sql.DB
}
