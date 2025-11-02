package projections

import (
	"context"
	"encoding/json"

	"backend-challenge-guinea/internal/contexts/users/domain"
	shared "backend-challenge-guinea/internal/shared/domain"
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


func (p *UserProjector) ProjectUserCreated(ctx context.Context, event shared.DomainEvent) error {

	var userCreated domain.UserCreatedEvent
	eventData, _ := json.Marshal(event)
	if err := json.Unmarshal(eventData, &userCreated); err != nil {
		p.log.Error("failed to unmarshal event", map[string]interface{}{
			"error":    err.Error(),
			"event_id": event.EventID(),
		})
		return err
	}

	userView := &domain.UserView{
		ID:          userCreated.UserID,
		Name:        userCreated.Name,
		Email:       userCreated.Email,
		DisplayName: userCreated.DisplayName,
		TenantID:    userCreated.TenantID(),
		CreatedAt:   userCreated.OccurredOn().Format("2006-01-02T15:04:05Z"),
	}

	if err := p.readModelRepo.Save(ctx, userView); err != nil {
		p.log.Error("failed to save user view", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userCreated.UserID,
		})
		return err
	}

	p.log.Info("user view projected", map[string]interface{}{
		"user_id":        userCreated.UserID,
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