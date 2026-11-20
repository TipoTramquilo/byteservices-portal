# Puertos y adaptadores

> **Estado:** En revisión

`core/` tiene la lógica y los puertos. Cada adaptador vive en su paquete y no conoce a los demás.

```
apps/api/
├── core/                 lógica: validación, código OTP, límites, reset
│   ├── service.go        Send y Verify
│   ├── otp.go            generación y hash del código
│   ├── store.go          puerto Store (y los tipos Visitor y Code)
│   └── sender.go         puerto Sender
└── adapters/
    ├── httpapi/          HTTP: /api/send, /api/verify, /api/health y el portal
    ├── postgres/         Store sobre Postgres + migraciones
    └── whatsapp/         Sender: Twilio, o consola si no hay credenciales
```

| Puerto | Adaptador |
|---|---|
| `Store` | [postgres](/apps/api/adapters/postgres/store.go) |
| `Sender` | [Twilio](/apps/api/adapters/whatsapp/twilio.go) o [consola](/apps/api/adapters/whatsapp/console.go) |

Dependencias directas: `github.com/jackc/pgx/v5` y `github.com/twilio/twilio-go`.

Ver [api-rest.md](/docs/explanation/architecture/api-rest.md) y [postgres.md](/docs/explanation/database/postgres.md).
