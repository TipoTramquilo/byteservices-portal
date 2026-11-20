# Prueba de carga

[k6](https://k6.io) en un contenedor, conectado a la red de compose, contra el portal que ya está
corriendo. Sirve para saber cuántas peticiones simultáneas aguanta la API y cuánto consume.

## Antes

```sh
make up                                        # portal en :8080
podman pull docker.io/grafana/k6
```

El contenedor de k6 se llama `k6run` y entra por `http://portal:8080`, así que no pasa por el
reenvío de puertos del Mac.

## Comandos

```sh
./loadtest/vus-1000.sh          # 1000 VUs en bucle, sin pausa (30 s) + CPU/RAM
./loadtest/escalera.sh          # 50 → 1000 VUs, nivel por nivel, con tabla resumen
```

Para verlo en vivo, en dos terminales:

```sh
podman stats byte-portal        # consumo del portal
podman run --rm --name k6run --network byte-portal_default \
  -v "$PWD/loadtest:/scripts" docker.io/grafana/k6 run \
  -e MODE=send -e VUS=1000 -e DUR=30s /scripts/k6.js
```

k6 dibuja la barra de progreso; al terminar imprime `rps`, `p50/p95/p99`, `fail`, `r429` y `r500`.

## Modos

| Variable | Qué hace |
| --- | --- |
| `MODE=send` | `POST /api/send` con teléfono y cédula únicos por iteración (por defecto) |
| `MODE=health` | `GET /api/health`, el techo del HTTP sin tocar la base |
| `MODE=index` | `GET /`, el portal compilado (28 KB) |
| `VUS=<n>` `DUR=<n>s` | VUs simultáneos y duración (por defecto 10 y 20s) |
| `RATE=<n>` | modelo abierto: fija `n` peticiones por segundo en vez de VUs |
| `SLEEP=<n>` | pausa entre iteraciones por VU (0 = sin pausa) |

`SLEEP` no se usa en los scripts: la idea es la peor carga, todos pidiendo a la vez.

## Qué mide y qué no

- Los VUs son **peticiones en vuelo a la vez**: con 1000 VUs, 1000 peticiones simultáneas sin esperar.
- k6 y la API viven en la **misma VM** de Podman (8 vCPU, 3.6 GB), así que a altas tasas k6 es parte
  del cuello de botella. Para el techo real haría falta un cliente fuera de la VM.
- La base SQLite usa **una sola conexión** (`SetMaxOpenConns(1)`), así que todas las consultas se
  serializan: es el límite del `send`.
- Con Twilio configurado, cada `send` hace una llamada HTTPS síncrona a Twilio; eso domina el
  tiempo de respuesta mucho antes que la base.

Los resultados quedan en `loadtest/out/` (ignorado por git).
