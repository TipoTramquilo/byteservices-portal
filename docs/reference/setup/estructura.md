# Estructura del repo

> **Estado:** Aprobado

Monorepo con dos apps: `apps/web` (portal) y `apps/api` (backend Go)

```
cautive-portal/
├── apps/
│   ├── web/                    portal (HTML/CSS/JS, sin dependencias)
│   │   ├── src/                index.html, byte.css, app.js, fonts/ (recortadas)
│   │   ├── fonts-src/          fuentes originales completas
│   │   ├── dist/               compilado por build.mjs (no va en git)
│   │   ├── build.mjs           minificador propio: mete byte.css y app.js en un index.html autocontenido
│   │   ├── subset-fonts.mjs    recorta fonts-src/ a src/fonts/ (make fonts)
│   │   └── design-system.html  muestrario de tokens y del portal en cada estado
│   └── api/                    backend Go
│       ├── main.go             configuración desde variables de entorno
│       ├── core/               lógica (Send) y puertos (Store, Router)
│       └── adapters/
│           ├── httpapi/        JSON sobre HTTP, portal compilado y redirección cautiva
│           ├── sqlite/         Store + migrations/
│           └── router/         OPNsense (API REST) y log
├── docker/                     Dockerfile y docker-compose.yaml (contexto: la raíz)
├── docs/                       esta documentación (Diátaxis)
├── brand/                      Manual Byte, logos y fuentes de referencia
├── Makefile                    comandos de todo el repo
├── PRODUCT.md                  usuarios, restricciones y compromisos de marca
├── DESIGN.md                   tokens de diseño y su razón
└── CLAUDE.md
```

## `Makefile`

| Comando        | Qué hace                                                    |
| -------------- | ----------------------------------------------------------- |
| `make web`     | Compila `apps/web/dist`                                     |
| `make web-dev` | Sirve `apps/web/src` en `:8000` (modo demo, sin API)        |
| `make fonts`   | Recorta las fuentes con `subset-font` vía `npx`; volver a correrlo si cambian los textos |
| `make api`     | Compila el portal y levanta la API en `:8080`               |
| `make build`   | Construye la imagen                                         |
| `make up`      | Levanta el contenedor                                       |
| `make down`    | Para el contenedor                                          |
| `make logs`    | Sigue los logs                                              |

Un solo `.gitignore` en la raíz. El `.dockerignore` deja pasar solo `apps/api/`, `apps/web/build.mjs` y `apps/web/src/`.
