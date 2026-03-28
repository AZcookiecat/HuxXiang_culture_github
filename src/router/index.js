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
import CreatePostPage from '../views/CreatePostPage.vue'
import { isAuthenticated, isAdmin } from '../services/authService.js'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomePage,
    meta: { title: '首页 - 湖湘文化数字化平台' }
  },
  {
    path: '/cultural-resources',
    name: 'cultural-resources',
    component: CulturalResourcesPage,
    meta: { title: '文化资源 - 湖湘文化数字化平台' }
  },
  {
    path: '/resource-detail/:id',
    name: 'resource-detail',
    component: ResourceDetailPage,
    props: true,
    meta: { title: '资源详情 - 湖湘文化数字化平台' }
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
    meta: { title: '数字化展示 - 湖湘文化数字化平台' }
  },
  {
    path: '/ai-assistant',
    name: 'ai-assistant',
    component: AiAssistantPage,
    meta: { title: 'AI 文化助手 - 湖湘文化数字化平台' }
  },
  {
    path: '/about',
    name: 'about',
    component: AboutPage,
    meta: { title: '关于我们 - 湖湘文化数字化平台' }
  },
  {
    path: '/contact',
    name: 'contact',
    component: ContactPage,
    meta: { title: '联系我们 - 湖湘文化数字化平台' }
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
    path: '/knowledge-graph',
    name: 'knowledge-graph',
    component: KnowledgeGraphPage,
    meta: { title: '知识图谱 - 湖湘文化数字化平台' }
  },
  {
    path: '/unity-webgl',
    name: 'unity-webgl',
    component: UnityWebGL,
    meta: { title: 'Unity WebGL - 湖湘文化数字化平台' }
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
    props: true,
    meta: {
      title: '帖子详情 - 湖湘文化数字化平台',
      requiresAuth: true
    }
  },
  {
    path: '/edit-post/:id',
    name: 'edit-post',
    component: EditPostPage,
    props: true,
    meta: {
      title: '编辑帖子 - 湖湘文化数字化平台',
      requiresAuth: true
    }
  },
  {
    path: '/create-post',
    name: 'create-post',
    component: CreatePostPage,
    meta: {
      title: '发布帖子 - 湖湘文化数字化平台',
      requiresAuth: true
    }
  },
  {
    path: '/poetry-digitalization',
    name: 'poetry-digitalization',
    component: PoetryDigitalizationPage,
    meta: { title: '诗词数字化 - 湖湘文化数字化平台' }
  },
  {
    path: '/architecture-3d',
    name: 'architecture-3d',
    component: Architecture3DPage,
    meta: { title: '建筑 3D - 湖湘文化数字化平台' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound,
    meta: { title: '页面未找到 - 湖湘文化数字化平台' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }

  if (to.meta.requiresAuth && !isAuthenticated()) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.requiresAdmin && !isAdmin()) {
    next({ name: 'home' })
    return
  }

  if (to.meta.guest && isAuthenticated()) {
    next({ name: 'home' })
    return
  }

  next()
})

export default router
