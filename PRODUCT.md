# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Static HTML/CSS with minimal inline JS, no framework. Derived from the user's requirement "lo más ligero posible"; the hosting target (router vs external server) is still open, so every page is designed for the strictest case: single file, inline CSS/SVG, no external fonts or requests.

## Users

People connecting to a Byte Wi-Fi network (offices, events, visitors) from their own phone. They are standing somewhere, often one-handed, with no internet yet, and want to get online in seconds.

## Product Purpose

Captive portal that grants Wi-Fi access after verifying a phone number with a one-time code sent by WhatsApp. Success: the visitor goes from opening the portal to connected in under 30 seconds without confusion.

## Positioning

Built and operated by Byte Services, a Venezuelan cybersecurity / Zero Trust infrastructure company. The portal itself is a small demonstration of that: verified, accountable access instead of an open network or a shared password.

## Operating Context

- Rendered inside the OS captive-portal mini browser (iOS CNA, Android captive portal WebView), which has limited features and no internet before authentication.
- Flow: full name + cédula (V/E prefix) + phone number → code sent by WhatsApp → enter 6-digit code → connected.
- Venezuelan numbers are the default (+58), but other country codes must be possible.

## Capabilities and Constraints

- OTP channel: WhatsApp only.
- Must be as light as possible; no third-party assets (fonts, CDNs, images) because they are unreachable before authentication.
- Hosting undecided: router (MikroTik/OpenWrt) or external server/controller.
- Undecided: terms-of-use text, session duration, resend cooldown values.

## Brand Commitments

Byte brand only (no white-label). Binding source: `brand/Manual Byte.pdf` (Manual de Identidad Rev. 2, Nov 2026), with the PAYMUL deck as the reference for how it looks on screen. byteservices.net is older and is not the authority where they differ.

- Type: Satoshi Bold for titles; IBM Plex Mono for buttons, UI and short text; IBM Plex Sans for long text.
- Colour: Azul Byte #309AE7 (primary), Violeta #5036EA, Azul Profundo #1A6FB8, Celeste #5FD6FF, Negro Byte #0A0F1C (never pure black). Greys #161D2E, #3A4358, #8A93A8, #D3D8E3, Blanco lila #EDEAF5, white for text only. States: success #22C55E (alt: Azul Byte), warning #F5A523, error #EF4444.
- Elements: the "byte" square pattern (azul→violeta), slightly softened corners, high contrast; no logo distortion, glows or low contrast.
- Logo: isotipo + divider + "Byte" wordmark (vectors in `Recursos/Iconotipos_Isotipos/Group 21.svg` and `Recursos/Logotipos/Byte.svg`).

## Evidence on Hand

Brand assets in `brand/Recursos` (fonts Satoshi / IBM Plex, vector logos). The gradient isotipo SVGs embed 1–2 MB PNGs and are unusable in the portal. No legal/terms copy provided; do not invent it.

## Product Principles

1. Connected fast: one field per step, nothing in the way.
2. Weightless: every byte must justify itself; system fonts and inline assets.
3. Trustworthy: it should read as a security company's front door, clear about why the number is asked.
4. Works in the worst mini browser, on the cheapest phone.

## Accessibility & Inclusion

Large tap targets, numeric keypads (`inputmode="numeric"`, `autocomplete="one-time-code"`), WCAG AA contrast, works without custom fonts.
