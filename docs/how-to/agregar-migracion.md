# Agregar una migración

> **Estado:** Aprobado

Las migraciones son forward-only. La decisión está en [sqlite.md](/docs/explanation/database/sqlite.md)

## Pasos

1. Crear `apps/api/adapters/sqlite/migrations/NNNN_nombre.sql` con el número siguiente al último (`0002_…`, `0003_…`)
2. Escribir solo el cambio: `ALTER TABLE`, `CREATE TABLE`, `CREATE INDEX`. No hace falta `BEGIN`/`COMMIT`: cada
   archivo corre en su propia transacción
3. Respetar las reglas del esquema: tablas `STRICT`, `CHECK` con el estándar del dato (E.164, ISO 8601)
4. Actualizar [esquema.md](/docs/reference/database/esquema.md) y, si cambia el contrato,
   [openapi.yaml](/docs/reference/api/openapi.yaml)
5. Arrancar la API: la migración se aplica sola y `PRAGMA user_version` pasa al número del archivo

## Reglas

- **Nunca editar una migración publicada.** Las bases que ya la aplicaron no la vuelven a correr; el cambio va en una nueva
- **El número sale del nombre.** Un archivo sin número detiene el arranque con `migración sin número`
- **Si una falla, nada se aplica a medias:** la transacción se revierte y la versión no cambia. La API no arranca
  hasta corregirla

## Ejemplo

```sql
-- 0002_visitor_email.sql
ALTER TABLE verified ADD COLUMN email TEXT CHECK (email IS NULL OR email LIKE '%_@_%');
```
