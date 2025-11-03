package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authDomain "backend-challenge-guinea/internal/contexts/auth/domain"
	userDomain "backend-challenge-guinea/internal/contexts/users/domain"
	vo "backend-challenge-guinea/internal/shared/domain/value_objects"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *userDomain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id, tenantID string) (*userDomain.User, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email, tenantID string) (*userDomain.User, error) {
	args := m.Called(ctx, email, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email, tenantID string) (bool, error) {
	args := m.Called(ctx, email, tenantID)
	return args.Bool(0), args.Error(1)
}

func TestAuthenticateCommandHandler_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	handler := NewAuthenticateCommandHandler(mockRepo)

	email, _ := vo.NewEmail("test@example.com")
	password, _ := vo.NewPassword("SecurePass123!")
	user, _ := userDomain.NewUser("Test User", email, password, "tenant-1", nil)

	mockRepo.On("FindByEmail", ctx, "test@example.com", "tenant-1").Return(user, nil)

	cmd := AuthenticateCommand{
		Email:    "test@example.com",
		Password: "SecurePass123!",
		TenantID: "tenant-1",
	}

	response, err := handler.Handle(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, user.ID(), response.UserID)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateCommandHandler_InvalidCredentials(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	handler := NewAuthenticateCommandHandler(mockRepo)

	email, _ := vo.NewEmail("test@example.com")
	password, _ := vo.NewPassword("SecurePass123!")
	user, _ := userDomain.NewUser("Test User", email, password, "tenant-1", nil)

	mockRepo.On("FindByEmail", ctx, "test@example.com", "tenant-1").Return(user, nil)

	cmd := AuthenticateCommand{
		Email:    "test@example.com",
		Password: "WrongPassword123!",
		TenantID: "tenant-1",
	}

	response, err := handler.Handle(ctx, cmd)

	assert.Error(t, err)
	assert.Equal(t, authDomain.ErrInvalidCredentials, err)
	assert.Nil(t, response)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateCommandHandler_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	handler := NewAuthenticateCommandHandler(mockRepo)

	mockRepo.On("FindByEmail", ctx, "nonexistent@example.com", "tenant-1").Return(nil, userDomain.ErrUserNotFound)

	cmd := AuthenticateCommand{
		Email:    "nonexistent@example.com",
		Password: "SecurePass123!",
		TenantID: "tenant-1",
	}

	response, err := handler.Handle(ctx, cmd)

	assert.Error(t, err)
	assert.Equal(t, authDomain.ErrInvalidCredentials, err)
	assert.Nil(t, response)
	mockRepo.AssertExpectations(t)
}