# Documentación del portal cautivo Byte

Organizada según [Diátaxis](https://diataxis.fr/) en cuatro tipos de documento, más los pendientes sin decisión
cerrada. Producto y diseño viven aparte: [PRODUCT.md](/PRODUCT.md) y [DESIGN.md](/DESIGN.md)

## Explicación (por qué se decidió)

- **Core**
  - [core/acceso-wifi.md](/docs/explanation/core/acceso-wifi.md) - registro con OTP por WhatsApp, límites de envío y reset diario (2:00 AM)
- **Arquitectura**
  - [architecture/puertos-y-adaptadores.md](/docs/explanation/architecture/puertos-y-adaptadores.md) - `core/` con sus puertos y un paquete por adaptador
  - [architecture/api-rest.md](/docs/explanation/architecture/api-rest.md) - Go, REST con JSON y mismo origen que el portal
- **Base de datos**
  - [database/postgres.md](/docs/explanation/database/postgres.md) - Postgres en Supabase, sin ORM, migraciones forward-only
- **Integraciones**
  - [integrations/twilio.md](/docs/explanation/integrations/twilio.md) - código por WhatsApp con la librería oficial de Twilio
- **Operación**
  - [ops/contenedores.md](/docs/explanation/ops/contenedores.md) - imagen sobre `scratch` con Podman y dónde usarla

## Referencia (datos exactos, para consultar)

- **API**
  - [api/openapi.yaml](/docs/reference/api/openapi.yaml) - contrato OpenAPI 3.2.0
  - [api/contrato-portal.md](/docs/reference/api/contrato-portal.md) - `window.BYTE_PORTAL`, cómo se identifica al dispositivo y cómo reacciona el portal
- **Base de datos**
  - [database/esquema.md](/docs/reference/database/esquema.md) - tablas, columnas y reglas
- **Setup**
  - [setup/estructura.md](/docs/reference/setup/estructura.md) - árbol del monorepo y comandos del `Makefile`
  - [setup/variables-de-entorno.md](/docs/reference/setup/variables-de-entorno.md) - variables que lee la API

## How-to (pasos para lograr algo)

- [correr-el-proyecto.md](/docs/how-to/correr-el-proyecto.md) - API en el host, en contenedor o solo el portal
- [agregar-migracion.md](/docs/how-to/agregar-migracion.md) - agregar una migración de Postgres
- [actualizar-openapi.md](/docs/how-to/actualizar-openapi.md) - cambiar el contrato OpenAPI
- [configurar-twilio.md](/docs/how-to/configurar-twilio.md) - mandar el código por WhatsApp de verdad

## Tutorial (aprender haciendo)

- [recorrido-flujo-completo.md](/docs/tutorials/recorrido-flujo-completo.md) - seguir una visita del portal a la base, en local

## Pendiente (sin decisión cerrada)

- [web/portal.md](/docs/pending/web/portal.md) - términos de uso y textos

## Estados

Cada documento lleva una línea `> **Estado:** ...` debajo del título (exactamente una):

- **Pendiente:** no revisado contra código. Funcionalidad planificada o en veremos
- **En revisión:** escrito o corregido, esperando revisión o validación
- **Aprobado:** revisado, vigente, alineado con el código actual
- **Outdated:** el contexto o el código cambió; conservado como referencia, con una nota `> **Nota:** ...`
- **Deprecated:** superado; incluye enlace al reemplazo

## Convenciones de escritura

- **Un tipo por documento (Diátaxis):** `explanation/` el porqué, `reference/` datos exactos, `how-to/` pasos,
  `tutorials/` aprender haciendo. Lo planificado sin implementar va en `pending/<área>/`
- **Estado actual, sin históricos.** Describe cómo está hoy el sistema
- **Decisiones con motivo** y, cuando aporta, la alternativa descartada
- **Diagramas en Mermaid**; árboles de carpetas en ASCII
- **Referencias con ruta absoluta** y líneas cuando es código: `[service.go](/apps/api/core/service.go) (líneas 58-118)`
- **Tono impersonal**, en español, sin punto final en listas y descripciones
- **Nombres reales:** lo citado existe en el repo o se marca como planificado
