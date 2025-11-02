package queries

import (
	"context"
	"errors"

	"backend-challenge-guinea/internal/contexts/users/domain"
)

type GetUserQuery struct {
	UserID   string
	TenantID string
}

type GetUserQueryHandler struct {
	readModel domain.UserReadModel 
}

func NewGetUserQueryHandler(readModel domain.UserReadModel) *GetUserQueryHandler {
	return &GetUserQueryHandler{
		readModel: readModel,
	}
}

func (h *GetUserQueryHandler) Handle(ctx context.Context, query GetUserQuery) (*domain.UserView, error) {

	if query.UserID == "" {
		return nil, errors.New("user ID is required")
	}
	
	user, err := h.readModel.FindByID(ctx, query.UserID, query.TenantID)
	if err != nil {
		return nil, err
	}

	return user, nil
}