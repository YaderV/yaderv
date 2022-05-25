package models

import (
	"errors"
)

var (
	// ErrNoRecord is raised when a query does not find a record
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials is raised when the user tries to log in with wrong credentials
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail is raised when an user tries to sign up with a duplicated email
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
