# ADR 001: Por qué usamos Arquitectura Hexagonal

**Fecha**: 2024-10-31

## El Problema

Necesitábamos una forma de organizar el código para que:
- El dominio no dependa de Gin, PostgreSQL o RabbitMQ
- Podamos hacer tests sin levantar bases de datos
- Si mañana cambiamos de PostgreSQL a MySQL, no tengamos que reescribir todo

## Lo que decidimos

Usar Arquitectura Hexagonal (también llamada Ports & Adapters).

Básicamente significa:
- **Dominio**: La lógica de negocio pura (User, Email, Password)
- **Puertos**: Interfaces que dice qué necesita el dominio (ej: `UserRepository`)
- **Adaptadores**: Las implementaciones reales (ej: `PostgresUserRepository`)

### Ejemplo
````go
// Puerto (en dominio)
type UserRepository interface {
    Save(user *User) error
}

// Adaptador (en infraestructura)
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(user *User) error {
    // código de PostgreSQL aquí
}
````

## Por qué está bien

- Puedo testear el dominio con mocks
- Si cambio de base de datos, solo cambio el adaptador
- El código está más ordenado

## Por qué es un poco molesto

- Escribo más código (interfaces + implementaciones)
- Al principio cuesta entenderlo
- Para features simples se siente como overkill

## Otras cosas que pensé

- **MVC normal**: Más simple pero todo queda muy acoplado
- **Clean Architecture**: Parecido pero con más capas (para este proyecto es mucho)

---

**Conclusión**: Vale la pena el esfuerzo extra porque el código queda más testeable y mantenible.