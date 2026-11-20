// Recorta las fuentes de fonts-src/ a los caracteres que puede mostrar el portal y las deja en src/fonts/
// Se corre con `make fonts`, que trae subset-font con npx (no queda como dependencia del repo)
import { readFile, writeFile, readdir, mkdir } from 'node:fs/promises';
import { createRequire } from 'node:module';
import { join } from 'node:path';

// npx instala subset-font en su caché y pone su node_modules/.bin en el PATH: se carga desde ahí
const bin = process.env.PATH.split(':').find((p) => p.includes('_npx') && p.endsWith(join('node_modules', '.bin')));
if (!bin) throw new Error('correr con: npx -y -p subset-font -- node apps/web/subset-fonts.mjs');
const subsetFont = createRequire(join(bin, '..', 'x.js'))('subset-font');

const DIR = import.meta.dirname, SRC = join(DIR, 'fonts-src'), OUT = join(DIR, 'src', 'fonts');

const uniq = (s) => [...new Set(s)].filter((c) => c.codePointAt(0) >= 0x20).join('');
const range = (a, b) => Array.from({ length: b - a + 1 }, (_, i) => String.fromCharCode(a + i)).join('');

// Lo que muestran el portal y el design system, más ASCII imprimible
let shown = range(0x20, 0x7e);
for (const f of ['src/index.html', 'src/app.js', 'design-system.html']) shown += await readFile(join(DIR, f), 'utf8');
shown = uniq(shown);

// Plex Sans es la fuente de los campos de nombre y cédula: el visitante puede escribir cualquier acento,
// así que lleva Latin-1 completo. Satoshi (títulos, teléfono, código) y Plex Mono (etiquetas, botones)
// solo muestran texto del portal
const textFor = (name) => (name.startsWith('plex-sans') ? uniq(shown + range(0xa0, 0xff)) : shown);

await mkdir(OUT, { recursive: true });
const kb = (n) => (n / 1024).toFixed(1) + ' KB';
for (const name of (await readdir(SRC)).filter((f) => f.endsWith('.woff2'))) {
  const original = await readFile(join(SRC, name));
  const subset = await subsetFont(original, textFor(name), { targetFormat: 'woff2' });
  await writeFile(join(OUT, name), subset);
  console.log(`${name.padEnd(20)} ${kb(original.length).padStart(8)} -> ${kb(subset.length)}`);
}
