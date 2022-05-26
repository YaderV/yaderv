package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX store the regular expression that we use to validate a email address
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator contains the form field error messages
type Validator struct {
	FieldErrors map[string]string
}

// Valid checks is a Validator (attached to a form) if valid
func (v Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError ties a message to a given key (input)
func (v Validator) AddFieldError(key, message string) {
	// We have to initialize the map if it does not exists
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField add an error message to the field if it is not valid
func (v Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.FieldErrors[key] = message
	}
}

// NotBlank checks if a field value is not blank
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxCharacters returns true is the valen length is less than max
func MaxCharacters(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// Min returns true if the value field length is grater than min
func Min(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min

}

// Matches return true if a given value match the rx regular expresion
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
