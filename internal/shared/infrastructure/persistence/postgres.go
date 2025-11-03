package persistence

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" 
	"backend-challenge-guinea/internal/shared/infrastructure/config"
)

func NewPostgresConnection(cfg config.DatabaseConfig) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)        
	db.SetMaxIdleConns(10)       
	db.SetConnMaxLifetime(time.Hour) 

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}