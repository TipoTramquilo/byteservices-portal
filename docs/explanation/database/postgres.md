# Base de datos: Postgres (Supabase)

> **Estado:** En revisión

## Decisión

La base es **Postgres**, en Supabase u otro proveedor, con el driver `pgx` a través de `database/sql`. Sin ORM.

## Por qué no SQLite

SQLite necesita un archivo en disco que sobreviva reinicios. Los hostings gratis (por ejemplo, el plan de Render) no
ofrecen disco persistente: si el servicio se reinicia, se pierden los registros. Postgres en Supabase guarda los datos
fuera del servicio que corre la API.

## Migraciones

- Viven en `apps/api/adapters/postgres/migrations/` como `NNNN_nombre.sql`, embebidas en el binario.
- Se aplican al arrancar la API, en orden, cada una en su propia transacción.
- Cada migración aplicada queda en la tabla `schema_migrations`.
- Son **forward-only**: no se edita una migración ya publicada; el cambio va en la siguiente.

## Restricciones

- Las columnas de fecha y hora son `TIMESTAMPTZ`. La API guarda siempre en UTC.
- Los formatos (teléfono E.164, cédula V/E) se validan en la base con `CHECK`, además de en la API.
- El hash del código se guarda; nunca el código.

## Conexión

La API lee `DATABASE_URL` (ver [variables-de-entorno.md](/docs/reference/setup/variables-de-entorno.md)). En Supabase,
usar la cadena de conexión del pooler. No va en el repositorio.

Esquema exacto: [esquema.md](/docs/reference/database/esquema.md)
