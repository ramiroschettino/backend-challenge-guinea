package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verificamos que las contraseñas validas se hashean bien
func TestNewPassword_ValidPassword(t *testing.T) {
	plainPassword := "SecurePass123!"
	
	password, err := NewPassword(plainPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, password.Hash())
	
	// La contraseña hasheada no tiene que ser igual a la contraseña original
	assert.NotEqual(t, plainPassword, password.Hash())
}

// Verificamos que la comparación funciona
func TestPassword_Compare(t *testing.T) {
	plainPassword := "SecurePass123!"
	
	password, _ := NewPassword(plainPassword)

	assert.True(t, password.Compare(plainPassword))
	assert.False(t, password.Compare("WrongPassword!"))
}

//Acá abajo las que deberian ser rechazadas por diferentes motivos.

// Contraseña demasiado corta
func TestNewPassword_TooShort(t *testing.T) {
	_, err := NewPassword("Short1!")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Al menos 8 caracteres")
}

// Contraseñas sin mayúsculas
func TestNewPassword_NoUppercase(t *testing.T) {
	_, err := NewPassword("lowercase123!")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uppercase")
}

// Contraseñas sin minusculas
func TestNewPassword_NoLowercase(t *testing.T) {
	_, err := NewPassword("UPPERCASE123!")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lowercase")
}

// Sin números
func TestNewPassword_NoNumber(t *testing.T) {
	_, err := NewPassword("NoNumbers!")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "number")
}

// Sin caracteres especiales
func TestNewPassword_NoSpecial(t *testing.T) {
	_, err := NewPassword("NoSpecial123")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "special")
}

// Verificamos que podemos reconstruir desde un hash
func TestFromHash(t *testing.T) {
	original, _ := NewPassword("SecurePass123!")
	hash := original.Hash()
	
	reconstructed := FromHash(hash)
	assert.True(t, reconstructed.Compare("SecurePass123!"))
}