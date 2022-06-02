package models

import (
	"database/sql"
	"errors"
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

	stmt := "INSERT INTO users (name, email, password_hash, activated) VALUES ($1, $2, $3, $4)"
	args := []interface{}{name, email, string(hashedPassword), false}

	_, err = m.DB.Exec(stmt, args...)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate query the database for a given email and check if the password
// matched with the hash password store in the db
func (m UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var passwordHash []byte

	stmt := "SELECT id, password_hash  FROM users WHERE email = $1 AND activated = $2"

	args := []interface{}{email, true}
	err := m.DB.QueryRow(stmt, args...).Scan(&id, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil

}
