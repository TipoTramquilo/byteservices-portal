# Motor de contenedores: podman si está, si no docker (forzar con `make ENGINE=docker`)
ENGINE ?= $(shell command -v podman 2>/dev/null || command -v docker 2>/dev/null)
ENGINE := $(if $(strip $(ENGINE)),$(strip $(ENGINE)),docker)

DC = "$(ENGINE)" compose --env-file .env -f docker/docker-compose.yaml

# El Go que pide apps/api/go.mod va primero en el PATH: si el del sistema es más viejo, `go`
# cambia de toolchain al vuelo, pero govulncheck y golangci-lint llaman a `go list` y verían el viejo
GOROOT_BIN := $(shell cd apps/api && go env GOROOT 2>/dev/null)/bin
export PATH := $(GOROOT_BIN):$(PATH)

# Herramientas de Go fijadas en tools/go.mod y tools/actionlint/go.mod (Dependabot las actualiza)
BIN := $(CURDIR)/.bin

# Herramientas del sistema: en local, si faltan, se avisa y se sigue; en CI (CI=true) faltar es un error
define need
	@if ! command -v $(1) >/dev/null 2>&1; then \
		if [ -n "$$CI" ]; then echo "✖ falta $(1)"; exit 1; fi; \
		echo "⚠ $(1) no está instalado: se salta ($(2))"; exit 0; \
	fi
endef

.PHONY: help web web-dev fonts api build up down logs \
	install tools check fmt lint web-check go-check go-lint go-fmt go-test go-vuln go-mod \
	shell-lint docker-lint actions-lint tf-check secrets

help:
	@echo "  make web       Compila apps/web/dist (index.html autocontenido + fuentes)"
	@echo "  make web-dev   Sirve apps/web/src en :8000 (modo demo, sin API)"
	@echo "  make fonts     Recorta las fuentes de apps/web/fonts-src a src/fonts"
	@echo "  make api       Compila el portal y levanta la API en :8080"
	@echo "  make build     Construye la imagen"
	@echo "  make up        Levanta el contenedor (construye si hace falta)"
	@echo "  make down      Para el contenedor"
	@echo "  make logs      Sigue los logs"
	@echo ""
	@echo "  make install   pnpm install (instala también los hooks de git) y herramientas de Go"
	@echo "  make check     Todo lo que valida CI salvo la imagen (lo corre el hook pre-push)"
	@echo "  make fmt       Formatea JS, CSS, JSON, Markdown, YAML y Go"
	@echo "  make go-check  Solo la API: lint, formato, módulos, tests y vulnerabilidades"
	@echo "  make web-check Solo el portal: lint, formato, ES2020, OpenAPI y build con presupuesto"

web:
	node apps/web/build.mjs

# Volver a correr si cambian los textos del portal
fonts:
	pnpm fonts

web-dev:
	python3 -m http.server -d apps/web/src 8000

# La API sirve apps/web/dist, así que se compila el portal antes
api: web
	cd apps/api && go run .

build: web
	$(DC) build

up: web
	$(DC) up -d --build

down:
	$(DC) down

logs:
	$(DC) logs -f

# ---------- calidad ----------

install: tools
	pnpm install

tools:
	cd tools && go build -o $(BIN)/ github.com/golangci/golangci-lint/v2/cmd/golangci-lint golang.org/x/vuln/cmd/govulncheck
	cd tools/actionlint && go build -o $(BIN)/ github.com/rhysd/actionlint/cmd/actionlint

check: web-check go-check shell-lint docker-lint actions-lint tf-check
	@echo "✔ check completo"

fmt:
	pnpm format
	cd apps/api && $(BIN)/golangci-lint fmt ./...
	terraform fmt -recursive infra

web-check:
	pnpm check

go-check: go-mod go-fmt go-lint go-test go-vuln

go-mod:
	cd apps/api && go mod tidy -diff && go mod verify
	cd tools && go mod tidy -diff
	cd tools/actionlint && go mod tidy -diff

go-fmt: tools
	cd apps/api && $(BIN)/golangci-lint fmt --diff ./...

go-lint: tools
	cd apps/api && go vet ./... && $(BIN)/golangci-lint run ./...

go-test:
	cd apps/api && go test -race -count=1 ./...

go-vuln: tools
	cd apps/api && $(BIN)/govulncheck ./...

shell-lint:
	$(call need,shellcheck,lint de scripts de shell)
	@if command -v shellcheck >/dev/null 2>&1; then shellcheck loadtest/*.sh && shellcheck -s sh -x .husky/pre-commit .husky/commit-msg .husky/pre-push .husky/post-checkout .husky/branch-name.sh; fi

docker-lint:
	$(call need,hadolint,lint del Dockerfile)
	@if command -v hadolint >/dev/null 2>&1; then hadolint docker/Dockerfile; fi

actions-lint: tools
	$(BIN)/actionlint

tf-check:
	$(call need,terraform,formato y validación de infra/local/terraform)
	@if command -v terraform >/dev/null 2>&1; then \
		terraform fmt -check -recursive infra && \
		terraform -chdir=infra/local/terraform init -backend=false -input=false >/dev/null && \
		terraform -chdir=infra/local/terraform validate; \
	fi

# Escanea todo el historial de git en busca de secretos verificados (CI lo hace en cada push)
secrets:
	$(call need,trufflehog,escaneo de secretos)
	@if command -v trufflehog >/dev/null 2>&1; then trufflehog git file://. --only-verified --fail; fi
