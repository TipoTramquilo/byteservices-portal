# Diagnosticar el portal: errores y pruebas

> **Estado:** Aprobado. Samsung y Pixel abren el portal desde el 24/11/2026, tras quitar el override de `gstatic`

Errores reales que aparecieron al poner el portal en OPNsense, con la causa comprobada y lo que se hizo. Las
comprobaciones se hacen en la shell de OPNsense (`tcsh`, ver [opnsense.md](/docs/reference/setup/opnsense.md)).

## Errores de la API

| Síntoma en el navegador                                                                         | Causa comprobada                                                              | Cómo se comprobó                                                  | Solución                                                                                 |
| ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------- | ----------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `500` en `POST /api/send`                                                                       | La API escuchaba en `:8080`; `lighttpd` reenvía `/api` a `127.0.0.1:8085`     | `listening on :8080` en el log                                    | Arrancar con `ADDR=127.0.0.1:8085`                                                       |
| `500` y log con `certificate is valid for fwportal.byteservices.net, not wifi.byteservices.net` | `OPNSENSE_SERVER_NAME` no coincidía con el certificado de la GUI              | Log de `getArp`                                                   | `OPNSENSE_SERVER_NAME=fwportal.byteservices.net`                                         |
| Log con `falta certificado CA local`                                                            | Se creía que hacía falta un PEM de CA                                         | El certificado de la GUI es de Let's Encrypt (`openssl s_client`) | No hace falta `OPNSENSE_CA_FILE`: el sistema ya confía en Let's Encrypt                  |
| `getArp: HTTP 403`                                                                              | La API key no tiene el privilegio de diagnóstico ARP                          | Log del 20/11 entre 21:52 y 21:58                                 | No vuelve a aparecer desde el reinicio del 24/11; confirmar en _System → Access → Users_ |
| Sesión autorizada para `127.0.0.1` y sin internet                                               | `lighttpd` reenvía desde su propio proceso; la API tomaba `127.0.0.1` como IP | Log `send desde 127.0.0.1`                                        | La API toma la IP real de `X-Forwarded-For` cuando la conexión viene de loopback         |
| Log con `sesión de 127.0.0.1 en estado ""`                                                      | Consecuencia de lo anterior: OPNsense recibía `127.0.0.1`                     | Log del 24/11 15:32 a 15:34                                       | Corregido con el cambio de IP; no aparece en los envíos posteriores                      |

Para ver la IP que llega a la API:

```sh
tail -20 /var/log/byte-api.log
```

Cada envío registra `send desde <IP> (peer=... X-Real-Ip=... X-Forwarded-For=...)`. La IP debe ser la del dispositivo.

## Errores del dispositivo

| Dispositivo | Síntoma                                                            | Estado                                                                         |
| ----------- | ------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| iPhone      | Abre el portal al conectarse                                       | Funciona                                                                       |
| Pixel       | Abre el portal al conectarse                                       | Funciona, también sin el override de `gstatic` (24/11/2026)                    |
| Samsung     | Dice "sin conexión", ofrece conectarse, al tocar no abre el portal | **Resuelto** el 24/11/2026: conecta después de quitar el override de `gstatic` |

Lo que se comprobó en el Samsung:

- **Abrir el portal a mano sí funciona** (escribir la dirección del portal en el navegador).
- **`http://neverssl.com` no abre el portal** (tampoco `http://1.1.1.1`). Una prueba anterior se hizo con
  `neverrssl.com` (error de escritura): ese dominio no existe, así que esa respuesta no sirvió de prueba.
- **La IP del Samsung está en la red del portal** (`172.16.1.200`, puerta de enlace `172.16.0.1`), y no tenía sesión
  autorizada.

Hipótesis en prueba:

1. **DNS de la prueba de conexión (confirmada).** El override de `connectivitycheck.gstatic.com` devolvía `172.16.0.1`
   (privada). Un tip de Fortinet indica que los Samsung con Android reciente solo intentan la detección si la respuesta
   DNS es una IP pública ([fuente](https://community.fortinet.com/fortinac-28/technical-tip-samsung-phones-with-newer-android-versions-not-automatically-launching-browser-97100)).
   Al quitar el override (24/11, 12:08) el Samsung conectó, y el Pixel, que usaba ese mismo override, también abre el portal sin él.
2. **HTTPS automático del navegador (sin probar).** Si el navegador sube `http://` a `https://`, la intercepción del portal
   (que solo cubre el puerto 80) no lo alcanza. Ya no hace falta para el Samsung; queda como causa posible de otros casos.

Procedimiento para repetir la prueba en un dispositivo nuevo:

1. Olvidar la red en el dispositivo y reconectarse.
2. Abrir `http://neverssl.com`. Si redirige al portal, la detección funciona y el fallo era solo la prueba automática.
3. Si no redirige, probar `http://1.1.1.1`. Si tampoco, la conexión HTTP no llega al portal en ese dispositivo.
4. Revisar el DNS del dispositivo (debe ser `172.16.0.1`) y que el DNS privado esté desactivado.

## Problemas de DNS (Unbound)

- **`unbound-control` no funciona en el equipo** (`control-enable is 'no'`). No hay forma de vaciar la caché desde la shell.
- **Después de borrar un override, la respuesta vieja sigue por su TTL.** Un `drill` dio `172.16.0.1` con TTL 3600 justo después de
  borrarlo; al repetir dio la IP pública con TTL 300. Esperar el TTL o reiniciar Unbound desde la GUI
  (_Services → Unbound DNS → General_).

Comprobación:

```sh
drill connectivitycheck.gstatic.com @172.16.0.1
```

La respuesta debe ser una IP pública si el override ya no existe.

## Sesiones que no se cortan o no aparecen

- La columna **End** de _Services → Captive Portal → Sessions_ es la **última actividad**, no el vencimiento: sube mientras
  el dispositivo navega. Una sesión activa puede verse con fin de 3 minutos después del inicio.
- Para confirmar la duración, la zona debe tener `hardtimeout 30` e `idletimeout 0` (ver [opnsense.md](/docs/reference/setup/opnsense.md)).
  Una sesión nueva debe cortarse 30 minutos después de su inicio.

## Ver el estado de la sesión de un dispositivo

```sh
sh -c 'set -a; . /usr/local/etc/byte/apikey.env; set +a; curl -sk -u "$OPNSENSE_KEY:$OPNSENSE_SECRET" https://127.0.0.1/api/captiveportal/session/list' | grep 172.16.1.200
```

Si no sale nada, el dispositivo no tiene sesión. La clave se lee del archivo de permisos 600 y no se imprime.

## Errores de shell en OPNsense

| Mensaje                                   | Causa                                         | Solución                                                                   |
| ----------------------------------------- | --------------------------------------------- | -------------------------------------------------------------------------- |
| `Illegal variable name.`                  | `procstat -e` mezcla el PID con las variables | No usar `procstat` para recuperar el entorno; usar el archivo de variables |
| `Ambiguous output redirect.`              | `tcsh` no acepta `2>/dev/null` con comodines  | Envolver en `sh -c '...'`                                                  |
| `Bad -c option`                           | Comillas sin cerrar al pegar el comando       | Pegar el comando completo en una línea                                     |
| `Unmatched '''.`                          | Lo mismo: el pegado cortó la línea            | Pegar de nuevo                                                             |
| `Ambiguous output redirect.` con `$(...)` | `tcsh` no soporta `$(...)`                    | Usar `sh -c`                                                               |
