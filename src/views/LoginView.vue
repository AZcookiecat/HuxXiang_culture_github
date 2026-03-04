<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h2>湖湘文化数字化平台</h2>
        <p>探索湖湘文化的数字世界</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="usernameOrEmail">用户名或邮箱</label>
          <input
            type="text"
            id="usernameOrEmail"
            v-model="usernameOrEmail"
            placeholder="请输入您的用户名或邮箱"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <input
            type="password"
            id="password"
            v-model="password"
            placeholder="请输入您的密码"
            required
          />
        </div>
        
        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="loading" class="loading">登录中...</div>
        
        <button type="submit" class="login-button" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
        
        <div class="register-link">
          <span>还没有账号？</span>
          <router-link to="/register">立即注册</router-link>
        </div>
      </form>
    </div>
    
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import authService from '../services/authService.js'

export default {
  name: 'LoginView',
  setup() {
    const router = useRouter()
    const usernameOrEmail = ref('')
    const password = ref('')
    const error = ref('')
    const loading = ref(false)

    // 处理登录
    const handleLogin = async () => {
      error.value = ''
      
      if (loading.value) return; // 防止重复提交
      
      // 简单验证
      if (!usernameOrEmail.value || !password.value) {
        error.value = '请填写所有必填字段'
        return
      }
      
      loading.value = true
      
      try {
        const response = await authService.login(usernameOrEmail.value, password.value)
        
        if (response.success) {
          // 触发登录成功事件，通知App.vue更新状态
          window.dispatchEvent(new Event('login-success'))
          
          // 获取重定向路径，如果有的话
          const redirect = router.currentRoute.value.query.redirect || '/'
          // 登录成功后跳转
          router.push(redirect)
        } else {
          error.value = response.message || '登录失败'
        }
      } catch (err) {
        error.value = err.message || '登录失败，请稍后重试'
      } finally {
        loading.value = false
      }
    }

    return {
      usernameOrEmail,
      password,
      error,
      loading,
      handleLogin
    }
  }
}
</script>

<style scoped>
/* 登录容器 - PC优先设计 */
.login-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 2rem;
  background: linear-gradient(135deg, #C8102E 0%, #8B0000 100%);
  position: relative;
  overflow: hidden;
}

/* 背景装饰元素 */
.login-container::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 70%);
  z-index: 0;
}

/* 登录框 - PC优先设计 */
.login-box {
  background: white;
  border-radius: 16px;
  padding: 3rem;
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 500px;
  margin-bottom: 2rem;
  position: relative;
  z-index: 1;
  transition: transform 0.3s ease;
}

.login-box:hover {
  transform: translateY(-5px);
}

/* 登录头部 */
.login-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.login-header h2 {
  color: #2c3e50;
  font-size: 2.4rem;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.login-header p {
  color: #7f8c8d;
  font-size: 1.1rem;
  margin: 0;
}

/* 登录表单 */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* 表单组 */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  color: #34495e;
  font-weight: 600;
  font-size: 1rem;
}

.form-group input {
  padding: 1rem 1.2rem;
  border: 2px solid #ecf0f1;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #C8102E;
  box-shadow: 0 0 0 3px rgba(200, 16, 46, 0.1);
}

/* 错误消息 */
.error-message {
  background-color: #fee;
  color: #e74c3c;
  padding: 0.8rem 1rem;
  border-radius: 8px;
  text-align: center;
  font-size: 0.95rem;
  border-left: 4px solid #e74c3c;
}

/* 加载指示器 */
.loading {
  background-color: #e3f2fd;
  color: #2196f3;
  padding: 0.8rem 1rem;
  border-radius: 8px;
  text-align: center;
  font-size: 0.95rem;
}

/* 登录按钮 */
.login-button {
  background: linear-gradient(135deg, #C8102E 0%, #8B0000 100%);
  color: white;
  border: none;
  padding: 1.1rem;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.3s;
  margin-top: 1rem;
}

.login-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

.login-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.login-button:active {
  transform: translateY(0);
}

/* 注册链接 */
.register-link {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 1rem;
}

.register-link span {
  color: #7f8c8d;
}

.register-link a {
  color: #C8102E;
  text-decoration: none;
  font-weight: 600;
  margin-left: 0.3rem;
  transition: color 0.3s;
}

.register-link a:hover {
  text-decoration: underline;
  color: #8B0000;
}

/* 响应式设计 - 增强多端适配 */

/* 大屏幕桌面（1400px及以上） */
@media (min-width: 1400px) {
  .login-box {
    max-width: 550px;
    padding: 3.5rem;
  }
  
  .login-header h2 {
    font-size: 2.6rem;
  }
}

/* 小屏幕桌面（992px-1199px） */
@media (max-width: 1199px) {
  .login-box {
    max-width: 480px;
  }
  
  .login-header h2 {
    font-size: 2.2rem;
  }
}

/* 平板设备（768px-991px） */
@media (max-width: 991px) {
  .login-box {
    padding: 2.5rem;
    max-width: 450px;
  }
  
  .login-header {
    margin-bottom: 2rem;
  }
  
  .login-header h2 {
    font-size: 2rem;
  }
  
  .login-form {
    gap: 1.2rem;
  }
}

/* 小平板/大屏手机（576px-767px） */
@media (max-width: 767px) {
  .login-container {
    padding: 1.5rem;
  }
  
  .login-box {
    padding: 2rem;
    border-radius: 12px;
  }
  
  .login-header h2 {
    font-size: 1.8rem;
  }
  
  .login-header p {
    font-size: 1rem;
  }
  
  .form-group input {
    padding: 0.9rem 1rem;
  }
  
  .login-button {
    padding: 1rem;
    font-size: 1rem;
  }
}

/* 手机设备（小于576px） */
@media (max-width: 575px) {
  .login-container {
    padding: 1rem;
    justify-content: flex-start;
    padding-top: 3rem;
  }
  
  .login-box {
    padding: 1.5rem;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  }
  
  .login-header h2 {
    font-size: 1.6rem;
  }
  
  .form-group {
    gap: 0.4rem;
  }
  
  .form-group label {
    font-size: 0.95rem;
  }
  
  .form-group input {
    padding: 0.8rem 0.9rem;
    font-size: 0.95rem;
  }
  
  .login-button {
    padding: 0.9rem;
  }
  
  .register-link {
    font-size: 0.95rem;
  }
}
</style>