<template>
  <div class="register-container">
    <div class="register-box">
      <div class="register-header">
        <h2>湖湘文化数字化平台</h2>
        <p>创建您的账号，开启文化探索之旅</p>
      </div>
      
      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label for="name">用户名</label>
          <input
            type="text"
            id="name"
            v-model="name"
            placeholder="请输入您的用户名"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="email">邮箱</label>
          <input
            type="email"
            id="email"
            v-model="email"
            placeholder="请输入您的邮箱"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <input
            type="password"
            id="password"
            v-model="password"
            placeholder="请设置您的密码（至少6位）"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input
            type="password"
            id="confirmPassword"
            v-model="confirmPassword"
            placeholder="请再次输入密码"
            required
          />
        </div>
        
        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="success" class="success-message">{{ success }}</div>
        <div v-if="loading" class="loading">注册中...</div>
        
        <button type="submit" class="register-button" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
        
        <div class="login-link">
          <span>已有账号？</span>
          <router-link to="/login">返回登录</router-link>
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
  name: 'RegisterView',
  setup() {
    const router = useRouter()
    const name = ref('')
    const email = ref('')
    const password = ref('')
    const confirmPassword = ref('')
    const error = ref('')
    const success = ref('')
    const loading = ref(false)

    // 处理注册
    const handleRegister = async () => {
      error.value = ''
      success.value = ''
      
      if (loading.value) return; // 防止重复提交
      
      // 简单验证
      if (!name.value || !email.value || !password.value || !confirmPassword.value) {
        error.value = '请填写所有必填字段'
        return
      }
      
      if (password.value !== confirmPassword.value) {
        error.value = '两次输入的密码不一致'
        return
      }
      
      if (password.value.length < 6) {
        error.value = '密码长度至少为6位'
        return
      }
      
      loading.value = true
      
      try {
        const response = await authService.register(name.value, email.value, password.value)
        
        if (response.success) {
          // 显示成功消息
          success.value = response.message || '注册成功！即将跳转到登录页面...'
          
          // 2秒后跳转到登录页面
          setTimeout(() => {
            router.push('/login')
          }, 2000)
        } else {
          error.value = response.message || '注册失败'
        }
      } catch (err) {
        error.value = err.message || '注册失败，请稍后重试'
      } finally {
        loading.value = false
      }
    }

    return {
      name,
      email,
      password,
      confirmPassword,
      error,
      success,
      loading,
      handleRegister
    }
  }
}
</script>

<style scoped>
/* 注册容器 - PC优先设计 */
.register-container {
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
.register-container::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 70%);
  z-index: 0;
}

/* 注册框 - PC优先设计 */
.register-box {
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

.register-box:hover {
  transform: translateY(-5px);
}

/* 注册头部 */
.register-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.register-header h2 {
  color: #2c3e50;
  font-size: 2.4rem;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.register-header p {
  color: #7f8c8d;
  font-size: 1.1rem;
  margin: 0;
}

/* 注册表单 */
.register-form {
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
  border-color: #1E40AF;
  box-shadow: 0 0 0 3px rgba(30, 64, 175, 0.1);
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

/* 成功消息 */
.success-message {
  background-color: #efe;
  color: #27ae60;
  padding: 0.8rem 1rem;
  border-radius: 8px;
  text-align: center;
  font-size: 0.95rem;
  border-left: 4px solid #27ae60;
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

/* 注册按钮 */
.register-button {
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

.register-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

.register-button:active {
  transform: translateY(0);
}

.register-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 登录链接 */
.login-link {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 1rem;
}

.login-link span {
  color: #7f8c8d;
}

.login-link a {
  color: #C8102E;
  text-decoration: none;
  font-weight: 600;
  margin-left: 0.3rem;
  transition: color 0.3s;
}

.login-link a:hover {
  text-decoration: underline;
  color: #8B0000;
}

/* 响应式设计 - 增强多端适配 */

/* 大屏幕桌面（1400px及以上） */
@media (min-width: 1400px) {
  .register-box {
    max-width: 550px;
    padding: 3.5rem;
  }
  
  .register-header h2 {
    font-size: 2.6rem;
  }
}

/* 小屏幕桌面（992px-1199px） */
@media (max-width: 1199px) {
  .register-box {
    max-width: 480px;
  }
  
  .register-header h2 {
    font-size: 2.2rem;
  }
}

/* 平板设备（768px-991px） */
@media (max-width: 991px) {
  .register-box {
    padding: 2.5rem;
    max-width: 450px;
  }
  
  .register-header {
    margin-bottom: 2rem;
  }
  
  .register-header h2 {
    font-size: 2rem;
  }
  
  .register-form {
    gap: 1.2rem;
  }
}

/* 小平板/大屏手机（576px-767px） */
@media (max-width: 767px) {
  .register-container {
    padding: 1.5rem;
  }
  
  .register-box {
    padding: 2rem;
    border-radius: 12px;
  }
  
  .register-header h2 {
    font-size: 1.8rem;
  }
  
  .register-header p {
    font-size: 1rem;
  }
  
  .form-group input {
    padding: 0.9rem 1rem;
  }
  
  .register-button {
    padding: 1rem;
    font-size: 1rem;
  }
}

/* 手机设备（小于576px） */
@media (max-width: 575px) {
  .register-container {
    padding: 1rem;
    justify-content: flex-start;
    padding-top: 3rem;
  }
  
  .register-box {
    padding: 1.5rem;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  }
  
  .register-header h2 {
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
  
  .register-button {
    padding: 0.9rem;
  }
  
  .login-link {
    font-size: 0.95rem;
  }
}
</style>
}