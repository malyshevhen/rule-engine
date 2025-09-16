## Essential Development Commands

### Build & Run
- `go build -o rule-engine cmd/main.go` - Build binary
- `go run cmd/main.go` - Run service
- `go run cmd/main.go migrate` - Run DB migrations
- `go clean` - Clean build artifacts

### Testing
- `go test ./...` - Run all unit tests
- `go test -tags=integration ./...` - Run integration tests
- `go test -tags=performance ./api` - Run performance tests
- `go test -run TestName ./path/to/package` - Run specific test
- `go test -v ./...` - Verbose test output
- `go test -race ./...` - Run with race detection

### Code Quality
- `gofmt -w .` - Format code
- `go vet ./...` - Vet code for issues
- `go mod tidy` - Clean dependencies
- `golangci-lint run` - Run linter (must pass)
- `sqruff lint internal/storage/db/migrations/` - Lint SQL migrations

### Docker
- `docker build -f containers/Containerfile -t rule-engine .` - Build container
- `docker-compose -f containers/compose.yaml up` - Run with compose

### Database
- Local Postgres: `docker run -d --name rule-engine-db -e POSTGRES_DB=rule_engine -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5433:5432 postgres:15-alpine`

### Documentation
- Swagger UI: `http://localhost:8080/swagger/` (when running)
- Health check: `GET /health`
- Metrics: `GET /metrics`