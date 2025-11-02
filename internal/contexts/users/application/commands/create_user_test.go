package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
	"backend-challenge-guinea/internal/contexts/users/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id, tenantID string) (*domain.User, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email, tenantID string) (*domain.User, error) {
	args := m.Called(ctx, email, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email, tenantID string) (bool, error) {
	args := m.Called(ctx, email, tenantID)
	return args.Bool(0), args.Error(1)
}

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

type MockIdempotencyRepository struct {
	mock.Mock
}

func (m *MockIdempotencyRepository) IsProcessed(ctx context.Context, key, tenantID string) (bool, string, error) {
	args := m.Called(ctx, key, tenantID)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockIdempotencyRepository) Store(ctx context.Context, key, tenantID, result string) error {
	args := m.Called(ctx, key, tenantID, result)
	return args.Error(0)
}


func TestCreateUserCommandHandler_Success(t *testing.T) {

	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockEventBus := new(MockEventBus)
	mockIdempotency := new(MockIdempotencyRepository)

	handler := NewCreateUserCommandHandler(mockRepo, mockEventBus, mockIdempotency)

	cmd := CreateUserCommand{
		Name:           "John Doe",
		Email:          "john@example.com",
		Password:       "SecurePass123!",
		TenantID:       "tenant-1",
		CorrelationID:  "corr-123",
		IdempotencyKey: "idem-key-1",
	}

	mockIdempotency.On("IsProcessed", ctx, cmd.IdempotencyKey, cmd.TenantID).Return(false, "", nil)
	mockRepo.On("ExistsByEmail", ctx, "john@example.com", cmd.TenantID).Return(false, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	mockIdempotency.On("Store", ctx, cmd.IdempotencyKey, cmd.TenantID, mock.AnythingOfType("string")).Return(nil)
	mockEventBus.On("Publish", ctx, mock.AnythingOfType("domain.UserCreatedEvent")).Return(nil)

	userID, err := handler.Handle(ctx, cmd)

	assert.NoError(t, err)
	assert.NotEmpty(t, userID)
	
	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
	mockIdempotency.AssertExpectations(t)
}

func TestCreateUserCommandHandler_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockEventBus := new(MockEventBus)
	mockIdempotency := new(MockIdempotencyRepository)

	handler := NewCreateUserCommandHandler(mockRepo, mockEventBus, mockIdempotency)

	cmd := CreateUserCommand{
		Name:          "John Doe",
		Email:         "john@example.com",
		Password:      "SecurePass123!",
		TenantID:      "tenant-1",
		CorrelationID: "corr-123",
	}

	mockRepo.On("ExistsByEmail", ctx, "john@example.com", cmd.TenantID).Return(true, nil)

	userID, err := handler.Handle(ctx, cmd)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserAlreadyExists, err)
	assert.Empty(t, userID)
	mockRepo.AssertExpectations(t)
}

func TestCreateUserCommandHandler_Idempotency(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockEventBus := new(MockEventBus)
	mockIdempotency := new(MockIdempotencyRepository)

	handler := NewCreateUserCommandHandler(mockRepo, mockEventBus, mockIdempotency)

	cmd := CreateUserCommand{
		Name:           "John Doe",
		Email:          "john@example.com",
		Password:       "SecurePass123!",
		TenantID:       "tenant-1",
		CorrelationID:  "corr-123",
		IdempotencyKey: "idem-key-1",
	}

	existingUserID := "user-123"

	mockIdempotency.On("IsProcessed", ctx, cmd.IdempotencyKey, cmd.TenantID).Return(true, existingUserID, nil)

	userID, err := handler.Handle(ctx, cmd)

	assert.NoError(t, err)
	assert.Equal(t, existingUserID, userID)
	mockIdempotency.AssertExpectations(t)
	
	mockRepo.AssertNotCalled(t, "Save")
	mockEventBus.AssertNotCalled(t, "Publish")
}