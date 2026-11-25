# SQLite en el mismo proceso, sin ORM

> **Estado:** Aprobado

## Decisión

- **SQLite** como base, dentro del proceso de la API, con el driver `modernc.org/sqlite` (Go puro, sin cgo)
- **SQL escrito a mano** con `database/sql`; sin ORM
- **Migraciones forward-only** numeradas, embebidas en el binario
- **Restricciones en la base:** tablas `STRICT` y `CHECK` por columna

Esquema exacto: [esquema.md](/docs/reference/database/esquema.md)

## Razones

### SQLite y no un servidor de base de datos

SQLite no es un servidor: es una librería que lee y escribe un archivo, dentro del proceso. No hay red entre la app
y la base ni otro contenedor que mantener. El respaldo es copiar el archivo (o `VACUUM INTO` en caliente).

Límites asumidos, irrelevantes para este uso:

- **Un solo escritor, lectores en paralelo.** En WAL: escrituras con una conexión, lecturas en un pool aparte
  (`query_only`) que no espera al escritor
- **Una sola instancia.** Dos procesos de la API no pueden compartir el archivo

### `modernc.org/sqlite` y no `mattn/go-sqlite3`

`mattn/go-sqlite3` usa cgo: rompe la compilación cruzada simple para el router (FreeBSD) y el binario estático.

### Sin ORM

Dos tablas y dos escrituras (`SaveVerified` y `SaveAccess`). Un ORM agregaría una dependencia y SQL que no se ve. Las consultas están en
[store.go](/apps/api/adapters/sqlite/store.go) y usan parámetros `?` en todos los valores.

### Restricciones según estándares

SQLite ignora los largos tipo `VARCHAR(100)`: las reglas van como `CHECK` y las tablas son `STRICT`, que rechazan
tipos equivocados en vez de convertirlos.

| Dato       | Estándar     | Regla                                                               |
| ---------- | ------------ | ------------------------------------------------------------------- |
| Teléfono   | E.164        | `+`, sin 0 inicial, hasta 15 dígitos                                |
| IP         | IPv4 o IPv6  | La normaliza la API (`netip`)                                       |
| MAC        | EUI-48       | `AA:BB:CC:DD:EE:FF`, mayúsculas con `:`; NULL si el router no la da |
| Fecha-hora | ISO 8601 UTC | `2026-11-24T18:03:00Z`; `datetime()` no puede dar NULL              |
| Día        | ISO 8601     | `2026-11-24`, en la zona horaria del portal                         |

SQLite no tiene tipo fecha: se eligió texto ISO 8601 sobre enteros Unix, porque se lee directo en el registro. El núcleo valida los mismos formatos antes de llegar a la base, así los errores salen como `400`.

### Migraciones con `PRAGMA user_version`

El adaptador aplica en orden los archivos `migrations/NNNN_nombre.sql` con número mayor a `PRAGMA user_version`, cada
uno en una transacción junto con la nueva versión: una migración que falla no deja la base a medias. Sin librería
aparte, porque la base vive en el mismo proceso.

Una migración publicada no se edita: se agrega la siguiente. Ver
[agregar-migracion.md](/docs/how-to/agregar-migracion.md).
