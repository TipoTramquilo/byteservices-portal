# Configurar Twilio

> **Estado:** En revisión

Pasos para que el código salga por WhatsApp. La decisión está en [twilio.md](/docs/explanation/integrations/twilio.md)

## En Twilio

1. Un remitente de WhatsApp: el sandbox para pruebas o un número aprobado para producción
2. Una plantilla en Content Template Builder, categoría autenticación, con la variable `{{1}}` para el código, aprobada
   por WhatsApp. Su ID empieza con `HX`
3. Con el sandbox, cada teléfono de prueba debe unirse antes mandando al número del sandbox el mensaje `join …` que
   indica Twilio

## En la API

```bash
TWILIO_ACCOUNT_SID=AC… \
TWILIO_AUTH_TOKEN=… \
TWILIO_WHATSAPP_FROM=+14155238886 \
TWILIO_CONTENT_SID=HX… \
make api
```

Con el contenedor, las mismas variables en el entorno de `make up` (las pasa el compose)

## Comprobar

Enviar datos desde el portal con un teléfono real. Si falla, `POST /api/send` responde `500` y el log trae el error
de Twilio, por ejemplo `ApiError 20003: Authentication Error` con credenciales inválidas
