package commands

import (
	"context"

	"backend-challenge-guinea/internal/contexts/users/domain"
	vo "backend-challenge-guinea/internal/shared/domain/value_objects"
)


type CreateUserCommand struct {
	Name           string
	Email          string
	Password       string
	DisplayName    *string
	TenantID       string
	CorrelationID  string
	IdempotencyKey string 
}


type CreateUserCommandHandler struct {
	repository      domain.UserRepository    
	eventBus        EventBus                 
	idempotencyRepo IdempotencyRepository    
}

func NewCreateUserCommandHandler(
	repo domain.UserRepository,
	eventBus EventBus,
	idempotencyRepo IdempotencyRepository,
) *CreateUserCommandHandler {
	return &CreateUserCommandHandler{
		repository:      repo,
		eventBus:        eventBus,
		idempotencyRepo: idempotencyRepo,
	}
}

func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) (string, error) {

	if cmd.IdempotencyKey != "" {
		processed, result, err := h.idempotencyRepo.IsProcessed(ctx, cmd.IdempotencyKey, cmd.TenantID)
		if err != nil {
			return "", err
		}
		if processed {
	
			return result, nil
		}
	}

	email, err := vo.NewEmail(cmd.Email)
	if err != nil {
		return "", err 
	}

	exists, err := h.repository.ExistsByEmail(ctx, email.Value(), cmd.TenantID)
	if err != nil {
		return "", err
	}
	if exists {
		return "", domain.ErrUserAlreadyExists
	}

	password, err := vo.NewPassword(cmd.Password)
	if err != nil {
		return "", err 
	}

	user, err := domain.NewUser(cmd.Name, email, password, cmd.TenantID, cmd.DisplayName)
	if err != nil {
		return "", err
	}

	if err := h.repository.Save(ctx, user); err != nil {
		return "", err
	}


	if cmd.IdempotencyKey != "" {
		if err := h.idempotencyRepo.Store(ctx, cmd.IdempotencyKey, cmd.TenantID, user.ID()); err != nil {
		}
	}

	event := domain.NewUserCreatedEvent(
		user.ID(),
		user.Name(),
		user.Email().Value(),
		cmd.TenantID,
		cmd.CorrelationID,
		user.DisplayName(),
	)

	if err := h.eventBus.Publish(ctx, event); err != nil {

		return user.ID(), err
	}

	return user.ID(), nil
}

type EventBus interface {
	Publish(ctx context.Context, event interface{}) error
}

type IdempotencyRepository interface {
	IsProcessed(ctx context.Context, key, tenantID string) (bool, string, error)
	Store(ctx context.Context, key, tenantID, result string) error
}