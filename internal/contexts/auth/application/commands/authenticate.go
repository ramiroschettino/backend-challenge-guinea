package commands

import (
	"time"
	"context"
	"backend-challenge-guinea/internal/contexts/auth/domain"
	userDomain "backend-challenge-guinea/internal/contexts/users/domain"
)

type AuthenticateCommand struct {
	Email    string
	Password string
	TenantID string
}

type AuthenticateResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	ExpiresAt string `json:"expires_at"`
}

type AuthenticateCommandHandler struct {
	userRepository userDomain.UserRepository 
}

func NewAuthenticateCommandHandler(userRepo userDomain.UserRepository) *AuthenticateCommandHandler {
	return &AuthenticateCommandHandler{
		userRepository: userRepo,
	}
}

func (h *AuthenticateCommandHandler) Handle(ctx context.Context, cmd AuthenticateCommand) (*AuthenticateResponse, error) {

	user, err := h.userRepository.FindByEmail(ctx, cmd.Email, cmd.TenantID)
	if err != nil {
		return nil, domain.ErrInvalidCredentials 
	}

	if !user.Password().Compare(cmd.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	session := domain.NewSession(user.ID(), cmd.TenantID, 24*time.Hour)

	return &AuthenticateResponse{
		Token:     session.Token(),
		UserID:    session.UserID(),
		ExpiresAt: session.ExpiresAt().Format(time.RFC3339),
	}, nil
}