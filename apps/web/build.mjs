// Builds dist/ from src/: one self-contained index.html (byte.css and app.mjs inlined, minified) plus fonts/.
// No dependencies: `node apps/web/build.mjs` (paths resolve from this file, any cwd works)
import { readFile, writeFile, mkdir, cp, rm } from 'node:fs/promises';
import { gzipSync } from 'node:zlib';

const SRC = `${import.meta.dirname}/src`,
  OUT = `${import.meta.dirname}/dist`;

const minCss = (css) =>
  css
    .replaceAll(/\/\*[\s\S]*?\*\//g, '')
    .replaceAll(/\s+/g, ' ')
    .replaceAll(/\s*([{}:;,>])\s*/g, '$1')
    .replaceAll(';}', '}')
    .trim();

const minJs = (js) =>
  js
    .replaceAll(/\/\*[\s\S]*?\*\//g, '')
    .split('\n')
    .map((l) => l.trim())
    .filter(Boolean)
    .join('\n');

const minHtml = (html) =>
  html
    .replaceAll(/<!--[\s\S]*?-->/g, '')
    .replaceAll(/>\s+</g, '><')
    .replaceAll(/\n\s*/g, ' ');

let html = await readFile(`${SRC}/index.html`, 'utf8');
const css = await readFile(`${SRC}/byte.css`, 'utf8');
const js = await readFile(`${SRC}/app.mjs`, 'utf8');

// The CSS and JS are inlined by exact string match: if a tag changes (e.g. a formatter turns it
// into `<link ... />`), fail instead of shipping a portal without styles or behaviour
for (const tag of ['<link rel="stylesheet" href="byte.css">', '<script src="app.mjs"></script>']) {
  if (!html.includes(tag)) throw new Error(`src/index.html must contain exactly: ${tag}`);
}

// keep scripts and styles out of the HTML whitespace pass, minify each on its own
const blocks = [];
const keep = (s) => `\u0000${blocks.push(s) - 1}\u0000`;
html = html
  .replace('<link rel="stylesheet" href="byte.css">', () => keep(`<style>${minCss(css)}</style>`))
  .replace('<script src="app.mjs"></script>', () => keep(`<script>${minJs(js)}</script>`))
  .replaceAll(/<style>([\s\S]*?)<\/style>/g, (_, c) => keep(`<style>${minCss(c)}</style>`))
  .replaceAll(/<script>([\s\S]*?)<\/script>/g, (_, j) => keep(`<script>${minJs(j)}</script>`));
// oxlint-disable-next-line no-control-regex -- \u0000 delimita los bloques reservados (ver keep)
html = minHtml(html).replaceAll(/\u0000(\d+)\u0000/g, (_, i) => blocks[i]);

await rm(OUT, { recursive: true, force: true });
await mkdir(OUT, { recursive: true });
await writeFile(`${OUT}/index.html`, html);
await cp(`${SRC}/fonts`, `${OUT}/fonts`, { recursive: true });

const kb = (n) => `${(n / 1024).toFixed(1)} KB`;
const gzip = gzipSync(html, { level: 9 }).length;
console.log(`dist/index.html  ${kb(html.length)}  (gzip ${kb(gzip)})`);

// Presupuesto de peso: la página viaja antes de que el visitante tenga internet, a veces con
// mala señal. Subirlo es una decisión, no un efecto secundario (CI y pre-push corren este build)
const BUDGET_GZIP = 12 * 1024;
if (gzip > BUDGET_GZIP) {
  throw new Error(`dist/index.html pesa ${kb(gzip)} gzip y el presupuesto es ${kb(BUDGET_GZIP)}`);
}
