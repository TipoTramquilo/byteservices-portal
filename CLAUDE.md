# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Wi-Fi captive portal for Byte Services: visitor enters nombre + cédula (V/E) + WhatsApp number → receives a 6-digit OTP by WhatsApp → enters it → connected. Static HTML/CSS with a small inline vanilla-JS IIFE; no framework, no package.json, no dependencies, no tests.

`PRODUCT.md` (users, constraints, brand commitments) and `DESIGN.md` (design tokens + rationale, in Spanish) are the source of truth for product and visual decisions. Read them before changing UI.

## Commands

```sh
node build.mjs                      # build dist/ (single self-contained index.html + fonts/)
python3 -m http.server -d portal    # serve source for development (any static server works)
```

There is no lint or test setup. Verify changes by opening the page in a browser (the `.mcp.json` registers `chrome-devtools` MCP for this), ideally at 390×844.

## Layout

- `portal/index.html` — the whole portal: markup for the three steps (`#s-datos`, `#s-codigo`, `#s-listo`) plus the terms sheet, and one inline `<script>` with all behaviour.
- `portal/byte.css` — all styles/tokens. Also linked by `design-system.html`.
- `portal/fonts/` — self-hosted woff2 (Satoshi 700, Plex Sans 400, Plex Mono 400/500). Not yet subsetted.
- `design-system.html` — living specimen of tokens and components; embeds the portal in each state via `?step=`.
- `build.mjs` — hand-rolled minifier: inlines `byte.css` into the HTML (it replaces the exact string `<link rel="stylesheet" href="byte.css">` — keep that tag verbatim), minifies inline `<style>`/`<script>` separately, writes `dist/`, prints raw and gzip size. `dist/` is committed; rebuild after editing `portal/`.
- `Adaptador PAYMUL DIRDEB para Bangente (1)/` — brand reference material (Manual Byte PDF, fonts, logos). `Manual Byte.pdf` (Rev. 2) is the binding brand authority; byteservices.net is not.
- `.claude/skills/impeccable`, `.impeccable/` — the Impeccable design workflow used to build the UI; `.impeccable/surfaces/portal-index-html.md` holds the direction contract and the list of open decisions.

## Backend contract

The portal has no real backend yet. It reads `window.BYTE_PORTAL` (set before the inline script):

- `send` — `POST {nombre, cedula:"V18456789", telefono:"+584125551234"}` → 2xx, optional `{verificado:true, redirect}` when the number was already verified on an earlier day (skips the code step, goes straight to "conectado"). Also used for resend. The network grants 1h of internet on `send` so the WhatsApp OTP can arrive; verifying extends it to the end of the day.
- `verify` — `POST {telefono, codigo}` → 2xx, optional `{redirect}` used for the "Empezar a navegar" link.

If `send` is unset the portal runs in **demo mode** (fake 700 ms latency, code `123456`, a "modo demo" badge).

Preview states for screenshots/QA: `portal/index.html?step=codigo|error|listo|datos-error` prefills sample data and jumps to that state.

## Hard constraints

- Runs inside the OS captive mini-browser (iOS CNA, Android WebView) **before** the device has internet: no CDNs, external fonts, images, or third-party requests. Everything inline or self-hosted in `portal/`.
- Weight matters ("lo más ligero posible"): every byte must justify itself; check the size printed by `node build.mjs`.
- Must work without the custom fonts (system fallbacks are in the font stacks) and on low-end phones; keep the JS ES5-style to match the existing code.
- Brand: only the palette/typography in `PRODUCT.md`/`DESIGN.md` (Negro Byte `#0A0F1C`, never pure black; Satoshi for titles, Plex Mono for UI/buttons, Plex Sans for text). Fields are underlines, not boxes. The gradient isotipo SVGs in the brand folder embed MB-size PNGs — don't use them.
- Don't invent terms-of-use/legal copy; it hasn't been provided.

Still undecided (don't assume): hosting (router MikroTik/OpenWrt vs external server), real API endpoints, session duration, resend cooldown (currently `WAIT=60`s in the script), font subsetting.
