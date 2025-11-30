<template>
  <header>
    <div class="container">
      <nav class="navbar">
        <a href="#" class="navbar-brand" @click.prevent="navigateTo('home')">
          <img src="/src/assets/imgs/logo1.jpg" alt="平台Logo" class="navbar-logo" loading="lazy">
          湖湘文化数字化平台
        </a>
        
        <!-- 移动端菜单按钮 -->
        <button class="navbar-toggle" @click="toggleMenu">
          <i class="fas fa-bars"></i>
        </button>
        
        <!-- 导航菜单 -->
        <ul class="navbar-nav" :class="{ show: isMenuOpen }">
          <li class="nav-item">
            <a href="#" class="nav-link" @click.prevent="navigateTo('home')">首页</a>
          </li>
          <li class="nav-item">
            <a href="#" class="nav-link" @click.prevent="navigateTo('cultural-resources')">文化资源</a>
          </li>
          <li class="nav-item">
            <a href="#" class="nav-link" @click.prevent="navigateTo('digital-showcase')">数字化展示</a>
          </li>
          <li class="nav-item">
            <a href="#" class="nav-link" @click.prevent="navigateTo('community')">互动社区</a>
          </li>
          <li class="nav-item">
            <a href="#" class="nav-link" @click.prevent="navigateTo('about')">关于我们</a>
          </li>
          <!-- 管理员页面导航项，只对管理员显示 -->
          <li class="nav-item" v-if="isLoggedIn && user?.isAdmin">
            <a href="#" class="nav-link admin-link" @click.prevent="navigateTo('admin')">管理中心</a>
          </li>
        </ul>
        
        <!-- 认证按钮 -->
        <div class="auth-buttons" v-if="!isLoggedIn">
          <button class="btn btn-outline" @click="navigateToRegister">注册</button>
          <button class="btn btn-primary" @click="navigateToLogin">登录</button>
        </div>
        
        <!-- 已登录用户信息 -->
        <div class="user-profile" v-else>
          <div class="user-avatar-container">
            <img 
              :src="userAvatar" 
              alt="用户头像" 
              class="user-avatar"
              @click="navigateTo('profile')"
              style="cursor: pointer;"
              loading="lazy"
            >
          </div>
          <span class="user-name">{{ user?.username || '用户' }}</span>
          <button class="btn btn-outline" @click="handleLogout">退出</button>
        </div>
      </nav>
    </div>
  </header>
</template>

<script>
import { ref, computed } from 'vue'

export default {
  name: 'Navbar',
  props: {
    isLoggedIn: {
      type: Boolean,
      required: true
    },
    user: {
      type: Object,
      default: null
    },
    navigateTo: {
      type: Function,
      required: true
    },
    navigateToLogin: {
      type: Function,
      required: true
    },
    navigateToRegister: {
      type: Function,
      required: true
    },
    handleLogout: {
      type: Function,
      required: true
    }
  },
  setup(props) {
    const isMenuOpen = ref(false)
    
    const toggleMenu = () => {
      isMenuOpen.value = !isMenuOpen.value
    }
    
    // 计算用户头像URL，使用用户名首字母或默认头像
    const userAvatar = computed(() => {
      if (!props.user?.username) {
        return `https://picsum.photos/seed/default/100`;
      }
      // 使用用户名首字母创建头像
      const initial = props.user.username.charAt(0).toUpperCase();
      return `https://picsum.photos/seed/${initial}${props.user.id}/100`;
    })
    
    return {
      isMenuOpen,
      toggleMenu,
      userAvatar
    }
  }
}
</script>

<style scoped>
:root {
  --primary-color: #C8102E; /* 湘红 */
  --secondary-color: #1E40AF; /* 湘蓝 */
  --accent-color: #D97706; /* 湘金 */
  --bg-color: #F9FAFB;
  --text-color: #1F2937;
  --light-text: #6B7280;
  --heading-font: 'Microsoft YaHei', 'STXihei', 'SimHei', sans-serif;
  --body-font: 'Microsoft YaHei', 'STXihei', 'SimSun', serif;
}

/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: var(--body-font);
}

/* 头部样式 */
header {
  background-color: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
}

.navbar-brand {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.8rem;
    font-weight: bold;
    color: var(--primary-color);
    font-family: var(--heading-font);
    text-decoration: none;
  }

.navbar-logo {
    height: 50px;
    width: auto;
    border-radius: 4px;
  }

.navbar-nav {
  display: flex;
  gap: 2rem;
  list-style: none;
}

.nav-link {
  text-decoration: none;
  color: var(--text-color);
  font-weight: 500;
  transition: color 0.3s;
}

.nav-link:hover {
  color: var(--primary-color);
}

/* 管理员导航链接样式 */
.nav-link.admin-link {
  font-weight: 600;
  color: var(--accent-color);
}

.nav-link.admin-link:hover {
  color: var(--primary-color);
}

/* 移动端菜单按钮 */
.navbar-toggle {
  display: none;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: var(--text-color);
}

/* 认证按钮样式 */
.auth-buttons {
  display: flex;
  gap: 1rem;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-avatar-container {
  position: relative;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--accent-color);
  transition: transform 0.3s ease;
}

.user-avatar:hover {
  transform: scale(1.1);
}

.user-name {
  font-weight: 600;
  color: var(--text-color);
  white-space: nowrap;
}

.btn {
  padding: 0.5rem 1.5rem;
  border-radius: 4px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  text-decoration: none;
}

.btn-outline {
  background-color: transparent;
  border: 1px solid var(--primary-color);
  color: var(--primary-color);
}

.btn-outline:hover {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
  border: none;
}

.btn-primary:hover {
  background-color: #A80D27;
}



/* 响应式调整 */
@media (max-width: 768px) {
  .navbar {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .navbar-nav {
    flex-direction: column;
    gap: 1rem;
    margin-top: 1rem;
    display: none;
    width: 100%;
  }
  
  .navbar-nav.show {
    display: flex;
  }
  
  .auth-buttons,
  .user-profile {
    display: flex;
    margin-top: 1rem;
    width: 100%;
    justify-content: center;
  }
  
  .navbar-toggle {
    display: block;
    position: absolute;
    top: 1rem;
    right: 1rem;
  }
}
</style>