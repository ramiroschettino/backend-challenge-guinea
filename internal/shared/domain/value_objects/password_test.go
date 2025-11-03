package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassword_ValidPassword(t *testing.T) {
	plainPassword := "SecurePass123!"
	
	password, err := NewPassword(plainPassword)
	
	assert.NoError(t, err)
	
	assert.NotEmpty(t, password.Hash())
	
	assert.NotEqual(t, plainPassword, password.Hash())
}

func TestPassword_Compare(t *testing.T) {
	plainPassword := "SecurePass123!"
	
	password, _ := NewPassword(plainPassword)
	
	assert.True(t, password.Compare(plainPassword))
	
	assert.False(t, password.Compare("WrongPassword!"))
}

func TestNewPassword_TooShort(t *testing.T) {
	_, err := NewPassword("Short1!")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "8 characters") 
}

func TestNewPassword_NoUppercase(t *testing.T) {
	_, err := NewPassword("lowercase123!")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uppercase") 
}

func TestNewPassword_NoLowercase(t *testing.T) {
	_, err := NewPassword("UPPERCASE123!")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lowercase") 
}

func TestNewPassword_NoNumber(t *testing.T) {
	_, err := NewPassword("NoNumbers!")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "number") 
}

func TestNewPassword_NoSpecial(t *testing.T) {
	_, err := NewPassword("NoSpecial123")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "special") 
}

func TestFromHash(t *testing.T) {
	original, _ := NewPassword("SecurePass123!")
	hash := original.Hash()
	
	reconstructed := FromHash(hash)
	
	assert.True(t, reconstructed.Compare("SecurePass123!"))
}