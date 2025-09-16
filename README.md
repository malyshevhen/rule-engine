# Rule Engine Microservice

A robust, scalable rule engine microservice for IoT automation built with Go. The service allows users to create and manage custom automation rules with Lua script execution in a secure sandboxed environment.

## Features

- **Rule Management**: Full CRUD operations for automation rules
- **Secure Lua Execution**: Sandboxed Lua script execution with platform API bindings
- **Trigger System**: Support for conditional and scheduled triggers
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
- Docker & Docker Compose (for local development)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd rule-engine
   ```

2. **Start PostgreSQL database**
   ```bash
   docker run -d \
     --name rule-engine-db \
     -e POSTGRES_DB=rule_engine \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=password \
     -p 5433:5432 \
     postgres:15-alpine
   ```

3. **Set environment variables**
   ```bash
   export DATABASE_URL="postgres://postgres:password@localhost:5433/rule_engine?sslmode=disable"
   export API_KEY="your-api-key-here"
   export JWT_SECRET="your-jwt-secret-here"
   ```

4. **Run database migrations**
   ```bash
   go run cmd/main.go migrate
   ```

5. **Start the service**
   ```bash
   go run cmd/main.go
   ```

The service will be available at `http://localhost:8080`

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

### Authentication

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
- `GET /api/v1/rules` - List all rules
- `GET /api/v1/rules/{id}` - Get rule by ID
- `PUT /api/v1/rules/{id}` - Update rule
- `DELETE /api/v1/rules/{id}` - Delete rule

#### Triggers

- `POST /api/v1/triggers` - Create a new trigger
- `GET /api/v1/triggers` - List all triggers
- `GET /api/v1/triggers/{id}` - Get trigger by ID

#### Actions

- `POST /api/v1/actions` - Create a new action
- `GET /api/v1/actions` - List all actions
- `GET /api/v1/actions/{id}` - Get action by ID

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

### Running Tests

```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# Performance tests
go test -tags=performance ./api
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter
golangci-lint run
```

### Database Migrations

```bash
# Run migrations
go run cmd/main.go migrate

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
docker build -f containers/Containerfile -t rule-engine:latest .

# Run in production
docker run -d \
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