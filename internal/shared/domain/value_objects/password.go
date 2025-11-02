package valueobjects

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 8
	BcryptCost       = 10 
)

type Password struct {
	hashedValue string
}

func NewPassword(plainPassword string) (Password, error) {
	if err := validatePassword(plainPassword); err != nil {
		return Password{}, err
	}
	
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), BcryptCost)
	if err != nil {
		return Password{}, err
	}
	
	return Password{hashedValue: string(hashedBytes)}, nil
}

func FromHash(hash string) Password {
	return Password{hashedValue: hash}
}

func (p Password) Hash() string {
	return p.hashedValue
}

func (p Password) Compare(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainPassword))
	return err == nil
}

func validatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return errors.New("password must be at least 8 characters long")
	}
	
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return errors.New("la contraseña debe al menos una mayuscula")
	}
	if !hasLower {
		return errors.New("la contraseña debe al menos una minuscula")
	}
	if !hasNumber {
		return errors.New("la contraseña debe tener al menos un número")
	}
	if !hasSpecial {
		return errors.New("La contraseña debe tener al menos un caracter especial")
	}
	
	return nil
}