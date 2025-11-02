.PHONY: help up down restart build test test-coverage migrate migrate-down logs clean

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

up: ## Levanta todos los servicios
	docker-compose up -d
	@echo "Esperando a que los servicios estén listos..."
	@timeout /t 8 /nobreak > nul
	@echo "Ejecutando migraciones..."
	@make migrate

down: ## Frena todos los servicios
	docker-compose down

restart: ## Reinicia todos los servicios
	docker-compose restart

build: ## Construye las imágenes
	docker-compose build

logs: ## Muestra logs de todos los servicios
	docker-compose logs -f

logs-api: ## Muestra logs solo de la API
	docker-compose logs -f api

logs-consumer: ## Muestra logs solo del consumer
	docker-compose logs -f consumer

test: ## Ejecuta los tests
	go test -v -race ./...

test-coverage: ## Ejecuta tests con reporte de cobertura
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Reporte de cobertura generado: coverage.html"

migrate: ## Ejecuta las migraciones
	migrate -path ./migrations -database "postgresql://backend_user:backend_pass@localhost:5432/backend_db?sslmode=disable" up

migrate-down: ## Revierte la última migración
	migrate -path ./migrations -database "postgresql://backend_user:backend_pass@localhost:5432/backend_db?sslmode=disable" down 1

migrate-create: ## Crea una nueva migración (uso: make migrate-create NAME=nombre)
	migrate create -ext sql -dir migrations -seq $(NAME)

clean: ## Limpia todo (containers, volumes, cache)
	docker-compose down -v
	del /Q coverage.out coverage.html 2>nul
	go clean -testcache

run-api: ## Ejecuta la API localmente (sin Docker)
	go run cmd/api/main.go

run-consumer: ## Ejecuta el consumer localmente (sin Docker)
	go run cmd/consumer/main.go

deps: ## Descarga dependencias
	go mod download
	go mod tidy

health: ## Verifica el health de los servicios
	@curl -s http://localhost:8080/health
	@echo.
	@curl -s http://localhost:8080/ready