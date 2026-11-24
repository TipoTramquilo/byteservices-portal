# Recorrido: el flujo completo en local

> **Estado:** Aprobado

Seguir una visita de punta a punta, del portal a la base, sin router real. Para el equipo real, ver
[diagnosticar-portal.md](/docs/how-to/diagnosticar-portal.md).

## 1. Levantar la API

```bash
make api
```

En la terminal: `listening on :8080`. Sin `OPNSENSE_URL` no hay router: la API solo registra el acceso en el log.

## 2. Enviar los datos

Abrir `http://localhost:8080` a 390×844 y completar nombre y apellido, cédula (V o E), empresa y un teléfono
(por ejemplo `+584125551234`). Al enviar:

```
[router] grant <ip> until <hora UTC>      ← sin OPNsense, el acceso solo se registra en el log
```

Qué pasó: [Send](/apps/api/core/service.go) validó los datos, guardó al visitante en `verified`, pidió la MAC al
router y registró la sesión en `access` con fin 30 minutos después.

## 3. Ver la base

```bash
sqlite3 portal.db 'select name,phone,company from verified;'
sqlite3 portal.db 'select ip,phone,start_at,end_at from access;'
```

El fin de cada sesión debe quedar 30 minutos después del inicio.

## 4. Volver a enviar

Enviar el mismo teléfono de nuevo: la respuesta sigue siendo `200 {"verified": true}`. Cada envío abre una sesión nueva;
no hay límite de envíos ni espera entre ellos.

## 5. Probar los errores

- Dejar el nombre con una sola palabra: `400 {"error": "invalid"}`
- Pedir `GET /api/send`: `405` con la cabecera `Allow: POST`

Contrato exacto: [openapi.yaml](/docs/reference/api/openapi.yaml). Diseño del flujo: [acceso-wifi.md](/docs/explanation/core/acceso-wifi.md).
