# Repository Guidelines

## Project Structure & Module Organization
This repository combines a Vue 3 frontend with a split backend. Frontend source lives in `src/`: page views in `src/views/`, shared components in `src/components/`, routing in `src/router/`, and API helpers in `src/services/`. Static assets are stored in `public/` and `src/assets/`. Flask code lives in `backend/`, with app setup in `backend/app.py` and `backend/app/__init__.py`, models in `backend/models/`, and Flask routes in `backend/routes/`. Community post APIs are served separately by the Go service in `backend/go_post_service/internal/`.

## Build, Test, and Development Commands
- `npm install`: install frontend dependencies.
- `npm run dev`: start the Vite frontend dev server.
- `npm run build`: create a production frontend build.
- `cd backend && pip install -r requirements.txt`: install Flask dependencies.
- `cd backend && python init_db.py`: create Flask tables and seed default data.
- `cd backend && python app.py`: run the Flask API on `http://127.0.0.1:5000`.
- `cd backend/go_post_service && go test ./...`: run Go unit and router tests.
- `cd backend/go_post_service && go run ./cmd/server`: run the community service on `:8080`.

## Coding Style & Naming Conventions
Use 2-space indentation in Vue, JavaScript, and CSS. Prefer PascalCase for Vue single-file components such as `PostDetailPage.vue`. Keep service, route, and helper modules lowercase, for example `src/services/api.js` and `backend/routes/main.py`. Follow existing Flask and Go idioms instead of introducing new patterns. No formatter or linter config is checked in, so keep edits minimal and consistent with nearby code.

## Testing Guidelines
There is no unified frontend or Flask test suite yet. For changes, run `npm run build` and `go test ./...` at minimum, then manually smoke-test login, cultural resources, and community flows. Add Go tests next to the affected package, for example `internal/server/router_test.go`. If Python tests are added later, use `test_<feature>.py`.

## Commit & Pull Request Guidelines
Recent commits use short, direct messages in Chinese or English, such as `backend: add auth guard` or `完成Gin项目重构,实现限流,熔断等功能`. Keep each commit scoped to one change. Pull requests should include a short summary, affected paths, local verification steps, linked issues when available, and screenshots for visible UI updates.

## Security & Configuration Tips
Do not commit secrets, production `.env` files, or real database credentials. During development, keep `/api/community` routed to the Go service and other `/api` endpoints routed to Flask.
