# Rule Engine Microservice

A robust, scalable rule engine microservice for IoT automation built with Go. The service allows users to create and manage custom automation rules with Lua script execution in a secure sandboxed environment.

## Features

- **Rule Management**: Full CRUD operations for automation rules
- **Secure Lua Execution**: Sandboxed Lua script execution with platform API bindings
- **Trigger System**: Support for conditional and scheduled triggers
- **Analytics Dashboard**: Real-time metrics visualization with historical trends and rule performance insights
- **RESTful API**: Complete REST API with OpenAPI/Swagger documentation
- **Authentication**: JWT and API key authentication
- **Observability**: Structured logging, Prometheus metrics, and health checks
- **Performance**: Comprehensive performance testing and optimization with Redis caching
- **Container Ready**: Production-ready Docker containerization

## Architecture

The service follows a clean architecture with clear separation of concerns:

- **API Layer**: HTTP handlers, middleware, and DTOs
- **Core Layer**: Business logic and domain models
- **Storage Layer**: Database repositories and migrations
- **Engine Layer**: Lua script execution and platform bindings

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 13+
- Podman & Podman Compose (for local development)

### Local Development Setup

1. **Clone the repository**
    ```bash
    git clone <repository-url>
    cd rule-engine
    ```

2. **Start the full development stack**
    ```bash
    make dev-up
    ```
    This starts all services using Podman Compose:
    - **PostgreSQL** database (port 5433)
    - **Redis** cache (port 6379)
    - **NATS** message bus (port 4222)
    - **Rule Engine** app (port 8080)

    The app will automatically run migrations and be ready to use!

3. **Access the application**
    - **API**: http://localhost:8080
    - **Analytics Dashboard**: http://localhost:8080/dashboard
    - **API Documentation**: http://localhost:8080/swagger/
    - **Health Check**: http://localhost:8080/health

4. **Stop the development stack**
    ```bash
    make dev-down
    ```

### Alternative Manual Setup

If you prefer manual setup or need more control:

```bash
# Start individual services
make db-up          # PostgreSQL only
make run-local      # Run app locally (requires DB to be running)

# Or use legacy workflow
make dev           # Manual DB + migrations + app
```

**Manual Podman Commands:**
```bash
# Start PostgreSQL manually
podman run -d \
  --name rule-engine-db \
  -e POSTGRES_DB=rule_engine \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5433:5432 \
  postgres:15-alpine
```

The service will be available at `http://localhost:8080`

**Access Points:**
- **API**: `http://localhost:8080/api/v1/`
- **API Documentation**: `http://localhost:8080/swagger/`
- **Analytics Dashboard**: `http://localhost:8080/dashboard`
- **Health Check**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:8080/metrics`

### Docker Deployment

1. **Build the container**
   ```bash
   docker build -f containers/Containerfile -t rule-engine .
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose -f containers/compose.yaml up
   ```

## API Documentation

### Swagger UI

Access the interactive API documentation at:
```
http://localhost:8080/swagger/
```

### Analytics Dashboard

The service includes a built-in analytics dashboard for monitoring rule execution metrics and performance:

**Dashboard URL:**
```
http://localhost:8080/dashboard
```

**Features:**
- Real-time execution statistics (total, successful, failed executions)
- Success rate and average latency metrics
- Historical trend charts for executions, success rates, and latency
- Rule-specific performance analytics
- Time range filtering (1 hour, 24 hours, 7 days, 30 days)
- Auto-refreshing data with 30-second intervals

**API Endpoint:**
- `GET /api/v1/analytics/dashboard?timeRange=24h` - Returns JSON data for dashboard metrics

**Dashboard Data Structure:**
```json
{
  "overall_stats": {
    "total_executions": 1250,
    "successful_executions": 1187,
    "failed_executions": 63,
    "success_rate": 94.96,
    "average_latency_ms": 45.2
  },
  "rule_stats": [
    {
      "rule_id": "rule-uuid",
      "rule_name": "Temperature Alert",
      "total_executions": 450,
      "successful_executions": 432,
      "failed_executions": 18,
      "success_rate": 96.0,
      "average_latency_ms": 42.1,
      "last_executed": "2024-01-01T12:00:00Z"
    }
  ],
  "execution_trend": {
    "metric": "executions_per_hour",
    "data": [
      {"timestamp": "2024-01-01T10:00:00Z", "value": 25},
      {"timestamp": "2024-01-01T11:00:00Z", "value": 30}
    ]
  },
  "success_rate_trend": {
    "metric": "success_rate_percent",
    "data": [...]
  },
  "latency_trend": {
    "metric": "average_latency_ms",
    "data": [...]
  },
  "time_range": "24h"
}
```

## Generated API Clients

The Rule Engine provides automatically generated REST API clients for Go and Python based on the OpenAPI specification.

### Client Generation

To generate the latest API clients:

```bash
# Generate both Go and Python clients
./scripts/generate-clients.sh all

# Generate only Go client
./scripts/generate-clients.sh go

# Generate only Python client
./scripts/generate-clients.sh python
```

### Go Client

The Go client is generated using OpenAPI Generator and provides type-safe access to all API endpoints.

**Installation:**
```bash
# The client is generated in clients/go/
# Add it as a dependency in your project
go get github.com/malyshevhen/rule-engine/clients/go
```

**Usage Example:**
```go
package main

import (
    "context"
    "log"

    ruleengine "github.com/malyshevhen/rule-engine/clients/go"
)

func main() {
    config := ruleengine.NewConfiguration()
    config.Host = "localhost:8080"
    config.Scheme = "http"
    config.AddDefaultHeader("Authorization", "ApiKey your-api-key")

    client := ruleengine.NewAPIClient(config)
    ctx := context.Background()

    // Create a rule
    createReq := ruleengine.ApiCreateRuleRequest{
        Name:      "Temperature Alert",
        LuaScript: "if event.temperature > 25 then return true end",
        Priority:  &[]int32{0}[0],
        Enabled:   &[]bool{true}[0],
    }

    rule, _, err := client.RulesAPI.RulesPost(ctx).ApiCreateRuleRequest(createReq).Execute()
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Created rule: %s", *rule.Name)
}
```

See `examples/go/main.go` for a complete usage example.

### Python Client

The Python client is generated using OpenAPI Generator and provides a Pythonic interface to the API.

**Installation:**
```bash
# Install the generated client
pip install ./clients/python

# Or install dependencies manually
pip install -r clients/python/requirements.txt
```

**Usage Example:**
```python
from rule_engine_client import ApiClient, Configuration
from rule_engine_client.api import RulesApi
from rule_engine_client.models import ApiCreateRuleRequest

# Configure client
config = Configuration()
config.host = "http://localhost:8080/api/v1"
config.api_key['Authorization'] = 'ApiKey your-api-key'

client = ApiClient(config)
rules_api = RulesApi(client)

# Create a rule
create_request = ApiCreateRuleRequest(
    name="Temperature Alert Rule",
    lua_script="if event.temperature > 25 then return true end",
    priority=0,
    enabled=True
)

rule = rules_api.rules_post(create_request)
print(f"Created rule: {rule.name} (ID: {rule.id})")
```

See `examples/python/example_usage.py` for a complete usage example.

### Authentication

Both clients support the same authentication methods as the REST API:

The API supports two authentication methods:

1. **API Key Authentication**
   ```
   Authorization: ApiKey your-api-key
   ```

2. **JWT Bearer Authentication**
   ```
   Authorization: Bearer your-jwt-token
   ```

### Core Endpoints

#### Rules

- `POST /api/v1/rules` - Create a new rule
- `GET /api/v1/rules?limit=50&offset=0` - List all rules (with pagination support)
- `GET /api/v1/rules/{id}` - Get rule by ID
- `PATCH /api/v1/rules/{id}` - Update rule (JSON Patch RFC 6902)
- `DELETE /api/v1/rules/{id}` - Delete rule

**Rule Update (PATCH) with JSON Patch:**

The PATCH endpoint supports JSON Patch (RFC 6902) for partial updates. Send a JSON array of patch operations:

```json
[
  {"op": "replace", "path": "/name", "value": "Updated Rule Name"},
  {"op": "replace", "path": "/lua_script", "value": "if event.temp > 30 then return true end"},
  {"op": "replace", "path": "/priority", "value": 10},
  {"op": "replace", "path": "/enabled", "value": false}
]
```

Supported operations: `add`, `remove`, `replace`, `test`

**Rules List Response with Pagination:**

```json
{
  "rules": [
    {
      "id": "uuid",
      "name": "Temperature Alert",
      "lua_script": "if event.temperature > 25 then return true end",
      "priority": 0,
      "enabled": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "limit": 50,
  "offset": 0,
  "count": 1
}
```

#### Triggers

- `POST /api/v1/triggers` - Create a new trigger
- `GET /api/v1/triggers` - List all triggers
- `GET /api/v1/triggers/{id}` - Get trigger by ID

#### Actions

- `POST /api/v1/actions` - Create a new action
- `GET /api/v1/actions` - List all actions
- `GET /api/v1/actions/{id}` - Get action by ID

#### Analytics

- `GET /api/v1/analytics/dashboard` - Get analytics dashboard data
- `GET /dashboard` - Analytics dashboard web interface

#### System

- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics

## Lua API

Rules and actions are defined using Lua scripts with access to a secure platform API. The API provides functions for device interaction, logging, and temporary data storage.

### Available Platform Functions

#### Device Management

**`get_device_state(device_id)`**
- **Parameters**: `device_id` (string) - The unique identifier of the device
- **Returns**: Table containing device state information, or `nil` and error message
- **Example**:
```lua
local state, err = get_device_state("temperature_sensor_1")
if err then
    log_message("error", "Failed to get device state: " .. err)
    return false
end

if state.online and state.temperature > 25 then
    -- Device is online and temperature is high
    return true
end
```

**`send_command(device_id, command, params)`**
- **Parameters**:
  - `device_id` (string) - The target device identifier
  - `command` (string) - The command to send
  - `params` (table, optional) - Additional command parameters
- **Returns**: `nil` on success, error message on failure
- **Example**:
```lua
-- Simple command
local err = send_command("ac_unit", "turn_on")
if err then
    log_message("error", "Failed to turn on AC: " .. err)
end

-- Command with parameters
local err = send_command("thermostat", "set_temperature", {
    temperature = 22,
    mode = "cool"
})
```

#### Logging

**`log_message(level, message)`**
- **Parameters**:
  - `level` (string) - Log level: "debug", "info", "warn", "error"
  - `message` (string) - The message to log
- **Returns**: Nothing
- **Example**:
```lua
log_message("info", "Temperature rule triggered")
log_message("debug", "Current temperature: " .. event.temperature)
log_message("error", "Failed to send command to device")
```

#### Time Functions

**`get_current_time()`**
- **Parameters**: None
- **Returns**: Current Unix timestamp (number)
- **Example**:
```lua
local now = get_current_time()
log_message("info", "Current timestamp: " .. now)
```

#### Data Storage

**`store_data(key, value)`**
- **Parameters**:
  - `key` (string) - Storage key
  - `value` (any) - Value to store (supports tables, strings, numbers, booleans)
- **Returns**: Nothing
- **Notes**: Data is stored per rule execution and automatically cleaned up
- **Example**:
```lua
-- Store sensor reading
store_data("last_temperature", event.temperature)
store_data("alert_count", 0)

-- Store complex data
store_data("device_info", {
    id = event.device_id,
    type = "temperature_sensor",
    last_reading = event.temperature
})
```

**`get_stored_data(key)`**
- **Parameters**: `key` (string) - Storage key
- **Returns**: Stored value, or `nil` if key doesn't exist
- **Example**:
```lua
local last_temp = get_stored_data("last_temperature")
if last_temp and event.temperature > last_temp + 5 then
    log_message("warn", "Temperature increased by more than 5 degrees")
end

local alert_count = get_stored_data("alert_count") or 0
alert_count = alert_count + 1
store_data("alert_count", alert_count)
```

### Rule Scripts

Rule scripts evaluate conditions and return a boolean indicating whether associated actions should execute.

**Return Values:**
- `true` - Execute associated actions
- `false` or `nil` - Skip actions

**Example Rule Scripts:**

```lua
-- Temperature threshold rule
if event.temperature > 25 then
    log_message("info", "High temperature detected: " .. event.temperature)
    return true  -- Execute actions
end
return false  -- Normal temperature, skip actions
```

```lua
-- Motion detection with cooldown
local last_motion = get_stored_data("last_motion_time")
local now = get_current_time()

if not last_motion or (now - last_motion) > 300 then  -- 5 minutes cooldown
    store_data("last_motion_time", now)
    log_message("info", "Motion detected after cooldown period")
    return true
end

return false
```

```lua
-- Device offline detection
local state, err = get_device_state(event.device_id)
if err then
    log_message("error", "Failed to get device state: " .. err)
    return false
end

if not state.online then
    log_message("warn", "Device " .. event.device_id .. " is offline")
    return true  -- Trigger offline alert
end

return false
```

### Action Scripts

Action scripts perform operations when rules evaluate to true. They don't need to return values.

**Example Action Scripts:**

```lua
-- Send notification
log_message("info", "Sending temperature alert")
local err = send_command("notification_service", "send_alert", {
    title = "High Temperature Alert",
    message = "Temperature exceeded threshold: " .. event.temperature,
    priority = "high",
    device_id = event.device_id
})
if err then
    log_message("error", "Failed to send notification: " .. err)
end
```

```lua
-- Control multiple devices
log_message("info", "Activating cooling system")

-- Turn on AC
local err1 = send_command("ac_unit", "turn_on")
if err1 then
    log_message("error", "Failed to turn on AC: " .. err1)
end

-- Adjust thermostat
local err2 = send_command("thermostat", "set_mode", {
    mode = "cool",
    temperature = 22
})
if err2 then
    log_message("error", "Failed to adjust thermostat: " .. err2)
end
```

### Security Notes

- Lua scripts run in a sandboxed environment with restricted access
- File system operations (`io`, `os`) are disabled
- Network operations are not allowed
- Scripts have a timeout to prevent infinite loops
- All platform API calls are logged for monitoring

## Configuration

The service is configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `REDIS_URL` | Redis connection string for caching | `localhost:6379` |
| `API_KEY` | API key for authentication | Required |
| `JWT_SECRET` | Secret for JWT token signing | Required |
| `PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |

## Development

### Makefile Commands

The project includes a comprehensive Makefile for development tasks:

```bash
# 🚀 Quick Start (Recommended)
make dev-up                # Start full development stack (PostgreSQL + Redis + NATS + App)
make dev-down              # Stop development stack
make dev-logs              # View logs from all services
make dashboard             # Open analytics dashboard in browser

# 🔧 Development
make run-local             # Run app locally with dev config
make migrate               # Run database migrations
make test                  # Run all unit tests
make test-integration      # Run integration tests
make quality               # Run all code quality checks

# 🐳 Container
make docker-build    # Build production container
make docker-run      # Run production container

# 🔍 Utilities
make health                # Check application health
make metrics               # Show Prometheus metrics
make clean                 # Clean build artifacts
```

### Running Tests

```bash
# Unit tests
make test
# Or: go test ./...

# Integration tests
make test-integration
# Or: go test -tags=integration ./...

# Performance tests
go test -tags=performance ./api
```

### Code Quality

```bash
# Run all quality checks
make quality

# Individual checks
make format    # Format code
make vet       # Vet code
make lint      # Run linter
```

### Database Migrations

```bash
# Run migrations
make migrate
# Or: go run cmd/main.go migrate

# Create new migration
# Add SQL files to internal/storage/db/migrations/
```

## Deployment

### Production Container

The production container includes:

- Non-root user execution
- Health checks
- Security hardening
- Minimal attack surface

```bash
# Build production image
podman build -f containers/Containerfile -t rule-engine:latest .

# Run in production
podman run -d \
  --name rule-engine \
  -p 8080:8080 \
  -e DATABASE_URL="..." \
  -e API_KEY="..." \
  -e JWT_SECRET="..." \
  rule-engine:latest
```

### Kubernetes Deployment

Use the included health checks and metrics endpoints for Kubernetes liveness/readiness probes.

## Monitoring

### Metrics

Prometheus metrics are available at `/metrics`:

- `rule_engine_requests_total` - Total HTTP requests
- `rule_engine_request_duration_seconds` - Request duration histogram
- `rule_engine_rules_executed_total` - Total rule executions
- `rule_engine_lua_execution_duration_seconds` - Lua execution time

### Health Checks

- `GET /health` - Basic health check
- Container health checks configured for orchestration

### Logging

Structured JSON logging with configurable levels:

```json
{
  "level": "info",
  "timestamp": "2024-01-01T12:00:00Z",
  "message": "Rule executed",
  "rule_id": "uuid",
  "execution_time_ms": 150
}
```

## Security

- **Lua Sandboxing**: Restricted execution environment
- **Input Validation**: Comprehensive input sanitization
- **Authentication**: JWT and API key support
- **Rate Limiting**: Configurable request rate limits
- **HTTPS Ready**: TLS termination support

## Performance

The service is optimized for high-throughput IoT workloads:

- **Concurrent Processing**: Handles thousands of concurrent requests
- **Efficient Lua Execution**: Fast script execution with caching
- **Redis Caching**: High-performance caching for rule and trigger data with automatic invalidation
- **Database Optimization**: Indexed queries and connection pooling
- **Load Testing**: Comprehensive performance test suite

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.