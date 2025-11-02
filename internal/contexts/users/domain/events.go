package domain

import (
	shared "backend-challenge-guinea/internal/shared/domain"
)

const (
	UserCreatedEventType = "user.created"
)


type UserCreatedEvent struct {
	shared.BaseEvent        
	UserID      string      `json:"user_id"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	DisplayName *string     `json:"display_name,omitempty"`
}

func NewUserCreatedEvent(userID, name, email, tenantID, correlationID string, displayName *string) UserCreatedEvent {
	return UserCreatedEvent{
		BaseEvent:   shared.NewBaseEvent(UserCreatedEventType, userID, tenantID, correlationID),
		UserID:      userID,
		Name:        name,
		Email:       email,
		DisplayName: displayName,
	}
}