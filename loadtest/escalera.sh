#!/bin/bash
# Escalera de concurrencia real contra el portal: 50 -> 1000 VUs, todos pidiendo a la vez
# (sin pausa entre peticiones). En cada nivel mide el consumo máximo del contenedor.
# Uso: ./loadtest/escalera.sh
set -u
IMG=docker.io/grafana/k6
NET=byte-portal_default
DIR=$(cd "$(dirname "$0")" && pwd)
OUTDIR=$DIR/out
OUT=$OUTDIR/escalera.txt
mkdir -p "$OUTDIR"
: > "$OUT"

run() { # mode vus dur etiqueta
  local mode=$1 vus=$2 dur=$3 label=$4
  local stats=$OUTDIR/stats-$label.log
  echo
  echo "### $label — mode=$mode vus=$vus dur=${dur}s (sin pausa)"
  ( for i in $(seq 1 200); do timeout 3 podman stats --no-stream \
      --format '{{.Name}} {{.CPUPerc}} {{.MemUsage}}' byte-portal k6run 2>/dev/null; sleep 1; done > "$stats" ) &
  local sp=$!
  podman run --rm --name k6run --network "$NET" -v "$DIR:/scripts" "$IMG" run \
      -e MODE="$mode" -e VUS="$vus" -e DUR="${dur}s" /scripts/k6.js
  kill $sp 2>/dev/null; wait $sp 2>/dev/null

  local cpu mem
  cpu=$(awk '$1=="byte-portal"{gsub("%","",$2); if($2+0>m)m=$2+0} END{printf "%.0f", m}' "$stats")
  mem=$(awk '$1=="byte-portal"{print $3}' "$stats" | sort -h | tail -1)
  {
    echo "[$label] $(jq -r '.metrics | "reqs=\(.http_reqs.values.count) rps=\(.http_reqs.values.rate|floor) p50=\(.http_req_duration.values["p(50)"]|floor)ms p95=\(.http_req_duration.values["p(95)"]|floor)ms p99=\(.http_req_duration.values["p(99)"]|floor)ms fail=\(.http_req_failed.values.rate*100|floor)% r429=\(.resp_429.values.count // 0) r500=\(.resp_500.values.count // 0)"' "$OUTDIR/last.json")"
    echo "  byte-portal: cpu_max=${cpu}% mem_max=${mem}"
  } | tee -a "$OUT"
}

echo "### base: health e index" | tee -a "$OUT"
run health 500 20 health-500
run index 200 20 index-200
echo "### send, 1000 peticiones en vuelo como máximo" | tee -a "$OUT"
for v in 50 100 200 400 700 1000; do run send "$v" 25 "send-$v"; done
echo
echo "Resumen en $OUT"
