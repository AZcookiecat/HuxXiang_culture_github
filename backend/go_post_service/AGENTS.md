# Repository Guidelines

## Project Structure & Module Organization

This service is a small Go Gin backend for `/api/community`.

- `cmd/server/main.go`: application entrypoint and startup wiring.
- `internal/community`: business code in `handler`, `service`, and `repository` layers.
- `internal/app`: shared infrastructure such as config, logging, middleware, cache, auth, and metrics.
- `internal/server`: router setup and route registration.
- `configs/app.env.example`: local environment template.
- `docs/`: migration notes and supporting documentation.

Keep new business logic inside `internal/community`. Put cross-cutting utilities in `internal/app`.

## Build, Test, and Development Commands

- `go run ./cmd/server`: run the service locally.
- `go test ./...`: run all unit tests.
- `go test ./... -cover`: run tests with coverage.
- `gofmt -w ./...`: format Go files before committing.
- `go mod tidy`: sync module dependencies after adding imports.

Run commands from the repository root: `D:\project_hutb\HuXiang_culture\backend\go_post_service`.

## Coding Style & Naming Conventions

Use standard Go style and tabs via `gofmt`. Keep functions short and explicit.

- Exported names use `CamelCase`; unexported names use `camelCase`.
- Handlers should parse input and return responses only.
- Services should hold validation, normalization, cache, and orchestration logic.
- Repositories should contain SQL and database mapping only.

Prefer simple structs and helper functions over abstractions or design patterns.

## Testing Guidelines

Use Go’s built-in `testing` package. Repository tests currently use `github.com/DATA-DOG/go-sqlmock`.

- Name tests like `TestListPostsAppliesKeywordFilter`.
- Add tests for new query branches, validation rules, and error paths.
- Keep tests close to the package they cover, for example `internal/community/service_test.go`.

## Commit & Pull Request Guidelines

Recent history mixes short English and Chinese messages. Keep commits short, scoped, and descriptive.

- Good examples: `community: add category stats`, `app: unify error responses`
- Avoid vague messages like `update` or `fix`

PRs should include:

- a short summary of behavior changes
- affected endpoints or config keys
- test results, for example `go test ./...`
- sample request/response when API behavior changes

## Security & Configuration Tips

Do not commit real secrets or production DSNs. Use `configs/app.env.example` as the template. Keep `JWT_SECRET_KEY`, `DATABASE_URL`, and CORS settings environment-driven.
