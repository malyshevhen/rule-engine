## Project Purpose

The Rule Engine is a robust, scalable microservice for IoT automation platforms. It enables users to create and manage custom automation rules through Lua scripts executed in a secure, sandboxed environment. The service handles:

- **Rule Management**: CRUD operations for automation rules via REST API
- **Trigger Evaluation**: Conditional triggers from NATS message bus and scheduled CRON triggers
- **Secure Lua Execution**: Sandboxed script execution with platform API bindings
- **Action Execution**: Performing operations based on rule evaluation results
- **Observability**: Structured logging, Prometheus metrics, and health checks

Key technologies: Go 1.24, PostgreSQL, NATS, Lua (gopher-lua), JWT authentication, Prometheus monitoring.