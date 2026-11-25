# Estructura del repo

> **Estado:** Aprobado

Monorepo con dos apps: `apps/web` (portal) y `apps/api` (backend Go)

```
byteservices-portal/
├── apps/
│   ├── web/                    portal (HTML/CSS/JS en ES2020, sin dependencias en runtime)
│   │   ├── src/                index.html, byte.css, app.mjs, fonts/ (recortadas)
│   │   ├── fonts-src/          fuentes originales completas
│   │   ├── dist/               compilado por build.mjs (no va en git)
│   │   ├── build.mjs           minificador propio: mete byte.css y app.mjs en un index.html autocontenido
│   │   ├── subset-fonts.mjs    recorta fonts-src/ a src/fonts/ (task web:fonts)
│   │   └── design-system.html  muestrario de tokens y del portal en cada estado
│   └── api/                    backend Go
│       ├── main.go             configuración desde variables de entorno
│       ├── core/               lógica (Send) y puertos (Store, Router)
│       ├── adapters/
│       │   ├── httpapi/        JSON sobre HTTP, portal compilado y redirección cautiva
│       │   ├── sqlite/         Store + migrations/
│       │   └── router/         OPNsense (API REST) y log
│       └── .golangci.yaml      linters y formateadores de Go
├── tools/                      herramientas de Go (golangci-lint, govulncheck, task; actionlint/ aparte)
├── .github/                    workflow de CI y Dependabot
├── .husky/                     hooks de git: pre-commit, commit-msg, pre-push, post-checkout
├── docker/                     Dockerfile y docker-compose.yaml (contexto: la raíz)
├── infra/local/                OPNsense en una VM + AP Byte (tareas `local:`)
├── loadtest/                   k6 contra el portal
├── docs/                       esta documentación (Diátaxis)
├── Taskfile.yml                tareas de todo el repo
├── package.json                herramientas de JS (pnpm): oxlint, oxfmt, es-check, husky, commitlint
├── .oxlintrc.json              reglas de lint de JS
└── .oxfmtrc.json               formato de JS, CSS, JSON, Markdown y YAML
```

## `Taskfile.yml`

`task --list` muestra todas. Las principales:

| Tarea            | Qué hace                                                                                       |
| ---------------- | ---------------------------------------------------------------------------------------------- |
| `task install`   | `pnpm install` (instala los hooks de git) y compila las herramientas de Go                     |
| `task web`       | Compila `apps/web/dist`; falla si pasa de 12 KB gzip                                           |
| `task web:dev`   | Sirve `apps/web/src` en `:8000` (modo demo, sin API)                                           |
| `task web:fonts` | Recorta las fuentes con `subset-font`; volver a correrlo si cambian los textos                 |
| `task api`       | Compila el portal y levanta la API en `:8080`                                                  |
| `task build`     | Construye la imagen                                                                            |
| `task up`        | Levanta el contenedor                                                                          |
| `task down`      | Para el contenedor                                                                             |
| `task logs`      | Sigue los logs                                                                                 |
| `task check`     | Todo lo que valida CI (lo corre el hook pre-push)                                              |
| `task fmt`       | Formatea JS, CSS, JSON, Markdown, YAML, Go y Terraform                                         |
| `task secrets`   | Escanea el historial de git con trufflehog                                                     |
| `task local:up`  | Entorno local con OPNsense en una VM (ver [levantar-local.md](/docs/how-to/levantar-local.md)) |

Detalle de cada chequeo: [ci.md](/docs/explanation/ops/ci.md)

Un solo `.gitignore` en la raíz. El `.dockerignore` deja pasar solo `apps/api/`, `apps/web/build.mjs` y `apps/web/src/`.
