# Correr el proyecto

> **Estado:** Aprobado

Todos los comandos salen del [Makefile](/Makefile), desde la raíz

## Requisitos

- **Node.js** para compilar el portal (`build.mjs` no tiene dependencias)
- **Go** para la API
- **Podman o Docker** para el contenedor. El `Makefile` detecta el motor; para forzarlo, `make ENGINE=docker`

## Opción A: API en el host

```bash
make api
```

Compila el portal y levanta la API en `:8080`. Abrir `http://localhost:8080`. Sin Twilio configurado, el código aparece en la terminal

## Opción B: todo en contenedor

```bash
make up     # construye y levanta; la base queda en el volumen portal_data
make logs   # el código aparece aquí sin Twilio
make down
```

Otro puerto: `PORTAL_PORT=8090 make up`

## Solo el portal

```bash
make web-dev
```

Sirve `apps/web/src` en `:8000` en modo demo (código `123456`), sin API. Para ver el design system, servir
`apps/web` en su lugar (`python3 -m http.server -d apps/web`) y abrir `/design-system.html`

## Antes de un commit

- `make web` para correr la API local con el portal nuevo: `dist/` no va en git, la imagen lo compila sola
- `cd apps/api && go vet ./...`
- Probar el flujo en el navegador a 390×844
