<script>
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from './components/Navbar.vue'
import authService from './services/authService.js'

export default {
  name: 'App',
  components: {
    Navbar
  },
  setup() {
    const router = useRouter()
    const isLoggedIn = ref(false)
    const user = ref(null)
    const alert = ref({ show: false, message: '', type: 'success' })
    let alertTimer = null

    const syncAuthState = () => {
      const loggedIn = authService.isAuthenticated()
      isLoggedIn.value = loggedIn
      user.value = loggedIn ? authService.getCurrentUser() : null
    }

    const showAlert = (message, type = 'success') => {
      if (alertTimer) {
        clearTimeout(alertTimer)
      }

      alert.value = {
        show: true,
        message,
        type
      }

      alertTimer = setTimeout(() => {
        alert.value.show = false
      }, 3000)
    }

    const navigateTo = (page) => {
      router.push({ name: page })
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    const navigateToLogin = () => {
      router.push('/login')
    }

    const navigateToRegister = () => {
      router.push('/register')
    }

    const handleLogout = async () => {
      try {
        const result = await authService.logout()
        syncAuthState()
        showAlert(result.message || '已退出登录', 'info')
      } finally {
        router.push('/')
      }
    }

    const handleAuthChanged = () => {
      syncAuthState()
    }

    onMounted(() => {
      syncAuthState()
      window.addEventListener('login-success', handleAuthChanged)
      window.addEventListener('logout', handleAuthChanged)
      window.addEventListener('storage', handleAuthChanged)
    })

    onUnmounted(() => {
      if (alertTimer) {
        clearTimeout(alertTimer)
      }
      window.removeEventListener('login-success', handleAuthChanged)
      window.removeEventListener('logout', handleAuthChanged)
      window.removeEventListener('storage', handleAuthChanged)
    })

    return {
      alert,
      handleLogout,
      isLoggedIn,
      navigateTo,
      navigateToLogin,
      navigateToRegister,
      showAlert,
      user
    }
  }
}
</script>

<template>
  <div id="app">
    <Navbar
      :is-logged-in="isLoggedIn"
      :user="user"
      :navigate-to="navigateTo"
      :navigate-to-login="navigateToLogin"
      :navigate-to-register="navigateToRegister"
      :handle-logout="handleLogout"
    />

    <main>
      <router-view v-slot="{ Component }">
        <component :is="Component" :show-alert="showAlert" />
      </router-view>
    </main>

    <footer class="site-footer">
      <div class="footer-container">
        <div class="footer-column">
          <h3><i class="fas fa-info-circle"></i> 关于我们</h3>
          <ul>
            <li><a href="#" @click.prevent="navigateTo('about')">平台介绍</a></li>
            <li><a href="#" @click.prevent="navigateTo('about')">团队成员</a></li>
            <li><a href="#" @click.prevent="navigateTo('about')">发展历程</a></li>
            <li><a href="#" @click.prevent="navigateTo('about')">合作伙伴</a></li>
          </ul>
        </div>
        <div class="footer-column">
          <h3><i class="fas fa-book"></i> 资源中心</h3>
          <ul>
            <li><a href="#" @click.prevent="navigateTo('cultural-resources')">文化资源</a></li>
            <li><a href="#" @click.prevent="navigateTo('cultural-resources')">资源分类</a></li>
            <li><a href="#" @click.prevent="navigateTo('cultural-resources')">资源上传</a></li>
            <li><a href="#" @click.prevent="navigateTo('cultural-resources')">资源检索</a></li>
          </ul>
        </div>
        <div class="footer-column">
          <h3><i class="fas fa-laptop-code"></i> 数字化展示</h3>
          <ul>
            <li><a href="#" @click.prevent="navigateTo('digital-showcase')">核心体验</a></li>
            <li><a href="#" @click.prevent="navigateTo('digital-showcase')">数字博物馆</a></li>
            <li><a href="#" @click.prevent="navigateTo('digital-showcase')">互动展示</a></li>
            <li><a href="#" @click.prevent="navigateTo('digital-showcase')">专题内容</a></li>
          </ul>
        </div>
        <div class="footer-column">
          <h3><i class="fas fa-users"></i> 互动社区</h3>
          <ul>
            <li><a href="#" @click.prevent="navigateTo('community')">文化论坛</a></li>
            <li><a href="#" @click.prevent="navigateTo('community')">文化活动</a></li>
            <li><a href="#" @click.prevent="navigateTo('community')">内容贡献</a></li>
            <li><a href="#" @click.prevent="navigateTo('community')">意见反馈</a></li>
          </ul>
        </div>
        <div class="footer-column">
          <h3><i class="fas fa-envelope"></i> 联系我们</h3>
          <ul>
            <li><a href="#" @click.prevent="navigateTo('contact')">联系方式</a></li>
            <li><a href="#" @click.prevent="navigateTo('contact')">合作咨询</a></li>
            <li><a href="#" @click.prevent="navigateTo('contact')">常见问题</a></li>
            <li><a href="#" @click.prevent="navigateTo('contact')">隐私政策</a></li>
          </ul>
        </div>
      </div>
      <div class="copyright">
        <p>&copy; 2025 湖湘文化数字化开发与传播平台 版权所有</p>
      </div>
    </footer>

    <div v-if="alert.show" :class="['alert-container', 'alert-' + alert.type]">
      <i v-if="alert.type === 'success'" class="fas fa-check-circle"></i>
      <i v-else-if="alert.type === 'error'" class="fas fa-exclamation-circle"></i>
      <i v-else class="fas fa-info-circle"></i>
      <span>{{ alert.message }}</span>
    </div>
  </div>
</template>

<style>
@import './assets/css/style.css';

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

main {
  flex: 1;
}

.site-footer {
  margin-top: auto;
}

.alert-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 260px;
  padding: 14px 18px;
  border-radius: 8px;
  color: #fff;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.18);
  animation: slide-in 0.25s ease-out;
}

.alert-success {
  background: #2e7d32;
}

.alert-error {
  background: #c62828;
}

.alert-info {
  background: #1565c0;
}

@keyframes slide-in {
  from {
    opacity: 0;
    transform: translateX(24px);
  }

  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>
