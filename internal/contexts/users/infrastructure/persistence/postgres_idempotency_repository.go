package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostgresIdempotencyRepository struct {
	db *sql.DB
}

func NewPostgresIdempotencyRepository(db *sql.DB) *PostgresIdempotencyRepository {
	return &PostgresIdempotencyRepository{db: db}
}

func (r *PostgresIdempotencyRepository) IsProcessed(ctx context.Context, key, tenantID string) (bool, string, error) {
	query := `SELECT result FROM idempotency_keys WHERE key = $1 AND tenant_id = $2`

	var result string
	err := r.db.QueryRowContext(ctx, query, key, tenantID).Scan(&result)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", nil 
		}
		return false, "", err
	}

	return true, result, nil 
}

func (r *PostgresIdempotencyRepository) Store(ctx context.Context, key, tenantID, result string) error {
	query := `
		INSERT INTO idempotency_keys (key, tenant_id, result, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key, tenant_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, key, tenantID, result, time.Now())
	return err
}