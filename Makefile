# Rule Engine Makefile
# This Makefile provides convenient commands for development, testing, and deployment

.PHONY: help build run run-local migrate clean test test-integration test-performance test-race test-verbose test-specific lint format vet tidy sql-lint docs quality docker-build docker-run docker-compose-up docker-compose-down dev-up dev-down dev-logs dev-restart db-up db-wait db-down dev ci setup logs health metrics dashboard

# Default target
help: ## Show this help message
	@echo "🚀 Rule Engine Development Commands"
	@echo ""
	@echo "📦 Quick Start (Recommended):"
	@echo "  make dev-up          Start full development stack"
	@echo "  make dev-down        Stop development stack"
	@echo "  make dashboard       Open analytics dashboard"
	@echo ""
	@echo "🔧 Development Commands:"
	@grep -E '^(run|run-local|migrate|dev-|test|lint|format|vet|tidy|quality|clients|docs|prepare-e2e):.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "🐳 Docker Commands:"
	@grep -E '^(docker-|dev-):.*?## .*$$' $(MAKEFILE_LIST) | grep -v "dev-up\|dev-down\|dev-logs\|dev-restart" | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "🗄️  Legacy Database Commands:"
	@grep -E '^(db-):.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "🔍 Utility Commands:"
	@grep -E '^(logs|health|metrics|dashboard|setup|clean|build|docs-clean):.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Build commands
build: ## Build the application binary
	go build -o rule-engine cmd/main.go

run: ## Run the application (requires environment variables to be set)
	go run cmd/main.go

run-local: ## Run the application locally with default development settings
	@echo "Starting rule engine with local development configuration..."
	@export DATABASE_URL="postgres://postgres:password@localhost:5433/rule_engine?sslmode=disable" && \
	export API_KEY="dev-api-key-12345" && \
	export JWT_SECRET="dev-jwt-secret-67890" && \
	export PORT="8080" && \
	export NATS_URL="nats://localhost:4222" && \
	export REDIS_URL="localhost:6379" && \
	export ALERTING_ENABLED="false" && \
	go run cmd/main.go

migrate: ## Run database migrations
	@export DATABASE_URL="postgres://postgres:password@localhost:5433/rule_engine?sslmode=disable" && \
	export API_KEY="dev-api-key-12345" && \
	export JWT_SECRET="dev-jwt-secret-67890" && \
	export PORT="8080" && \
	export NATS_URL="nats://localhost:4222" && \
	export REDIS_URL="localhost:6379" && \
	export ALERTING_ENABLED="false" && \
	go run cmd/main.go migrate

clean: ## Clean build artifacts
	go clean
	rm -f rule-engine

# Testing commands
test: ## Run all unit tests
	go test ./...

test-integration: ## Run integration tests
	go test ./tests/e2e/...

test-performance: ## Run performance tests
	go test -tags=performance ./api

test-race: ## Run tests with race detection
	go test -race $(go list ./... | grep -v tests)

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-specific: ## Run specific test (usage: make test-specific TEST=TestName)
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-specific TEST=TestName"; \
		exit 1; \
	fi
	go test -run $(TEST) ./...

test-e2e: ## Run e2e tests (rebuilds Docker image first)
	@echo "Preparing for e2e tests..."
	make prepare-e2e
	go test -v ./tests/e2e/ -timeout=10m

test-e2e-only: ## Run e2e tests without rebuilding (assumes image is up to date)
	go test -v ./tests/e2e/ -timeout=10m

# Code quality commands
lint: ## Run linter (golangci-lint)
	golangci-lint run

vet: ## Vet code for potential issues
	go vet ./...

tidy: ## Clean and update Go module dependencies
	go mod tidy

sql-lint: ## Lint SQL migration files
	sqruff lint internal/storage/db/migrations/

docs: ## Generate OpenAPI/Swagger documentation
	swag init -g cmd/main.go -o docs/ --ot json,yaml

docs-clean: ## Clean generated documentation files
	rm -f docs/swagger.json docs/swagger.yaml

clients-go: ## Generate Go HTTP client from OpenAPI spec using openapi-generator-cli
	oapi-codegen -package client -generate types,client -o pkg/client/rule_engine.go docs/swagger.yaml

quality: format vet tidy sql-lint ## Run all code quality checks

# Container commands
docker-build: ## Build container (includes latest code changes)
	@echo "Building Docker image with latest changes..."
	docker build -f containers/Containerfile -t rule-engine:local .
	@echo "Docker image built successfully!"

docker-run: ## Run container
	docker run -d --name rule-engine -p 8080:8080 \
		-e DATABASE_URL="postgres://postgres:password@localhost:5433/rule_engine?sslmode=disable" \
		-e API_KEY="your-api-key-here" \
		-e JWT_SECRET="your-jwt-secret-here" \
		rule-engine:local

docker-compose-up: ## Start services with Docker Compose
	docker-compose -f containers/compose.yaml up

docker-compose-down: ## Stop services with Docker Compose
	docker-compose -f containers/compose.yaml down

# Local development stack
dev-up: ## Start local development stack with podman-compose
	podman-compose -f docker-compose.dev.yml up -d

dev-down: ## Stop local development stack
	podman-compose -f docker-compose.dev.yml down

dev-logs: ## Show logs from development stack
	podman-compose -f docker-compose.dev.yml logs -f

dev-restart: ## Restart development stack
	podman-compose -f docker-compose.dev.yml restart

# Legacy database commands (for manual setup)
db-up: ## Start local PostgreSQL database (legacy)
	podman run -d --name rule-engine-db \
		-e POSTGRES_DB=rule_engine \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-p 5433:5432 \
		postgres:15-alpine

db-wait: ## Wait for database to be ready (legacy)
	@echo "Waiting for database to be ready..."
	@until podman exec rule-engine-db pg_isready -U postgres -d rule_engine >/dev/null 2>&1; do \
		echo "Database not ready, waiting..."; \
		sleep 2; \
	done
	@echo "Database is ready!"

db-down: ## Stop local PostgreSQL database (legacy)
	podman stop rule-engine-db
	podman rm rule-engine-db

# Development workflow (legacy - use dev-up instead)
dev: db-up db-wait migrate run-local ## Start development environment (legacy - use dev-up instead)

# E2E test preparation
prepare-e2e: ## Prepare for e2e tests (rebuild everything)
	@echo "Preparing for e2e tests..."
	make docs
	@echo "E2E test preparation complete!"

# CI/CD simulation
ci:  ## Run CI pipeline locally
	gh act push

# All-in-one setup for new developers
setup: ## Initial project setup
	go mod download
	@echo "Project setup complete. Run 'make dev' to start development environment."

# Utility commands
logs: ## Show application logs (if running in container)
	docker logs -f rule-engine

health: ## Check application health
	curl -f http://localhost:8080/health || echo "Health check failed"

metrics: ## Show application metrics
	curl -f http://localhost:8080/metrics || echo "Metrics endpoint failed"

dashboard: ## Open analytics dashboard in browser (requires app to be running)
	@echo "Analytics dashboard available at: http://localhost:8080/dashboard"
	@echo "API documentation generated in docs/ directory (run 'make docs' to regenerate)"
	@echo ""
	@echo "Opening dashboard..."
	@which xdg-open >/dev/null 2>&1 && xdg-open http://localhost:8080/dashboard || \
	which open >/dev/null 2>&1 && open http://localhost:8080/dashboard || \
	echo "Please open http://localhost:8080/dashboard in your browser"
	@echo ""
	@echo "To view API docs: check docs/swagger.json or docs/swagger.yaml"
