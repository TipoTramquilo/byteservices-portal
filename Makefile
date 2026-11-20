# Motor de contenedores: podman si está, si no docker (forzar con `make ENGINE=docker`)
ENGINE ?= $(shell command -v podman 2>/dev/null || command -v docker 2>/dev/null)
ENGINE := $(if $(strip $(ENGINE)),$(strip $(ENGINE)),docker)

DC = "$(ENGINE)" compose --env-file .env -f docker/docker-compose.yaml

.PHONY: help web web-dev fonts api build up down logs

help:
	@echo "  make web       Compila apps/web/dist (index.html autocontenido + fuentes)"
	@echo "  make web-dev   Sirve apps/web/src en :8000 (modo demo, sin API)"
	@echo "  make fonts     Recorta las fuentes de apps/web/fonts-src a src/fonts (npx subset-font)"
	@echo "  make api       Compila el portal y levanta la API en :8080"
	@echo "  make build     Construye la imagen"
	@echo "  make up        Levanta el contenedor (construye si hace falta)"
	@echo "  make down      Para el contenedor"
	@echo "  make logs      Sigue los logs"

web:
	node apps/web/build.mjs

# Volver a correr si cambian los textos del portal; subset-font llega con npx y no queda en el repo
fonts:
	npx -y -p subset-font -- node apps/web/subset-fonts.mjs

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
