## Post-Task Checklist

After completing any code changes, always run:

1. **Tests**: `go test ./...` (unit tests) and `go test -tags=integration ./...` (integration)
2. **Linting**: `golangci-lint run` (must pass with zero issues)
3. **Formatting**: `gofmt -w .` and `go mod tidy`
4. **Vet**: `go vet ./...`
5. **SQL Lint**: `sqruff lint internal/storage/db/migrations/`
6. **Build**: `go build -o rule-engine cmd/main.go` (ensure compiles)
7. **Race Detection**: `go test -race ./...` for concurrent code

All checks must pass before considering the task complete. Target >85% test coverage.