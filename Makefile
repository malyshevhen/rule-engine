# Rule Engine Makefile
# This Makefile provides convenient commands for development, testing, and deployment

.PHONY: help build run migrate clean test test-integration test-performance test-race lint format vet tidy sql-lint docker-build docker-run docker-compose-up db-up db-down

# Default target
help: ## Show this help message
	@echo "Rule Engine Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

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
	go test -tags=integration ./...

test-performance: ## Run performance tests
	go test -tags=performance ./api

test-race: ## Run tests with race detection
	go test -race ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-specific: ## Run specific test (usage: make test-specific TEST=TestName)
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-specific TEST=TestName"; \
		exit 1; \
	fi
	go test -run $(TEST) ./...

# Code quality commands
lint: ## Run linter (golangci-lint)
	golangci-lint run

format: ## Format code with gofmt
	gofmt -w .

vet: ## Vet code for potential issues
	go vet ./...

tidy: ## Clean and update Go module dependencies
	go mod tidy

sql-lint: ## Lint SQL migration files
	sqruff lint internal/storage/db/migrations/

docs: ## Generate OpenAPI/Swagger documentation
	swag init -g cmd/main.go -o docs/

quality: format vet tidy sql-lint ## Run all code quality checks

# Docker commands
docker-build: ## Build Docker container
	docker build -f containers/Containerfile -t rule-engine .

docker-run: ## Run Docker container
	docker run -d --name rule-engine -p 8080:8080 \
		-e DATABASE_URL="postgres://postgres:password@localhost:5433/rule_engine?sslmode=disable" \
		-e API_KEY="your-api-key-here" \
		-e JWT_SECRET="your-jwt-secret-here" \
		rule-engine

docker-compose-up: ## Start services with Docker Compose
	docker-compose -f containers/compose.yaml up

docker-compose-down: ## Stop services with Docker Compose
	docker-compose -f containers/compose.yaml down

# Database commands
db-up: ## Start local PostgreSQL database
	docker run -d --name rule-engine-db \
		-e POSTGRES_DB=rule_engine \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-p 5433:5432 \
		postgres:15-alpine

db-wait: ## Wait for database to be ready
	@echo "Waiting for database to be ready..."
	@until docker exec rule-engine-db pg_isready -U postgres -d rule_engine >/dev/null 2>&1; do \
		echo "Database not ready, waiting..."; \
		sleep 2; \
	done
	@echo "Database is ready!"

db-down: ## Stop local PostgreSQL database
	docker stop rule-engine-db
	docker rm rule-engine-db

# Development workflow
dev: db-up db-wait migrate run-local ## Start development environment (DB + migrations + app)

# CI/CD simulation
ci: quality test test-integration lint ## Run CI pipeline locally

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
	@echo "API documentation available at: http://localhost:8080/swagger/index.html"
	@which xdg-open >/dev/null 2>&1 && xdg-open http://localhost:8080/dashboard || \
	which open >/dev/null 2>&1 && open http://localhost:8080/dashboard || \
	echo "Please open http://localhost:8080/dashboard in your browser"