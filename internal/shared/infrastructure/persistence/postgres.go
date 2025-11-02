package persistence

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Driver de PostgreSQL
	"backend-challenge-guinea/internal/shared/infrastructure/config"
)

// NewPostgresConnection crea y configura la conexión a PostgreSQL
func NewPostgresConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	// Construir connection string
	// Formato: host=localhost port=5432 user=user password=pass dbname=db sslmode=disable
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	// Abrir conexión
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)        // Máximo 25 conexiones simultáneas
	db.SetMaxIdleConns(10)        // Mantener 10 conexiones idle
	db.SetConnMaxLifetime(time.Hour) // Cerrar conexiones después de 1 hora

	// Verificar que la conexión funciona
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}