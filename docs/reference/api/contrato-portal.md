# Contrato entre el portal y la API

> **Estado:** Aprobado

Los endpoints, cuerpos y errores exactos están en [openapi.yaml](/docs/reference/api/openapi.yaml). Este documento
cubre lo que el OpenAPI no describe: cómo el portal encuentra la API y cómo se identifica al dispositivo

## `window.BYTE_PORTAL`

El portal ([app.js](/apps/web/src/app.js), línea 8) lee la configuración, que debe existir antes de cargar el script

| Campo    | Uso                                                          |
| -------- | ------------------------------------------------------------ |
| `send`   | URL de `POST /api/send`; también se usa para reenviar        |
| `verify` | URL de `POST /api/verify`                                    |

La API lo inyecta al servir `/`: `window.BYTE_PORTAL={send:"/api/send",verify:"/api/verify"}`. Sin `send`, el
portal corre en **modo demo**: latencia falsa de 700 ms, código `123456` y una marca "modo demo"

## Dispositivo

La página no manda ningún identificador del dispositivo, y la API no usa la IP ni la MAC. El registro se hace con los
datos del formulario y el teléfono verificado.

## Respuestas que cambian el flujo

| Respuesta                         | Qué hace el portal                                       |
| --------------------------------- | -------------------------------------------------------- |
| `send` 2xx con `verified: true`   | Salta directo a "conectado" |
| `send` 2xx sin `verified`         | Pasa al paso del código y arranca el contador de 60 s    |
| `verify` 2xx                      | "Conectado"; `redirect` reemplaza el enlace de navegar   |
| Cualquier no 2xx                  | Mensaje de error del paso                                |

## Estados de vista previa

`apps/web/src/index.html?step=codigo|error|listo|datos-error` precarga datos de ejemplo y salta a ese estado; lo
usa el [design system](/apps/web/design-system.html)
