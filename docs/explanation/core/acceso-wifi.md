# Acceso al Wi-Fi: registro y sesión de 30 minutos

> **Estado:** vigente en la rama `captive-opnsense` (desde `aeaca09`). Por probar en teléfonos reales. La duración de OPNsense está pendiente de verificar (ver Riesgos).

## Contexto

El portal es para un **evento de un día**. El visitante llena un formulario y recibe internet durante 30 minutos. No hay
código de verificación ni envío por WhatsApp: el teléfono no necesita datos móviles para entrar.

## Decisión

1. Al enviar el formulario (`POST /api/send`), la API registra al visitante y autoriza la IP en OPNsense por **30 minutos** (`SessionTime`).
2. Cada envío abre una sesión nueva. No hay límite de envíos por teléfono ni espera entre envíos.
3. El walled garden, el WhatsApp y la tabla de códigos se retiraron.

## Flujo

1. **Conexión.** El visitante se conecta al Wi-Fi. El router no le da internet y el sistema operativo abre el portal cautivo.
2. **Datos.** El visitante escribe nombre, cédula, empresa y teléfono. El portal llama a `POST /api/send`.
3. **Autorización.** La API toma la **IP de la conexión** (o la IP real de `X-Forwarded-For` cuando llega desde el proxy local),
   pide la **MAC al router** por ARP (si no la conoce, sigue sin MAC) y guarda el acceso.
4. **Sesión.** La API llama a `GrantAccess(ip, ahora + 30 min)`: OPNsense recibe `sessionTimeout` en segundos.
5. **Respuesta.** `200 {verified: true, redirect}`. El portal muestra "Conectado" y el enlace para navegar.

La lógica vive en [service.go](/apps/api/core/service.go) (función `Send`) y el adaptador en [opnsense.go](/apps/api/adapters/router/opnsense.go) (`GrantAccess`).

## Razones

- **IP de la petición y no MAC en la página.** El navegador no puede leer la MAC. La IP sale de la conexión
  ([clientIP](/apps/api/adapters/httpapi/handlers.go)). Requiere que el AP esté en modo puente: con NAT, todos los
  visitantes llegarían con la IP del router.
- **MAC del router, no de la página.** El backend se la pide al router con la IP. Si el router no la conoce, el acceso
  queda con IP, teléfono y cédula.
- **30 minutos por envío.** Es la duración por visita; el visitante vuelve a registrarse si necesita más.

## Límites

| Límite                | Valor           | Dónde                               |
| --------------------- | --------------- | ----------------------------------- |
| Duración de la sesión | 30 min          | `SessionTime` (`apps/api/main.go`)  |
| Envíos por teléfono   | sin límite      | `POST /api/send`                    |
| Zona del portal       | `OPNSENSE_ZONE` | OPNsense: Services → Captive Portal |

## Riesgos asumidos

- **Duración en OPNsense:** la zona tiene `hardtimeout 30` e `idletimeout 0` (verificado en `fwportal`). La columna
  _End_ de _Sessions_ es la última actividad, no el vencimiento: la prueba es ver si una sesión nueva se corta a los 30 minutos.
- **Detección en Samsung:** el override DNS de `gstatic` impedía que el Samsung abriera el portal; se quitó el 24/11 y
  conectó, y el Pixel también abre el portal sin él. Ver [diagnosticar-portal.md](/docs/how-to/diagnosticar-portal.md).
- **Número ajeno:** quien escriba el teléfono de otra persona queda con acceso; el sistema no verifica la identidad.

## Verificación

- En **Services → Captive Portal → Sessions**, compara la hora de inicio y de fin de una sesión nueva: debe ser de 30 minutos.
- En el log de la API (`/var/log/byte-api.log`), cada envío registra `send desde <IP>` con la IP usada para autorizar.

Contrato exacto: [openapi.yaml](/docs/reference/api/openapi.yaml)
