#!/bin/bash
# Rampa de carga contra /api/send: 50 -> 1000 VUs, un nivel tras otro, todos pidiendo a la vez
# (sin pausa entre peticiones). Cada nivel deja una fila en results/ramp.csv con el resumen de k6
# y el consumo máximo del portal. Después: node loadtest/report.mjs
# Uso: ./loadtest/ramp.sh
set -u
IMG=docker.io/grafana/k6
NET=byte-portal_default
DIR=$(cd "$(dirname "$0")" && pwd)
RESULTS=$DIR/results
CSV=$RESULTS/ramp.csv
mkdir -p "$RESULTS"
echo "etapa,modo,vus,rps,p50_ms,p95_ms,p99_ms,fallos_pct,r429,r500,cpu_max_pct,mem_max_mb" > "$CSV"

run() { # modo vus dur etiqueta
  local mode=$1 vus=$2 dur=$3 label=$4
  local stats=$RESULTS/stats-$label.log
  echo
  echo "### $label: modo=$mode vus=$vus duración=${dur}s (sin pausa)"
  ( for i in $(seq 1 200); do timeout 3 podman stats --no-stream \
      --format '{{.Name}} {{.CPUPerc}} {{.MemUsage}}' byte-portal k6run 2>/dev/null; sleep 1; done > "$stats" ) &
  local sp=$!
  podman run --rm --name k6run --network "$NET" -v "$DIR:/scripts" "$IMG" run \
      -e MODE="$mode" -e VUS="$vus" -e DUR="${dur}s" /scripts/scripts/portal.js
  kill $sp 2>/dev/null; wait $sp 2>/dev/null

  # El consumo llega como "39.62MB / 7.8GB": se pasa todo a MB para poder comparar
  local cpu mem
  cpu=$(awk -v n=byte-portal '$1==n{gsub("%","",$2); if($2+0>m)m=$2+0} END{printf "%.1f", m}' "$stats")
  mem=$(awk -v n=byte-portal '
    function mb(v) {
      x = v + 0
      if (v ~ /GiB/) return x * 1024
      if (v ~ /MiB/) return x * 1.048576
      if (v ~ /GB/) return x * 1000
      if (v ~ /kB|KB/) return x / 1000
      return x
    }
    $1==n { if (mb($3) > m) m = mb($3) } END { printf "%.1f", m }' "$stats")

  local rps p50 p95 p99 fail r429 r500
  read -r rps p50 p95 p99 fail r429 r500 < <(jq -r '.metrics | [
      (.http_reqs.values.rate | floor),
      (.http_req_duration.values["p(50)"] | floor),
      (.http_req_duration.values["p(95)"] | floor),
      (.http_req_duration.values["p(99)"] | floor),
      ((.http_req_failed.values.rate * 1000 | round) / 10),
      (.resp_429.values.count // 0),
      (.resp_500.values.count // 0)] | @tsv' "$RESULTS/last.json")
  echo "$label,$mode,$vus,$rps,$p50,$p95,$p99,$fail,$r429,$r500,$cpu,$mem" >> "$CSV"
  echo "  rps=$rps p50=${p50}ms p95=${p95}ms p99=${p99}ms fallos=${fail}% 429=$r429 500=$r500 | portal cpu=${cpu}% mem=${mem}MB"
}

run health 500 20 health-500
run index 200 20 index-200
for v in 50 100 200 400 700 1000; do run send "$v" 25 "send-$v"; done
echo
echo "Resumen en $CSV"
