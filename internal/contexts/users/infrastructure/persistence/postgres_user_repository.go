package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"backend-challenge-guinea/internal/contexts/users/domain"
	vo "backend-challenge-guinea/internal/shared/domain/value_objects"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users_write (id, name, email, password_hash, display_name, tenant_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID(),
		user.Name(),
		user.Email().Value(),
		user.Password().Hash(),
		user.DisplayName(),
		user.TenantID(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id, tenantID string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, display_name, tenant_id, created_at, updated_at
		FROM users_write
		WHERE id = $1 AND tenant_id = $2
	`

	var (
		userID       string
		name         string
		email        string
		passwordHash string
		displayName  *string
		tenantId     string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&userID, &name, &email, &passwordHash, &displayName, &tenantId, &createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	emailVO, _ := vo.NewEmail(email)
	passwordVO := vo.FromHash(passwordHash)

	return domain.Reconstitute(
		userID,
		name,
		emailVO,
		passwordVO,
		tenantId,
		displayName,
		createdAt,
		updatedAt,
	), nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email, tenantID string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, display_name, tenant_id, created_at, updated_at
		FROM users_write
		WHERE email = $1 AND tenant_id = $2
	`

	var (
		userID       string
		name         string
		emailStr     string
		passwordHash string
		displayName  *string
		tenantId     string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.db.QueryRowContext(ctx, query, email, tenantID).Scan(
		&userID, &name, &emailStr, &passwordHash, &displayName, &tenantId, &createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	emailVO, _ := vo.NewEmail(emailStr)
	passwordVO := vo.FromHash(passwordHash)

	return domain.Reconstitute(
		userID,
		name,
		emailVO,
		passwordVO,
		tenantId,
		displayName,
		createdAt,
		updatedAt,
	), nil
}


func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email, tenantID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users_write WHERE email = $1 AND tenant_id = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email, tenantID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}