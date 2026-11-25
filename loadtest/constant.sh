#!/bin/bash
# Carga constante: 1000 VUs en bucle, sin pausa, contra /api/send durante DUR segundos.
# Mide el consumo máximo del portal y de k6. Uso: ./loadtest/constant.sh [segundos]
set -u
IMG=docker.io/grafana/k6
NET=byte-portal_default
DIR=$(cd "$(dirname "$0")" && pwd)
RESULTS=$DIR/results
DUR=${1:-30}
mkdir -p "$RESULTS"
rm -f "$RESULTS/last.json"

echo
echo "### 1000 VUs simultáneos, ${DUR} s, sin pausa"
( for _ in $(seq 1 $((DUR + 60))); do timeout 3 podman stats --no-stream \
    --format '{{.Name}} {{.CPUPerc}} {{.MemUsage}}' byte-portal k6run 2>/dev/null; sleep 1; done > "$RESULTS/stats-constant-1000.log" ) &
SP=$!
podman run --rm --name k6run --network "$NET" -v "$DIR:/loadtest" "$IMG" run \
    -e MODE=send -e VUS=1000 -e DUR="${DUR}s" /loadtest/k6/portal.mjs
kill $SP 2>/dev/null; wait $SP 2>/dev/null

echo
echo "### resumen"
jq -r '.metrics | "reqs=\(.http_reqs.values.count) rps=\(.http_reqs.values.rate|floor) p50=\(.http_req_duration.values["p(50)"]|floor)ms p95=\(.http_req_duration.values["p(95)"]|floor)ms p99=\(.http_req_duration.values["p(99)"]|floor)ms fail=\(.http_req_failed.values.rate*100|floor)% r429=\(.resp_429.values.count // 0) r500=\(.resp_500.values.count // 0)"' "$RESULTS/last.json"
for c in byte-portal k6run; do
  cpu=$(awk -v n="$c" '$1==n{gsub("%","",$2); if($2+0>m)m=$2+0} END{printf "%.0f", m}' "$RESULTS/stats-constant-1000.log")
  mem=$(awk -v n="$c" '$1==n{print $3}' "$RESULTS/stats-constant-1000.log" | sort -h | tail -1)
  echo "  $c: cpu_max=${cpu}% mem_max=${mem}"
done
