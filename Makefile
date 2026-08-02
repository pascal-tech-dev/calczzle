.PHONY: help up down build logs ps cover cover-backend cover-frontend

COMPOSE ?= docker compose
GO_IMAGE ?= golang:1.26-alpine
NODE_IMAGE ?= node:22-alpine
UID := $(shell id -u)
GID := $(shell id -g)

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

up: ## Build and start the production stack (http://localhost)
	$(COMPOSE) up --build -d

down: ## Stop and remove the production stack
	$(COMPOSE) down

build: ## Build production images without starting
	$(COMPOSE) build

logs: ## Follow compose logs
	$(COMPOSE) logs -f

ps: ## Show compose service status
	$(COMPOSE) ps

cover: cover-backend cover-frontend ## Run backend + frontend coverage in ephemeral containers

cover-backend: ## Go tests + coverage → backend/coverage.out and coverage.html
	docker run --rm \
		--user "$(UID):$(GID)" \
		-e GOCACHE=/tmp/go-cache \
		-e GOMODCACHE=/tmp/go-mod \
		-v "$(CURDIR)/backend:/src" \
		-w /src \
		$(GO_IMAGE) \
		sh -c 'go test ./... -coverprofile=coverage.out -covermode=atomic \
			&& echo \
			&& go tool cover -func=coverage.out | grep -v "^total:" \
			&& echo \
			&& go tool cover -func=coverage.out | grep "^total:" \
			&& go tool cover -html=coverage.out -o coverage.html'
	@echo "Wrote backend/coverage.out and backend/coverage.html"

cover-frontend: ## Vitest coverage → frontend/coverage/
	docker run --rm \
		--user "$(UID):$(GID)" \
		-e npm_config_cache=/tmp/npm-cache \
		-v "$(CURDIR)/frontend:/app" \
		-w /app \
		$(NODE_IMAGE) \
		sh -c 'npm ci && npm run test:coverage'
	@echo "Wrote frontend/coverage/"
