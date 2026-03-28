# Gin Post Service

This directory contains a parallel Go implementation of the existing post system under `/api/community`. It is intended to run alongside the Flask service during shadow, canary, and rollback-safe migration.

## Commands

- `go run ./cmd/server`: start the Gin service locally.
- `go test ./... -cover`: run unit and handler tests.
- `go test ./... -bench .`: run micro-benchmarks.

## Environment

- `DATABASE_URL`: SQLAlchemy-style MySQL DSN from the current Flask app.
- `READ_DATABASE_URL`: optional read-replica DSN for read/write split.
- `JWT_SECRET_KEY`: must match the current Flask JWT secret.
- `GO_POST_SERVICE_ADDR`: listen address, default `:8080`.

## Migration Notes

- Existing Flask routes stay online while this service is deployed behind a header or path-based canary rule.
- Cache invalidation is event-driven inside the process today; the `EventBus` boundary is where Kafka/NATS can be attached later without changing handlers.
- No existing WebSocket endpoints were found in the current post system, so only REST compatibility is implemented here.
