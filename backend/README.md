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


### 帖子系统的实际业务逻辑在 internal/community/ 目录下，采用 三层架构 ：
handler.go HTTP 请求处理层 - 解析请求参数、调用 service、返回响应 service.go 业务逻辑层 - 核心业务规则、验证、限流、熔断 (gobreaker) repository.go 数据访问层 - MySQL 数据库读写操作



go_post_service 完整目录结构及分工如下
go_post_service/
├── cmd/server/main.go          # 程序入口 - 启动服务器、加载配置
│
├── internal/
│   ├── app/                    # 通用基础设施 (非业务相关)
│   │   ├── auth.go            # JWT 认证
│   │   ├── cache.go           # 缓存管理
│   │   ├── config.go          # 配置加载
│   │   ├── db.go              # 数据库连接
│   │   ├── errors.go          # 统一错误处理
│   │   ├── events.go          # 事件系统
│   │   ├── logger.go          # 日志
│   │   ├── metrics.go         # 指标监控
│   │   ├── middleware.go      # Gin 中间件 (限流、CORS 等)
│   │   └── response.go        # 统一响应格式
│   │
│   ├── community/             # 帖子系统业务模块
│   │   ├── handler.go         # HTTP 处理器 (路由注册、参数解析)
│   │   ├── service.go         # 业务逻辑 (CRUD、验证、熔断)
│   │   ├── repository.go      # 数据访问 (SQL 查询)
│   │   ├── service_test.go    # 业务层单元测试
│   │   └── repository_test.go # 数据层单元测试
│   │
│   └── server/                # 服务器框架
│       ├── router.go          # Gin 路由配置
│       └── router_test.go     # 路由测试
│
├── configs/app.env.example    # 环境变量配置示例
├── docs/migration-plan.md     # 迁移文档
├── AGENTS.md                  # Agent 规范
├── README.md                  # 项目说明
├── go.mod / go.sum            # Go 依赖
└── run.pid                    # 运行时 PID 文件