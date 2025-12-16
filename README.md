# 湖湘文化数字化平台

## 项目概述

湖湘文化数字化平台是一个致力于通过数字化手段展示、传播和传承湖湘文化精髓的综合性Web应用。该项目旨在利用现代信息技术，将湖湘地区丰富的历史文化遗产、民俗风情、传统艺术等资源进行数字化呈现，促进文化的保护与传承，同时为用户提供便捷的文化资源获取和互动交流平台。

## 项目结构

项目采用前后端分离的架构设计，分为两个主要部分：

1. **HuXiang_github**：基于Vue 3的前端项目
2. **HX_project**：基于Flask的后端项目

## 技术栈

### 前端技术栈

| 技术/框架 | 版本 | 用途 |
|---------|------|------|
| Vue.js | 3.5.21 | 前端框架，用于构建用户界面和组件化开发 |
| Vue Router | 4.5.1 | 前端路由管理，实现页面间的导航和参数传递 |
| Vite | 7.1.7 | 现代化前端构建工具，提供快速的开发体验和优化的构建输出 |
| Font Awesome | 7.0.1 | 图标库，提供丰富的图标资源用于界面设计 |
| 原生CSS | - | 样式设计，实现响应式布局和界面美化 |
| Unity WebGL | - | 3D内容展示，用于数字化展示部分 |

### 后端技术栈

| 技术/框架 | 版本 | 用途 |
|---------|------|------|
| Flask | - | Python Web框架，提供后端API服务 |
| Flask-SQLAlchemy | 3.0.5 | ORM框架，简化数据库操作 |
| PyMySQL | 1.1.0 | MySQL数据库驱动，用于数据库连接和操作 |
| Werkzeug | 2.3.7 | WSGI工具库，提供Web服务器网关接口实现 |

## 功能模块

### 1. 首页模块

作为平台的门户，展示平台的核心功能入口和最新文化资源推荐，提供全局导航和搜索功能。

### 2. 文化资源库模块

集中展示湖湘地区的各类文化资源，包括历史遗迹、传统艺术、文学作品、民俗风情、饮食文化、建筑风格、宗教文化、历史伟人事迹等。支持按分类浏览、关键词搜索和分页展示功能。

资源数据结构：
- 标题、描述、分类、图片URL、创建日期

### 3. 互动社区模块

提供用户交流互动的平台，支持用户发布、浏览和评论文化相关的帖子和讨论。该模块需要用户登录才能访问。

### 4. 数字化展示模块

利用Unity WebGL技术实现文化资源的3D展示和互动体验，为用户提供沉浸式的文化体验。

### 5. 知识图谱模块

可视化展示湖湘文化元素之间的关联关系，帮助用户更好地理解和探索湖湘文化的内在联系。

### 6. 用户认证与管理模块

- **登录/注册**：支持用户通过用户名/邮箱和密码进行登录和注册
- **个人资料**：用户可以查看和管理个人信息
- **权限控制**：区分普通用户和管理员权限

### 7. 管理员模块

提供后台管理功能，支持管理员对文化资源、用户和社区内容进行管理和维护

### 3. 互动社区模块

提供用户交流互动的平台，支持用户发布、浏览和评论文化相关的帖子和讨论。该模块需要用户登录才能访问。

### 4. 数字化展示模块

利用Unity WebGL技术实现文化资源的3D展示和互动体验，为用户提供沉浸式的文化体验。

### 5. 知识图谱模块

可视化展示湖湘文化元素之间的关联关系，帮助用户更好地理解和探索湖湘文化的内在联系。

### 6. 用户认证与管理模块

- **登录/注册**：支持用户通过用户名/邮箱和密码进行登录和注册
- **个人资料**：用户可以查看和管理个人信息
- **权限控制**：区分普通用户和管理员权限



### 7. 管理员模块

提供后台管理功能，支持管理员对文化资源、用户和社区内容进行管理和维护。

## 前端项目结构

```
src/
├── assets/         # 静态资源文件
│   ├── css/        # 样式文件
│   └── imgs/       # 图片文件
├── components/     # Vue组件
├── views/          # 页面组件
├── services/       # 服务层，处理API请求和业务逻辑
├── router/         # 路由配置
├── App.vue         # 根组件
└── main.js         # 应用入口文件
```
<mcfile name="README.md" path="e:/Project_huxiangwenhua/HuXiang_github/README.md"></mcfile>

## 核心页面组件

| 页面组件 | 路径 | 功能描述 |
|--------|------|---------|
| HomePage | src/views/HomePage.vue | 平台首页，展示核心功能入口和推荐资源 |
| CulturalResourcesPage | src/views/CulturalResourcesPage.vue | 文化资源库，展示各类湖湘文化资源 |
| CommunityPage | src/views/CommunityPage.vue | 互动社区，提供用户交流平台 |
| DigitalShowcasePage | src/views/DigitalShowcasePage.vue | 数字化展示，提供3D文化体验 |
| AboutPage | src/views/AboutPage.vue | 关于我们，介绍平台的背景和宗旨 |
| ContactPage | src/views/ContactPage.vue | 联系我们，提供联系方式和反馈渠道 |
| ResourceDetailPage | src/views/ResourceDetailPage.vue | 资源详情，展示单个文化资源的详细信息 |
| KnowledgeGraphPage | src/views/KnowledgeGraphPage.vue | 知识图谱，可视化展示文化元素关联 |
| UnityWebGL | src/views/UnityWebGL.vue | Unity WebGL内容展示 |
| LoginView/RegisterView | src/views/LoginView.vue/src/views/RegisterView.vue | 用户登录和注册页面 |
| ProfilePage | src/views/ProfilePage.vue | 用户个人资料页面 |
| AdminPage | src/views/AdminPage.vue | 管理员后台管理页面 |

<mcfile name="index.js" path="e:/Project_huxiangwenhua/HuXiang_github/src/router/index.js"></mcfile>

## 特色功能

### 1. 响应式设计

平台采用响应式设计，支持在不同尺寸的设备上提供良好的用户体验，包括桌面端、平板和移动设备。

### 2. 模拟数据服务

前端项目包含模拟数据服务，用于在后端API未完全开发时提供数据支持，确保前端功能可以正常演示和测试。

### 3. 全局事件总线

实现了全局事件总线机制，方便组件间的通信和数据传递。

```javascript
// 全局事件总线 - 用于组件间通信
app.config.globalProperties.$eventBus = app
```
<mcfile name="main.js" path="e:/Project_huxiangwenhua/HuXiang_github/src/main.js"></mcfile>

## 项目意义

湖湘文化数字化平台通过现代信息技术手段，为湖湘文化的保护、传承和传播提供了新的途径。该平台不仅可以促进湖湘文化的研究和交流，还可以增强公众对湖湘文化的认知和认同，对于推动湖湘文化的创新发展具有重要意义。

## 开发与部署

### 前端开发环境设置

1. 安装Node.js 14.18+ 或 16+
2. 克隆项目仓库
3. 安装依赖：`npm install`
4. 启动开发服务器：`npm run dev`
5. 构建生产版本：`npm run build`

### 后端开发环境设置

1. 安装Python和pip
2. 安装依赖：`pip install -r requirements.txt`
3. 运行应用：`python run.py`

## 结语

湖湘文化数字化平台是一个融合了现代信息技术和传统文化的创新项目，旨在通过数字化手段推动湖湘文化的保护、传承和发展。随着项目的不断完善和优化，相信该平台将为湖湘文化的传播和弘扬做出更大的贡献。