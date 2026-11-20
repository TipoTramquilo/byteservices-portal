# Recorrido: el flujo completo en local

> **Estado:** En revisión

Seguir una visita de punta a punta, del portal a la base, sin Twilio. Necesitas una base Postgres (Supabase o una local).

## 1. Levantar la API

```bash
export DATABASE_URL='postgres://usuario:clave@host:5432/base'
make api
```

En la terminal: `listening on :8080`. Las migraciones se aplican al arrancar. Sin `OTP_SECRET` la API usa una clave
aleatoria, así que los códigos pendientes dejan de valer al reiniciar (esperado en local).

## 2. Enviar los datos

Abrir `http://localhost:8080` a 390×844 y completar nombre, cédula, empresa y un WhatsApp venezolano.
En la terminal aparece el código, porque sin Twilio el envío es por consola:

```
[whatsapp] code for +58…: 123456           ← el código (adaptador de consola)
```

Qué pasó: [Send](/apps/api/core/service.go) validó los datos, guardó el hash del código en `codes` y lo mandó con el `Sender`.

## 3. Equivocarse

Escribir `000000`: el portal muestra "Código incorrecto" y `codes.attempts` sube a 1.

## 4. Verificar

Escribir el código de la terminal: "Listo, ya estás conectado". En la base, la fila de `codes` desaparece y aparece
una en `verified` con el nombre, la cédula, la empresa y la hora de verificación (en UTC).
[Verify](/apps/api/core/service.go) hizo ese movimiento.

## 5. Volver otro día

Recargar la página y mandar el mismo teléfono: el portal salta directo a "conectado", sin código.
`send` respondió `{"verified": true}`.

## 6. Mirar los límites

Pedir el código de nuevo antes de 60 s responde `429 too_soon`; ver la tabla de límites en
[acceso-wifi.md](/docs/explanation/core/acceso-wifi.md).
