package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	
	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	
	if !emailRegex.MatchString(email) {
		return Email{}, errors.New("invalid email format")
	}
	
	return Email{value: email}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}