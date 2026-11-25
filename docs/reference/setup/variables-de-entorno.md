# Variables de entorno de la API

> **Estado:** Aprobado

Lo que lee `apps/api` al arrancar. Fuente: [main.go](/apps/api/main.go) y [handlers.go](/apps/api/adapters/httpapi/handlers.go)

| Variable               | Por defecto       | Uso                                                                                                               |
| ---------------------- | ----------------- | ----------------------------------------------------------------------------------------------------------------- |
| `ADDR`                 | `:8080`           | Dirección de escucha. En OPNsense: `127.0.0.1:8085` (lo que espera el `lighttpd` del portal)                      |
| `DB`                   | `portal.db`       | Archivo SQLite. En el contenedor: `/data/portal.db`                                                               |
| `PORTAL_DIR`           | `../web/dist`     | Portal compilado que se sirve en `/`                                                                              |
| `TIMEZONE`             | `America/Caracas` | Zona horaria para calcular el día (`day`) y el reset                                                              |
| `RESET_AT`             | `02:00`           | Hora local (`HH:MM`) en que cambia el día. Define la columna `day`; no corta sesiones                             |
| `REDIRECT`             | vacío             | Enlace de "Empezar a navegar". Vacío deja el del portal                                                           |
| `NETWORK_NAME`         | `Byte Invitados`  | Nombre de la red que muestra la página (título y pastilla). Se escapa como HTML                                   |
| `CAPTIVE_HOME`         | vacío             | Si está definida, cualquier host distinto del portal recibe `302` a esta URL                                      |
| `CORS_ORIGIN`          | vacío             | Origen que puede llamar a la API desde otra página. Vacío: sin encabezados CORS                                   |
| `OPNSENSE_URL`         | vacío             | GUI de OPNsense, por ejemplo `https://127.0.0.1`. **Sin ella, el acceso no se abre** (solo se registra en el log) |
| `OPNSENSE_KEY`         | vacío             | Key de la API de un usuario de OPNsense                                                                           |
| `OPNSENSE_SECRET`      | vacío             | Secret de esa key                                                                                                 |
| `OPNSENSE_ZONE`        | `0`               | ID de la zona del portal cautivo                                                                                  |
| `OPNSENSE_CA_FILE`     | vacío             | PEM de la CA de la GUI, si no es un certificado público. Con Let's Encrypt no hace falta                          |
| `OPNSENSE_SERVER_NAME` | vacío             | Nombre que trae el certificado de la GUI (SAN). Hace falta si `OPNSENSE_URL` usa una IP                           |

Ejemplo de la configuración en uso en `fwportal`: ver [opnsense.md](/docs/reference/setup/opnsense.md).

## Contenedor

En el compose se leen del archivo `.env` de la raíz. Ver [contenedores.md](/docs/explanation/ops/contenedores.md).
