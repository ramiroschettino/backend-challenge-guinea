package domain

import "context"

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id, tenantID string) (*User, error)
	FindByEmail(ctx context.Context, email, tenantID string) (*User, error)
	ExistsByEmail(ctx context.Context, email, tenantID string) (bool, error)
}

type UserReadModel interface {
	FindByID(ctx context.Context, id, tenantID string) (*UserView, error)
	FindAll(ctx context.Context, tenantID string) ([]UserView, error)
}

type UserView struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name,omitempty"`
	TenantID    string  `json:"tenant_id"`
	CreatedAt   string  `json:"created_at"`
}