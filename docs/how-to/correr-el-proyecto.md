# Correr el proyecto

> **Estado:** Aprobado

Todos los comandos salen del [Taskfile](/Taskfile.yml), desde la raíz (`task --list` los muestra)

## Requisitos

- **[Task](https://taskfile.dev)** para las tareas del repo
- **Node.js** 26+ (`engines` en `package.json`) y **pnpm**, que se instala aparte porque Node ya no trae Corepack
  (`npm install -g pnpm`). Cualquier versión sirve: en el repo pnpm usa la fijada en `packageManager`.
  `build.mjs` no tiene dependencias, pnpm solo instala las herramientas de calidad
- **Go** para la API (la versión de `apps/api/go.mod`; con un Go más viejo se descarga sola)

Primera vez: `task install` (dependencias de JS, hooks de git y herramientas de Go)

- **Podman o Docker** para el contenedor. El `Taskfile.yml` detecta el motor; para forzarlo, `ENGINE=docker task up`

## Opción A: API en el host

```bash
task api
```

Compila el portal y levanta la API en `:8080`. Abrir `http://localhost:8080`. Sin `OPNSENSE_URL`, el acceso no se abre
en ningún equipo: la API solo lo registra en el log, así que sirve para probar el formulario y la base.

## Opción B: todo en contenedor

```bash
task up     # construye y levanta; la base queda en el volumen portal_data
task logs
task down
```

Otro puerto: `PORTAL_PORT=8090 task up`

## Solo el portal

```bash
task web:dev
```

Sirve `apps/web/src` en `:8000` en modo demo (sin API). Para ver el design system, servir `apps/web` en su lugar
(`python3 -m http.server -d apps/web`) y abrir `/design-system.html`.

Estados de vista previa: `index.html?step=listo` y `index.html?step=datos-error`.

## Antes de un commit

Los hooks de git validan solos: el pre-commit formatea y lintea lo que está en stage y el pre-push corre
`task check` (ver [ci.md](/docs/explanation/ops/ci.md)). Lo que no automatizan:

- Probar el flujo en el navegador a 390×844
- `task web` para que `dist/` esté al día si se va a correr la API (no va en git; la imagen lo compila sola)

## Compilar para OPNsense

```bash
cd apps/api && GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -o byte-api-freebsd .
```

Subir el binario y reiniciar la API: ver [opnsense.md](/docs/reference/setup/opnsense.md#arrancar-la-api).
