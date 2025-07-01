package auth

import (
	"fmt"
	"regexp"
)

// Username is a unique identifier for a user.
// A username can only contain letters and numbers, and must be between 4 and 30 characters long.
type Username struct {
	value string
}

func (u Username) String() string {
	return u.value
}

func (u Username) Equal(u2 Username) bool {
	return u.value == u2.value
}

func (u Username) MarshalText() ([]byte, error) {
	return []byte(u.value), nil
}

var usernameRegex = regexp.MustCompile(`^[A-Za-z\d]{4,30}$`)

func ParseUsername(username string) (Username, error) {
	if !usernameRegex.MatchString(username) {
		return Username{}, fmt.Errorf("invalid username %q", username)
	}
	return Username{value: username}, nil
}

func MustParseUsername(username string) Username {
	u, err := ParseUsername(username)
	if err != nil {
		panic(err)
	}
	return u
}
