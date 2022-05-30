package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
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

// Insert a user
func (m UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO USER (name, email, password_hash) VALUES (?, ?, ?)"
	args := []interface{}{name, email, string(hashedPassword)}

	_, err = m.DB.Exec(stmt, args...)

	if err != nil {
		return err
	}

	return nil
}
