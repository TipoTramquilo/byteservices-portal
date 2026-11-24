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

Compila el portal y levanta la API en `:8080`. Abrir `http://localhost:8080`. Sin `OPNSENSE_URL`, el acceso no se abre
en ningún equipo: la API solo lo registra en el log, así que sirve para probar el formulario y la base.

## Opción B: todo en contenedor

```bash
make up     # construye y levanta; la base queda en el volumen portal_data
make logs
make down
```

Otro puerto: `PORTAL_PORT=8090 make up`

## Solo el portal

```bash
make web-dev
```

Sirve `apps/web/src` en `:8000` en modo demo (sin API). Para ver el design system, servir `apps/web` en su lugar
(`python3 -m http.server -d apps/web`) y abrir `/design-system.html`.

Estados de vista previa: `index.html?step=listo` y `index.html?step=datos-error`.

## Antes de un commit

- `cd apps/api && go vet ./... && go test ./...`
- `make web` para que `dist/` esté al día (no va en git; la imagen lo compila sola)
- Probar el flujo en el navegador a 390×844

## Compilar para OPNsense

```bash
cd apps/api && GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -o byte-api-freebsd .
```

Subir el binario y reiniciar la API: ver [opnsense.md](/docs/reference/setup/opnsense.md#arrancar-la-api).
