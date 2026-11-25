---
name: Portal cautivo Byte
description: Acceso Wi-Fi verificado por WhatsApp, escrito en líneas sobre Negro Byte.
colors:
  azul-byte: '#309AE7'
  azul-hover: '#4AA8EC'
  violeta: '#5036EA'
  azul-profundo: '#1A6FB8'
  celeste: '#5FD6FF'
  negro-byte: '#0A0F1C'
  panel: '#0E1424'
  gris-900: '#161D2E'
  linea-panel: '#232B40'
  gris-700: '#3A4358'
  gris-400: '#8A93A8'
  gris-200: '#D3D8E3'
  blanco-lila: '#EDEAF5'
  blanco: '#FFFFFF'
  exito: '#22C55E'
  aviso: '#F5A523'
  error: '#EF4444'
  error-texto: '#F87171'
typography:
  display:
    fontFamily: 'Satoshi, Avenir Next, Segoe UI, system-ui, sans-serif'
    fontSize: 'clamp(30px, 10.4vw, 44px)'
    fontWeight: 700
    lineHeight: 1.1
    letterSpacing: '-0.02em'
    fontFeature: 'tnum'
  display-otp:
    fontFamily: 'Satoshi, Avenir Next, Segoe UI, system-ui, sans-serif'
    fontSize: '40px'
    fontWeight: 700
    lineHeight: 1
  headline:
    fontFamily: 'Satoshi, Avenir Next, Segoe UI, system-ui, sans-serif'
    fontSize: '30px'
    fontWeight: 700
    lineHeight: 1.12
    letterSpacing: '-0.02em'
  title:
    fontFamily: 'Satoshi, Avenir Next, Segoe UI, system-ui, sans-serif'
    fontSize: '22px'
    fontWeight: 700
    lineHeight: 1.2
  body:
    fontFamily: 'Plex Sans, system-ui, -apple-system, Segoe UI, Roboto, sans-serif'
    fontSize: '15px'
    fontWeight: 400
    lineHeight: 1.5
  field:
    fontFamily: 'Plex Sans, system-ui, -apple-system, Segoe UI, Roboto, sans-serif'
    fontSize: '18px'
    fontWeight: 400
    lineHeight: 1.3
  label:
    fontFamily: 'Plex Mono, ui-monospace, SF Mono, Menlo, Consolas, monospace'
    fontSize: '12.5px'
    fontWeight: 400
    lineHeight: 1
  button:
    fontFamily: 'Plex Mono, ui-monospace, SF Mono, Menlo, Consolas, monospace'
    fontSize: '16px'
    fontWeight: 500
    lineHeight: 1
  small:
    fontFamily: 'Plex Sans, system-ui, -apple-system, Segoe UI, Roboto, sans-serif'
    fontSize: '12px'
    fontWeight: 400
    lineHeight: 1.5
rounded:
  px: '1.5px'
  focus: '4px'
  seg-item: '5px'
  seg: '7px'
  btn: '10px'
  panel: '14px'
  pill: '999px'
spacing:
  s2: '8px'
  s3: '12px'
  s4: '16px'
  s5: '24px'
  s6: '32px'
  gutter: '24px'
components:
  button-primary:
    backgroundColor: '{colors.azul-byte}'
    textColor: '{colors.negro-byte}'
    typography: '{typography.button}'
    rounded: '{rounded.btn}'
    height: '56px'
    width: '100%'
  button-primary-hover:
    backgroundColor: '#4AA8EC'
  button-primary-disabled:
    backgroundColor: '{colors.gris-900}'
    textColor: '{colors.gris-400}'
  link:
    textColor: '{colors.azul-byte}'
    typography: '{typography.label}'
  link-disabled:
    textColor: '{colors.gris-400}'
  field-line:
    textColor: '{colors.blanco}'
    typography: '{typography.field}'
    padding: '6px 0 10px'
  field-label:
    textColor: '{colors.gris-400}'
    typography: '{typography.label}'
  field-error:
    textColor: '{colors.error-texto}'
  network-pill:
    textColor: '{colors.gris-200}'
    typography: '{typography.label}'
    rounded: '{rounded.pill}'
    padding: '6px 12px'
  segmented-option-checked:
    backgroundColor: '{colors.gris-700}'
    textColor: '{colors.blanco}'
  otp-cell:
    textColor: '{colors.blanco}'
    typography: '{typography.display-otp}'
    height: '68px'
  terms-sheet:
    backgroundColor: '{colors.panel}'
    rounded: '{rounded.panel}'
    padding: '24px 24px 28px'
---

# Design System: Portal cautivo Byte

## Overview

**Creative North Star: "El número es la interfaz"**

El portal es la puerta de entrada de una empresa de seguridad, y se comporta como tal: un fondo Negro Byte casi absoluto, una sola acción en Azul Byte y los datos del visitante escritos sobre líneas, no encerrados en cajas. Lo que el visitante teclea es lo más grande de la pantalla. El número de WhatsApp se compone a escala de cartel en Satoshi Bold y, debajo, una regla de píxeles se llena de azul a violeta con cada dígito; en el paso del código, la misma regla se vacía como temporizador de reenvío.

La densidad es baja y vertical: una columna de 440px como máximo, una pregunta por paso, el botón anclado al pie bajo el pulgar. El mundo sale del Manual de Identidad Byte Rev. 2 y de su aplicación en pantalla en el deck PAYMUL (de ahí el panel #0E1424 y la línea #232B40). El patrón de cuadrados "byte" aparece agrupado en una sola esquina, nunca repartido como textura.

La superficie es plana: profundidad por escalón tonal y filetes de 1px, sin sombras ni brillos, con contraste alto en todo el texto. El manual prohíbe distorsionar el logo, los brillos y el bajo contraste; el portal lo respeta.

**Key Characteristics:**

- Fondo Negro Byte (#0A0F1C), nunca negro puro.
- Campos como líneas de 1px; el dato escrito en blanco sobre la regla.
- Número de teléfono y código OTP en Satoshi Bold a tamaño monumental, cifras tabulares.
- Regla de píxeles azul→violeta como firma: progreso al escribir, cuenta atrás al esperar.
- Tres voces tipográficas con funciones fijas: Satoshi (títulos y números), Plex Mono (etiquetas, botones, datos), Plex Sans (texto).
- Azul Byte como única acción; plano, sin sombras.

## Colors

Una noche azul-tinta casi monocroma, con un único azul eléctrico que actúa y un violeta que solo cierra la transición de los píxeles.

### Primary

- **Azul Byte** (#309AE7): la acción, el foco activo de un campo, el subrayado activo de la celda OTP, la palabra destacada del titular y el cursor de texto. Fondo del botón primario con texto Negro Byte encima.

### Secondary

- **Violeta** (#5036EA): extremo final de la regla de píxeles (cada píxel interpola de Azul Byte a Violeta según su posición) y algunos cuadrados del patrón. Nunca es una acción ni un fondo.
- **Azul Profundo** (#1A6FB8): solo dentro del patrón byte de esquina.

### Tertiary

- **Celeste** (#5FD6FF): anillo de foco por teclado (2px, desplazado 3px), punto del pill de red y cuadrados del patrón.

### Neutral

- **Negro Byte** (#0A0F1C): fondo de todas las pantallas; texto del botón primario.
- **Panel** (#0E1424): superficie elevada un escalón (hoja de términos, control V/E, tarjetas del espécimen).
- **Gris 900** (#161D2E): fondo del botón deshabilitado.
- **Línea de panel** (#232B40): filetes sobre panel, borde del pill, píxeles apagados de la regla.
- **Gris 700** (#3A4358): la línea de los campos en reposo, separador del lockup, filetes de la lista de confirmación, opción V/E seleccionada.
- **Gris 400** (#8A93A8): etiquetas, texto secundario, placeholders, texto deshabilitado.
- **Gris 200** (#D3D8E3): valores en lecturas de datos, texto del pill, enlaces del texto legal.
- **Blanco lila** (#EDEAF5): texto corrido.
- **Blanco** (#FFFFFF): titulares, lockup y lo que el visitante escribe.

### Estados

- **Éxito** (#22C55E): punto del pill cuando la red está conectada.
- **Aviso** (#F5A523): estado del manual, reservado; el portal aún no lo usa.
- **Error** (#EF4444): línea del campo, píxeles de la regla y subrayado OTP en error.
- **Error texto** (#F87171): mensajes de error en texto pequeño sobre Negro Byte; existe para que el texto de error cumpla AA, el rojo del manual se queda en líneas y píxeles.

### Named Rules

**La regla de la única acción.** Azul Byte es lo único que se pulsa o que indica foco activo. Ningún otro color de la marca se usa como botón o enlace.

**La regla del negro que no es negro.** El fondo es Negro Byte (#0A0F1C). Nunca #000.

**La regla del blanco escrito.** El blanco puro es para titulares, el logo y los datos que teclea el visitante; el texto corrido va en Blanco lila.

## Typography

**Display Font:** Satoshi Bold (con Avenir Next, Segoe UI, system-ui)
**Body Font:** IBM Plex Sans Regular (con system-ui, -apple-system, Segoe UI, Roboto)
**Label/Mono Font:** IBM Plex Mono Regular y Medium (con ui-monospace, SF Mono, Menlo, Consolas)

**Character:** Satoshi da la voz grande y geométrica de la marca; Plex Mono convierte etiquetas, botones y datos en lecturas de terminal, propias de una empresa de seguridad; Plex Sans queda neutral para lo que se lee de corrido. Solo se cargan esos cuatro pesos, como subconjuntos WOFF2 locales (latín básico, Latin-1, guiones, comillas, puntos suspensivos y flechas): `satoshi-700.woff2` (17.1 KB), `plex-mono-400.woff2` (13.6 KB), `plex-mono-500.woff2` (14.0 KB) y `plex-sans-400.woff2` (21.0 KB), 66 KB en total en `portal/fonts/`, con `font-display: swap`. Satoshi se precarga porque compone el titular y el número. Los originales completos quedan en `Recursos` del manual. Las pilas de respaldo mantienen el portal legible sin fuentes propias.

### Hierarchy

- **Display** (Satoshi 700, clamp(30px, 10.4vw, 44px), 1.1, -0.02em, cifras tabulares): el número de WhatsApp. Si no cabe, baja en pasos de 2px hasta 22px.
- **Display OTP** (Satoshi 700, 40px, 1): los seis dígitos del código.
- **Headline** (Satoshi 700, 30px, 1.12, -0.02em, text-wrap balance): la pregunta de cada paso. Una palabra puede ir en Azul Byte (sin cursiva) para nombrar el canal o el resultado.
- **Title** (Satoshi 700, 22px, 1.2): títulos de hojas secundarias (términos).
- **Body** (Plex Sans 400, 15px, 1.5): texto de apoyo, máximo unos 34ch bajo el titular.
- **Field** (Plex Sans 400, 18px, 1.3): nombre y cédula escritos sobre la línea.
- **Label** (Plex Mono 400, 12.5px, 1): etiquetas de campo, temporizador, filas de confirmación; Plex Mono 500 para el pill y los enlaces.
- **Button** (Plex Mono 500, 16px, 1): botón primario, en caja de oración.
- **Small** (12px, 1.4–1.5): texto legal en Plex Sans; mensajes de error en Plex Mono.

### Named Rules

**La regla de las tres voces.** Satoshi solo en titulares y números; Plex Mono en todo lo que es interfaz o dato (etiquetas, botones, pill, valores); Plex Sans en el texto que se lee. No se intercambian.

**La regla del número monumental.** El dato principal de cada paso (teléfono, código) se compone en Satoshi Bold a escala display, en blanco. Es el elemento más grande de la pantalla.

## Layout

Una sola columna centrada de 440px máximo, con margen lateral de 24px y respeto de las zonas seguras (`env(safe-area-inset-*)`). Cada paso es una columna flexible: barra superior (lockup a la izquierda, pill de red a la derecha), titular a clamp(28px, 7vh, 52px) del borde superior, campos separados 32px, y el bloque del pie empujado al fondo con `margin-top: auto`, de modo que el botón queda bajo el pulgar con el texto legal debajo.

Nombre y cédula comparten una fila en rejilla de dos columnas (1fr y 1.08fr, separación 16px). La escala de espaciado es 8 / 12 / 16 / 24 / 32px.

Desde 700px el contenido se centra verticalmente y cada paso tiene una altura mínima de 720px. Por debajo de 1000px el patrón de esquina vive dentro de la columna; desde 1000px sale de ella y se fija en la esquina inferior derecha de la ventana (240px, a 40px de los bordes).

## Elevation & Depth

Sistema plano. No hay sombras ni brillos en ningún elemento; la profundidad se expresa con un escalón tonal (Negro Byte → Panel) y filetes de 1px en Línea de panel o Gris 700. La única superficie superpuesta, la hoja de términos, sube desde el pie con fondo Panel, filete superior y esquinas superiores de 14px.

### Named Rules

**La regla plana.** Nada proyecta sombra ni resplandor. Si algo necesita separarse, sube un escalón tonal y gana un filete de 1px.

## Shapes

Esquinas ligeramente suavizadas, como pide el manual: 10px en botones, 14px en paneles, píldora completa en el indicador de red, y 1.5px en los cuadrados del patrón y el punto del pill (1px en los píxeles de la regla). El control V/E usa 7px por fuera y 5px en la opción.

La forma recurrente es el cuadrado pequeño: píxeles de 6px en la regla, cuadrados de 13px con 2px de separación en el patrón, punto de 7px en el pill. Los campos no tienen forma propia: son una línea horizontal.

### Named Rules

**La regla de la línea, no la caja.** Un campo de texto es una línea de 1px bajo el dato (2px en las celdas OTP). No se dibujan cajas alrededor de campos.

## Components

### Buttons

Firme y anclado, a lo ancho de la columna.

- **Shape:** esquinas suavizadas (10px), 56px de alto, ancho completo.
- **Primary:** fondo Azul Byte, texto Negro Byte en Plex Mono 500 16px.
- **Hover / Active:** el fondo aclara a #4AA8EC; al pulsar escala a 0.985. Transiciones de 140ms con la curva de salida `cubic-bezier(.16,1,.3,1)`.
- **Disabled:** fondo Gris 900, texto Gris 400, sin transformación.
- **Ocupado:** el texto se vuelve transparente y aparece un giro de 18px con trazo Negro Byte de 2px.
- **Enlace:** Plex Mono 500 12.5px en Azul Byte, subrayado con desplazamiento de 3px; deshabilitado pasa a Gris 400 sin subrayado. "Reenviar código" permanece oculto mientras corre la cuenta atrás y aparece cuando la regla se vacía.

### Chips

- **Pill de red:** borde de 1px en Línea de panel, forma de píldora, Plex Mono 500 12.5px en Gris 200, con un punto cuadrado de 7px en Celeste (en Éxito cuando la red está conectada).

### Cards / Containers

- **Hoja de términos:** hoja inferior que se abre por `:target` (`#terminos`), sin JS; fondo Panel, filete superior en Línea de panel, esquinas superiores de 14px, relleno 24px 24px 28px. Sin sombra.
- **Lista de confirmación:** filas separadas por filetes Gris 700, etiqueta en Plex Mono Gris 400 a la izquierda y valor en Gris 200 a la derecha.

### Inputs / Fields

- **Style:** línea inferior de 1px en Gris 700, sin fondo ni caja; dato en Plex Sans 18px blanco; placeholder Gris 400; etiqueta Plex Mono 12.5px en Gris 400, 8px por encima.
- **Focus:** la línea pasa a Azul Byte (140ms). El foco por teclado en controles añade un anillo Celeste de 2px desplazado 3px.
- **Error:** la línea pasa a Error y aparece debajo un mensaje en Plex Mono 12px en Error texto.
- **Cédula V/E:** control segmentado sobre Panel con borde de Línea de panel; opción en Plex Mono 500 13px Gris 400, seleccionada en Gris 700 con texto blanco. Va dentro de la misma línea que el número.

### Navigation

- **Lockup:** isotipo (26px) | separador de 1×22px en Gris 700 | wordmark "Byte" (19px), en blanco, separación 12px. Siempre arriba a la izquierda, sin distorsión.

### Regla de píxeles (signature)

Cuarenta cuadrados de hasta 6×6px, separados 3px, esquinas de 1px. Apagados en Línea de panel; encendidos, cada uno toma el color interpolado entre Azul Byte y Violeta según su posición. Se llena en proporción a los dígitos escritos del teléfono, se vacía segundo a segundo como cuenta atrás de reenvío y aparece llena en la confirmación. En error, todos los encendidos pasan a Error.

### Número monumental (signature)

Selector de prefijo (+58 por defecto) en Plex Mono 500 17px Gris 400 con un chevron dibujado en CSS, alineado por la línea base junto al número en Display blanco. El número se formatea en grupos (412 555 1234) mientras se escribe. La regla de píxeles va 12px debajo.

### Celdas OTP

Seis celdas en rejilla con separación de 8px que reflejan un único input real. Cada celda: 68px de alto, línea inferior de 2px en Gris 700, dígito en Display OTP. La celda actual subraya en Azul Byte y muestra un cursor de 2×34px que parpadea a 1s. En error, línea en Error y dígitos en Error texto.

### Patrón byte de esquina

Cuadrados de 13px con paso de 15px y esquinas de 1.5px, en Azul Byte, Celeste, Azul Profundo y Violeta con opacidades de 0.57 a 0.88, agrupados en una sola esquina: dentro de la columna (180px, detrás del contenido) por debajo de 1000px, fijo en la esquina inferior derecha de la ventana (240px) desde 1000px. Es decoración: `aria-hidden`, sin interacción.

### Movimiento

Entrada de cada paso: aparece y sube 10px en 320ms con la curva de salida. Cambios de estado: 140ms. Con `prefers-reduced-motion: reduce` se desactivan todas las transiciones y animaciones.

## Do's and Don'ts

### Do:

- **Do** usar Negro Byte (#0A0F1C) como fondo de toda pantalla.
- **Do** reservar Azul Byte (#309AE7) para la acción, el foco activo y la única palabra destacada del titular.
- **Do** escribir los campos sobre una línea de 1px (Gris 700 en reposo, Azul Byte con foco, Error en error).
- **Do** componer el dato principal del paso en Satoshi Bold a escala display, con cifras tabulares.
- **Do** expresar el paso de azul a violeta en píxeles discretos, como en la regla y el patrón.
- **Do** agrupar el patrón byte en una sola esquina.
- **Do** mantener el botón primario a ancho completo, 56px de alto y anclado al pie del paso.
- **Do** usar Error texto (#F87171) para mensajes de error en texto pequeño, para que cumplan AA.

### Don't:

- **Don't** usar negro puro (#000) como fondo.
- **Don't** añadir sombras, brillos o resplandores; el manual los prohíbe y el sistema es plano.
- **Don't** dibujar cajas alrededor de los campos de texto.
- **Don't** distorsionar el logo ni alterar las proporciones del lockup (isotipo | separador | wordmark).
- **Don't** usar Violeta, Celeste o Azul Profundo como color de acción.
- **Don't** usar texto de bajo contraste; el texto secundario no baja de Gris 400 sobre Negro Byte.
- **Don't** usar Satoshi para etiquetas o botones ni Plex Mono para texto corrido.

## Pendientes abiertos

- **Texto de términos de uso.** La hoja de términos es un marcador provisional; la redacción es una decisión abierta de Byte.
