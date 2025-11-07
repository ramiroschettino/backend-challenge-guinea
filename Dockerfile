# Etapa 1: Build (compila los binarios de Go)
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Compila los binarios para la API y el consumer
RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer cmd/consumer/main.go

# Etapa 2: Development (entorno para desarrollo local)
FROM golang:1.24-alpine AS development
WORKDIR /app
RUN apk add --no-cache git make bash curl
# Instala la herramienta de migraciones
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# Copia el código y los binarios compilados desde la etapa anterior
COPY --from=builder /app .

# Etapa 3: Production (imagen final optimizada) esta no la vamos a usar

# FROM alpine:latest AS production
# RUN apk --no-cache add ca-certificates bash curl
# WORKDIR /root/
# # Copia los binarios y migraciones necesarios
# COPY --from=builder /app/api .
# COPY --from=builder /app/consumer .
# COPY --from=builder /app/migrations ./migrations
# COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
# # Expone el puerto 8080 y define el punto de entrada
# EXPOSE 8080
# CMD ["./api"]
