<script>
import { ref, onMounted } from 'vue'
import Navbar from './components/Navbar.vue'
import { useRouter } from 'vue-router'

export default {
  name: 'App',
  components: {
    Navbar
  },
  setup() {
    // 状态管理
    const isLoggedIn = ref(false)
    const user = ref(null)
    const alert = ref({ show: false, message: '', type: 'success' })

    // 获取路由实例
    const router = useRouter()

    // 页面导航方法（与组件内部导航保持兼容）
    const navigateTo = (page) => {
      router.push({ name: page })
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    const goToResourceDetail = (resourceId) => {
      router.push({ name: 'resource-detail', params: { id: resourceId } })
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    const goToResources = (tag = null) => {
      router.push({ name: 'cultural-resources', query: { tag } })
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    const goBack = () => {
      router.back()
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    // 导航到登录页
    const navigateToLogin = () => {
      router.push('/login')
    }

    // 导航到注册页
    const navigateToRegister = () => {
      router.push('/register')
    }

    // 登出处理
    const handleLogout = async () => {
      try {
        // 在实际环境中使用真实API
        // await userAPI.logout();
        
        // 清除本地存储中的用户信息
        localStorage.removeItem('user')
        
        // 更新状态
        isLoggedIn.value = false
        user.value = null
        
        // 显示提示消息
        showAlert('已成功登出！', 'info')
        
        // 重定向到首页
        router.push('/')
      } catch (error) {
        console.error('登出失败:', error)
        // 即使API调用失败，也将本地状态设置为未登录
        localStorage.removeItem('user')
        isLoggedIn.value = false
        user.value = null
      }
    }

    // 显示提示消息
    const showAlert = (message, type = 'success') => {
      alert.value = { show: true, message, type }
      setTimeout(() => {
        alert.value.show = false
      }, 3000)
    }

    // 监听导航事件（保持向后兼容）
    onMounted(() => {
      const handleNavigate = (event) => {
        if (event.detail && event.detail.page) {
          navigateTo(event.detail.page)
        }
      }
      window.addEventListener('navigate', handleNavigate)
      
      // 监听登录成功事件
      const handleLoginSuccess = () => {
        checkUserLoginStatus()
      }
      window.addEventListener('login-success', handleLoginSuccess)
      
      // 尝试检查用户登录状态
      checkUserLoginStatus()
      
      // 清理事件监听器
      return () => {
        window.removeEventListener('navigate', handleNavigate)
        window.removeEventListener('login-success', handleLoginSuccess)
      }
    })

    // 检查用户登录状态
    const checkUserLoginStatus = () => {
      try {
        // 从本地存储读取用户信息
        const storedUser = localStorage.getItem('user')
        if (storedUser) {
          const userInfo = JSON.parse(storedUser)
          isLoggedIn.value = true
          user.value = userInfo
        }
      } catch (error) {
        console.error('检查用户登录状态失败:', error)
      }
    }

    return {
      isLoggedIn,
      user,
      alert,
      navigateTo,
      goToResourceDetail,
      goToResources,
      goBack,
      navigateToLogin,
      navigateToRegister,
      handleLogout,
      showAlert
    }
  }
}
</script>

<template>
  <div id="app">
    <!-- 头部导航 -->
    <Navbar 
      :is-logged-in="isLoggedIn" 
      :user="user"
      :navigate-to="navigateTo"
      :navigate-to-login="navigateToLogin"
      :navigate-to-register="navigateToRegister"
      :handle-logout="handleLogout"
    />

    <!-- 主要内容区域 -->
    <main>
      <router-view />
    </main>

    <!-- 页脚 -->
    <footer>
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
            <li><a href="#" @click.prevent="navigateTo('digital-showcase')">核心数字化体验</a></li>
            <li><a href="#" @click.prevent="navigateTo('digital-showcase')">数字博物馆</a></li>
            <li><a href="#" @click.prevent="navigateTo('digital-showcase')">互动文化体验</a></li>
            <li><a href="#" @click.prevent="navigateTo('digital-showcase')">特色专题</a></li>
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



    <!-- 消息提示 -->
    <div v-if="alert.show" :class="['alert', 'alert-' + alert.type]" class="alert-container">
      <i v-if="alert.type === 'success'" class="fas fa-check-circle"></i>
      <i v-else-if="alert.type === 'error'" class="fas fa-exclamation-circle"></i>
      <i v-else-if="alert.type === 'info'" class="fas fa-info-circle"></i>
      <span>{{ alert.message }}</span>
    </div>
  </div>
</template>

<style>
/* 导入全局样式 */
@import './assets/css/style.css';

/* 自定义全局样式 */
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

main {
  flex: 1;
}

/* 增强的提示框样式 */
.alert-container {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 15px 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 500;
  z-index: 1000;
  animation: slideIn 0.3s ease-out;
  min-width: 280px;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.alert-success {
  background-color: #4caf50;
  color: white;
  border-left: 4px solid #388e3c;
}

.alert-error {
  background-color: #f44336;
  color: white;
  border-left: 4px solid #d32f2f;
}

.alert-info {
  background-color: #2196f3;
  color: white;
  border-left: 4px solid #1976d2;
}

.alert-container i {
  font-size: 20px;
}


</style>
