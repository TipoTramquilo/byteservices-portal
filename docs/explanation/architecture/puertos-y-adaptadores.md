# Puertos y adaptadores en `apps/api`

> **Estado:** Aprobado

## Contexto

La API depende de tres cosas externas, dos sin decidir al escribirla: el router que da internet (MikroTik u
OpenWrt) y el proveedor del mensaje (finalmente Twilio). La lógica (hora gratis, código, reset diario) no debe
cambiar cuando cambien

## Decisión

Hexagonal ligero: un núcleo `core/` con la lógica y sus puertos, y un paquete por adaptador en `adapters/`

```
apps/api/
├── main.go             arma todo: elige adaptadores según las variables de entorno
├── core/
│   ├── service.go      Send(), Verify(): toda la lógica
│   ├── otp.go          generar, hashear y comparar el código
│   ├── store.go        puerto Store + tipos (Visitor, Code, Access)
│   ├── sender.go       puerto Sender
│   └── router.go       puerto Router
└── adapters/
    ├── httpapi/        entrada: JSON sobre HTTP y el portal compilado
    ├── sqlite/         Store con SQL directo + migraciones
    ├── whatsapp/       Sender: Twilio o consola
    └── router/         Router: solo log hasta elegir router
```

Regla de dependencias: `core` no importa nada de `adapters`; los adaptadores importan `core`

| Puerto   | Definición                                     | Adaptadores                                                                 |
| -------- | ---------------------------------------------- | --------------------------------------------------------------------------- |
| `Store`  | [store.go](/apps/api/core/store.go) (línea 47) | [sqlite](/apps/api/adapters/sqlite/store.go)                                |
| `Sender` | [sender.go](/apps/api/core/sender.go)          | [Twilio](/apps/api/adapters/whatsapp/twilio.go), [consola](/apps/api/adapters/whatsapp/console.go) |
| `Router` | [router.go](/apps/api/core/router.go): `MAC(ip)` y `GrantAccess(ip, until)` | [log](/apps/api/adapters/router/log.go) (OPNsense: planificado)   |

## Razones

- **Puertos dentro de `core`, no en `ports/`.** En Go la interfaz la define quien la usa. Con un paquete `ports/`
  aparte, los tipos que comparten núcleo y puertos tendrían que vivir en `ports/` para evitar imports circulares
- **Un archivo por puerto, no un paquete por puerto.** En Go los archivos ordenan y los paquetes aíslan: un paquete
  por puerto obligaría a otro paquete más para los tipos compartidos
- **Sin Clean Architecture completa.** Capas `domain/application/infrastructure` con DTOs entre ellas serían
  ceremonia para dos endpoints y tres tablas
- **`httpapi` y no `http`.** Un paquete `http` choca con `net/http` en cada import

## Convenciones

- Identificadores, archivos y tablas en inglés; comentarios en español, sin punto final
- Sin framework HTTP: `net/http` enruta por método desde Go 1.22 (`POST /api/send`)
- Dependencias directas: `modernc.org/sqlite` y `github.com/twilio/twilio-go`

Ver [api-rest.md](/docs/explanation/architecture/api-rest.md) y [sqlite.md](/docs/explanation/database/sqlite.md)
