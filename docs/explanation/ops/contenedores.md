# Contenedor sobre `scratch` con Podman

> **Estado:** Aprobado

## Decisión

Una imagen con la API y el portal compilado, construida con [Dockerfile](/docker/Dockerfile) y levantada con
[docker-compose.yaml](/docker/docker-compose.yaml). El `Makefile` detecta Podman o Docker.

- **Portal:** etapa `node:alpine` que corre `apps/web/build.mjs` (sin dependencias); el resultado va a `/app/dist`
- **Build:** `golang:alpine`, binario estático (`CGO_ENABLED=0`) con caché de módulos y de compilación
- **Runtime:** `scratch` con el binario, `apps/web/dist` y los certificados CA (para hablar por HTTPS con la GUI del router)
- **Usuario:** `65532:65532`, sin root. `scratch` no tiene `/etc/passwd`: va por número
- **Datos:** volumen `portal_data` en `/data` con la base SQLite

## Razones

- **`scratch` y no distroless.** El binario es estático y solo necesita los certificados CA
- **Tags sin fijar.** `golang:alpine` toma la última versión en cada build, por decisión del proyecto
- **Zona horaria embebida.** El binario incluye `time/tzdata`; `TIMEZONE` por defecto es `America/Caracas`. El log sale en UTC
- **El binario hace de healthcheck.** `scratch` no tiene shell ni curl: `captiveportal health` pide `/api/health` y sale con 0 o 1
- **Healthcheck repetido en el compose.** Podman construye en formato OCI, que descarta el `HEALTHCHECK` del Dockerfile
- **`in_pod: false`.** Sin pod, el puerto se publica con `rootlessport` y no depende del contenedor infra

## Dónde usarlo

| Destino          | Recomendación                                              |
| ---------------- | ---------------------------------------------------------- |
| Equipo en la LAN | Contenedor                                                 |
| OPNsense         | Binario solo: OPNsense es FreeBSD y no corre Podman. Ver [opnsense.md](/docs/reference/setup/opnsense.md) |
| Router pequeño   | Binario solo; Podman pesa demasiado para el router         |

## Variables

Se leen del archivo `.env` de la raíz (ignorado por git). Lista completa en
[variables-de-entorno.md](/docs/reference/setup/variables-de-entorno.md).
