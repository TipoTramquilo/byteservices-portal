# Código por WhatsApp con Twilio

> **Estado:** En revisión

## Decisión

El código se manda por WhatsApp con la API de mensajes de Twilio, usando la librería oficial
`github.com/twilio/twilio-go`. Adaptador: [twilio.go](/apps/api/adapters/whatsapp/twilio.go)

Sin `TWILIO_ACCOUNT_SID` la API usa el adaptador de [consola](/apps/api/adapters/whatsapp/console.go), que imprime
el código en el log, para desarrollar sin cuenta

## Razones

- **Twilio** ya estaba contratado
- **Librería oficial y no `net/http`.** La llamada es un solo POST a `/Messages.json` y cabía en ~30 líneas con la
  librería estándar; se eligió la oficial por tipos para cada parámetro y errores ya parseados. Costo: la librería se
  genera para todas las APIs de Twilio y es buena parte del tamaño del binario
- **Plantilla aprobada.** El mensaje lo inicia la empresa, así que WhatsApp exige una plantilla aprobada (Content
  Template, `HX…`) con la variable `{{1}}` para el código
- **El código lo genera la API, no Twilio Verify.** Verify generaría y comprobaría el código por su cuenta; se
  mantiene en el núcleo para que el cambio de proveedor solo toque el adaptador `Sender`

## Limitaciones

- `twilio-go` no acepta `context`: un timeout de 10 s en el cliente lo reemplaza
- Si Twilio falla, `POST /api/send` responde `500`; el código ya quedó guardado y la hora gratis ya se dio

## Pendiente de validar

El envío real no se probó con credenciales; sí la conexión TLS desde la imagen `scratch` (Twilio respondió
`401 Authentication Error` con credenciales falsas). Pasos en [configurar-twilio.md](/docs/how-to/configurar-twilio.md)
