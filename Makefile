COMPOSE_FILE := deploy/docker-compose.yml
PROJECT_NAME := fullstack-engineering-lab

.PHONY: up down logs restart dev build clean test help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## Start all services with Docker Compose
	docker compose -f $(COMPOSE_FILE) --env-file .env up -d

down: ## Stop all services
	docker compose -f $(COMPOSE_FILE) down

logs: ## Show logs for all services
	docker compose -f $(COMPOSE_FILE) logs -f

restart: ## Restart all services
	docker compose -f $(COMPOSE_FILE) restart

dev: ## Start local development (without Docker)
	@bash scripts/dev.sh

build: ## Build all Docker images
	docker compose -f $(COMPOSE_FILE) build

clean: ## Stop services and remove volumes
	docker compose -f $(COMPOSE_FILE) down -v --remove-orphans

test: ## Run all tests
	cd apps/server && go test ./...
	cd apps/web && npm run type-check 2>/dev/null || true

init: ## Initialize project (copy .env, check dependencies)
	@bash scripts/init.sh

build-server: ## Build Go server binary
	cd apps/server && go build -o bin/server ./cmd/server/

build-web: ## Build frontend
	cd apps/web && npm install && npm run build

build-docs: ## Build docs site
	cd apps/docs && npm install && npm run build
