## Codebase Structure

Follows Standard Go Project Layout:

- `cmd/main.go` - Application entry point
- `api/` - HTTP handlers, middleware, DTOs, routing
- `internal/core/` - Business logic domain models and services
  - `action/`, `rule/`, `trigger/` - Domain entities
  - `manager/` - Orchestration logic
- `internal/engine/executor/` - Lua script execution engine
  - `context/` - Execution context
  - `platform/` - Platform API bindings
- `internal/storage/` - Data persistence layer
  - `db/` - Database connection and migrations
  - `action/`, `rule/`, `trigger/` - Repositories
- `internal/metrics/` - Prometheus metrics
- `pkg/logger/` - Shared logging package
- `containers/` - Docker configuration
- `lua/` - Lua type stubs
- `docs/` - Generated API documentation