// Recorta las fuentes de fonts-src/ a los caracteres que puede mostrar el portal y las deja en src/fonts/
// Se corre con `task web:fonts` (`pnpm fonts`); subset-font es una devDependency, no entra al portal
import { readFile, writeFile, readdir, mkdir } from 'node:fs/promises';
import { join } from 'node:path';
import subsetFont from 'subset-font';

const DIR = import.meta.dirname,
  SRC = join(DIR, 'fonts-src'),
  OUT = join(DIR, 'src', 'fonts');

const uniq = (s) => [...new Set(s)].filter((c) => c.codePointAt(0) >= 0x20).join('');
const range = (a, b) =>
  Array.from({ length: b - a + 1 }, (_, i) => String.fromCharCode(a + i)).join('');

// Lo que muestran el portal y el design system, más ASCII imprimible
let shown = range(0x20, 0x7e);
const sources = ['src/index.html', 'src/app.mjs', 'design-system.html'];
shown = uniq(
  shown + (await Promise.all(sources.map((f) => readFile(join(DIR, f), 'utf8')))).join('')
);

// Plex Sans es la fuente de los campos de nombre y cédula: el visitante puede escribir cualquier acento,
// así que lleva Latin-1 completo. Satoshi (títulos, teléfono, código) y Plex Mono (etiquetas, botones)
// solo muestran texto del portal
const textFor = (name) => (name.startsWith('plex-sans') ? uniq(shown + range(0xa0, 0xff)) : shown);

await mkdir(OUT, { recursive: true });
const kb = (n) => `${(n / 1024).toFixed(1)} KB`;
// Se recortan en paralelo; el resumen sale en el orden de los archivos
const fonts = (await readdir(SRC)).filter((f) => f.endsWith('.woff2'));
const results = await Promise.all(
  fonts.map(async (name) => {
    const original = await readFile(join(SRC, name));
    const subset = await subsetFont(original, textFor(name), { targetFormat: 'woff2' });
    await writeFile(join(OUT, name), subset);
    return `${name.padEnd(20)} ${kb(original.length).padStart(8)} -> ${kb(subset.length)}`;
  })
);
for (const line of results) console.log(line);
