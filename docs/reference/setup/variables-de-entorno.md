# Variables de entorno de la API

> **Estado:** Aprobado

Las lee [main.go](/apps/api/main.go). Ninguna es obligatoria para correr en local

| Variable               | Por defecto              | Uso                                                                              |
| ---------------------- | ------------------------ | -------------------------------------------------------------------------------- |
| `ADDR`                 | `:8080`                  | Dirección de escucha                                                             |
| `DB`                   | `portal.db`              | Archivo SQLite. En el contenedor: `/data/portal.db`                              |
| `PORTAL_DIR`           | `../web/dist`            | Portal compilado que se sirve en `/`. En el contenedor: `/app/dist`              |
| `TIMEZONE`             | `America/Caracas`        | Zona horaria para ubicar el reset y la columna `day`. Las horas se guardan en UTC |
| `RESET_AT`             | `02:00`                  | Hora local (`HH:MM`) del reset diario: vencen los accesos y se reinician los límites diarios |
| `REDIRECT`             | vacío                    | Enlace de "Empezar a navegar"; vacío deja el del portal                          |
| `OTP_SECRET`           | aleatoria al arrancar    | Clave HMAC de los códigos. **Fijarla en producción**: si no, un reinicio invalida los códigos pendientes |
| `TWILIO_ACCOUNT_SID`   | vacío                    | Sin ella el código se imprime en el log en vez de mandarse                        |
| `TWILIO_AUTH_TOKEN`    | vacío                    | Token de la cuenta de Twilio                                                     |
| `TWILIO_WHATSAPP_FROM` | vacío                    | Número de WhatsApp de Twilio en E.164 (`+14155238886`)                           |
| `TWILIO_CONTENT_SID`   | vacío                    | ID de la plantilla aprobada (`HX…`) con `{{1}}` para el código                   |

El compose ([docker-compose.yaml](/docker/docker-compose.yaml)) pasa `OTP_SECRET`, `REDIRECT`, `RESET_AT` y las de Twilio, y
publica el puerto en `PORTAL_PORT` (por defecto `8080`)
