# Contenedor sobre `scratch` con Podman

> **Estado:** Aprobado

## Decisión

Una imagen con la API y el portal compilado, construida con [Dockerfile](/docker/Dockerfile) y levantada con
[docker-compose.yaml](/docker/docker-compose.yaml). El `Makefile` detecta Podman o Docker

- **Build:** `golang:alpine`, binario estático (`CGO_ENABLED=0`, `-trimpath -ldflags="-s -w"`) con caché de
  módulos y de compilación
- **Runtime:** `scratch` con el binario, `apps/web/dist` y los certificados CA. ~12 MB
- **Usuario:** `65532:65532`, sin root. `scratch` no tiene `/etc/passwd`: va por número
- **Datos:** volumen `portal_data` en `/data` con la base SQLite

## Razones

- **`scratch` y no distroless.** El binario es estático y solo necesita los certificados CA para hablar con Twilio por
  HTTPS; distroless sumaba ~3 MB sin uso. Además `distroless/static-debian12` está deprecado en el registro
- **Tags sin fijar.** `golang:alpine` toma la última versión en cada build, por decisión del proyecto
- **Sin zona horaria en la imagen.** `TIMEZONE` no se fija en el Dockerfile; el binario embebe `time/tzdata` y
  su valor por defecto es `America/Caracas`. El log sale en UTC
- **El binario hace de healthcheck.** `scratch` no tiene shell ni curl: `captiveportal health` pide `/api/health` y
  sale con 0 o 1
- **Healthcheck repetido en el compose.** Podman construye en formato OCI, que descarta el `HEALTHCHECK` del Dockerfile
- **`in_pod: false`.** Sin pod, el puerto se publica con `rootlessport` y no depende del
  contenedor infra

## Dónde usarlo

| Destino           | Recomendación                                                              |
| ----------------- | -------------------------------------------------------------------------- |
| Equipo en la LAN  | Contenedor                                                                 |
| Mini PC OPNsense  | Binario solo: OPNsense es FreeBSD y no corre Podman                        |
| OpenWrt chico     | Binario solo; Podman pesa demasiado para el router                         |

## En producción

`OTP_SECRET` debe fijarse: sin ella se genera una clave al arrancar y los códigos pendientes dejan de servir tras
un reinicio. Variables: [variables-de-entorno.md](/docs/reference/setup/variables-de-entorno.md)
