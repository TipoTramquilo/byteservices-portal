// Builds dist/: one self-contained index.html (CSS inlined, minified) plus fonts/.
// No dependencies: `node build.mjs`
import { readFile, writeFile, mkdir, cp, rm } from 'node:fs/promises';
import { gzipSync } from 'node:zlib';

const SRC = 'portal', OUT = 'dist';

const minCss = (css) => css
  .replace(/\/\*[\s\S]*?\*\//g, '')
  .replace(/\s+/g, ' ')
  .replace(/\s*([{}:;,>])\s*/g, '$1')
  .replace(/;}/g, '}')
  .trim();

const minJs = (js) => js
  .replace(/\/\*[\s\S]*?\*\//g, '')
  .split('\n').map((l) => l.trim()).filter(Boolean).join('\n');

const minHtml = (html) => html
  .replace(/<!--[\s\S]*?-->/g, '')
  .replace(/>\s+</g, '><')
  .replace(/\n\s*/g, ' ');

let html = await readFile(`${SRC}/index.html`, 'utf8');
const css = await readFile(`${SRC}/byte.css`, 'utf8');

// keep scripts and styles out of the HTML whitespace pass, minify each on its own
const blocks = [];
const keep = (s) => `\u0000${blocks.push(s) - 1}\u0000`;
html = html
  .replace('<link rel="stylesheet" href="byte.css">', () => keep(`<style>${minCss(css)}</style>`))
  .replace(/<style>([\s\S]*?)<\/style>/g, (_, c) => keep(`<style>${minCss(c)}</style>`))
  .replace(/<script>([\s\S]*?)<\/script>/g, (_, j) => keep(`<script>${minJs(j)}</script>`));
html = minHtml(html).replace(/\u0000(\d+)\u0000/g, (_, i) => blocks[i]);

await rm(OUT, { recursive: true, force: true });
await mkdir(OUT, { recursive: true });
await writeFile(`${OUT}/index.html`, html);
await cp(`${SRC}/fonts`, `${OUT}/fonts`, { recursive: true });

const kb = (n) => (n / 1024).toFixed(1) + ' KB';
console.log(`dist/index.html  ${kb(html.length)}  (gzip ${kb(gzipSync(html, { level: 9 }).length)})`);
