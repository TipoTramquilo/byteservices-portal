// Genera results/report.html a partir de results/ramp.csv: una página con gráficos SVG, sin dependencias.
// Uso: node loadtest/report.mjs
import { readFileSync, writeFileSync } from 'node:fs';

const DIR = import.meta.dirname;
const CSV = `${DIR}/results/ramp.csv`,
  OUT = `${DIR}/results/report.html`;

const num = (n, d = 0) =>
  n.toLocaleString('es-ES', { minimumFractionDigits: d, maximumFractionDigits: d });

const rows = readFileSync(CSV, 'utf8')
  .trim()
  .split('\n')
  .slice(1)
  .map((l) => {
    const [etapa, modo, vus, rps, p50, p95, p99, fallos, r429, r500, cpu, mem] = l.split(',');
    return {
      etapa,
      modo,
      vus: +vus,
      rps: +rps,
      p50: +p50,
      p95: +p95,
      p99: +p99,
      fallos: +fallos,
      r429: +r429,
      r500: +r500,
      cpu: +cpu,
      mem: +mem,
    };
  });
const send = rows.filter((r) => r.modo === 'send').toSorted((a, b) => a.vus - b.vus);
const base = rows.filter((r) => r.modo !== 'send');
const last = send.at(-1);
let peak = send[0];
for (const r of send) if (r.rps > peak.rps) peak = r;

// Eje Y con 4 divisiones de paso "redondo" (1, 2, 2.5, 5 o 10 por potencia de diez):
// el máximo es 4 pasos, así que cada marca cae en un número limpio
function niceMax(v) {
  const exp = 10 ** Math.floor(Math.log10(v / 4));
  return 4 * [1, 2, 2.5, 5, 10].map((m) => m * exp).find((step) => 4 * step >= v);
}

const W = 640,
  H = 280,
  P = { l: 56, r: 16, t: 16, b: 40 };
const IW = W - P.l - P.r,
  IH = H - P.t - P.b;

// Ejes y rejilla comunes; las categorías (los niveles de VUs) van igual espaciadas
function frame(labels, yMax, fmtY) {
  const n = labels.length;
  const x = (i) => P.l + ((i + 0.5) * IW) / n;
  const y = (v) => P.t + IH - (v / yMax) * IH;
  const grid = [0, 1, 2, 3, 4]
    .map((k) => {
      const v = (yMax * k) / 4;
      return (
        `<line x1="${P.l}" x2="${W - P.r}" y1="${y(v)}" y2="${y(v)}" class="grid"/>` +
        `<text x="${P.l - 8}" y="${y(v) + 4}" class="tick" text-anchor="end">${fmtY(v)}</text>`
      );
    })
    .join('');
  const xs = labels
    .map((l, i) => `<text x="${x(i)}" y="${H - 16}" class="tick" text-anchor="middle">${l}</text>`)
    .join('');
  return { x, y, base: grid + xs };
}

function barChart(labels, values, { color, yMax, fmtY, fmtV }) {
  const f = frame(labels, yMax, fmtY);
  const bw = Math.min(56, (IW / labels.length) * 0.55);
  const bars = values
    .map((v, i) => {
      const top = f.y(v),
        cx = f.x(i);
      return (
        `<rect x="${cx - bw / 2}" y="${top}" width="${bw}" height="${P.t + IH - top}" rx="3" fill="${color}"/>` +
        `<text x="${cx}" y="${top - 6}" class="val" text-anchor="middle">${fmtV(v)}</text>`
      );
    })
    .join('');
  return `<svg viewBox="0 0 ${W} ${H}" role="img">${f.base}${bars}</svg>`;
}

function lineChart(labels, series, { yMax, fmtY }) {
  const f = frame(labels, yMax, fmtY);
  const lines = series
    .map(({ values, color }) => {
      const pts = values.map((v, i) => `${f.x(i)},${f.y(v)}`);
      return `<polyline points="${pts.join(' ')}" fill="none" stroke="${color}" stroke-width="2.5" stroke-linejoin="round"/>${values
        .map((v, i) => `<circle cx="${f.x(i)}" cy="${f.y(v)}" r="4" fill="${color}"/>`)
        .join('')}`;
    })
    .join('');
  return `<svg viewBox="0 0 ${W} ${H}" role="img">${f.base}${lines}</svg>`;
}

const legend = (items) =>
  `<div class="legend">${items
    .map(([name, color]) => `<span><i style="background:${color}"></i>${name}</span>`)
    .join('')}</div>`;

const labels = send.map((r) => num(r.vus));
const maxRps = niceMax(Math.max(...send.map((r) => r.rps)));
const maxLat = niceMax(Math.max(...send.map((r) => r.p99)));
const maxCpu = niceMax(Math.max(...send.map((r) => r.cpu), 1));
const maxMem = niceMax(Math.max(...send.map((r) => r.mem)));

const charts = [
  {
    title: 'Peticiones por segundo',
    sub: 'Por VUs simultáneos. Más alto, mejor: el techo es donde deja de subir',
    legend: '',
    svg: barChart(
      labels,
      send.map((r) => r.rps),
      { color: '#309AE7', yMax: maxRps, fmtY: (v) => num(v), fmtV: (v) => num(v) }
    ),
  },
  {
    title: 'Latencia (ms)',
    sub: 'Por VUs simultáneos. Tiempo de respuesta de /api/send por percentil',
    legend: legend([
      ['p50', '#5FD6FF'],
      ['p95', '#309AE7'],
      ['p99', '#5036EA'],
    ]),
    svg: lineChart(
      labels,
      [
        { values: send.map((r) => r.p50), color: '#5FD6FF' },
        { values: send.map((r) => r.p95), color: '#309AE7' },
        { values: send.map((r) => r.p99), color: '#5036EA' },
      ],
      { yMax: maxLat, fmtY: (v) => num(v) }
    ),
  },
  {
    title: 'CPU del portal (%)',
    sub: 'Por VUs simultáneos. Pico de CPU del contenedor',
    legend: '',
    svg: barChart(
      labels,
      send.map((r) => r.cpu),
      { color: '#5FD6FF', yMax: maxCpu, fmtY: (v) => `${num(v)}%`, fmtV: (v) => `${num(v)}%` }
    ),
  },
  {
    title: 'Memoria del portal (MB)',
    sub: 'Por VUs simultáneos. Pico de RAM del contenedor',
    legend: '',
    svg: lineChart(labels, [{ values: send.map((r) => r.mem), color: '#5036EA' }], {
      yMax: maxMem,
      fmtY: (v) => num(v),
    }),
  },
];

const cell = (v) => `<td>${v}</td>`;
const tableRows = rows
  .map(
    (r) =>
      `<tr>${[
        r.etapa,
        r.modo,
        num(r.vus),
        num(r.rps),
        num(r.p50),
        num(r.p95),
        num(r.p99),
        `${num(r.fallos, 1)}%`,
        num(r.r429),
        num(r.r500),
        `${num(r.cpu)}%`,
        num(r.mem, 1),
      ]
        .map(cell)
        .join('')}</tr>`
  )
  .join('');

const generated = new Date().toLocaleString('es-ES', {
  timeZone: 'America/Caracas',
  dateStyle: 'medium',
  timeStyle: 'short',
});

const html = `<!doctype html>
<html lang="es">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Prueba de carga · Portal Byte</title>
<style>
  :root { --bg: #0A0F1C; --card: #121A2B; --line: #232C40; --text: #EDEAF5; --muted: #8B93A7;
    --blue: #309AE7; --cyan: #5FD6FF; --violet: #5036EA; }
  * { box-sizing: border-box; }
  body { margin: 0; padding: 32px 16px; background: var(--bg); color: var(--text);
    font: 15px/1.5 "Plex Sans", system-ui, -apple-system, Segoe UI, Roboto, sans-serif; }
  main { max-width: 1100px; margin: 0 auto; }
  h1 { font: 700 28px/1.2 Satoshi, Avenir Next, Segoe UI, system-ui, sans-serif; margin: 0 0 4px; }
  .lede { color: var(--muted); margin: 0 0 24px; }
  .kpis { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 12px; margin-bottom: 24px; }
  .kpi { background: var(--card); border: 1px solid var(--line); border-radius: 12px; padding: 16px; }
  .kpi b { display: block; font: 500 26px/1.1 "Plex Mono", ui-monospace, SF Mono, Menlo, Consolas, monospace; color: var(--cyan); }
  .kpi span { color: var(--muted); font-size: 13px; }
  .charts { display: grid; grid-template-columns: 1fr; gap: 16px; }
  @media (min-width: 900px) { .charts { grid-template-columns: 1fr 1fr; } }
  .card { background: var(--card); border: 1px solid var(--line); border-radius: 12px; padding: 18px; }
  .card h2 { font: 700 17px/1.3 Satoshi, Avenir Next, Segoe UI, system-ui, sans-serif; margin: 0; }
  .card p { color: var(--muted); margin: 2px 0 10px; font-size: 13px; }
  svg { width: 100%; height: auto; display: block; }
  .grid { stroke: var(--line); stroke-width: 1; }
  .tick, .val { fill: var(--muted); font: 12px "Plex Mono", ui-monospace, SF Mono, Menlo, Consolas, monospace; }
  .val { fill: var(--text); }
  .unit { fill: var(--muted); font: 11px "Plex Sans", system-ui, -apple-system, Segoe UI, Roboto, sans-serif; }
  .legend { display: flex; gap: 14px; margin-bottom: 8px; font-size: 13px; color: var(--muted); }
  .legend i { display: inline-block; width: 10px; height: 10px; border-radius: 2px; margin-right: 6px; }
  h3 { font: 700 17px/1.3 Satoshi, Avenir Next, Segoe UI, system-ui, sans-serif; margin: 32px 0 10px; }
  .table-wrap { overflow-x: auto; background: var(--card); border: 1px solid var(--line); border-radius: 12px; }
  table { width: 100%; border-collapse: collapse; font: 13px "Plex Mono", ui-monospace, SF Mono, Menlo, Consolas, monospace; }
  th, td { padding: 9px 12px; text-align: right; white-space: nowrap; border-bottom: 1px solid var(--line); }
  th { color: var(--muted); font-weight: 500; }
  th:first-child, td:first-child, th:nth-child(2), td:nth-child(2) { text-align: left; }
  tr:last-child td { border-bottom: 0; }
  footer { margin-top: 28px; color: var(--muted); font-size: 13px; }
</style>
</head>
<body>
<main>
  <h1>Prueba de carga · /api/send</h1>
  <p class="lede">Portal cautivo Byte en un contenedor Podman, k6 en la red de compose. Generado ${generated}.</p>

  <section class="kpis">
    <div class="kpi"><b>${num(peak.rps)}</b><span>rps máximo (a ${num(peak.vus)} VUs)</span></div>
    <div class="kpi"><b>${num(last.p95)} ms</b><span>p95 a ${num(last.vus)} VUs</span></div>
    <div class="kpi"><b>${num(last.p99)} ms</b><span>p99 a ${num(last.vus)} VUs</span></div>
    <div class="kpi"><b>${num(Math.max(...send.map((r) => r.cpu)))}%</b><span>CPU máxima del portal</span></div>
    <div class="kpi"><b>${num(send.reduce((a, r) => a + r.r429 + r.r500, 0))}</b><span>respuestas 429 y 500</span></div>
  </section>

  <section class="charts">
    ${charts.map((c) => `<article class="card"><h2>${c.title}</h2><p>${c.sub}</p>${c.legend}${c.svg}</article>`).join('\n    ')}
  </section>

  <h3>Todas las etapas</h3>
  <div class="table-wrap">
    <table>
      <thead><tr><th>Etapa</th><th>Modo</th><th>VUs</th><th>rps</th><th>p50 ms</th><th>p95 ms</th><th>p99 ms</th>
        <th>Fallos</th><th>429</th><th>500</th><th>CPU máx</th><th>RAM máx MB</th></tr></thead>
      <tbody>${tableRows}</tbody>
    </table>
  </div>

  <footer>Los VUs son peticiones en vuelo a la vez. k6 y el portal comparten la VM de Podman, así que a altas tasas
    k6 también pesa en el resultado. Etapas de referencia: ${base.map((r) => r.etapa).join(', ')}.</footer>
</main>
</body>
</html>
`;

writeFileSync(OUT, html);
console.log(`Reporte en ${OUT}`);
