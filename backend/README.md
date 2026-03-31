# HuXiang Culture Backend

This backend now has a split responsibility:

- Flask backend in the repository root serves user/authentication and cultural resource APIs.
- Gin post service in `go_post_service/` serves all community post APIs under `/api/community`.

## Current Backend Structure

```text
backend/
|-- app/
|   `-- __init__.py
|-- go_post_service/
|   |-- cmd/
|   |-- configs/
|   |-- docs/
|   `-- internal/
|-- models/
|   |-- cultural_resource.py
|   |-- test_belong_hai.py
|   `-- user.py
|-- routes/
|   |-- auth.py
|   |-- cultural_resources.py
|   `-- main.py
|-- app.py
|-- config.py
|-- init_db.py
`-- requirements.txt
```

## API Ownership

### Flask

- `GET /`
- `GET /health`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/auth/profile`
- `PUT /api/auth/profile`
- `POST /api/auth/upload-avatar`
- `POST /api/auth/logout`
- `GET /api/resources/`
- `GET /api/resources/<id>`
- `POST /api/resources/`
- `PUT /api/resources/<id>`
- `DELETE /api/resources/<id>`
- `POST /api/resources/<id>/like`

### Gin

- All community post APIs under `/api/community`

See [go_post_service/README.md](D:\project_hutb\HuXiang_culture\backend\go_post_service\README.md) for the Gin service routes and configuration.

## Running The Services

### Flask backend

```bash
pip install -r requirements.txt
python init_db.py
python app.py
```

### Gin post service

```bash
cd go_post_service
go test ./...
go run ./cmd/server
```

## Configuration

### Flask

- `DATABASE_URL`
- `SECRET_KEY`
- `JWT_SECRET_KEY`
- `AVATAR_UPLOAD_PATH`

### Gin

- `DATABASE_URL`
- `READ_DATABASE_URL`
- `JWT_SECRET_KEY`
- `GO_POST_SERVICE_ADDR`
- `POST_CACHE_TTL`
- `POST_SERVICE_READ_TIMEOUT`
- `POST_SERVICE_WRITE_TIMEOUT`
- `POST_DB_MAX_OPEN_CONNS`
- `POST_DB_MAX_IDLE_CONNS`
- `POST_DB_CONN_MAX_LIFETIME`
- `POST_RATE_LIMIT_RPS`
- `POST_RATE_LIMIT_BURST`
- `POST_CORS_ALLOW_ORIGINS`
- `POST_ENABLE_HTTPS_REDIRECT`
- `POST_LOG_JSON`

## Migration Note

The legacy Flask community post implementation has been removed from this backend. If you run a reverse proxy or frontend dev proxy, route `/api/community` to the Gin service instead of Flask.
