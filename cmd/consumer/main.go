package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend-challenge-guinea/internal/contexts/users/application/projections"
	"backend-challenge-guinea/internal/contexts/users/domain"
	usersPersistence "backend-challenge-guinea/internal/contexts/users/infrastructure/persistence"
	"backend-challenge-guinea/internal/shared/infrastructure/bus"
	"backend-challenge-guinea/internal/shared/infrastructure/config"
	"backend-challenge-guinea/internal/shared/infrastructure/persistence"
	"backend-challenge-guinea/internal/shared/logger"
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Inicializar logger
	appLogger, err := logger.NewLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	appLogger.Info("starting consumer", map[string]interface{}{
		"env": cfg.Env,
	})

	// 3. Conectar a PostgreSQL
	db, err := persistence.NewPostgresConnection(cfg.Database)
	if err != nil {
		appLogger.Error("failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	appLogger.Info("connected to database", nil)

	// 4. Conectar a RabbitMQ
	eventBus, err := bus.NewRabbitMQBus(cfg.RabbitMQ.URL, cfg.RabbitMQ.Exchange, appLogger)
	if err != nil {
		appLogger.Error("failed to connect to rabbitmq", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}
	defer eventBus.Close()

	appLogger.Info("connected to rabbitmq", nil)

	// 5. Inicializar repositorios
	userReadModelRepo := usersPersistence.NewPostgresUserReadModel(db)

	// 6. Inicializar projector
	userProjector := projections.NewUserProjector(userReadModelRepo, appLogger)

	// 7. Suscribir el projector al evento UserCreated
	err = eventBus.Subscribe(domain.UserCreatedEventType, func(ctx context.Context, event interface{}) error {
		// Convertir el evento de map a UserCreatedEvent
		eventMap, ok := event.(map[string]interface{})
		if !ok {
			appLogger.Error("invalid event format", nil)
			return nil
		}

		// Serializar y deserializar para obtener el struct correcto
		eventBytes, err := json.Marshal(eventMap)
		if err != nil {
			appLogger.Error("failed to marshal event", map[string]interface{}{
				"error": err.Error(),
			})
			return err
		}

		var userCreatedEvent domain.UserCreatedEvent
		if err := json.Unmarshal(eventBytes, &userCreatedEvent); err != nil {
			appLogger.Error("failed to unmarshal to UserCreatedEvent", map[string]interface{}{
				"error": err.Error(),
			})
			return err
		}

		// Proyectar usando el evento tipado
		return userProjector.ProjectUserCreated(ctx, userCreatedEvent)
	})

	if err != nil {
		appLogger.Error("failed to subscribe to events", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// 8. Iniciar el consumo de mensajes
	ctx := context.Background()
	if err := eventBus.Start(ctx); err != nil {
		appLogger.Error("failed to start event bus", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("Failed to start consumer: %v", err)
	}

	appLogger.Info("consumer started, waiting for events...", nil)

	// 9. Esperar señal de terminación
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("shutting down consumer...", nil)

	// 10. Cerrar conexiones
	if err := eventBus.Close(); err != nil {
		appLogger.Error("error closing event bus", map[string]interface{}{
			"error": err.Error(),
		})
	}

	appLogger.Info("consumer stopped", nil)
}