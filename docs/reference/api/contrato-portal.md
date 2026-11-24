# Contrato entre el portal y la API

> **Estado:** Aprobado

Los endpoints, cuerpos y errores exactos están en [openapi.yaml](/docs/reference/api/openapi.yaml). Este documento
cubre lo que el OpenAPI no describe: cómo el portal encuentra la API y cómo se identifica al dispositivo

## `window.BYTE_PORTAL`

El portal ([app.js](/apps/web/src/app.js)) lee la configuración, que debe existir antes de cargar el script

| Campo  | Uso                                              |
| ------ | ------------------------------------------------ |
| `send` | URL de `POST /api/send`                          |

La API la inyecta al servir `/` y `/index.html`: `window.BYTE_PORTAL={send:"/api/send"}`. Sin `send`, el portal
corre en **modo demo**: latencia falsa de 700 ms y una marca "modo demo"

## Dispositivo

La página no manda ningún identificador del dispositivo: el navegador no puede leer la MAC. La API usa la IP de la
conexión ([clientIP](/apps/api/adapters/httpapi/handlers.go)) y, con ella, le pide la MAC al router y le abre el
acceso. Solo cuando la conexión llega desde `127.0.0.1` (el proxy `lighttpd` de OPNsense) toma la IP real de
`X-Real-Ip` o `X-Forwarded-For`; desde cualquier otra IP ignora esas cabeceras

## Respuestas que cambian el flujo

| Respuesta                    | Qué hace el portal                                    |
| ---------------------------- | ----------------------------------------------------- |
| `send` 2xx                   | Muestra "Conectado". Si trae `redirect`, el botón de navegar lo usa |
| `send` 400                   | Marca los campos con error                            |
| `send` 500 u otro no 2xx     | Muestra "No pudimos conectarte. Revisa tus datos"     |

Solo hay dos vistas: el formulario de datos y "Conectado". No hay paso de código.

## Estados de vista previa

`apps/web/src/index.html?step=listo|datos-error` precarga datos de ejemplo y salta a ese estado; lo usa el
[design system](/apps/web/design-system.html)
