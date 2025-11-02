package persistence

import (
	"context"
	"database/sql"
	"errors"

	"backend-challenge-guinea/internal/contexts/users/domain"
)


type PostgresUserReadModel struct {
	db *sql.DB
}

func NewPostgresUserReadModel(db *sql.DB) *PostgresUserReadModel {
	return &PostgresUserReadModel{db: db}
}

func (r *PostgresUserReadModel) Save(ctx context.Context, view *domain.UserView) error {
	query := `
		INSERT INTO users_read (id, name, email, display_name, tenant_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id, tenant_id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			display_name = EXCLUDED.display_name
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		view.ID,
		view.Name,
		view.Email,
		view.DisplayName,
		view.TenantID,
		view.CreatedAt,
	)

	return err
}

func (r *PostgresUserReadModel) FindByID(ctx context.Context, id, tenantID string) (*domain.UserView, error) {
	query := `
		SELECT id, name, email, display_name, tenant_id, created_at
		FROM users_read
		WHERE id = $1 AND tenant_id = $2
	`

	var view domain.UserView
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&view.ID,
		&view.Name,
		&view.Email,
		&view.DisplayName,
		&view.TenantID,
		&view.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &view, nil
}

func (r *PostgresUserReadModel) FindAll(ctx context.Context, tenantID string) ([]domain.UserView, error) {
	query := `
		SELECT id, name, email, display_name, tenant_id, created_at
		FROM users_read
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []domain.UserView
	for rows.Next() {
		var view domain.UserView
		if err := rows.Scan(
			&view.ID,
			&view.Name,
			&view.Email,
			&view.DisplayName,
			&view.TenantID,
			&view.CreatedAt,
		); err != nil {
			return nil, err
		}
		views = append(views, view)
	}

	return views, rows.Err()
}