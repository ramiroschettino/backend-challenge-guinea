# Backend Challenge - Arquitectura Hexagonal + CQRS

Sistema backend con Arquitectura Hexagonal, CQRS, multi-tenant y event sourcing construido con Go.

## 🏗️ Arquitectura

- **Hexagonal (Ports & Adapters)**: Dominio independiente de frameworks
- **CQRS**: Separación de escritura (commands) y lectura (queries)
- **Event-Driven**: RabbitMQ para eventos de dominio
- **Multi-tenant**: Aislamiento por tenant a nivel de aplicación
- **TDD**: Desarrollo guiado por tests con cobertura >80%

## 📁 Estructura del Proyecto
```
backend-challenge-guinea/
├── cmd/
│   ├── api/          # Servidor HTTP (REST API)
│   └── consumer/     # Consumidor de eventos (proyecciones)
├── internal/
│   ├── contexts/     # Bounded contexts (users, auth)
│   │   ├── users/
│   │   │   ├── domain/           # Entidades, value objects, eventos
│   │   │   ├── application/      # Casos de uso (commands, queries)
│   │   │   └── infrastructure/   # Adaptadores (PostgreSQL, HTTP)
│   │   └── auth/
│   │       ├── domain/           # Sesiones, autenticación
│   │       ├── application/      # Login
│   │       └── infrastructure/   # HTTP handlers
│   └── shared/       # Código compartido entre contexts
│       ├── domain/              # Events base, value objects
│       ├── infrastructure/      # Config, logger, bus, middleware
│       └── logger/
├── migrations/       # Migraciones SQL
├── docs/            # Documentación y ADRs
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── .env.sample
```

## 🚀 Inicio Rápido

### Prerrequisitos

- Go 1.21+
- Docker y Docker Compose
- Make (opcional pero recomendado)

### Instalación

1. **Clonar el repositorio**
```bash
git clone https://github.com/ramiroschettino/backend-challenge-guinea.git
cd backend-challenge-guinea
```

2. **Configurar variables de entorno**
```bash
copy .env.example .env
# Editar .env si es necesario (los valores por defecto funcionan)
```

3. **Levantar todos los servicios**
```bash
make up
```

Esto va a:
- Levantar PostgreSQL
- Levantar RabbitMQ
- Ejecutar migraciones automáticamente
- Iniciar la API en http://localhost:8080
- Iniciar el consumer de eventos

### Comandos Disponibles
```bash
make help              # Ver todos los comandos disponibles
make up                # Levantar servicios
make down              # Detener servicios
make logs              # Ver logs de todos los servicios
make logs-api          # Ver logs solo de la API
make logs-consumer     # Ver logs del consumer
make test              # Ejecutar tests
make test-coverage     # Ejecutar tests con reporte de cobertura
make migrate           # Ejecutar migraciones manualmente
make migrate-down      # Revertir última migración
make clean             # Limpiar todo (containers, volumes, cache)
```

# 📡 Endpoints

## Health Checks

### Verificar que la API está viva
```
GET http://localhost:8080/health
```

### Verificar que las dependencias están listas
```
GET http://localhost:8080/ready
```

---

## Users API

### Crear Usuario

```
POST http://localhost:8080/api/v1/users

Headers:
Content-Type: application/json
X-Tenant-Id: tenant-1
X-Idempotency-Key: unique-key-123

Body:
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "display_name": "Johnny"
}
```

### Obtener Usuario

```
GET http://localhost:8080/api/v1/users/{user_id}

Headers:
X-Tenant-Id: tenant-1
```

---

## Auth API

### Login

**Paso 1: Primero crea un usuario**

```
POST http://localhost:8080/api/v1/users

Headers:
Content-Type: application/json
X-Tenant-Id: tenant-1

Body:
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Paso 2: Luego haz login**

```
POST http://localhost:8080/api/v1/auth/login

Headers:
Content-Type: application/json
X-Tenant-Id: tenant-1

Body:
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Respuesta:**
```json
{
  "token": "uuid-token-here",
  "user_id": "user-id-here",
  "expires_at": "2025-11-02T12:00:00Z"
}
```

### Headers Requeridos

- `X-Tenant-Id`: Identificador del tenant (requerido)
- `X-Correlation-Id`: ID de correlación para trazabilidad (opcional, se genera automáticamente)
- `X-Idempotency-Key`: Clave de idempotencia para evitar duplicados (opcional)

## 🧪 Testing

### Ejecutar todos los tests
```bash
make test
```

### Ejecutar tests con cobertura
```bash
make test-coverage
```

Esto genera un reporte HTML en `coverage.html`

### Ejecutar tests de un paquete específico
```bash
go test -v ./internal/contexts/users/domain/...
go test -v ./internal/contexts/users/application/commands/...
go test -v ./internal/contexts/auth/application/commands/...
```

## 🗄️ Base de Datos

### Migraciones

Las migraciones se ejecutan automáticamente al hacer `make up`.

**Crear nueva migración:**
```bash
make migrate-create NAME=add_user_status
```

**Ejecutar migraciones manualmente:**
```bash
make migrate
```

**Revertir última migración:**
```bash
make migrate-down
```

### Esquema

- **users_write**: Tabla de escritura (write model)
- **users_read**: Tabla de lectura optimizada (read model / proyección)
- **idempotency_keys**: Gestión de idempotencia

## 🐰 RabbitMQ

**Management UI**: http://localhost:15672
- Usuario: `guest`
- Contraseña: `guest`

### Eventos

- `user.created`: Se publica cuando se crea un usuario
  - El consumer escucha este evento y actualiza el read model

## 🔧 Configuración

Todas las configuraciones se gestionan mediante variables de entorno (archivo `.env`).

Ver `.env.sample` para la lista completa de variables disponibles.

## 🎯 Features

### Implementados

- ✅ Arquitectura Hexagonal (Ports & Adapters)
- ✅ CQRS (Command Query Responsibility Segregation)
- ✅ Event Sourcing con RabbitMQ
- ✅ Multi-tenant (X-Tenant-Id header)
- ✅ Idempotencia en comandos
- ✅ Rate limiting por tenant
- ✅ Feature flags por tenant
- ✅ Correlation ID para trazabilidad
- ✅ Logs estructurados (JSON)
- ✅ Health checks
- ✅ Migraciones automáticas
- ✅ Autenticación básica (login)
- ✅ Tests unitarios (>80% cobertura en dominio)
- ✅ Docker & Docker Compose
- ✅ Graceful shutdown

### Multi-Tenant

Cada request debe incluir el header `X-Tenant-Id`. Los datos están aislados por tenant:
- Usuarios de `tenant-1` no pueden ver usuarios de `tenant-2`
- Rate limiting por tenant
- Feature flags por tenant

### Idempotencia

Los comandos soportan idempotencia mediante el header `X-Idempotency-Key`:
- Si envías el mismo comando con la misma clave 2 veces, solo se ejecuta una vez
- Útil para reintentos en caso de timeouts o errores de red

### Feature Flags

Ejemplo: `display_name` solo está habilitado para algunos tenants.

Configuración actual:
- `tenant-1`: ✅ display_name habilitado
- `tenant-2`: ❌ display_name deshabilitado

## 📊 Monitoreo

### Logs

Los logs están estructurados en formato JSON con campos:
- `level`: debug, info, warn, error
- `msg`: mensaje del log
- `correlation_id`: para seguir el flujo de una request
- `tenant_id`: identificador del tenant
- Campos custom según el contexto

**Ver logs:**
```bash
make logs              # Todos los servicios
make logs-api          # Solo API
make logs-consumer     # Solo consumer
```

## 🔒 Seguridad

- Contraseñas hasheadas con bcrypt (costo 10)
- Validación de email
- Política de contraseñas:
  - Mínimo 8 caracteres
  - Al menos 1 mayúscula
  - Al menos 1 minúscula
  - Al menos 1 número
  - Al menos 1 carácter especial

## 🏗️ Decisiones Arquitectónicas

Ver los ADRs (Architecture Decision Records) en `/docs/adr/`:
- [ADR-001: Arquitectura Hexagonal](docs/adr/001-hexagonal-architecture.md)
- [ADR-002: CQRS Pattern](docs/adr/002-cqrs-pattern.md)
- [ADR-003: Estrategia Multi-Tenant](docs/adr/003-multi-tenant-strategy.md)

## 📈 Diagramas

Ver [flujo de comandos y eventos](docs/diagrams/command-flow.md) para entender cómo fluye una request por el sistema.

## 🛠️ Tecnologías

- **Go 1.21+**: Lenguaje principal
- **Gin**: Framework HTTP
- **PostgreSQL 15**: Base de datos
- **RabbitMQ 3**: Message broker
- **Docker**: Containerización
- **golang-migrate**: Migraciones de DB
- **Zap**: Logger estructurado
- **Testify**: Framework de testing

## 📝 Notas

- El proyecto sigue principios SOLID y Clean Architecture
- El dominio es completamente independiente de frameworks
- La separación CQRS permite escalar lectura y escritura independientemente
- Los eventos permiten eventual consistency entre write y read models
- El sistema está preparado para escalar horizontalmente

### Agregar nuevo contexto

1. Crear estructura en `internal/contexts/{nuevo_contexto}/`
2. Definir dominio (entidades, eventos, puertos)
3. Implementar casos de uso en application
4. Crear adaptadores en infrastructure
5. Registrar handlers y rutas

### Agregar nuevo caso de uso

1. Crear command/query en application
2. Implementar handler con tests
3. Crear endpoint HTTP en infrastructure
4. Documentar en README
