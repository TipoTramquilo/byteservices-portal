# ============================================================
#  Portal cautivo Byte - apps/api (Go) + apps/web/dist (portal compilado)
#  Contexto de build: la raíz del repo (necesita apps/api/ y apps/web/dist/)
#  podman build -f docker/Dockerfile -t byte-portal .
# ============================================================

# ---------- Build: binario estático ----------
FROM docker.io/library/golang:alpine AS build
WORKDIR /src

COPY apps/api/go.mod apps/api/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY apps/api/ ./
# CGO_ENABLED=0: modernc.org/sqlite es Go puro, así el binario no depende de libc
# y corre sobre scratch. -s -w quita los símbolos de debug
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/captiveportal . \
    && mkdir -p /out/data

# ---------- Runtime: scratch (solo el binario, el portal y los certificados) ----------
FROM scratch

# Certificados CA para la API de Twilio por HTTPS: es lo único del sistema que usa el binario
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/captiveportal /app/captiveportal
COPY apps/web/dist/ /app/dist/
# Volumen de la base SQLite: nace con dueño 65532 para que la app escriba sin ser root
COPY --from=build --chown=65532:65532 /out/data /data

ENV ADDR=:8080 \
    DB=/data/portal.db \
    PORTAL_DIR=/app/dist
EXPOSE 8080
VOLUME /data

# scratch no tiene /etc/passwd: el usuario va por número
USER 65532:65532

# El mismo binario hace de healthcheck: pide /api/health y sale con 0 o 1
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/app/captiveportal", "health"]

ENTRYPOINT ["/app/captiveportal"]
