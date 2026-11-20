# SQLite en el mismo proceso, sin ORM

> **Estado:** Aprobado

## Decisión

- **SQLite** como base, dentro del proceso de la API, con el driver `modernc.org/sqlite` (Go puro, sin cgo)
- **SQL escrito a mano** con `database/sql`; sin ORM
- **Migraciones forward-only** numeradas, embebidas en el binario
- **Restricciones en la base** según estándares: tablas `STRICT` y `CHECK` por columna

Esquema exacto: [esquema.md](/docs/reference/database/esquema.md)

## Razones

### SQLite y no Postgres

SQLite no es un servidor: es una librería que lee y escribe un archivo, y lo normal es usarla dentro del proceso.
No hay red entre la app y la base ni otro contenedor que mantener; el respaldo es copiar el archivo (o
`VACUUM INTO` en caliente)

Límites asumidos, irrelevantes para un torneo:

- **Un solo escritor, lectores en paralelo.** En WAL: escrituras con una conexión (`SetMaxOpenConns(1)`),
  `synchronous=NORMAL`, `busy_timeout` y checkpoint cada ~16 MB; lecturas en un pool aparte (`query_only`,
  una conexión por núcleo hasta 8) que no espera al escritor
- **Una sola instancia.** Dos procesos de la API no pueden compartir el archivo; escalar horizontalmente exigiría Postgres

### `modernc.org/sqlite` y no `mattn/go-sqlite3`

`mattn/go-sqlite3` usa cgo: rompe la compilación cruzada simple y el binario estático que corre sobre `scratch`

### Sin ORM

Tres tablas y ocho consultas. Un ORM agregaría una dependencia y SQL que no se ve. Las consultas están en
[store.go](/apps/api/adapters/sqlite/store.go)

### Restricciones según estándares

SQLite ignora los largos tipo `VARCHAR(100)`: las reglas van como `CHECK` y las tablas son `STRICT`, que rechazan
tipos equivocados en vez de convertirlos

| Dato       | Estándar                                      | Regla                                                  |
| ---------- | --------------------------------------------- | ------------------------------------------------------ |
| Teléfono   | E.164                                         | `+`, sin 0 inicial, hasta 15 dígitos                   |
| IP         | IPv4 o IPv6 canónica                          | La normaliza la API (`netip`); IPv4 mapeada pasa a IPv4 |
| MAC        | EUI-48                                        | `AA:BB:CC:DD:EE:FF`, mayúsculas con `:`; NULL si el router no la da |
| Fecha-hora | ISO 8601 UTC                                  | `2026-11-19T18:03:00Z`; `datetime()` no puede dar NULL |
| Día        | ISO 8601                                      | `2026-11-19`, en la zona horaria del portal            |

SQLite no tiene tipo fecha: se eligió texto ISO 8601 (legible en el registro del torneo) sobre enteros Unix. El
núcleo valida los mismos formatos antes de llegar a la base, así los errores salen como `400` y no como fallas de SQL

### Migraciones con `PRAGMA user_version`

[migrate](/apps/api/adapters/sqlite/migrate.go) (líneas 21-48) aplica en orden los archivos
`migrations/NNNN_nombre.sql` con número mayor a `PRAGMA user_version`, cada uno en una transacción junto con la
nueva versión: una migración que falla no deja la base a medias. Sin librería ni servicio `migrate` aparte, porque la
base vive en el mismo proceso

Una migración publicada no se edita: se agrega la siguiente. Ver
[agregar-migracion.md](/docs/how-to/agregar-migracion.md)
