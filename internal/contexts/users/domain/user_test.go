package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vo "backend-challenge-guinea/internal/shared/domain/value_objects"
)

// Crear usuario valido
func TestNewUser_Success(t *testing.T) {

	email, _ := vo.NewEmail("test@example.com")
	password, _ := vo.NewPassword("SecurePass123!")
	displayName := "Test User"
	
	user, err := NewUser("John Doe", email, password, "tenant-1", &displayName)
	

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID())
	assert.Equal(t, "John Doe", user.Name())
	assert.Equal(t, "test@example.com", user.Email().Value())
	assert.Equal(t, "tenant-1", user.TenantID())
	assert.Equal(t, &displayName, user.DisplayName())
	assert.NotZero(t, user.CreatedAt())
	assert.NotZero(t, user.UpdatedAt())
}

// Sin nombre (tiene que fallar)
func TestNewUser_EmptyName(t *testing.T) {
	email, _ := vo.NewEmail("test@example.com")
	password, _ := vo.NewPassword("SecurePass123!")
	
	user, err := NewUser("", email, password, "tenant-1", nil)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrInvalidUserName, err)
}

// Recontruir desde database
func TestReconstitute(t *testing.T) {
	email, _ := vo.NewEmail("test@example.com")
	password := vo.FromHash("$2a$10$...")
	
	user := Reconstitute(
		"user-123",
		"John Doe",
		email,
		password,
		"tenant-1",
		nil,
		time.Now(),
		time.Now(),
	)
	
	assert.NotNil(t, user)
	assert.Equal(t, "user-123", user.ID())
	assert.Equal(t, "John Doe", user.Name())
}