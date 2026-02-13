import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../views/HomePage.vue'
import CulturalResourcesPage from '../views/CulturalResourcesPage.vue'
import CommunityPage from '../views/CommunityPage.vue'
import DigitalShowcasePage from '../views/DigitalShowcasePage.vue'
import AboutPage from '../views/AboutPage.vue'
import ContactPage from '../views/ContactPage.vue'
import ResourceDetailPage from '../views/ResourceDetailPage.vue'
import KnowledgeGraphPage from '../views/KnowledgeGraphPage.vue'
import UnityWebGL from '../views/UnityWebGL.vue'
import NotFound from '../views/NotFound.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ProfilePage from '../views/ProfilePage.vue'
import AdminPage from '../views/AdminPage.vue'
import PostDetailPage from '../views/PostDetailPage.vue'
import AiAssistantPage from '../views/AiAssistantPage.vue'
import PoetryDigitalizationPage from '../views/PoetryDigitalizationPage.vue'
import Architecture3DPage from '../views/Architecture3DPage.vue'
import EditPostPage from '../views/EditPostPage.vue'
import CreatePostPage from '../views/CreatePostPage.vue'  // 导入新的创建帖子页面
import { isAuthenticated, isAdmin } from '../services/authService.js' // 导入认证服务


const routes = [
  {
    path: '/',
    name: 'home',
    component: HomePage,
    meta: {
      title: '首页 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/cultural-resources',
    name: 'cultural-resources',
    component: CulturalResourcesPage,
    meta: {
      title: '文化资源 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/resource-detail/:id',
    name: 'resource-detail',
    component: ResourceDetailPage,
    meta: {
      title: '资源详情 - 湖湘文化数字化平台'
    },
    props: true
  },
  {
    path: '/community',
    name: 'community',
    component: CommunityPage,
    meta: {
      title: '互动社区 - 湖湘文化数字化平台',
      requiresAuth: true
    }
  },
  {
    path: '/digital-showcase',
    name: 'digital-showcase',
    component: DigitalShowcasePage,
    meta: {
      title: '数字化展示 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/ai-assistant',
    name: 'ai-assistant',
    component: AiAssistantPage,
    meta: {
      title: 'AI文化助手 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/about',
    name: 'about',
    component: AboutPage,
    meta: {
      title: '关于我们 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/contact',
    name: 'contact',
    component: ContactPage,
    meta: {
      title: '联系我们 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: {
      title: '登录 - 湖湘文化数字化平台',
      guest: true
    }
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView,
    meta: {
      title: '注册 - 湖湘文化数字化平台',
      guest: true
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound,
    meta: {
      title: '页面未找到 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/knowledge-graph',
    name: 'knowledge-graph',
    component: KnowledgeGraphPage,
    meta: {
      title: '知识图谱 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/unity-webgl',
    name: 'unity-webgl',
    component: UnityWebGL,
    meta: {
      title: '虚拟现实体验 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/profile',
    name: 'profile',
    component: ProfilePage,
    meta: {
      title: '个人中心 - 湖湘文化数字化平台',
      requiresAuth: true
    }
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminPage,
    meta: {
      title: '管理中心 - 湖湘文化数字化平台',
      requiresAuth: true,
      requiresAdmin: true
    }
  },
  {
    path: '/post-detail/:id',
    name: 'post-detail',
    component: PostDetailPage,
    meta: {
      title: '帖子详情 - 湖湘文化数字化平台',
      requiresAuth: true
    },
    props: true
  },
  {
    path: '/edit-post/:id',
    name: 'edit-post',
    component: EditPostPage,
    meta: {
      title: '编辑帖子 - 湖湘文化数字化平台',
      requiresAuth: true
    },
    props: true
  },
  {
    path: '/create-post',
    name: 'create-post',
    component: CreatePostPage,
    meta: {
      title: '发布新主题 - 湖湘文化数字化平台',
      requiresAuth: true
    }
  },
  {
    path: '/poetry-digitalization',
    name: 'poetry-digitalization',
    component: PoetryDigitalizationPage,
    meta: {
      title: '湖湘诗词数字化 - 湖湘文化数字化平台'
    }
  },
  {
    path: '/architecture-3d',
    name: 'architecture-3d',
    component: Architecture3DPage,
    meta: {
      title: '湖湘建筑3D模型 - 湖湘文化数字化平台'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由导航守卫 - 设置页面标题和用户认证
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title
  }

  // 检查用户是否已登录
  const userIsAuthenticated = isAuthenticated()
  const user = userIsAuthenticated ? JSON.parse(localStorage.getItem('user')) : null

  // 如果页面需要认证
  if (to.meta.requiresAuth) {
    if (!userIsAuthenticated) {
      // 未登录，跳转到登录页
      next({ name: 'login', query: { redirect: to.fullPath } })
    } else if (to.meta.requiresAdmin && !isAdmin()) {  // 使用authService的isAdmin函数
      // 需要管理员权限，但用户不是管理员
      next({ name: 'home' })
    } else {
      // 已登录且权限足够，继续访问
      next()
    }
  } else if (to.meta.guest) {
    // 如果是登录/注册页，已登录用户重定向到首页
    if (userIsAuthenticated) {
      next({ name: 'home' })
    } else {
      next()
    }
  } else {
    // 其他页面直接访问
    next()
  }
})

export default router