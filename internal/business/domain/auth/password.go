package auth

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Password represents a password that can be used to authenticate a user.
// A password must be between 8 and 64 characters long, and contain at least one letter and one number.
// Only "@$!%*#?&_-^" special characters are allowed.
type Password struct {
	value string
}

func (p Password) String() string {
	return p.value
}

func (p Password) Equal(p2 Password) bool {
	return p.value == p2.value
}

func (p Password) MarshalText() ([]byte, error) {
	// redact the password from logs
	return []byte("REDACTED"), nil
}

// Hash hashes the password using bcrypt.
func (p Password) Hash() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p.value), 12)
}

func ParsePassword(password string) (Password, error) {
	// Since golang uses RE2 regex engine, we can't use lookaheads.
	// So we'll drop regex and use a simple loop to validate.

	if len(password) < 8 || len(password) > 64 {
		return Password{}, fmt.Errorf("password must be 8-64 characters")
	}

	hasNumber := false
	hasLetter := false
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsLetter(c):
			hasLetter = true
		case strings.Contains("@$!%*#?&_-^", string(c)):
			// valid special character, no-op
		default:
			return Password{}, fmt.Errorf("invalid character %q in password", c)
		}
	}

	if !hasNumber {
		return Password{}, fmt.Errorf("password must contain at least one number")
	}
	if !hasLetter {
		return Password{}, fmt.Errorf("password must contain at least one letter")
	}

	return Password{value: password}, nil
}

func MustParsePassword(password string) Password {
	p, err := ParsePassword(password)
	if err != nil {
		panic(err)
	}
	return p
}
