# Repository Guidelines

## Project Structure & Module Organization
This repository is split into a Vue 3 frontend and a Flask backend. Frontend source lives in `src/`, with page views in `src/views/`, shared UI in `src/components/`, router setup in `src/router/`, and API helpers in `src/services/`. Static assets are under `public/` and `src/assets/`. Backend code lives in `backend/`, with routes in `backend/routes/`, SQLAlchemy models in `backend/models/`, app bootstrap in `backend/app.py`, and database initialization in `backend/init_db.py`.

## Build, Test, and Development Commands
- `npm install`: install frontend dependencies.
- `npm run dev`: start Vite dev server with `/api` proxied to `http://127.0.0.1:5000`.
- `npm run build`: produce the frontend production bundle.
- `cd backend && pip install -r requirements.txt`: install Flask dependencies.
- `cd backend && python init_db.py`: create or refresh database tables and seed data.
- `cd backend && python app.py`: run the backend locally on port `5000`.

## Coding Style & Naming Conventions
Use 2-space indentation in Vue, JavaScript, and CSS. Prefer Vue SFCs with PascalCase filenames such as `CommunityPage.vue` and `CommentsSection.vue`. Keep route and service files in lowercase where already established, for example `src/services/api.js` and `backend/routes/community.py`. Follow existing alias usage (`@/services/api.js`) on the frontend and PEP 8-style naming on the backend. No lint or formatter config is checked in, so keep edits minimal and consistent with surrounding code.

## Testing Guidelines
There is no standardized automated test suite configured yet. The existing `backend/models/test_belong_hai.py` is not a full project test harness. Before opening a PR, at minimum run `npm run build`, start `python app.py`, and manually smoke-test key flows such as login, community posts, and resource detail pages. If you add tests, place backend tests under `backend/tests/test_<feature>.py`.

## Commit & Pull Request Guidelines
Recent history uses short, direct messages in Chinese or English and frequent merge commits. Prefer concise, scoped commits such as `frontend: fix community pagination` or `backend: add post like guard`. For PRs, include a short summary, affected areas (`src/`, `backend/routes/`, etc.), manual verification steps, and screenshots for UI changes. Link related issues when available.

## Security & Configuration Tips
Do not commit secrets or database credentials. Keep environment-specific values in local config or `.env` files, and verify the Vite proxy target and backend database settings before running locally.
