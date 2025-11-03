# ADR 003: Cómo manejamos múltiples clientes (Multi-tenant)

**Fecha**: 2024-10-31

## El Problema

El sistema tiene que soportar varios clientes (tenants):
- Los datos de `tenant-1` no pueden verlos los de `tenant-2`
- Cada tenant puede tener configuraciones diferentes
- No queremos tener 100 bases de datos si tenemos 100 clientes

## Lo que decidimos

Usar **un solo PostgreSQL pero filtrar todo por `tenant_id`**.

### Cómo funciona

1. Cada request tiene que traer el header `X-Tenant-Id`
2. Todas las queries filtran por ese tenant_id
3. Cada tenant puede tener su propia config (rate limits, feature flags)
````http
POST /api/v1/users
X-Tenant-Id: tenant-1

→ Se guarda con tenant_id = 'tenant-1'
````
````sql
SELECT * FROM users_write 
WHERE id = 'user-123' AND tenant_id = 'tenant-1'
````

## Por qué está bien

- Una sola base de datos = más barato
- Deploy una vez y afecta a todos
- Más simple de desarrollar

## Por qué puede ser un problema

- Si me olvido un `WHERE tenant_id =` en alguna query, expongo datos de otros tenants (muy malo)
- Si un tenant hace muchas queries, puede hacer lento a todos
- No puedo escalar un tenant específico

## Qué más consideré

**DB por tenant**: 
- Mejor aislamiento pero muy caro
- Tener que gestionar 50 bases de datos es un dolor de cabeza

**Schema por tenant**:
- Intermedio pero PostgreSQL tiene límite de schemas

---

**Conclusión**: Para empezar, tenant_id funciona bien. Si algún cliente grande necesita aislamiento total, podemos migrarlo a su propia DB después.

**IMPORTANTE**: Tengo que tener mucho cuidado de SIEMPRE filtrar por tenant_id en todas las queries.