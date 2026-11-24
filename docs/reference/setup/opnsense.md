# OPNsense: configuración del portal

> **Estado:** En revisión. Valores verificados en `fwportal.byteservices.net` (OPNsense 26.7.9_3) el 24/11/2026

Lo que tiene que estar configurado en el equipo para que la API funcione. Los errores y su diagnóstico están en
[diagnosticar-portal.md](/docs/how-to/diagnosticar-portal.md).

## Red

| Elemento | Valor |
| --- | --- |
| LAN de visitantes (`vtnet1`) | `172.16.0.1/22` |
| WAN (`vtnet0`) | DHCP |
| Certificado de la GUI | Let's Encrypt, para `fwportal.byteservices.net` |
| Certificado del portal (`:8000`) | Let's Encrypt, SAN `wifi.byteservices.net`, vence 18/02/2027 |

## Portal cautivo

**Services → Captive Portal → (zona)**

| Campo | Valor | Por qué |
| --- | --- | --- |
| Hard timeout | `30` | Corta cada sesión a los 30 minutos |
| Idle timeout | `0` | Sin corte por inactividad |

Verificado con la API: `GET /api/captiveportal/settings/get` devuelve `"hardtimeout":"30"` y `"idletimeout":"0"`.

### Lighttpd del portal

`lighttpd` escucha en `:8000` y `:9000` y reenvía `/api` a `127.0.0.1:8085` (`proxy.server` en
`/var/etc/lighttpd-cp-zone-0.conf`). Por eso la API debe escuchar en `8085`. El módulo `mod_extforward` está activo con
la cabecera `X-Real-Ip`; la API además lee `X-Forwarded-For`.

## Overrides de DNS

**Services → Unbound DNS → Overrides**

| Host | Dominio | IP | Estado |
| --- | --- | --- | --- |
| `wifi` | `byteservices.net` | `172.16.0.4` | Activo: el portal la necesita |
| `captive` | `apple.com` | `172.16.0.1` | Activo: detección de iPhone |
| `www` | `msftconnecttest.com` | `172.16.0.1` | Activo: detección de Windows |
| `connectivitycheck` | `gstatic.com` | — | **Borrado el 24/11/2026**: devolvía `172.16.0.1` y el Samsung no disparaba la detección (era `172.16.0.1`) |

Si se vuelve a crear el de `gstatic.com`, hay que recordar que devuelve una IP privada y que los Samsung con Android
reciente no disparan la detección con ella (ver [diagnosticar-portal.md](/docs/how-to/diagnosticar-portal.md)).

## API key de la API de Byte

- Archivo: `/usr/local/etc/byte/apikey.env`, permisos `600`, dueño `root`. Contiene `OPNSENSE_KEY` y `OPNSENSE_SECRET`
- El usuario de la key necesita permisos para **diagnóstico ARP** (`getArp`) y **captive portal** (session connect/list/disconnect)
- Si aparece `HTTP 403` en el log, revisar *System → Access → Users* del usuario de la key

## Arrancar la API

El binario va en `/usr/local/sbin/byte-api` y se arranca con el archivo de variables (no con el entorno de un proceso):

```sh
sh -c 'pkill -x byte-api; sleep 2; set -a; . /usr/local/etc/byte/apikey.env; set +a; export OPNSENSE_URL=https://127.0.0.1 OPNSENSE_SERVER_NAME=fwportal.byteservices.net OPNSENSE_ZONE=0 ADDR=127.0.0.1:8085 DB=/var/db/byte/portal.db PORTAL_DIR=/usr/local/www/byte TIMEZONE=America/Caracas RESET_AT=02:00 NETWORK_NAME=Byte_Services; daemon -p /var/run/byte-api.pid -o /var/log/byte-api.log /usr/local/sbin/byte-api; sleep 3; tail -3 /var/log/byte-api.log'
```

Resultado correcto: `listening on 127.0.0.1:8085`.

Pendiente: que el arranque sobreviva a un reinicio de OPNsense (servicio `rc.d` o equivalente). Hoy no lo hace.

## Log y base

| Qué | Dónde |
| --- | --- |
| Log de la API | `/var/log/byte-api.log` |
| Base SQLite | `/var/db/byte/portal.db` |
| Certificado del portal | `/var/etc/cert-cp-zone0.pem` |

Consultas útiles:

```sh
sqlite3 /var/db/byte/portal.db 'select ip,mac,phone,start_at,end_at from access order by start_at desc limit 5;'
sqlite3 /var/db/byte/portal.db 'select name,phone,company,verified_at from verified order by verified_at desc limit 5;'
```

## Notas del equipo

- El shell de root es `tcsh`: envolver comandos con `$(...)`, `2>/dev/null` o comodines en `sh -c '...'`
- `unbound-control` no está habilitado (`control-enable: no`)
