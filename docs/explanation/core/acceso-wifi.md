# Acceso al Wi-Fi: hora gratis, código y reset diario

> **Estado:** Aprobado

## Contexto

El código de verificación llega por WhatsApp, pero el visitante no tiene internet hasta verificarlo:
sin datos móviles el código nunca llega. El portal es para un **torneo de un día**

## Decisión

1. Al enviar los datos (`POST /api/send`) el dispositivo recibe **1 h de internet completo**, una vez al día por IP,
   MAC, teléfono y cédula: si cualquiera de ellos ya la usó hoy, no hay otra
2. Al verificar el código (`POST /api/verify`) el teléfono queda verificado y el dispositivo tiene internet **hasta el reset de las 2:00 AM**
3. Los días siguientes, un teléfono ya verificado entra **sin código**: `POST /api/send` responde `verified: true` y da internet hasta el reset

## Flujo

### Paso a paso

**1. Conexión.** El visitante se conecta al Wi-Fi. El router todavía no le da internet y el sistema operativo abre el
portal cautivo (en iOS y Android, un mini-navegador)

**2. Datos.** El visitante escribe nombre, cédula, empresa donde trabaja y WhatsApp, y el portal llama a `POST /api/send`. La página no manda
nada que identifique al dispositivo: la API toma la **IP de la conexión** y con ella le pide la **MAC al router**, que
la conoce por ARP. Si el router no la sabe, se sigue sin MAC

**3. ¿Ya verificado?** La API busca el teléfono en `verified`. Si lo verificó otro día, no hace falta código: el router
le da internet **hasta el reset**, se registra el acceso y el portal salta directo a "Conectado"

**4. Límites de envío.** Si el teléfono es nuevo, la API revisa que no se haya pedido un código hace menos de 60 s
(`429 too_soon`) ni se hayan hecho ya 5 envíos hoy a ese número (`429 too_many`). Las peticiones del **mismo teléfono**
se atienden de una en una (un lock en memoria por número, el despliegue es una sola instancia), así que dos pestañas o
dos móviles a la vez no generan dos códigos: la segunda recibe `too_soon` y el portal lo interpreta como enviado

**5. Hora gratis.** La API revisa en `access` si **la IP, la MAC, el teléfono o la cédula** ya recibieron su hora gratis
hoy. Si ninguno la usó, el router abre **1 h de internet**; si alguno ya la usó, no hay otra hora, pero el código se
manda igual (quien ya tiene internet lo recibe; quien no, puede usar sus datos móviles)

**6. Código.** La API genera un código de 6 dígitos, guarda solo su hash (vence en 10 min) y lo manda por WhatsApp con
Twilio. Gracias a la hora gratis, el teléfono tiene internet para recibirlo. El portal pasa al paso del código y arranca
el contador de reenvío

**7. Verificación.** El visitante escribe el código y el portal llama a `POST /api/verify`. La API busca el código
pendiente de ese teléfono y lo rechaza si no hay (`400 wrong_code`), si ya se usaron los 5 intentos (`429 too_many`:
hay que pedir otro), si venció o si el hash no coincide (`400 wrong_code`; un intento más)

**8. Conectado.** Con el código correcto, los datos pasan a `verified`, el código se borra y el router da internet
**hasta el reset**. A partir de ahí, ese teléfono entra sin código los días siguientes (paso 3)

### Secuencia completa

```mermaid
sequenceDiagram
    autonumber
    actor V as Visitante
    participant P as Portal (apps/web)
    participant A as API (apps/api)
    participant R as Router (OPNsense)
    participant D as SQLite
    participant W as Twilio (WhatsApp)

    V->>R: se conecta al Wi-Fi
    R-->>V: abre el portal cautivo
    V->>P: nombre, cédula, empresa, WhatsApp
    P->>A: POST /api/send {name, id_number, company, phone}
    Note over A: IP = la de la conexión
    A->>R: MAC(ip)
    R-->>A: MAC por ARP, o vacía
    A->>D: ¿teléfono verificado otro día?

    alt ya verificado
        A->>R: GrantAccess(ip, reset)
        A->>D: access (verified)
        A-->>P: 200 {verified: true}
        P-->>V: "Conectado", sin código
    else teléfono nuevo
        A->>D: ¿IP, MAC, teléfono o cédula ya usaron su hora hoy?
        opt ninguno la usó
            A->>R: GrantAccess(ip, +1 h)
            A->>D: access (free)
        end
        A->>D: codes: hash del código, vence en 10 min
        A->>W: código de 6 dígitos
        A-->>P: 200
        W-->>V: código por WhatsApp (llega gracias a la hora gratis)
        V->>P: escribe el código
        P->>A: POST /api/verify {phone, code}
        A->>D: compara el hash, suma intentos
        alt código correcto
            A->>D: verified + borra el código
            A->>R: GrantAccess(ip, reset)
            A->>D: access (verified)
            A-->>P: 200
            P-->>V: "Conectado"
        else incorrecto o vencido
            A-->>P: 400 wrong_code
            P-->>V: "Código incorrecto"
        end
    end
```

### Decisión de `POST /api/send`

```mermaid
flowchart TD
    S([POST /api/send]) --> VAL{"¿datos válidos?<br/>nombre, cédula V/E, E.164"}
    VAL -- no --> E400[400 invalid]
    VAL -- sí --> MAC["MAC(ip) al router"]
    MAC --> VER{"¿teléfono en verified?"}
    VER -- sí --> FULL["internet hasta el reset"] --> OKV(["200 {verified: true}"])
    VER -- no --> WAIT{"¿código enviado hace menos de 60 s?"}
    WAIT -- sí --> E429a[429 too_soon]
    WAIT -- no --> MAX{"¿5 envíos hoy a ese teléfono?"}
    MAX -- sí --> E429b[429 too_many]
    MAX -- no --> FREE{"¿IP, MAC, teléfono o cédula<br/>ya usaron su hora hoy?"}
    FREE -- no --> HOUR["1 h de internet"] --> CODE
    FREE -- sí --> CODE["guarda el hash y manda el código"]
    CODE --> OK([200])
```

### Decisión de `POST /api/verify`

```mermaid
flowchart TD
    S([POST /api/verify]) --> VAL{"¿teléfono y código válidos?"}
    VAL -- no --> E400[400 invalid]
    VAL -- sí --> PEND{"¿hay código pendiente?"}
    PEND -- no --> WRONG[400 wrong_code]
    PEND -- sí --> TRIES{"¿5 intentos usados?"}
    TRIES -- sí --> E429["429 too_many<br/>pedir otro con /api/send"]
    TRIES -- no --> EXP{"¿vencido?<br/>10 min"}
    EXP -- sí --> WRONG
    EXP -- no --> MATCH{"¿coincide el hash?"}
    MATCH -- no --> ADD["suma un intento"] --> WRONG
    MATCH -- sí --> DONE["verified + borra el código<br/>internet hasta el reset"] --> OK([200])
```

La lógica vive en [service.go](/apps/api/core/service.go) (líneas 58-161)

## Razones

- **Hora gratis y no walled garden.** Un walled garden solo para WhatsApp (dominios y rangos IP de Meta) dejaría
  usar WhatsApp sin registrarse; la hora gratis es más simple de configurar en cualquier router. En un torneo de un
  día, una hora de uso sin verificar es aceptable
- **Verificar da internet hasta el reset.** Sin un premio por verificar, mucha gente no escribiría el código: con internet
  activo, iOS da la red por conectada y cierra el portal
- **IP de la petición y no MAC.** El navegador no puede leer la MAC y, pasada por la URL, la escribe el visitante.
  La IP sale de la conexión ([clientIP](/apps/api/adapters/httpapi/handlers.go) (líneas 91-99)) y es lo que usa el
  portal cautivo de OPNsense para abrir el acceso. Requiere que el AP esté en modo puente: con NAT, todos los
  visitantes llegarían con la IP del router
- **MAC del router, no de la página.** El backend se la pide al router con la IP (`Router.MAC`; el router la sabe
  por ARP). Si el router no la conoce, el acceso queda con IP, teléfono y cédula
- **IP, MAC, teléfono y cédula para la hora gratis.** Los teléfonos modernos usan MAC privada y la rotan; la IP suele
  cambiar con ella; teléfono y cédula los escribe el visitante sin verificar. Ninguno basta solo: hay que cambiar
  todos a la vez
- **Saltar el código otros días.** El backend recuerda los teléfonos verificados (tabla `verified`); volver al día
  siguiente no obliga a esperar otro WhatsApp

## El reset diario y las horas

- **El día va de reset a reset (2:00 AM a 2:00 AM)**, no de medianoche a medianoche. A la 1:00 AM del 21 todavía es
  el día del 20, así que no hay otra hora gratis ni se reinician los límites, y todo acceso vence a las 2:00 AM. Lo calcula [window](/apps/api/core/service.go) (líneas 185-195)
- **Todas las horas se guardan en UTC** (ISO 8601 con `Z`): `verified_at`, `expires_at`, `sent_at`, `start_at`,
  `end_at`. La zona horaria (`TIMEZONE`) solo sirve para saber cuándo son las 2:00 AM locales; el log también sale en UTC
- La columna `day` no es una hora sino la etiqueta del día: la fecha local en que empezó (`2026-11-20` para
  todo lo que pasa entre las 2:00 AM del 20 y las 2:00 AM del 21)

## Riesgos asumidos

- **Uso sin verificar:** quien no escribe el código igual navega su hora
- **Rotación de MAC e IP con datos nuevos:** quien olvide la red (nueva MAC, nueva IP) y escriba otro teléfono y otra
  cédula obtiene otra hora gratis
- **Número ajeno verificado:** quien escriba el teléfono de alguien ya verificado entra sin código
- **Llegada cerca del reset:** la hora gratis se corta en el reset; a la 1:58 AM solo quedan 2 min

## Límites

| Límite                    | Valor  | Dónde                                                          |
| ------------------------- | ------ | -------------------------------------------------------------- |
| Hora gratis               | 1 h    | `FreeTime` en [main.go](/apps/api/main.go) (líneas 57-67)      |
| Vigencia del código       | 10 min | `CodeTTL`                                                      |
| Espera entre envíos       | 60 s   | `ResendWait`; igual a `WAIT` en el portal                      |
| Envíos por teléfono y día | 5      | `MaxSendsPerDay`                                               |
| Intentos por código       | 5      | `MaxAttempts`; agotados, se pide otro código con `/api/send`   |
| Reset diario              | 02:00  | `RESET_AT`, hora local de `TIMEZONE` (06:00 UTC en Venezuela)  |

Contrato exacto: [openapi.yaml](/docs/reference/api/openapi.yaml)
