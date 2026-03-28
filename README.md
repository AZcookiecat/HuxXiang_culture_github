# 湖湘文化数字化平台

[![Vue 3](https://img.shields.io/badge/Vue-3.5.21-42b883?style=flat&logo=vue.js)](https://vuejs.org)
[![Flask](https://img.shields.io/badge/Flask-2.3.3-000000?style=flat&logo=flask)](https://flask.palletsprojects.com)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-00758F?style=flat&logo=mysql)](https://www.mysql.com)

## 项目简介

湖湘文化数字化平台是一个致力于通过数字化手段展示、传播和传承湖湘文化精髓的综合性 Web 应用。项目采用前后端分离架构，前端使用 Vue 3 + Vite 构建，后端使用 Flask + SQLAlchemy，提供文化资源展示、社区互动、知识图谱、AI 助手等功能。

## 核心功能

| 功能模块 | 说明 |
|---------|------|
| 🏠 首页 | 平台门户，核心功能入口 |
| 📚 文化资源库 | 湖湘文化资源展示（历史遗迹、传统艺术、诗词、美食等） |
| 💬 互动社区 | 用户发帖、评论、点赞 |
| 🕸️ 知识图谱 | 文化元素关联可视化 |
| 🎮 数字化展示 | Unity WebGL 3D 体验 |
| 🤖 AI 助手 | 智能问答服务 |
| 👤 用户系统 | 登录注册、个人中心 |
| ⚙️ 管理后台 | 内容管理、用户管理 |

## 技术栈

### 前端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue.js | ^3.5.21 | 渐进式前端框架 |
| Vue Router | ^4.5.1 | 客户端路由 |
| Vite | ^7.1.7 | 现代化构建工具 |
| Font Awesome | ^7.0.1 | 图标库 |

### 后端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Flask | 2.3.3 | Python Web 框架 |
| Flask-SQLAlchemy | 3.0.5 | ORM 数据库操作 |
| Flask-CORS | 4.0.0 | 跨域资源共享 |
| Flask-JWT-Extended | 4.5.3 | JWT 身份认证 |
| PyMySQL | 1.1.0 | MySQL 数据库驱动 |
| Werkzeug | 2.3.7 | WSGI 工具库 |

## 项目结构

```
HuXiang_culture/
│
├── src/                          # 前端 Vue 项目
│   ├── assets/                   # 静态资源
│   │   ├── css/                  # 样式文件
│   │   └── imgs/                 # 图片资源
│   ├── components/               # 公共组件
│   │   └── CommentsSection.vue   # 评论组件
│   ├── views/                    # 页面组件
│   │   ├── HomePage.vue          # 首页
│   │   ├── CommunityPage.vue     # 社区页
│   │   ├── CreatePostPage.vue    # 发帖页
│   │   ├── PostDetailPage.vue    # 帖子详情
│   │   ├── CulturalResourcesPage.vue  # 文化资源
│   │   ├── KnowledgeGraphPage.vue     # 知识图谱
│   │   ├── LoginView.vue         # 登录
│   │   ├── RegisterView.vue      # 注册
│   │   └── ...                   # 其他页面
│   ├── router/
│   │   └── index.js              # 路由配置
│   ├── services/
│   │   ├── api.js                # API 服务
│   │   └── authService.js        # 认证服务
│   ├── main.js                   # 入口文件
│   ├── App.vue                   # 根组件
│   └── style.css                 # 全局样式
│
├── backend/                      # 后端 Flask 项目
│   ├── models/                   # 数据模型
│   │   ├── user.py               # 用户模型
│   │   ├── community_post.py     # 帖子/评论模型
│   │   └── cultural_resource.py  # 文化资源模型
│   ├── routes/                   # API 路由
│   │   ├── auth.py               # 认证路由
│   │   ├── community.py          # 社区路由
│   │   ├── cultural_resources.py  # 文化资源路由
│   │   └── main.py               # 通用路由
│   ├── app/
│   │   └── __init__.py           # 应用初始化
│   ├── app.py                    # 应用入口
│   ├── config.py                  # 配置文件
│   ├── init_db.py                 # 数据库初始化
│   └── requirements.txt           # Python 依赖
│
├── public/                       # 公共静态资源
│   └── unity-webgl/              # Unity WebGL 游戏
│
├── index.html                    # HTML 入口
├── package.json                  # 前端依赖
├── vite.config.js                # Vite 配置
└── README.md                     # 项目文档
```

## 环境要求

| 环境 | 要求 |
|------|------|
| Node.js | 16+ |
| Python | 3.8+ |
| MySQL | 8.0+ |
| npm / pip | 最新版本 |

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd HuXiang_culture
```

### 2. 前端配置

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

前端服务启动后访问：`http://localhost:5173`

### 3. 后端配置

```bash
# 进入后端目录
cd backend

# 创建虚拟环境（推荐）
python -m venv venv

# 激活虚拟环境
# Windows:
venv\Scripts\activate
# Linux/Mac:
source venv/bin/activate

# 安装依赖
pip install -r requirements.txt

# 配置数据库
# 修改 config.py 中的数据库连接信息
# 默认: mysql+pymysql://root:password@localhost/huxiang_culture

# 初始化数据库
python init_db.py

# 启动服务
python app.py
```

后端服务启动后访问：`http://localhost:5000`

### 4. 环境变量（可选）

在 `backend/` 目录下创建 `.env` 文件：

```env
DATABASE_URL=mysql+pymysql://root:password@localhost/huxiang_culture
SECRET_KEY=your-secret-key
JWT_SECRET_KEY=your-jwt-secret-key
```

## API 接口文档

### 认证接口 (auth)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/auth/register` | 用户注册 | 否 |
| POST | `/api/auth/login` | 用户登录 | 否 |
| GET | `/api/auth/profile` | 获取用户信息 | 是 |
| PUT | `/api/auth/profile` | 更新用户信息 | 是 |
| POST | `/api/auth/upload-avatar` | 上传头像 | 是 |
| POST | `/api/auth/logout` | 登出 | 是 |

### 文化资源接口 (cultural-resources)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/cultural-resources` | 获取资源列表 | 否 |
| GET | `/api/cultural-resources/<id>` | 获取资源详情 | 否 |
| POST | `/api/cultural-resources/<id>/like` | 点赞资源 | 是 |
| POST | `/api/cultural-resources` | 创建资源 | 是(管理员) |

### 社区帖子接口 (community)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/community/posts` | 获取帖子列表 | 否 |
| GET | `/api/community/posts/<id>` | 获取帖子详情 | 否 |
| POST | `/api/community/posts` | 发布帖子 | 是 |
| PUT | `/api/community/posts/<id>` | 更新帖子 | 是 |
| DELETE | `/api/community/posts/<id>` | 删除帖子 | 是 |
| POST | `/api/community/posts/<id>/like` | 点赞帖子 | 是 |
| GET | `/api/community/posts/<id>/comments` | 获取评论列表 | 否 |
| POST | `/api/community/posts/<id>/comments` | 添加评论 | 是 |
| DELETE | `/api/community/comments/<id>` | 删除评论 | 是 |

### 通用接口 (main)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/` | 首页 | 否 |
| GET | `/health` | 健康检查 | 否 |

## 数据模型

### User（用户）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | Integer | 主键 |
| username | String | 用户名 |
| email | String | 邮箱 |
| password_hash | String | 密码哈希 |
| avatar_url | String | 头像URL |
| role | String | 角色(user/admin) |
| created_at | DateTime | 创建时间 |

### CommunityPost（帖子）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | Integer | 主键 |
| title | String | 标题 |
| content | Text | 内容 |
| author_id | Integer | 作者ID |
| category | String | 分类 |
| status | String | 状态 |
| view_count | Integer | 浏览量 |
| like_count | Integer | 点赞数 |
| comment_count | Integer | 评论数 |
| created_at | DateTime | 创建时间 |

### Comment（评论）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | Integer | 主键 |
| content | Text | 内容 |
| author_id | Integer | 作者ID |
| post_id | Integer | 帖子ID |
| parent_id | Integer | 父评论ID（支持嵌套） |
| created_at | DateTime | 创建时间 |

### CulturalResource（文化资源）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | Integer | 主键 |
| title | String | 标题 |
| description | Text | 描述 |
| category | String | 分类 |
| image_url | String | 图片URL |
| created_at | DateTime | 创建时间 |

## 使用示例

### 前端 API 调用

```javascript
// 获取帖子列表
import { get } from './services/api';

const posts = await get('/community/posts');

// 登录
import { post } from './services/api';

const result = await post('/auth/login', {
  username: 'user123',
  password: 'password123'
});

// 获取 token
const token = result.access_token;
```

### 后端模型查询

```python
from app import db
from models.user import User
from models.community_post import CommunityPost, Comment

# 查询用户
user = User.query.filter_by(username='test').first()

# 查询帖子及评论
post = CommunityPost.query.get(1)
comments = Comment.query.filter_by(post_id=1, parent_id=None).all()

# 创建帖子
new_post = CommunityPost(
    title='新帖子',
    content='帖子内容',
    author_id=user.id,
    category='讨论'
)
db.session.add(new_post)
db.session.commit()
```

## 开发规范

### 前端规范

- 使用 Vue 3 Composition API
- 组件文件使用 PascalCase 命名
- API 请求统一通过 `services/api.js` 封装
- 样式使用原生 CSS 或 SCSS

### 后端规范

- 遵循 Flask 蓝图画分路由
- 使用 SQLAlchemy ORM 操作数据库
- API 返回统一 JSON 格式
- 需要认证的接口使用 JWT

## 常见问题

### Q: 前端无法连接后端？
A: 检查后端是否启动在 `http://localhost:5000`，确认 CORS 配置正确。

### Q: 数据库连接失败？
A: 确认 MySQL 服务已启动，配置文件中用户名密码正确，数据库已创建。

### Q: 如何切换到生产环境？
A: 修改 `config.py` 中的数据库 URL，使用生产级服务器（如 Gunicorn）。

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/xxx`)
3. 提交更改 (`git commit -m 'Add xxx'`)
4. 推送分支 (`git push origin feature/xxx`)
5. 创建 Pull Request

## 许可证

本项目基于 MIT 许可证开源，详见 [LICENSE](LICENSE) 文件。

## 联系方式

- 项目作者：[作者名称]
- 邮箱：[email@example.com]
- GitHub：[https://github.com/your-repo](https://github.com/your-repo)

---

<p align="center">
  湖湘文化数字化平台 | 致力于湖湘文化的保护与传承
</p>