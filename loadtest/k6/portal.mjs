// Prueba de carga del portal cautivo Byte (k6).
// MODE=health|index|send   VUS=<n>   DUR=<duración>   SLEEP=<pausa entre iteraciones>
// BASE=<url del portal, por defecto la red de compose>
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';
import exec from 'k6/execution';

// Variables de `k6 run -e CLAVE=valor`. k6 solo las expone en la global __ENV (no hay módulo
// equivalente): se lee una vez aquí y el resto del script usa `env`
const env = __ENV;
const BASE = env.BASE || 'http://portal:8080';
const MODE = env.MODE || 'send';
const VUS = Number(env.VUS || 10);
const DUR = env.DUR || '20s';
const PAUSE = Number(env.SLEEP || 0);
// RATE>0 cambia a modelo abierto: pide RATE peticiones/s sin esperar a que responda la anterior
const RATE = Number(env.RATE || 0);

const r429 = new Counter('resp_429');
const r500 = new Counter('resp_500');

// Sal distinta en cada corrida: el límite es por teléfono y día, y un reenvío dentro de
// los 60 s daría too_soon, así que los teléfonos no deben repetirse entre corridas
const SALT = String(Date.now() % 1000).padStart(3, '0');

export const options = {
  scenarios: {
    carga:
      RATE > 0
        ? {
            executor: 'constant-arrival-rate',
            rate: RATE,
            timeUnit: '1s',
            duration: DUR,
            preAllocatedVUs: Math.min(2000, Math.max(50, Math.ceil(RATE / 10))),
            maxVUs: Math.min(5000, Math.max(200, RATE)),
          }
        : {
            executor: 'constant-vus',
            vus: VUS,
            duration: DUR,
            gracefulStop: '10s',
          },
  },
  summaryTrendStats: ['avg', 'p(50)', 'p(95)', 'p(99)', 'max'],
};

export default function portal() {
  if (MODE === 'health') {
    const res = http.get(`${BASE}/api/health`);
    check(res, { 204: (r) => r.status === 204 });
  } else if (MODE === 'index') {
    const res = http.get(`${BASE}/`);
    check(res, { 200: (r) => r.status === 200 });
  } else {
    // teléfono y cédula únicos por iteración (+584 + sal + 8 dígitos)
    const n = String(exec.vu.idInTest * 10000 + (exec.vu.iterationInScenario % 10000)).padStart(
      8,
      '0'
    );
    const body = JSON.stringify({
      name: `Carga ${n}`,
      id_number: `V${n}`,
      company: 'Carga Byte',
      phone: `+584${SALT}${n}`,
    });
    const res = http.post(`${BASE}/api/send`, body, {
      headers: { 'Content-Type': 'application/json' },
    });
    if (res.status === 429) r429.add(1);
    if (res.status >= 500) r500.add(1);
    check(res, { 200: (r) => r.status === 200 });
  }
  if (PAUSE > 0) sleep(PAUSE);
}

export function handleSummary(data) {
  const m = data.metrics;
  const d = m.http_req_duration.values;
  const line = [
    `mode=${MODE}`,
    `vus=${VUS}`,
    `rate=${RATE}`,
    `reqs=${m.http_reqs.values.count}`,
    `rps=${m.http_reqs.values.rate.toFixed(1)}`,
    `p50=${(d['p(50)'] ?? 0).toFixed(1)}ms`,
    `p95=${(d['p(95)'] ?? 0).toFixed(1)}ms`,
    `p99=${(d['p(99)'] ?? 0).toFixed(1)}ms`,
    `max=${(d.max ?? 0).toFixed(1)}ms`,
    `fail=${((m.http_req_failed?.values.rate ?? 0) * 100).toFixed(2)}%`,
    `r429=${m.resp_429?.values.count ?? 0}`,
    `r500=${m.resp_500?.values.count ?? 0}`,
  ].join(' ');
  return {
    '/loadtest/results/last.json': JSON.stringify({
      mode: MODE,
      vus: VUS,
      rate: RATE,
      duration: DUR,
      pause: PAUSE,
      metrics: m,
    }),
    stdout: `${line}\n`,
  };
}
