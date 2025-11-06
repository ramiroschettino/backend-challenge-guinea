package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	authCommands "backend-challenge-guinea/internal/contexts/auth/application/commands"
	authHttp "backend-challenge-guinea/internal/contexts/auth/infrastructure/http"
	"backend-challenge-guinea/internal/contexts/users/application/commands"
	"backend-challenge-guinea/internal/contexts/users/application/queries"
	usersHttp "backend-challenge-guinea/internal/contexts/users/infrastructure/http"
	usersPersistence "backend-challenge-guinea/internal/contexts/users/infrastructure/persistence"
	"backend-challenge-guinea/internal/shared/infrastructure/bus"
	"backend-challenge-guinea/internal/shared/infrastructure/config"
	sharedHttp "backend-challenge-guinea/internal/shared/infrastructure/http"
	"backend-challenge-guinea/internal/shared/infrastructure/middleware"
	"backend-challenge-guinea/internal/shared/infrastructure/persistence"
	"backend-challenge-guinea/internal/shared/logger"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger, err := logger.NewLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	appLogger.Info("starting api server", map[string]interface{}{
		"env":  cfg.Env,
		"port": cfg.Port,
	})

	db, err := persistence.NewPostgresConnection(cfg.Database)
	if err != nil {
		appLogger.Error("failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	appLogger.Info("connected to database", nil)

	if err := runMigrations(db, appLogger); err != nil {
		appLogger.Error("failed to run migrations", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Migrations failed: %v", err)
	}

	eventBus, err := bus.NewRabbitMQBus(cfg.RabbitMQ.URL, cfg.RabbitMQ.Exchange, appLogger)
	if err != nil {
		appLogger.Error("failed to connect to rabbitmq", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}
	defer eventBus.Close()

	appLogger.Info("connected to rabbitmq", nil)

	userRepository := usersPersistence.NewPostgresUserRepository(db)
	userReadModel := usersPersistence.NewPostgresUserReadModel(db)
	idempotencyRepo := usersPersistence.NewPostgresIdempotencyRepository(db)

	createUserHandler := commands.NewCreateUserCommandHandler(
		userRepository,
		eventBus,
		idempotencyRepo,
	)
	getUserHandler := queries.NewGetUserQueryHandler(userReadModel)

	authenticateHandler := authCommands.NewAuthenticateCommandHandler(userRepository)

	featureFlags := middleware.NewFeatureFlags()
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	userHandlers := usersHttp.NewUserHandlers(createUserHandler, getUserHandler, featureFlags)
	healthHandlers := sharedHttp.NewHealthHandlers(db)

	authHandlers := authHttp.NewAuthHandlers(authenticateHandler)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	healthHandlers.RegisterRoutes(router)
	userHandlers.RegisterRoutes(router, rateLimiter)

	authHandlers.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		appLogger.Info("api server listening", map[string]interface{}{
			"port": cfg.Port,
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("server failed", map[string]interface{}{
				"error": err.Error(),
			})
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("shutting down server...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
	}

	appLogger.Info("server stopped", nil)
}

func runMigrations(db *sql.DB, logger logger.Logger) error {
	logger.Info("running database migrations...", nil)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	// Ejecutar migraciones
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	logger.Info("migrations completed successfully", nil)
	return nil
}
