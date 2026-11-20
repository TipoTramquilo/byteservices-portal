# Esquema de la base

> **Estado:** En revisión

Las tablas se crean con la migración [0001_init.sql](/apps/api/adapters/postgres/migrations/0001_init.sql).

## `verified`

Teléfonos que pasaron el código. Un teléfono aparece una sola vez.

| Columna | Tipo | Regla |
|---|---|---|
| `phone` | `TEXT` | Clave primaria. E.164: `+` y 8 a 15 dígitos |
| `name` | `TEXT` | 3 a 100 caracteres |
| `id_number` | `TEXT` | `V` o `E` seguido de 6 a 8 dígitos |
| `company` | `TEXT` | 2 a 100 caracteres |
| `verified_at` | `TIMESTAMPTZ` | Hora UTC de la verificación |

## `codes`

Código pendiente por teléfono. Cada envío reemplaza al anterior.

| Columna | Tipo | Regla |
|---|---|---|
| `phone` | `TEXT` | Clave primaria. Mismo formato que `verified` |
| `name`, `id_number`, `company` | `TEXT` | Mismas reglas que `verified` |
| `hash` | `TEXT` | 64 caracteres: HMAC-SHA256 en hex. Nunca el código |
| `expires_at` | `TIMESTAMPTZ` | Vence a los 10 minutos |
| `attempts` | `INTEGER` | Intentos fallidos, mínimo 0 |
| `sent_at` | `TIMESTAMPTZ` | Hora UTC del último envío |
| `day` | `TEXT` | Día del reset al que pertenece el envío (`YYYY-MM-DD`) |
| `sends_today` | `INTEGER` | Envíos del día, mínimo 1 |

Las horas se guardan en UTC.
