# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Wi-Fi captive portal for Byte Services: visitor enters nombre + cédula (V/E) + WhatsApp number → receives a 6-digit OTP by WhatsApp → enters it → connected. Monorepo with two apps: `apps/web` (static HTML/CSS with a small inline vanilla-JS IIFE; no framework, no package.json, no dependencies) and `apps/api` (Go). No test setup.

`PRODUCT.md` (users, constraints, brand commitments) and `DESIGN.md` (design tokens + rationale, in Spanish) are the source of truth for product and visual decisions. Read them before changing UI.

`docs/` (Diátaxis, in Spanish; conventions in `docs/README.md`) holds the decisions (`explanation/`), exact data (`reference/`, including the OpenAPI contract at `docs/reference/api/openapi.yaml`), how-tos, a tutorial and open items (`pending/`). Update the matching doc and the OpenAPI when behaviour or the API changes.

## Commands

Run from the repo root (`make help` lists them):

```sh
make web       # build apps/web/dist/ (single self-contained index.html + fonts/)
make web-dev   # serve apps/web/src/ for UI work (demo mode, no API)
make fonts     # re-subset apps/web/fonts-src into src/fonts (after changing portal copy)
make api       # API + built portal on :8080 (codes printed to the log without TWILIO_ACCOUNT_SID)
make up        # same, containerized (scratch image, SQLite in the portal_data volume); make down / make logs
```

There is no lint or test setup. Verify changes by opening the page in a browser (the `.mcp.json` registers `chrome-devtools` MCP for this), ideally at 390×844.

## Layout

Web (`apps/web/`):
- `src/index.html` — the portal markup: the three steps (`#s-datos`, `#s-codigo`, `#s-listo`) plus the terms sheet.
- `src/app.js` — all behaviour, one ES5 IIFE (includes the `window.BYTE_PORTAL` contract comment).
- `src/byte.css` — all styles/tokens. Also linked by `design-system.html`.
- `src/fonts/` — self-hosted woff2 (Satoshi 700, Plex Sans 400, Plex Mono 400/500), subsetted from `fonts-src/` by `subset-fonts.mjs` (`make fonts`, which fetches `subset-font` with npx; nothing is added to the repo). Plex Sans keeps full Latin-1 because it renders the name field; Satoshi and Plex Mono keep only the characters the portal shows. Rerun `make fonts` after changing portal copy, or new characters fall back to system fonts.
- `design-system.html` — living specimen of tokens and components; embeds the portal in each state via `?step=`.
- `build.mjs` — hand-rolled minifier: inlines `byte.css` and `app.js` into the HTML (it replaces the exact strings `<link rel="stylesheet" href="byte.css">` and `<script src="app.js"></script>` — keep both tags verbatim), minifies inline `<style>`/`<script>` separately, writes `dist/`, prints raw and gzip size. `dist/` is not committed (gitignored): the Dockerfile builds it, and `make web` rebuilds it for local runs (`make api` serves it).

API (`apps/api/`):
- Go API (`/api/send`, `/api/verify`, `/api/health`) that also serves `apps/web/dist/`. Ports and adapters: `core/` holds the logic and its ports (`Store`, `Sender`, `Router`); `adapters/` has `httpapi`, `sqlite` (plain SQL; numbered migrations in `migrations/`, embedded and applied on startup via `PRAGMA user_version`; never edit a published one, add the next), `whatsapp` (Twilio or console) and `router` (log only until the router is chosen). Dependencies: `modernc.org/sqlite`, `twilio-go`. Identifiers in English, comments in Spanish.

Root:
- `docs/` — documentation (Diátaxis); `docs/reference/api/openapi.yaml` is the API contract.
- `docker/` — `Dockerfile` (build context is the repo root: stage `node` compiles the portal with `build.mjs`, stage `golang` builds the API, runtime `scratch`) and `docker-compose.yaml`.
- `brand/` — brand reference material (Manual Byte PDF, fonts, logos). `Manual Byte.pdf` (Rev. 2) is the binding brand authority; byteservices.net is not.
- `.claude/skills/impeccable`, `.impeccable/` — the Impeccable design workflow used to build the UI; `.impeccable/surfaces/portal-index-html.md` holds the direction contract and the list of open decisions.

## Backend contract

The portal reads `window.BYTE_PORTAL` (set before the inline script; `apps/api` injects it when serving the page):

- `send` — `POST {name, id_number:"V18456789", phone:"+584125551234"}` → 2xx, optional `{verified:true, redirect}` when the number was already verified on an earlier day (skips the code step, goes straight to "conectado"). Also used for resend. The network grants 1h of internet on `send` so the WhatsApp OTP can arrive (once a day per IP, MAC, phone and cédula: any of them used blocks it); verifying extends it until the daily reset at `RESET_AT` (02:00 local; all times are stored in UTC). The device is identified by the request's IP; the backend asks the `Router` adapter for its MAC (the page never sends it).
- `verify` — `POST {phone, code}` → 2xx, optional `{redirect}` used for the "Empezar a navegar" link.

If `send` is unset the portal runs in **demo mode** (fake 700 ms latency, code `123456`, a "modo demo" badge).

Preview states for screenshots/QA: `apps/web/src/index.html?step=codigo|error|listo|datos-error` prefills sample data and jumps to that state.

## Hard constraints

- Runs inside the OS captive mini-browser (iOS CNA, Android WebView) **before** the device has internet: no CDNs, external fonts, images, or third-party requests. Everything inline or self-hosted in `apps/web/src/`.
- Weight matters ("lo más ligero posible"): every byte must justify itself; check the size printed by `make web`.
- Must work without the custom fonts (system fallbacks are in the font stacks) and on low-end phones; keep the JS ES5-style to match the existing code.
- Brand: only the palette/typography in `PRODUCT.md`/`DESIGN.md` (Negro Byte `#0A0F1C`, never pure black; Satoshi for titles, Plex Mono for UI/buttons, Plex Sans for text). Fields are underlines, not boxes. The gradient isotipo SVGs in the brand folder embed MB-size PNGs — don't use them.
- Don't invent terms-of-use/legal copy; it hasn't been provided.

Still undecided (don't assume): hardware (mini PC with OPNsense, likely, vs Raspberry), which `Router` adapter to write, resend cooldown (currently `WAIT=60`s in the script).
