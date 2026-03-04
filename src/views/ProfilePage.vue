<template>
  <div class="profile-container">
    <div class="profile-box">
      <div class="profile-header">
        <h2>个人中心</h2>
        <p>管理您的账户信息</p>
      </div>
      
      <div class="profile-content">
        <!-- 头像上传区域 -->
        <div class="avatar-section">
          <div class="avatar-upload">
            <div class="avatar-preview" :style="{ backgroundImage: `url(${avatarPreview})` }">
              <input type="file" ref="avatarInput" class="avatar-input" accept="image/*" @change="handleAvatarChange">
            </div>
            <label for="avatar-upload-btn" class="avatar-upload-label">
              更换头像
            </label>
            <button id="avatar-upload-btn" type="button" class="avatar-upload-button" @click="triggerAvatarUpload">
              选择图片
            </button>
          </div>
        </div>
        
        <!-- 用户信息表单 -->
        <form @submit.prevent="handleUpdateProfile" class="profile-form">
          <div class="form-group">
            <label for="username">昵称</label>
            <input
              type="text"
              id="username"
              v-model="userInfo.username"
              placeholder="请输入您的昵称"
              required
            />
          </div>
          
          <div class="form-group">
            <label for="email">邮箱</label>
            <input
              type="email"
              id="email"
              v-model="userInfo.email"
              placeholder="请输入您的邮箱"
              required
              readonly
            />
            <small class="form-hint">邮箱不可修改</small>
          </div>
          
          <div v-if="error" class="error-message">{{ error }}</div>
          <div v-if="success" class="success-message">{{ success }}</div>
          <div v-if="loading" class="loading">保存中...</div>
          
          <button type="submit" class="save-button" :disabled="loading">
            {{ loading ? '保存中...' : '保存修改' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import authService from '../services/authService.js'

export default {
  name: 'ProfilePage',
  setup() {
    const router = useRouter()
    const avatarInput = ref(null)
    const userInfo = ref({ username: '', email: '', id: '' })
    const avatarPreview = ref('')
    const error = ref('')
    const success = ref('')
    const loading = ref(false)

    // 初始化用户信息
    const initUserInfo = () => {
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        const parsedUser = JSON.parse(storedUser)
        userInfo.value = {
          username: parsedUser.username,
          email: parsedUser.email,
          id: parsedUser.id
        }
        // 设置头像预览，优先使用已上传的头像
        avatarPreview.value = parsedUser.avatar || getUserAvatar(parsedUser.username)
      } else {
        // 未登录，重定向到登录页
        router.push('/login')
      }
    }

    // 获取用户头像（根据用户名生成）
    const getUserAvatar = (username) => {
      // 使用与导航栏相同的头像生成服务
      const initial = username.charAt(0).toUpperCase()
      // 如果用户信息中有ID，则使用ID，否则仅使用用户名首字母
      const userId = userInfo.value?.id ? userInfo.value.id : 'default'
      return `https://picsum.photos/seed/${initial}${userId}/100`
    }

    // 触发头像上传
    const triggerAvatarUpload = () => {
      avatarInput.value?.click()
    }

    // 处理头像更改
    const handleAvatarChange = (event) => {
      const file = event.target.files[0]
      if (file) {
        // 检查文件大小
        if (file.size > 5 * 1024 * 1024) { // 5MB
          error.value = '文件大小不能超过5MB'
          return
        }
        
        // 使用FileReader预览图片
        const reader = new FileReader()
        reader.onload = (e) => {
          avatarPreview.value = e.target.result
          success.value = '头像已更新（保存后生效）'
          setTimeout(() => { success.value = '' }, 3000)
        }
        reader.readAsDataURL(file)
      }
    }

    // 处理更新个人资料
    const handleUpdateProfile = async () => {
      error.value = ''
      success.value = ''
      
      if (loading.value) return; // 防止重复提交
      
      // 简单验证
      if (!userInfo.value.username) {
        error.value = '昵称不能为空'
        return
      }
      
      loading.value = true
      
      try {
        // 准备要发送的数据
        const profileData = {
          username: userInfo.value.username,
        }
        
        // 检查头像是否是base64格式（意味着用户上传了新头像）
        if (avatarPreview.value.startsWith('data:image')) {
          // 将base64转换为文件并上传
          const base64Response = await fetch(avatarPreview.value)
          const blob = await base64Response.blob()
          
          // 创建一个文件对象
          const file = new File([blob], `${userInfo.value.username}_avatar.jpg`, { type: 'image/jpeg' })
          
          // 上传头像文件
          const uploadResponse = await authService.uploadAvatar(file)
          
          if (uploadResponse.success) {
            profileData.avatar = uploadResponse.avatar_url
          } else {
            throw new Error(uploadResponse.message || '头像上传失败')
          }
        } else {
          // 如果不是base64格式，说明是URL，直接使用当前值
          profileData.avatar = avatarPreview.value
        }
        
        // 使用真实的更新服务
        const response = await authService.updateProfile(profileData)
        
        if (response.success) {
          // 更新localStorage中的用户信息
          const updatedUser = {
            ...JSON.parse(localStorage.getItem('user')),
            username: userInfo.value.username,
            avatar: profileData.avatar
          }
          localStorage.setItem('user', JSON.stringify(updatedUser))
          
          // 触发全局状态更新
          window.dispatchEvent(new Event('login-success'))
          
          // 显示成功消息
          success.value = response.message || '个人资料更新成功'
          setTimeout(() => { success.value = '' }, 3000)
        } else {
          error.value = response.message || '更新失败，请稍后重试'
        }
      } catch (err) {
        error.value = err.message || '更新失败，请稍后重试'
      } finally {
        loading.value = false
      }
    }

    // 组件挂载时初始化用户信息
    onMounted(() => {
      initUserInfo()
    })

    return {
      userInfo,
      avatarInput,
      avatarPreview,
      error,
      success,
      loading,
      triggerAvatarUpload,
      handleAvatarChange,
      handleUpdateProfile
    }
  }
}
</script>

<style scoped>
/* 个人中心容器 - PC优先设计 */
.profile-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 2rem;
  background: linear-gradient(135deg, #C8102E 0%, #8B0000 100%); /* 使用湖湘文化特色红色系 */
  position: relative;
  overflow: hidden;
}

/* 背景装饰元素 */
.profile-container::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 70%);
  z-index: 0;
}

/* 个人中心框 - PC优先设计 */
.profile-box {
  background: white;
  border-radius: 16px;
  padding: 3rem;
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 600px;
  margin-bottom: 2rem;
  position: relative;
  z-index: 1;
  transition: transform 0.3s ease;
}

.profile-box:hover {
  transform: translateY(-5px);
}

/* 个人中心头部 */
.profile-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.profile-header h2 {
  color: #2c3e50;
  font-size: 2.4rem;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.profile-header p {
  color: #7f8c8d;
  font-size: 1.1rem;
  margin: 0;
}

/* 头像上传区域 */
.avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: 2.5rem;
}

.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.avatar-preview {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  border: 4px solid #ecf0f1;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
}

.avatar-preview:hover {
  border-color: #C8102E;
  transform: scale(1.05);
}

.avatar-input {
  opacity: 0;
  position: absolute;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.avatar-upload-label {
  color: #7f8c8d;
  font-size: 0.95rem;
  margin-bottom: 0.5rem;
}

.avatar-upload-button {
  background: linear-gradient(135deg, #C8102E 0%, #8B0000 100%); /* 与湖湘文化主题一致 */
  color: white;
  border: none;
  padding: 0.6rem 1.2rem;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.3s;
}

.avatar-upload-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(200, 16, 46, 0.4);
}

/* 个人资料表单 */
.profile-form {
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
  box-shadow: 0 0 0 3px rgba(30, 64, 175, 0.1);
}

.form-group input[readonly] {
  background-color: #f8f9fa;
  cursor: not-allowed;
}

.form-hint {
  color: #95a5a6;
  font-size: 0.85rem;
  margin-top: 0.2rem;
}

/* 错误和成功消息 */
.error-message {
  background-color: #fee;
  color: #e74c3c;
  padding: 0.8rem 1rem;
  border-radius: 8px;
  text-align: center;
  font-size: 0.95rem;
  border-left: 4px solid #e74c3c;
}

.success-message {
  background-color: #e6f7e6;
  color: #27ae60;
  padding: 0.8rem 1rem;
  border-radius: 8px;
  text-align: center;
  font-size: 0.95rem;
  border-left: 4px solid #27ae60;
}

/* 保存按钮 */
.save-button {
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

.save-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

.save-button:active {
  transform: translateY(0);
}

.save-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.save-button:active {
  transform: translateY(0);
}

/* 响应式设计 - 增强多端适配 */

/* 大屏幕桌面（1400px及以上） */
@media (min-width: 1400px) {
  .profile-box {
    max-width: 650px;
    padding: 3.5rem;
  }
  
  .profile-header h2 {
    font-size: 2.6rem;
  }
}

/* 小屏幕桌面（992px-1199px） */
@media (max-width: 1199px) {
  .profile-box {
    max-width: 580px;
  }
  
  .profile-header h2 {
    font-size: 2.2rem;
  }
}

/* 平板设备（768px-991px） */
@media (max-width: 991px) {
  .profile-box {
    padding: 2.5rem;
    max-width: 550px;
  }
  
  .profile-header {
    margin-bottom: 2rem;
  }
  
  .profile-header h2 {
    font-size: 2rem;
  }
  
  .profile-form {
    gap: 1.2rem;
  }
}

/* 小平板/大屏手机（576px-767px） */
@media (max-width: 767px) {
  .profile-container {
    padding: 1.5rem;
  }
  
  .profile-box {
    padding: 2rem;
    border-radius: 12px;
  }
  
  .profile-header h2 {
    font-size: 1.8rem;
  }
  
  .profile-header p {
    font-size: 1rem;
  }
  
  .form-group input {
    padding: 0.9rem 1rem;
  }
  
  .save-button {
    padding: 1rem;
    font-size: 1rem;
  }
}

/* 手机设备（小于576px） */
@media (max-width: 575px) {
  .profile-container {
    padding: 1rem;
    justify-content: flex-start;
    padding-top: 3rem;
  }
  
  .profile-box {
    padding: 1.5rem;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  }
  
  .profile-header h2 {
    font-size: 1.6rem;
  }
  
  .avatar-preview {
    width: 100px;
    height: 100px;
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
  
  .save-button {
    padding: 0.9rem;
  }
}

/* 确保所有CSS规则都正确闭合 */
</style>
