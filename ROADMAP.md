# Rule Engine Development Roadmap

This roadmap outlines the step-by-step development plan for building a robust, scalable rule engine microservice for the IoT platform. The service will manage custom automation rules, execute Lua scripts securely, and handle triggers from events and schedules.

## Current State

- Database schema is fully defined in `internal/storage/db/migrations/initial_schema.sql`
- Project structure follows Go standards
- All Go code files contain placeholder TODO comments
- Basic dependencies configured in `go.mod`

## Phase 1: Core Infrastructure Setup

### 1.1 Application Bootstrap

- [ ] Implement main entry point in `cmd/main.go`
- [ ] Create app initialization in `cmd/app/rule_engine.go`
- [ ] Set up configuration management (environment variables)
- [ ] Add graceful shutdown handling

### 1.2 Database Layer

- [ ] Implement PostgreSQL connection pool in `internal/storage/db/postgres.go`
- [ ] Add database migrations runner
- [ ] Create database health checks

### 1.3 Basic Models and Repositories

- [ ] Define storage models:
  - [ ] `internal/storage/rule/model.go`
  - [ ] `internal/storage/trigger/model.go`
  - [ ]`internal/storage/action/model.go`
- [ ] Implement repositories:
  - [ ] `internal/storage/rule/repository.go`
  - [ ] `internal/storage/trigger/repository.go`
  - [ ] `internal/storage/action/repository.go`

## Phase 2: Core Business Logic

### 2.1 Domain Models

- [ ] Implement core domain models:
  - [ ] `internal/core/rule/model.go`
  - [ ] `internal/core/trigger/model.go`
  - [ ] `internal/core/action/model.go`

### 2.2 Business Services

- [ ] Implement core services:
  - [ ] `internal/core/rule/service.go`
  - [ ] `internal/core/trigger/service.go`
  - [ ] `internal/core/action/service.go`

### 2.3 Execution Engine

- [ ] Define execution context model in `internal/engine/executor/context/model.go`
- [ ] Implement execution context service in `internal/engine/executor/context/service.go`
- [ ] Build Lua script executor in `internal/engine/executor/service.go`
- [ ] Implement secure Lua sandboxing (disable io, os, networking)
- [ ] Add execution timeout handling
- [ ] Create Lua API bindings for platform interaction

## Phase 3: API Layer

### 3.1 API Infrastructure

- [ ] Define DTOs in `api/dto.go` (Rule, Trigger, Action DTOs)
- [ ] Implement request helpers in `api/request.go`
- [ ] Implement response helpers in `api/response.go`
- [ ] Add middleware in `api/middleware.go` (logging, authentication)

### 3.2 HTTP Server and Routing

- [ ] Set up HTTP server in `api/server.go`
- [ ] Implement router in `api/router.go` with endpoints:
  - [ ] Rules: POST /rules, GET /rules, GET /rules/:id, PUT /rules/:id, DELETE /rules/:id
  - [ ] Triggers: POST /triggers, GET /triggers, GET /triggers/:id, PUT /triggers/:id, DELETE /triggers/:id
  - [ ] Actions: POST /actions, GET /actions, GET /actions/:id, PUT /actions/:id, DELETE /actions/:id

## Phase 4: Trigger System

### 4.1 Conditional Triggers

- [ ] Integrate NATS message bus client
- [ ] Implement event listener for conditional triggers
- [ ] Add event parsing and condition evaluation

### 4.2 Scheduled Triggers

- [ ] Implement CRON scheduler
- [ ] Add scheduled trigger execution logic

## Phase 5: Observability and Security

### 5.1 Logging and Monitoring

- [ ] Implement structured logging (JSON format) using `slog`
- [ ] Add Prometheus metrics (execution counts, latency, error rates)
- [ ] Integrate OpenTelemetry for distributed tracing

### 5.2 Security Enhancements

- [ ] Add JWT/API key authentication
- [ ] Implement input sanitization
- [ ] Add rate limiting
- [ ] Secure secret management (environment variables)

## Phase 6: Testing and Quality Assurance

### 6.1 Unit Testing

- [ ] Write unit tests for all models, repositories, and services
- [ ] Achieve >85% test coverage
- [ ] Add mock implementations for external dependencies

### 6.2 Integration Testing

- [ ] Implement integration tests for API endpoints
- [ ] Add database integration tests
- [ ] Test Lua script execution

### 6.3 Performance and Load Testing

- [ ] Benchmark Lua execution performance
- [ ] Load test API endpoints
- [ ] Optimize database queries

## Phase 7: Deployment and Operations

### 7.1 Containerization

- [ ] Optimize `Containerfile` for production
- [ ] Add health check endpoints
- [ ] Configure non-root user

### 7.2 CI/CD Pipeline

- [ ] Set up GitHub Actions for CI
- [ ] Add automated testing and linting
- [ ] Implement automated deployment

### 7.3 Documentation

- [ ] Generate OpenAPI/Swagger documentation
- [ ] Update README.md with setup and deployment instructions
- [ ] Document Lua API for users

## Phase 8: Advanced Features

### 8.1 Rule Dependencies and Orchestration

- [ ] Implement rule chaining and dependencies
- [ ] Add rule execution ordering
- [ ] Support complex trigger conditions

### 8.2 Performance Optimizations

- [ ] Add caching layer (Redis)
- [ ] Implement rule execution queuing
- [ ] Add horizontal scaling support

### 8.3 Monitoring and Alerting

- [ ] Add alerting for failed executions
- [ ] Implement execution analytics dashboard
- [ ] Add audit logging for rule changes

## Success Criteria

- [ ] All TODO comments resolved
- [ ] Comprehensive test suite with high coverage
- [ ] Secure Lua execution environment
- [ ] RESTful API with full CRUD operations
- [ ] Support for conditional and scheduled triggers
- [ ] Production-ready deployment configuration
- [ ] Complete observability stack
- [ ] Performance benchmarks met

