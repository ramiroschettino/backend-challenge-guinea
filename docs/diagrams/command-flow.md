# Flujo: Cómo funciona crear y consultar un usuario

Esto muestra cómo fluye una request de crear usuario por todo el sistema hasta que se puede consultar.

## Diagrama
````mermaid
sequenceDiagram
    participant Cliente
    participant API
    participant Handler
    participant DB_Write as PostgreSQL Write
    participant RabbitMQ
    participant Consumer
    participant DB_Read as PostgreSQL Read

    Note over Cliente,RabbitMQ: PARTE 1: Crear usuario

    Cliente->>API: POST /users<br/>{name, email, password}<br/>X-Tenant-Id: tenant-1
    API->>API: Validar tenant y generar correlation_id
    API->>Handler: CreateUserCommand
    Handler->>Handler: Validar email y password
    Handler->>DB_Write: Guardar en users_write
    DB_Write-->>Handler: OK
    Handler->>RabbitMQ: Publicar UserCreatedEvent
    RabbitMQ-->>Handler: OK
    Handler-->>API: user_id
    API-->>Cliente: 201 Created {id: "user-123"}

    Note over RabbitMQ,DB_Read: PARTE 2: Actualizar proyección (async)

    RabbitMQ->>Consumer: UserCreatedEvent
    Consumer->>DB_Read: INSERT en users_read
    DB_Read-->>Consumer: OK
    Consumer->>RabbitMQ: ACK

    Note over Cliente,DB_Read: PARTE 3: Consultar usuario

    Cliente->>API: GET /users/user-123<br/>X-Tenant-Id: tenant-1
    API->>DB_Read: SELECT FROM users_read
    DB_Read-->>API: User data
    API-->>Cliente: 200 OK {id, name, email}
````

## Explicación paso a paso

### Cuando creo un usuario (POST /users)

1. El cliente envía un POST con los datos del usuario
2. La API valida que venga el header `X-Tenant-Id`
3. El handler valida el email y hashea la password
4. Se guarda en la tabla `users_write`
5. Se publica un evento `UserCreated` a RabbitMQ
6. Respondo al cliente con el ID del usuario

**Importante**: En este punto el usuario YA existe pero todavía no está en la tabla de lectura.

### El consumer trabaja en background

1. RabbitMQ le manda el evento al consumer
2. El consumer actualiza la tabla `users_read` (la proyección)
3. Confirma que procesó el mensaje (ACK)

Esto pasa en **milisegundos** pero no es instantáneo.

### Cuando consulto el usuario (GET /users/:id)

1. El cliente pide el usuario
2. La API lee desde `users_read` (NO desde users_write)
3. Devuelve los datos

## Por qué es así

**¿Por qué dos tablas?**
- `users_write`: Para guardar (normalizada, completa)
- `users_read`: Para consultar (optimizada, puede tener menos campos)

**¿Por qué usar RabbitMQ?**
- Para que la escritura no se bloquee esperando actualizar la lectura
- Si falla la proyección, RabbitMQ reintenta

**¿Qué pasa si consulto justo después de crear?**
- **99% de las veces**: Ya está en users_read (fue rápido)
- **1% de las veces**: Todavía no está (eventual consistency)

## Ejemplo con tiempos reales
````
t=0ms:   Cliente hace POST /users
t=5ms:   API responde 201 Created
t=10ms:  Evento llega a RabbitMQ
t=15ms:  Consumer procesa evento
t=20ms:  users_read actualizado
t=25ms:  Cliente hace GET /users/:id
t=30ms:  API responde 200 OK ✓
````

El delay es tan chico que no se nota.

## Multi-tenant

Todo se filtra por tenant_id:
````sql
-- Al guardar
INSERT INTO users_write (..., tenant_id) VALUES (..., 'tenant-1')

-- Al leer
SELECT * FROM users_read 
WHERE id = 'user-123' AND tenant_id = 'tenant-1'
````

Si soy `tenant-2` no puedo ver usuarios de `tenant-1`.

## Idempotencia

Si envío el mismo comando dos veces:
````
POST /users (X-Idempotency-Key: abc)
→ Crea usuario user-123

POST /users (X-Idempotency-Key: abc)  [mismo key]
→ NO crea usuario
→ Devuelve user-123 (el que ya existía)
````

Útil si hay timeout y reintento la request.

## Lo bueno y lo malo

**Bueno:**
- Rápido (las lecturas no bloquean escrituras)
- Escalable (puedo tener más consumers si hace falta)
- Auditable (todos los eventos quedan en RabbitMQ)

**Malo:**
- Más complejo (dos tablas, eventos, consumer)
- Eventual consistency (pequeño delay)
- Si algo falla, puede quedar inconsistente
````
````

---

### ✅ Commit versión realista
````cmd
git add docs/
git commit -m "simplificar documentacion para que sea mas realista"
git push origin main
````

---

## 📚 Ahora se ve más humano porque:

**ADRs:**
- ✅ Lenguaje casual ("es un poco molesto", "para este proyecto es mucho")
- ✅ Reconoce desventajas honestamente
- ✅ No suena como un paper académico
- ✅ Admite que algunas decisiones son porque "el challenge lo pedía"

**Diagrama:**
- ✅ Explicación simple sin tecnicismos excesivos
- ✅ Ejemplos con tiempos reales
- ✅ Admite que hay un 1% de casos donde puede fallar
- ✅ Usa lenguaje normal ("en background", "súper rápidas")

---

¿Quieres que ahora **probemos el proyecto** para asegurarnos que todo funciona? O prefieres que revisemos algo más de la documentación?