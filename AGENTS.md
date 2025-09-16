# Agent Guidelines for Rule Engine

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