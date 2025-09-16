## Code Style and Conventions

### Go Conventions
- Formatting: `gofmt` with 4-space indentation
- Package names: lowercase, single word when possible
- Function names: PascalCase for exported, camelCase for unexported
- Variable names: camelCase, descriptive and concise
- Error handling: explicit returns, no panics for recoverable errors
- Concurrency: goroutines/channels with race detection, context.Context for timeouts

### SQL Style
- Keywords: UPPERCASE
- Indentation: 4 spaces
- Tables: snake_case, plural (e.g., `execution_logs`)
- Columns: snake_case (e.g., `created_at`, `rule_id`)
- Timestamps: TIMESTAMPTZ
- Foreign keys: CASCADE delete where appropriate
- Enums: for status fields

### Naming Conventions
- Go structs: PascalCase (e.g., `ExecutionLog`)
- JSON fields: snake_case in tags
- Database tables: snake_case, plural
- Internal packages only (no external deps on internal)

### Architecture
- Standard Go Project Layout
- Clear separation: API, core business logic, storage
- Dependency injection pattern
- 12-Factor App: env vars for config

### Imports
- Standard library first, then third-party, then internal
- Group with blank lines
- Aliases for conflicts