# ADR 002: Por qué separamos lectura y escritura (CQRS)

**Fecha**: 2025-11-02

## El Problema

En un CRUD normal:
- La misma tabla para escribir y leer
- Si hago muchos INSERTs, las lecturas se ponen lentas
- Las consultas complejas necesitan muchos JOINs

Y el challenge pedía específicamente CQRS.

## Lo que decidí
Separar en dos modelos:

**Para escribir:**
- Tabla `users_write` normalizada
- Guardo el usuario
- Publico un evento

**Para leer:**
- Tabla `users_read` optimizada para consultas
- Se actualiza cuando llega el evento
- Las consultas son súper rápidas (sin JOINs)

### El flujo
````
1. POST /users
   → Guardo en users_write
   → Publico evento UserCreated

2. Consumer escucha evento
   → Actualiza users_read

3. GET /users/:id
   → Leo desde users_read
````

## Por qué está bien

- Las lecturas no bloquean las escrituras
- Puedo optimizar cada tabla para su propósito
- Si tengo millones de lecturas, puedo escalar esa parte sola

## Por qué es complicado

- Tengo que mantener 2 tablas sincronizadas
- Hay un pequeño delay (milisegundos) entre escribir y leer
- Necesito RabbitMQ corriendo
- Si algo falla, los modelos pueden quedar inconsistentes

## Lo que pensé hacer en vez de esto

- **CRUD normal**: Más simple pero no escalaba bien
- **Event Sourcing completo**: Guardar todos los eventos sin tablas... demasiado complejo para esto

---

**Conclusión**: CQRS tiene sentido cuando las lecturas son mucho más frecuentes que las escrituras, que es este caso. Además el challenge lo pedía explícitamente.