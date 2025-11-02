package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verificamos que emails válidos se crean correctamente
func TestNewEmail_ValidEmails(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@example.com",
		"user+tag@example.co.uk",
		"TEST@EXAMPLE.COM", 
	}

	for _, emailStr := range validEmails {
		t.Run(emailStr, func(t *testing.T) {
			email, err := NewEmail(emailStr)
			
			assert.NoError(t, err)
			
			assert.NotEmpty(t, email.Value())
		})
	}
}

	// Verificamos que emails inválidos son rechazados
func TestNewEmail_InvalidEmails(t *testing.T) {
	invalidEmails := []string{
		"",                    
		"notanemail",          
		"@example.com",        
		"user@",               
		"user @example.com",   
		"user@example",        
	}

	for _, emailStr := range invalidEmails {
		t.Run(emailStr, func(t *testing.T) {
			_, err := NewEmail(emailStr)
			
			assert.Error(t, err)
		})
	}
}

// Verificamos que los emails se normalizan
func TestNewEmail_Normalization(t *testing.T) {
	// Email con espacios y mayúsculas
	email, err := NewEmail("  TEST@EXAMPLE.COM  ")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", email.Value())
}