# Puertos y adaptadores en `apps/api`

> **Estado:** Aprobado

## Contexto

La API depende de un router (OPNsense en el sitio real) que da o quita internet. La lógica (sesión de 30 minutos,
registro del visitante) no debe cambiar al cambiar el router ni la base.

## Decisión

Hexagonal ligero: un núcleo `core/` con la lógica y sus puertos, y un paquete por adaptador en `adapters/`

```
apps/api/
├── main.go             arma todo: elige adaptadores según las variables de entorno
├── core/
│   ├── service.go      Send(): toda la lógica
│   ├── store.go        puerto Store + tipos (Visitor, Access)
│   └── router.go       puerto Router
└── adapters/
    ├── httpapi/        entrada: JSON sobre HTTP y el portal compilado
    ├── sqlite/         Store con SQL directo + migraciones
    └── router/         Router: OPNsense (API REST) o solo log si no hay OPNSENSE_URL
```

Regla de dependencias: `core` no importa nada de `adapters`; los adaptadores importan `core`.

| Puerto   | Definición                                            | Adaptadores                                              |
| -------- | ----------------------------------------------------- | -------------------------------------------------------- |
| `Store`  | [store.go](/apps/api/core/store.go)                   | [sqlite](/apps/api/adapters/sqlite/store.go)             |
| `Router` | [router.go](/apps/api/core/router.go): `MAC(ip)` y `GrantAccess(ip, until)` | [opnsense](/apps/api/adapters/router/opnsense.go), [log](/apps/api/adapters/router/log.go) |

Sin WhatsApp: el código de verificación y su puerto `Sender` se retiraron.

## Razones

- **Puertos dentro de `core`, no en `ports/`.** En Go la interfaz la define quien la usa.
- **Un archivo por puerto, no un paquete por puerto.** En Go los archivos ordenan y los paquetes aíslan.
- **Sin Clean Architecture completa.** Capas `domain/application/infrastructure` con DTOs entre ellas serían
  ceremonia para dos endpoints y dos tablas.
- **`httpapi` y no `http`.** Un paquete `http` choca con `net/http` en cada import.

## Convenciones

- Identificadores, archivos y tablas en inglés; comentarios en español, sin punto final
- Sin framework HTTP: `net/http` enruta por método desde Go 1.22 (`POST /api/send`)
- Dependencias directas: `modernc.org/sqlite`

Ver [api-rest.md](/docs/explanation/architecture/api-rest.md) y [sqlite.md](/docs/explanation/database/sqlite.md).
