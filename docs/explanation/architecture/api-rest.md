# API REST con JSON, en Go

> **Estado:** Aprobado

## Decisión

- **Go** para `apps/api`, con la librería estándar para HTTP y JSON
- **REST con JSON**: dos endpoints de acceso y uno de salud. Contrato en [openapi.yaml](/docs/reference/api/openapi.yaml)
- El mismo proceso sirve el portal compilado (`apps/web/dist`) desde el mismo origen

## Razones

### Go y no Rust (ni Node)

- **Despliegue:** un binario estático; compilar para otro equipo es una línea (`GOOS=linux GOARCH=arm64 go build`).
  En Rust la compilación cruzada a MIPS/ARM de routers es más engorrosa
- **Velocidad de desarrollo:** HTTP, JSON y SQL son estándar y simples; el rendimiento de Rust no se nota con el
  tráfico de un torneo
- Node estaba instalado y servía, pero un binario sin runtime encaja mejor en un servidor pequeño

### REST y no gRPC

- **El cliente es un navegador**, y además el mini-navegador del portal cautivo: gRPC necesitaría gRPC-Web, un
  proxy (Envoy) y una librería JS de decenas de KB en un portal de 11 KB con gzip
- El portal ya llamaba con `fetch` y JSON; son dos endpoints con payloads diminutos

### Mismo origen para portal y API

[index](/apps/api/adapters/httpapi/handlers.go) (líneas 76-89) sirve `index.html` inyectando
`window.BYTE_PORTAL={send:"/api/send",verify:"/api/verify"}`, con las rutas de la misma API. Sin
otro origen no hay CORS

## Errores

Cuerpo `{"error": "<código>"}`; el portal solo distingue 2xx de error. El mapeo está en
[fail](/apps/api/adapters/httpapi/handlers.go) (líneas 110-124) y el catálogo en el OpenAPI
