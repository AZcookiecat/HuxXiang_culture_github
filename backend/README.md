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

- 用户注册与登录系统（支持用户名或邮箱登录）
- 文化资源管理（增删改查）
- 社区互动功能（帖子发布、评论、回复、点赞）
- 基于角色的权限控制（普通用户/管理员）
- 响应式 API 接口设计
- 分页功能支持
- 浏览量统计功能

## 目录结构

```
backend/
├── app/                    # 应用主模块
│   ├── __init__.py         # 应用工厂函数
│   └── models/             # 数据模型（已迁移至根目录/models）
├── config.py               # 配置文件
├── models/                 # 数据模型
│   ├── user.py             # 用户模型
│   ├── cultural_resource.py # 文化资源模型
│   └── community_post.py   # 社区帖子模型
├── routes/                 # API 路由
│   ├── main.py             # 主页和健康检查路由
│   ├── cultural_resources.py # 文化资源相关路由
│   ├── auth.py             # 认证相关路由
│   └── community.py        # 社区相关路由
├── app.py                  # 应用启动文件
├── init_db.py              # 数据库初始化脚本（包含点赞关联表创建）
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

#### 开发环境配置
开发环境默认使用 SQLite 数据库，数据文件存储在 `instance/huxiang_culture_dev.db`。
如需使用 SQLite，可在 [config.py](file:///c:/Users/Lenovo/Desktop/HuxXiang_culture_github/backend/config.py) 中修改配置：
```python
SQLALCHEMY_DATABASE_URI = os.environ.get('DEV_DATABASE_URL') or 'sqlite:///../instance/huxiang_culture_dev.db'
```

#### 生产环境配置
生产环境使用 MySQL 数据库，修改 [config.py](file:///c:/Users/Lenovo/Desktop/HuxXiang_culture_github/backend/config.py) 中的数据库配置：
```python
SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL') or 'mysql+pymysql://root:qq123123@localhost/huxiang_culture'
```

### 5. 初始化数据库

创建数据库表结构并添加初始数据：

```bash
python init_db.py
```

此脚本会创建所有必要的数据表，包括点赞关联表，以及管理员账户和示例数据。

## API 接口说明

### 认证接口

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录（支持用户名或邮箱登录）
- `GET /api/auth/profile` - 获取用户信息（需要认证）
- `PUT /api/auth/profile` - 更新用户信息（需要认证）
- `POST /api/auth/logout` - 用户登出

### 文化资源接口

- `GET /api/resources` - 获取文化资源列表
- `GET /api/resources/<id>` - 获取特定文化资源详情
- `POST /api/resources` - 创建新的文化资源（需要管理员权限）
- `PUT /api/resources/<id>` - 更新文化资源（需要管理员权限）
- `DELETE /api/resources/<id>` - 删除文化资源（需要管理员权限）

### 社区接口

- `GET /api/community/posts` - 获取社区帖子列表（支持分页）
- `GET /api/community/posts/<id>` - 获取帖子详情（增加浏览量）
- `POST /api/community/posts` - 发布新帖子（需要认证）
- `PUT /api/community/posts/<id>` - 编辑帖子（需要认证且为本人或管理员）
- `DELETE /api/community/posts/<id>` - 删除帖子（需要认证且为本人或管理员）
- `POST /api/community/posts/<id>/like` - 给帖子点赞/取消点赞（需要认证）
- `GET /api/community/posts/<post_id>/comments` - 获取帖子评论列表
- `POST /api/community/posts/<post_id>/comments` - 发表评论（需要认证）
- `DELETE /api/community/comments/<id>` - 删除评论（需要认证且为本人或管理员）
- `GET /api/community/posts/related/<post_id>` - 获取相关帖子推荐

## 环境变量配置

创建 `.env` 文件配置敏感信息：

```bash
SECRET_KEY=your-secret-key
JWT_SECRET_KEY=your-jwt-secret-key
DATABASE_URL=mysql+pymysql://username:password@localhost/database_name
DEV_DATABASE_URL=sqlite:///../instance/huxiang_culture_dev.db
AVATAR_UPLOAD_PATH=/path/to/avatar/uploads  # 生产环境头像上传路径（可选）
```

**部署环境必设变量**：
- `DATABASE_URL`: 生产环境数据库连接字符串
  - Linux: `export DATABASE_URL="mysql+pymysql://username:password@host:port/dbname"`
  - Windows: `set DATABASE_URL=mysql+pymysql://username:password@host:port/dbname`
- `SECRET_KEY`: 应用密钥，用于加密会话和JWT
  - Linux: `export SECRET_KEY="your-production-secret-key"`
  - Windows: `set SECRET_KEY=your-production-secret-key`

**部署环境可选变量**：
- `AVATAR_UPLOAD_PATH`: 指定头像文件存储路径，默认为项目根目录的public/static/avatars
  - Linux: `export AVATAR_UPLOAD_PATH="/path/to/avatar/uploads"`
  - Windows: `set AVATAR_UPLOAD_PATH=C:\path\to\avatar\uploads`
  - （若未设置则使用默认路径）

## 部署注意事项

### 1. 服务器环境准备
- 确保服务器已安装Python 3.10+、Node.js环境
- 配置MySQL数据库（使用utf8mb4字符集）
- 设置环境变量（DATABASE_URL、SECRET_KEY、AVATAR_UPLOAD_PATH）

### 2. 文件上传目录设置
- 如设置了AVATAR_UPLOAD_PATH，确保目录存在且有适当的写入权限
- Linux示例：`sudo mkdir -p /var/www/uploads/avatars && sudo chmod 755 /var/www/uploads/avatars`
- Windows示例：确保指定目录有写入权限

### 3. 数据库初始化
- 部署后运行 `python init_db.py` 初始化数据库表结构和示例数据
- 此脚本会自动创建所有必要的数据表，包括点赞关联表，无需单独执行

### 4. 服务部署
- 使用Gunicorn运行后端服务：`gunicorn -w 4 -b 0.0.0.0:5000 app:app`
- 前端使用Nginx提供静态文件服务
- 配置Nginx反向代理将 `/api/` 请求转发到后端

### 5. Nginx配置示例
``nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API请求代理到Flask后端
    location /api/ {
        proxy_pass http://127.0.0.1:5000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 头像等静态文件
    location /static/ {
        alias /var/www/uploads/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### 6. 进程管理
- 推荐使用systemd（Linux）或PM2（跨平台）管理进程
- 配置自动启动和故障恢复机制

### 7. 安全要求
- 生产环境禁用调试模式
- 不要在代码中硬编码敏感信息
- 定期轮换SECRET_KEY
- 限制文件上传类型和大小

### 8. 维护要求
- 定期备份数据库和上传文件目录
- 监控服务器磁盘空间使用情况
- 及时更新系统和软件包的安全补丁

## 运行项目

```bash
python app.py
```

应用将在 `http://localhost:5000` 上运行。

## 前后端联调配置

在前端开发环境中，需要配置代理将API请求转发到后端服务。在前端项目的 [vite.config.js](file:///c:/Users/Lenovo/Desktop/HuxXiang_culture_github/vite.config.js) 中添加：

```javascript
export default defineConfig({
  // ... 其他配置
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:5000',
        changeOrigin: true,
        secure: false
      }
    }
  }
})
```

## 数据模型

### 用户模型 (User)
- id: 用户唯一标识
- username: 用户名（唯一）
- email: 邮箱（唯一）
- password_hash: 加密后的密码
- bio: 个人简介
- avatar: 头像URL
- role: 用户角色（user/admin）
- created_at: 创建时间
- updated_at: 更新时间

### 文化资源模型 (CulturalResource)
- id: 资源唯一标识
- title: 标题
- description: 描述
- content: 详细内容
- type: 类型
- category: 分类
- tags: 标签
- image_url: 图片URL
- author: 作者
- status: 状态（draft/published）
- created_at: 创建时间
- updated_at: 更新时间
- views: 浏览量

### 社区帖子模型 (CommunityPost)
- id: 帖子唯一标识
- title: 标题
- content: 内容
- author_id: 作者ID（关联用户）
- category: 分类
- created_at: 创建时间
- updated_at: 更新时间
- view_count: 浏览量
- like_count: 点赞数
- comment_count: 评论数
- status: 状态（published/deleted）

### 评论模型 (Comment)
- id: 评论唯一标识
- content: 评论内容
- post_id: 关联帖子ID
- author_id: 作者ID（关联用户）
- created_at: 创建时间
- updated_at: 更新时间
- parent_id: 回复的父评论ID（支持回复评论）
- liked_users: 点赞用户关系（多对多关联用户表）

## 开发规范

- 使用 Flask 的应用工厂模式
- 所有 API 响应使用 JSON 格式
- 使用 JWT 进行身份验证
- 数据库操作使用 SQLAlchemy ORM
- 代码遵循 PEP 8 规范
- API 接口统一以 `/api` 为前缀
- 错误处理返回结构化的错误信息
- 在应用工厂模式下，SQLAlchemy实例通过`app.db = db`附加到应用实例，并在路由中使用`current_app.db`访问数据库实例

## 权限控制

- 认证路由使用 `@jwt_required()` 装饰器
- 管理员权限通过检查 `user.role === 'admin'` 来验证
- 用户只能编辑自己的内容（帖子、评论等）
- 删除操作会检查权限和关联关系

## 部署

生产环境中建议使用 Gunicorn 或 uWSGI 部署 Flask 应用：

```bash
gunicorn app:app
```

配合 Nginx 作为反向代理服务器。

## 项目维护

- 数据库变更时需要重新运行 [init_db.py](file:///c:/Users/Lenovo/Desktop/HuxXiang_culture_github/backend/init_db.py)
- 服务重启后需要确保数据库连接正常
- 定期备份数据库以防数据丢失

## 许可证

湖湘文化数字化平台采用 MIT 许可证。