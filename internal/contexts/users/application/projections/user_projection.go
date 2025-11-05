package projections

import (
	"context"

	"backend-challenge-guinea/internal/contexts/users/domain"
)

type UserProjector struct {
	readModelRepo UserReadModelRepository
	log           Logger
}

func NewUserProjector(repo UserReadModelRepository, log Logger) *UserProjector {
	return &UserProjector{
		readModelRepo: repo,
		log:           log,
	}
}


func (p *UserProjector) ProjectUserCreated(ctx context.Context, event domain.UserCreatedEvent) error {

	userView := &domain.UserView{
		ID:          event.UserID,
		Name:        event.Name,
		Email:       event.Email,
		DisplayName: event.DisplayName,
		TenantID:    event.TenantID(),
		CreatedAt:   event.OccurredOn().Format("2006-01-02T15:04:05Z"),
	}

	if err := p.readModelRepo.Save(ctx, userView); err != nil {
		p.log.Error("failed to save user view", map[string]interface{}{
			"error":   err.Error(),
			"user_id": event.UserID,
		})
		return err
	}

	p.log.Info("user view projected", map[string]interface{}{
		"user_id":        event.UserID,
		"correlation_id": event.CorrelationID(),
	})

	return nil
}

type UserReadModelRepository interface {
	Save(ctx context.Context, view *domain.UserView) error
	FindByID(ctx context.Context, id, tenantID string) (*domain.UserView, error)
}

type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}