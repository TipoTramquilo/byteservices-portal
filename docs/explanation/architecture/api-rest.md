# API REST con JSON, en Go

> **Estado:** Aprobado

## Decisión

- **Go** para `apps/api`, con la librería estándar para HTTP y JSON
- **REST con JSON:** dos endpoints (`/api/send`, `/api/health`). Contrato en [openapi.yaml](/docs/reference/api/openapi.yaml)
- El mismo proceso sirve el portal compilado (`apps/web/dist`) desde el mismo origen

## Razones

### Go y no Rust (ni Node)

- **Despliegue:** un binario estático; compilar para el router es una línea (`GOOS=freebsd GOARCH=amd64 go build`)
- **Velocidad de desarrollo:** HTTP, JSON y SQL son estándar y simples; el tráfico del portal no necesita más
- Un binario sin runtime encaja mejor en un router que Node

### REST y no gRPC

- **El cliente es un navegador**, y además el mini-navegador del portal cautivo: gRPC necesitaría gRPC-Web y una
  librería JS que pesa más que todo el portal
- El portal llama con `fetch` y JSON; son dos endpoints con payloads diminutos

### Mismo origen para portal y API

[index](/apps/api/adapters/httpapi/handlers.go) sirve `index.html` inyectando
`window.BYTE_PORTAL={send:"/api/send"}`, como lo haría una página alojada en el router. Sin otro origen no hace falta
CORS; `CORS_ORIGIN` solo existe para una página alojada en otro dominio.

### Redirección del portal cautivo

Para cualquier host distinto del portal, `CAPTIVE_HOME` (si está definido) responde con `302` al portal. Es lo que
hace que el sistema operativo del teléfono abra la página de acceso.

## Errores

Cuerpo `{"error": "<código>"}`. Códigos: `invalid` (400), `method_not_allowed` (405) e `internal` (500). El
mapeo está en [fail](/apps/api/adapters/httpapi/handlers.go) y el catálogo en el OpenAPI.
