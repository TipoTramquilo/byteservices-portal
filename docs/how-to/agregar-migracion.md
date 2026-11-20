# Agregar una migración

Las migraciones son forward-only. La decisión está en [postgres.md](/docs/explanation/database/postgres.md).

1. Crear `apps/api/adapters/postgres/migrations/NNNN_nombre.sql` con el número siguiente al último
   (`0002_…`, `0003_…`). No editar una migración ya publicada.
2. Escribir el SQL. Usar `CHECK` para las reglas de formato que ya existen en [esquema.md](/docs/reference/database/esquema.md).
3. Arrancar la API contra una base vacía y otra con la versión anterior, y confirmar que la migración se aplica una sola vez.
4. Actualizar [esquema.md](/docs/reference/database/esquema.md) y, si cambia el contrato, [openapi.yaml](/docs/reference/api/openapi.yaml).
