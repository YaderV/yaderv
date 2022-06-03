package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX store the regular expression that we use to validate a email address
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator handles field validators and erros messages
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid checks is a Validator (attached to a form) if valid
func (v Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddNonFieldError add an general error message that's no tied to a specific field
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// AddFieldError adds an error message to a given input field
func (v *Validator) AddFieldError(key, message string) {
	// We have to initialize the map if it does not exists
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField add an error message to the field if it is not valid
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		// We have to use AddFieldError in the case the map is nil (non-initialize)
		v.AddFieldError(key, message)
	}
}

// NotEmpty validate that a given string array is not empty
func NotEmpty(values []string) bool {
	return len(values) > 0
}

// NotBlank checks if a field value is not blank
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxCharacters returns true is the valen length is less than max
func MaxCharacters(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// MinChars returns true if the value field length is grater than min
func MinChars(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min

}

// Matches return true if a given value match the rx regular expresion
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// ConfirmPassword compares two password to check if they are the same
func ConfirmPassword(password, repeatedPassword string) bool {
	return password == repeatedPassword
}
