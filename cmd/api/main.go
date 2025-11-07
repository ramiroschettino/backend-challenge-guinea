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

	// Importo los distintos contextos y módulos de la aplicación
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

	// Cargo la configuración del proyecto (variables de entorno, puertos, etc.)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Inicializo el logger con el nivel y formato definidos en la configuración
	appLogger, err := logger.NewLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Log inicial indicando que el servidor está arrancando
	appLogger.Info("starting api server", map[string]interface{}{
		"env":  cfg.Env,
		"port": cfg.Port,
	})

	// Establezco la conexión con la base de datos PostgreSQL
	db, err := persistence.NewPostgresConnection(cfg.Database)
	if err != nil {
		appLogger.Error("failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close() // Cierro la conexión cuando el programa termine

	appLogger.Info("connected to database", nil)

	// Ejecuto las migraciones de la base de datos
	if err := runMigrations(db, appLogger); err != nil {
		appLogger.Error("failed to run migrations", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Migrations failed: %v", err)
	}

	// Conecto con RabbitMQ para manejar eventos de dominio
	eventBus, err := bus.NewRabbitMQBus(cfg.RabbitMQ.URL, cfg.RabbitMQ.Exchange, appLogger)
	if err != nil {
		appLogger.Error("failed to connect to rabbitmq", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}
	defer eventBus.Close()

	appLogger.Info("connected to rabbitmq", nil)

	// Inicializo los repositorios del contexto de usuarios
	userRepository := usersPersistence.NewPostgresUserRepository(db)
	userReadModel := usersPersistence.NewPostgresUserReadModel(db)
	idempotencyRepo := usersPersistence.NewPostgresIdempotencyRepository(db)

	// Handlers de comandos y consultas del contexto de usuarios
	createUserHandler := commands.NewCreateUserCommandHandler(
		userRepository,
		eventBus,
		idempotencyRepo,
	)
	getUserHandler := queries.NewGetUserQueryHandler(userReadModel)

	// Handler de autenticación
	authenticateHandler := authCommands.NewAuthenticateCommandHandler(userRepository)

	// Middlewares de control de features y rate limiting
	featureFlags := middleware.NewFeatureFlags()
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Inicializo los controladores HTTP de cada módulo
	userHandlers := usersHttp.NewUserHandlers(createUserHandler, getUserHandler, featureFlags)
	healthHandlers := sharedHttp.NewHealthHandlers(db)
	authHandlers := authHttp.NewAuthHandlers(authenticateHandler)

	// Si estamos en producción, desactivo el modo debug de Gin
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Creo el router principal y registro las rutas de la API
	router := gin.Default()
	healthHandlers.RegisterRoutes(router)
	userHandlers.RegisterRoutes(router, rateLimiter)
	authHandlers.RegisterRoutes(router)

	// Configuro el servidor HTTP
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Inicio el servidor en una goroutine para no bloquear el hilo principal
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

	// Canal para escuchar señales del sistema (Ctrl+C o kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Espera hasta que llegue una señal

	appLogger.Info("shutting down server...", nil)

	// Contexto con timeout para permitir un apagado controlado (graceful shutdown)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
	}

	appLogger.Info("server stopped", nil)
}

// Función auxiliar que ejecuta las migraciones de base de datos
func runMigrations(db *sql.DB, logger logger.Logger) error {
	logger.Info("running database migrations...", nil)

	// Configuro el driver de migraciones para PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	// Creo la instancia de migraciones leyendo los archivos de la carpeta /migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	// Ejecuto las migraciones pendientes (si no hay cambios, ignoro el error)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	logger.Info("migrations completed successfully", nil)
	return nil
}
