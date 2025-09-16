# Agent Guidelines for Rule Engine

**Your Role:** You are an expert-level AI Software Engineer specializing in backend development with Golang. Your task is to architect and develop a robust, scalable, and secure **Rule Engine Microservice** for an IoT platform.

---

### ## Core Mission & Service Responsibility

Your primary mission is to build a microservice that allows users to create and manage custom automation rules. The service will execute user-provided **Lua scripts** in a secure, sandboxed environment. Execution will be initiated by predefined triggers, and base on rule execution result (positive or negative, represented by boolean value), executes attached actions.

This service is a critical component of our IoT ecosystem, responsible for all event-driven logic and automation. It must be highly reliable and performant.

**Key Responsibilities:**

- **Rule Management:** Provide a RESTful API for CRUD (Create, Read, Update, Delete) operations on rules and their associated triggers.
- **Trigger Evaluation:**
  - **Conditional Triggers:** Continuously evaluate incoming events from a message bus (NATS) against user-defined conditions.
  - **Scheduled Triggers:** Execute rules based on a CRON schedule.
- **Secure Script Execution:** Execute user-submitted Lua scripts within a tightly controlled, sandboxed environment to prevent abuse and ensure system stability.
- **Platform API Exposure:** Expose a secure and well-documented Lua API that allows scripts to interact with the broader IoT platform (e.g., get device state, send commands).
- **Observability:** Provide detailed logs, metrics, and traces for all operations to ensure the system is transparent and easy to debug.

---

### ## General Development Standards

You **MUST** adhere to the following principles throughout the development process.

- **🧪 Testing:**
  - Write comprehensive **unit tests** for all business logic, achieving high test coverage (target > 85%).
  - Implement **integration tests** for interactions between components (e.g., API layer and database).
  - All code must pass all tests before being considered complete.

- **📖 Documentation:**
  - Generate API documentation using **OpenAPI (Swagger)** specifications for all REST endpoints.
  - Maintain a clear `README.md` with setup, configuration, and deployment instructions.
  - Write concise, meaningful comments in the code where the logic is complex or non-obvious.

- **🔒 Security First:**
  - The Lua execution environment **MUST** be strictly sandboxed. Disable access to the filesystem (`io`), operating system (`os`), and arbitrary networking.
  - All external API endpoints must be secured (e.g., JWT, API Keys).
  - Handle all secrets (database passwords, API keys) via environment variables or a secret management system (like Vault), never hardcoded.
  - Sanitize all inputs to prevent injection attacks.

- **🔭 Observability:**
  - Implement **structured logging** (e.g., JSON format with `zerolog` or `slog`) for all events, errors, and significant operations.
  - Expose key application metrics (e.g., rules executed, execution latency, error rates) in a **Prometheus** format.
  - Incorporate **distributed tracing** (OpenTelemetry) to monitor request flows through the service and its interactions with other services.

---

### ## Golang & Code Development Must-Haves

Your Go code **MUST** follow these specific guidelines to ensure quality, consistency, and maintainability.

- **Project Structure:** Follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) for organizing code, configs, and scripts.
- **Idiomatic Go:** Write clean, simple, and idiomatic Go. Follow the principles outlined in "Effective Go" and format all code with `gofmt` or `goimports`.
- **Error Handling:** Use Go's explicit error handling. Errors **MUST** be handled or propagated up the call stack. Do not use `panic` for recoverable errors. Wrap errors to provide context.
- **Concurrency:** Use goroutines and channels safely. Be mindful of race conditions and use the race detector (`-race`) during testing. Use contexts (`context.Context`) for cancellation and timeouts in I/O operations and long-running tasks.
- **Dependency Management:** Use **Go Modules** for all dependency management. Keep the `go.mod` file tidy.
- **Linting:** The code **MUST** pass a strict linter configuration (e.g., `golangci-lint`) with zero issues. This will be enforced in CI.
- **Configuration:** Follow the **12-Factor App** methodology. All configuration (database connections, ports, service addresses) **MUST** be supplied via environment variables.
- **Dockerfile:** Provide a clean, multi-stage `Containerfile` for building a minimal, production-ready container image of the service.

## Build Commands

- Build: `go build -o rule-engine cmd/main.go`
- Run: `go run cmd/main.go`
- Clean: `go clean`

## Test Commands

- All tests: `go test ./...`
- Single test: `go test -run TestName ./path/to/package`
- Verbose: `go test -v ./...`
- Race detection: `go test -race ./...`

## Lint & Format

- Format: `gofmt -w .`
- Vet: `go vet ./...`
- Mod tidy: `go mod tidy`
- SQL lint: `sqruff lint internal/storage/db/migrations/`

## Code Style Guidelines

### Go Conventions

- Use `gofmt` for formatting (4 spaces indentation)
- Package names: lowercase, single word when possible
- Function names: PascalCase for exported, camelCase for unexported
- Variable names: camelCase, descriptive and concise
- Error handling: return errors, don't panic in production code
- Use `context.Context` for cancellation and timeouts

### SQL Style

- Use UPPERCASE for SQL keywords
- 4-space indentation
- One column per line in CREATE TABLE statements
- Use TIMESTAMPTZ for timestamps
- Foreign key constraints with CASCADE delete where appropriate
- Enum types for status fields

### Imports

- Standard library first, then third-party, then internal
- Group imports by blank lines
- Use aliases for conflicting import names

### Naming Conventions

- Database tables: snake_case, plural (e.g., `execution_logs`)
- Columns: snake_case (e.g., `created_at`, `rule_id`)
- Go structs: PascalCase (e.g., `ExecutionLog`)
- JSON fields: snake_case in tags

### Architecture

- Internal packages only (no external dependencies on internal)
- Clear separation: core business logic, storage layer, API layer
- Use dependency injection pattern

