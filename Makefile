.PHONY: help up down restart build test test-coverage migrate migrate-down logs clean

help: ## Muestra esta ayuda
	@echo Comandos disponibles:
	@echo   make up          - Levanta todos los servicios
	@echo   make down        - Frena todos los servicios
	@echo   make migrate     - Ejecuta migraciones
	@echo   make test        - Ejecuta tests
	@echo   make logs        - Muestra logs

up: ## Levanta todos los servicios
	docker-compose up -d
	@echo Esperando a que los servicios esten listos...
	@timeout /t 8 /nobreak > nul 2>&1
	@echo Ejecutando migraciones...
	@$(MAKE) migrate

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
	@echo Reporte de cobertura generado: coverage.html

migrate: ## Ejecuta las migraciones
	migrate -path ./migrations -database "postgresql://backend_user:backend_pass@localhost:5432/backend_db?sslmode=disable" up

migrate-down: ## Revierte la última migración
	migrate -path ./migrations -database "postgresql://backend_user:backend_pass@localhost:5432/backend_db?sslmode=disable" down 1

migrate-create: ## Crea una nueva migración (uso: make migrate-create NAME=nombre)
	migrate create -ext sql -dir migrations -seq $(NAME)

clean: ## Limpia todo (containers, volumes, cache)
	docker-compose down -v
	@if exist coverage.out del /Q coverage.out 2>nul
	@if exist coverage.html del /Q coverage.html 2>nul
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