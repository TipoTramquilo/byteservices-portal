# Documentación del portal cautivo Byte

Organizada según [Diátaxis](https://diataxis.fr/) en cuatro tipos de documento, más los pendientes sin decisión
cerrada

## Explicación (por qué se decidió)

- **Core**
  - [core/acceso-wifi.md](/docs/explanation/core/acceso-wifi.md) - registro del visitante y sesión de 30 minutos, sin código
- **Arquitectura**
  - [architecture/puertos-y-adaptadores.md](/docs/explanation/architecture/puertos-y-adaptadores.md) - `core/` con sus puertos y un paquete por adaptador
  - [architecture/api-rest.md](/docs/explanation/architecture/api-rest.md) - Go, REST con JSON, mismo origen que el portal y redirección cautiva
- **Base de datos**
  - [database/sqlite.md](/docs/explanation/database/sqlite.md) - SQLite en proceso, sin ORM, tablas `STRICT` y migraciones forward-only
- **Portal**
  - [web/es2020.md](/docs/explanation/web/es2020.md) - JavaScript en ES2020: teléfonos objetivo, qué queda fuera y cómo se valida
- **Operación**
  - [ops/contenedores.md](/docs/explanation/ops/contenedores.md) - imagen sobre `scratch` con Podman
  - [ops/ci.md](/docs/explanation/ops/ci.md) - hooks de git, CI y Dependabot: qué valida cada uno y por qué

## Referencia (datos exactos, para consultar)

- **API**
  - [api/openapi.yaml](/docs/reference/api/openapi.yaml) - contrato OpenAPI 3.2.0
  - [api/contrato-portal.md](/docs/reference/api/contrato-portal.md) - `window.BYTE_PORTAL` y cómo reacciona el portal
- **Base de datos**
  - [database/esquema.md](/docs/reference/database/esquema.md) - tablas, columnas y reglas
- **Setup**
  - [setup/estructura.md](/docs/reference/setup/estructura.md) - árbol del monorepo y tareas del `Taskfile.yml`
  - [setup/variables-de-entorno.md](/docs/reference/setup/variables-de-entorno.md) - variables que lee la API
  - [setup/opnsense.md](/docs/reference/setup/opnsense.md) - configuración del equipo: zona, DNS, API key y arranque de la API

## How-to (pasos para lograr algo)

- [correr-el-proyecto.md](/docs/how-to/correr-el-proyecto.md) - API en el host o en contenedor, y compilar para OPNsense
- [levantar-local.md](/docs/how-to/levantar-local.md) - OPNsense en una VM y el AP `Byte` en el PC, para probar en un teléfono
- [diagnosticar-portal.md](/docs/how-to/diagnosticar-portal.md) - errores reales (500, 403, TLS, Samsung, DNS y shell) con su causa y solución
- [agregar-migracion.md](/docs/how-to/agregar-migracion.md) - agregar una migración de SQLite
- [actualizar-openapi.md](/docs/how-to/actualizar-openapi.md) - cambiar el contrato OpenAPI

## Tutorial (aprender haciendo)

- [recorrido-flujo-completo.md](/docs/tutorials/recorrido-flujo-completo.md) - seguir una visita del portal a la base, en local

## Pendiente (sin decisión cerrada)

- [ops/router.md](/docs/pending/ops/router.md) - estado del equipo real y lo que falta: arranque tras reinicio y nombre de la red
- [ops/estado-fwportal.html](/docs/pending/ops/estado-fwportal.html) - resumen visual del estado de `fwportal` y sus pendientes (se abre en el navegador)
- [web/portal.md](/docs/pending/web/portal.md) - términos de uso, nombre de la red y textos

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
- **Referencias con ruta absoluta** y líneas cuando es código
- **Tono impersonal**, en español, sin punto final en listas y descripciones
- **Nombres reales:** lo citado existe en el repo o se marca como planificado
