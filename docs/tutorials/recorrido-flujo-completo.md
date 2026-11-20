# Recorrido: el flujo completo en local

> **Estado:** Aprobado

Seguir una visita de punta a punta, del portal a la base, sin Twilio ni router

## 1. Levantar la API

```bash
make api
```

En la terminal: `listening on :8080` y un aviso de `OTP_SECRET vacío` (esperado en local)

## 2. Enviar los datos

Abrir `http://localhost:8080` a 390×844 y completar nombre, cédula y un WhatsApp venezolano.
En la terminal aparecen dos líneas:

```
[router] grant 127.0.0.1 until …     ← la hora gratis (adaptador de log)
[whatsapp] code for +58…: 123456           ← el código (adaptador de consola)
```

Qué pasó: [Send](/apps/api/core/service.go) (líneas 58-118) validó los datos, vio que ni la IP, ni el teléfono, ni la cédula habían usado su hora,
la dio con `GrantAccess`, guardó el hash del código en `codes` y lo mandó con el `Sender`

## 3. Equivocarse

Escribir `000000`: el portal muestra "Código incorrecto" y `codes.attempts` sube a 1

## 4. Verificar

Escribir el código de la terminal: "Listo, ya estás conectado". En la terminal, otro `grant` hasta las 2:00 AM (en el log, `06:00:00Z`: las horas van en UTC).
[Verify](/apps/api/core/service.go) (líneas 122-161) movió los datos a `verified`, borró el código y registró el
acceso `verified`

## 5. Volver otro día

Recargar la página y mandar el mismo teléfono (como si fuera otro día u otro dispositivo): el portal salta
directo a "conectado", sin código. `send` respondió `{"verified": true}`

## 6. Mirar los límites

Pedir el código de nuevo antes de 60 s responde `429 too_soon`; ver la tabla de límites en
[acceso-wifi.md](/docs/explanation/core/acceso-wifi.md)
