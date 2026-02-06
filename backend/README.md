# 湖湘文化数字化平台 - 后端

湖湘文化数字化平台后端服务，基于 Flask 框架构建，提供 RESTful API 接口，用于管理和展示湖湘地区的文化资源。

## 项目概述

湖湘文化数字化平台致力于通过现代信息技术手段，对湖湘地区丰富的历史文化遗产、民俗风情、传统艺术等进行数字化展示与传播，推动文化的保护与传承。

## 技术栈

- **Python**: 3.x
- **Flask**: Web 框架
- **Flask-SQLAlchemy**: ORM 框架
- **Flask-JWT-Extended**: JWT 认证
- **Flask-CORS**: 跨域资源共享
- **PyMySQL**: MySQL 数据库驱动
- **Werkzeug**: WSGI 工具库

## 功能特性

- 用户注册与登录系统
- 文化资源管理（增删改查）
- 社区互动功能（帖子发布、评论、回复）
- 基于角色的权限控制（普通用户/管理员）
- 响应式 API 接口设计

## 目录结构

```
backend/
├── app/                    # 应用主模块
│   ├── __init__.py         # 应用工厂函数
│   ├── models/             # 数据模型
│   │   ├── user.py         # 用户模型
│   │   ├── cultural_resource.py  # 文化资源模型
│   │   └── community_post.py     # 社区帖子和评论模型
│   ├── routes/             # API 路由
│   │   ├── main.py         # 主页和健康检查路由
│   │   ├── cultural_resources.py # 文化资源相关路由
│   │   ├── auth.py         # 认证相关路由
│   │   └── community.py    # 社区相关路由
│   └── utils/              # 工具函数
├── config/                 # 配置文件
├── models/                 # 数据模型（备用）
├── routes/                 # API 路由（备用）
├── static/                 # 静态文件
├── templates/              # 模板文件
├── migrations/             # 数据库迁移文件
├── app.py                  # 应用启动文件
├── init_db.py              # 数据库初始化脚本
├── requirements.txt        # 项目依赖
└── README.md               # 项目说明文档
```

## 安装与配置

### 1. 环境准备

确保系统已安装 Python 3.x 和 pip。

### 2. 克隆项目

```bash
git clone <repository-url>
cd backend
```

### 3. 创建虚拟环境并安装依赖

```bash
# 创建虚拟环境
python -m venv venv

# 激活虚拟环境
# Windows:
venv\Scripts\activate
# macOS/Linux:
source venv/bin/activate

# 安装依赖
pip install -r requirements.txt
```

### 4. 配置数据库

创建 MySQL 数据库，或修改 [config.py](file:///c%3A/Users/Lenovo/Desktop/HuxXiang_culture_github/backend/config.py) 中的数据库配置：

```python
# 默认配置
SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL') or 'mysql+pymysql://root:password@localhost/huxiang_culture'
```

### 5. 初始化数据库

创建数据库表结构：

```bash
python init_db.py
```

## API 接口说明

### 认证接口

- `POST /api/auth/register` - 用户注册
  - 请求体: `{"username": "用户名", "email": "邮箱", "password": "密码"}`
  - 响应: `{"success": true, "message": "注册成功", "user": {...}}`

- `POST /api/auth/login` - 用户登录（支持用户名或邮箱登录）
  - 请求体: `{"username": "用户名或邮箱", "password": "密码"}` 或 `{"usernameOrEmail": "用户名或邮箱", "password": "密码"}`
  - 响应: `{"success": true, "message": "登录成功", "access_token": "...", "user": {...}}`

- `GET /api/auth/profile` - 获取用户信息（需要认证）
  - 响应: `{"success": true, "data": {...}}`

- `PUT /api/auth/profile` - 更新用户信息（需要认证）
  - 请求体: `{"bio": "简介", "avatar": "头像URL"}`
  - 响应: `{"success": true, "message": "资料更新成功"}`

- `POST /api/auth/logout` - 用户登出（需要认证）
  - 响应: `{"success": true, "message": "登出成功"}`

### 文化资源接口

- `GET /api/resources` - 获取文化资源列表
- `GET /api/resources/<id>` - 获取特定文化资源详情
- `POST /api/resources` - 创建新的文化资源（需要管理员权限）

### 社区接口

- `GET /api/community/posts` - 获取社区帖子列表
- `GET /api/community/posts/<id>` - 获取帖子详情
- `POST /api/community/posts` - 发布新帖子（需要认证）
- `POST /api/community/posts/<post_id>/comments` - 添加评论（需要认证）

## 环境变量配置

创建 `.env` 文件配置敏感信息：

```bash
SECRET_KEY=your-secret-key
JWT_SECRET_KEY=your-jwt-secret-key
DATABASE_URL=mysql+pymysql://username:password@localhost/database_name
```

## 运行项目

```bash
python app.py
```

应用将在 `http://localhost:5000` 上运行。

## 数据库连接测试

你可以使用以下方法测试数据库连接：

1. 访问健康检查接口：
```bash
curl http://localhost:5000/health
```

## 开发规范

- 使用 Flask 的应用工厂模式
- 所有 API 响应使用 JSON 格式
- 使用 JWT 进行身份验证
- 数据库操作使用 SQLAlchemy ORM
- 代码遵循 PEP 8 规范

## 部署

生产环境中建议使用 Gunicorn 或 uWSGI 部署 Flask 应用：

```bash
gunicorn app:app
```

配合 Nginx 作为反向代理服务器。

## 许可证

湖湘文化数字化平台采用 MIT 许可证。