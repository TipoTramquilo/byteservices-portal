---
version: 1
slug: "portal-index-html"
primary_target: "apps/web/src/index.html"
related_targets: ["apps/web/design-system.html"]
---

# Portal cautivo Byte

Scope: the captive-portal flow (datos → código WhatsApp → conectado) plus the design-system specimen that documents its tokens and components. Visitor mode: Operate.

Audience and job: a guest on the Byte Wi-Fi, on their own phone inside the OS captive mini-browser, one-handed, offline until verified. Task: give nombre, cédula (V/E) and WhatsApp, type the 6-digit code, get online. Constraints: as light as possible, no third-party requests, works without JS-heavy features.

## Direction contract

THESIS: the number you type is the interface. Refuses the category default of a centred card with a logo and three stacked inputs; data is written on ruled lines and the WhatsApp number is set at poster scale.

OWN-WORLD: Manual Byte Rev. 2. Negro Byte #0A0F1C ground, lines #3A4358, Azul Byte #309AE7 action, azul→violeta (#5036EA) pixel rule, Satoshi Bold display, IBM Plex Mono for labels and buttons, IBM Plex Sans for text, the byte square pattern clustered in one corner. Fields are underlines, not boxes; buttons 10px radius.

STORY: the visitor sees who is asking (Byte lockup, network pill), understands why (code by WhatsApp), writes name and cédula, types the number which grows huge and fills the pixel rule, sends, types the code in six monumental cells, sees "Conectado".

FIRST VIEWPORT: 390×844. Top bar: lockup left, "Byte Invitados" pill right. Satoshi 30px question. Name + cédula on one ruled row. WhatsApp: +58 selector beside the number at 46px Satoshi, pixel rule below filling per digit. Primary button pinned to the bottom, legal line under it. Signature move: the pixel rule that fills azul→violeta as digits arrive (and counts down the resend timer on the code step).

FORM: "Número monumental", position 2 of my ordered list; seed:9c14c58e214d43c7.

FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance

## Decisions open
- Hosting (router vs server) and backend endpoints for send/verify: the build exposes hooks, no real API.
- Font subsetting (fonttools not installed; user not yet asked to approve install).
- Terms-of-use text.
