# Variables de entorno de la API

> **Estado:** En revisión

| Variable | Valor por defecto | Para qué |
|---|---|---|
| `DATABASE_URL` | (obligatoria) | Cadena de conexión a Postgres (Supabase u otro) |
| `ADDR` | `:8080` | Dirección y puerto donde escucha la API |
| `PORTAL_DIR` | `../web/dist` | Carpeta con el portal compilado (`index.html` y `fonts/`) |
| `TIMEZONE` | `America/Caracas` | Zona horaria del día y del reset |
| `RESET_AT` | `02:00` | Hora local del reset diario (`HH:MM`) |
| `REDIRECT` | (vacío) | Enlace de "Empezar a navegar" tras conectarse |
| `OTP_SECRET` | aleatorio al arrancar | Clave HMAC de los hashes de código. Si cambia, los códigos pendientes dejan de valer |
| `TWILIO_ACCOUNT_SID` | (vacío) | Sin esta variable, el código sale en el log |
| `TWILIO_AUTH_TOKEN` | (vacío) | Token de Twilio |
| `TWILIO_WHATSAPP_FROM` | (vacío) | Número de WhatsApp de Twilio en E.164 |
| `TWILIO_CONTENT_SID` | (vacío) | ID de la plantilla aprobada de WhatsApp (`HX…`) |

`DATABASE_URL` no va en el repositorio ni en la imagen: se pasa al contenedor o al entorno del servicio.
