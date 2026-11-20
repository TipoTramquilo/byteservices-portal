# ============================================================
#  Portal cautivo Byte - apps/api (Go) + apps/web (portal, compilado en la imagen)
#  Contexto de build: la raíz del repo (necesita apps/api/ y apps/web/)
#  podman build -f docker/Dockerfile -t byte-portal .
# ============================================================

# ---------- Portal: build.mjs no tiene dependencias, solo necesita Node ----------
FROM docker.io/library/node:alpine AS web
WORKDIR /src
COPY apps/web/build.mjs apps/web/build.mjs
COPY apps/web/src apps/web/src
RUN node apps/web/build.mjs

# ---------- Build: binario estático ----------
FROM docker.io/library/golang:alpine AS build
WORKDIR /src

COPY apps/api/go.mod apps/api/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY apps/api/ ./
# CGO_ENABLED=0: el binario no depende de libc y corre sobre scratch. -s -w quita los símbolos de debug
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/captiveportal .

# ---------- Runtime: scratch (solo el binario, el portal y los certificados) ----------
FROM scratch

# Certificados CA para la API de Twilio por HTTPS: es lo único del sistema que usa el binario
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/captiveportal /app/captiveportal
COPY --from=web /src/apps/web/dist/ /app/dist/
# La base es Postgres (DATABASE_URL se pasa en tiempo de ejecución); el contenedor no guarda datos
ENV ADDR=:8080 \
    PORTAL_DIR=/app/dist
EXPOSE 8080

# scratch no tiene /etc/passwd: el usuario va por número
USER 65532:65532

# El mismo binario hace de healthcheck: pide /api/health y sale con 0 o 1
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/app/captiveportal", "health"]

ENTRYPOINT ["/app/captiveportal"]
