<template>
  <div class="create-post-page">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <button class="back-button" @click="goBack">
        <i class="fas fa-arrow-left"></i> 返回社区
      </button>
    </div>

    <!-- 页面标题 -->
    <div class="page-header">
      <div class="container">
        <h1>发布新主题</h1>
        <p>分享您对湖湘文化的见解与热情</p>
      </div>
    </div>

    <!-- 创建表单 -->
    <div class="container">
      <div class="create-form-container">
        <form @submit.prevent="submitForm" class="create-form">
          <div class="form-group">
            <label for="title">标题 *</label>
            <input 
              type="text" 
              id="title" 
              v-model="formData.title" 
              placeholder="请输入帖子标题"
              required
              maxlength="200"
            />
          </div>

          <div class="form-group">
            <label for="category">分类 *</label>
            <select id="category" v-model="formData.category" required>
              <option value="">请选择分类</option>
              <option value="文化讨论">文化讨论</option>
              <option value="历史研究">历史研究</option>
              <option value="传统艺术">传统艺术</option>
              <option value="饮食文化">饮食文化</option>
            </select>
          </div>

          <div class="form-group">
            <label for="content">内容 *</label>
            <textarea 
              id="content" 
              v-model="formData.content" 
              placeholder="请输入帖子内容，分享您对湖湘文化的见解..."
              rows="15"
              required
              maxlength="5000"
            ></textarea>
          </div>

          <div class="form-actions">
            <button type="button" class="btn-cancel" @click="goBack">取消</button>
            <button type="submit" class="btn-submit" :disabled="isSubmitting">
              {{ isSubmitting ? '发布中...' : '发布帖子' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { request } from '@/services/api.js'

export default {
  name: 'CreatePostPage',
  props: {
    showAlert: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    const router = useRouter()
    
    const formData = ref({
      title: '',
      content: '',
      category: ''
    })
    
    const isSubmitting = ref(false)
    
    // 提交表单
    const submitForm = async () => {
      if (!formData.value.title.trim() || !formData.value.content.trim() || !formData.value.category) {
        if (props.showAlert) {
          props.showAlert('请填写所有必填项', 'error')
        }
        return
      }
      
      isSubmitting.value = true
      
      try {
        const response = await request('/community/posts', 'POST', formData.value)
        
        if (response.success) {
          if (props.showAlert) {
            props.showAlert('帖子发布成功', 'success')
          }
          router.push(`/post-detail/${response.data.id}`)
        } else {
          throw new Error(response.message || '发布帖子失败')
        }
      } catch (err) {
        console.error('发布帖子错误:', err)
        if (props.showAlert) {
          props.showAlert(err.message || '发布帖子失败', 'error')
        }
      } finally {
        isSubmitting.value = false
      }
    }
    
    // 返回上一页
    const goBack = () => {
      router.go(-1)
    }
    
    return {
      formData,
      isSubmitting,
      submitForm,
      goBack
    }
  }
}
</script>

<style scoped>
.create-post-page {
  padding: 2rem 0;
}

.back-button-container {
  position: fixed;
  top: 2rem;
  left: 2rem;
  z-index: 1000;
}

.back-button {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 25px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  font-weight: 500;
  font-size: 1rem;
}

.back-button:hover {
  background-color: #a60e24;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

.page-header {
  margin: 0 0 2rem;
  padding: 2rem 0;
  text-align: center;
  background: linear-gradient(to right, #f8f9fa, #e9ecef);
}

.page-header h1 {
  font-size: 2.5rem;
  color: var(--text-color);
  margin-bottom: 0.5rem;
  font-family: var(--heading-font);
  font-weight: 700;
}

.page-header p {
  color: var(--light-text);
  font-size: 1.2rem;
  margin: 0;
}

.container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 1rem;
}

.create-form-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.create-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 600;
  color: #2c3e50;
  font-size: 1rem;
  margin-bottom: 0.3rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 0.8rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(200, 16, 46, 0.1);
}

.form-group textarea {
  resize: vertical;
  min-height: 250px;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.btn-cancel,
.btn-submit {
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  border: none;
}

.btn-cancel {
  background: #f8f9fa;
  color: #6c757d;
  border: 1px solid #dee2e6;
}

.btn-cancel:hover {
  background: #e9ecef;
}

.btn-submit {
  background: var(--primary-color);
  color: white;
}

.btn-submit:hover:not(:disabled) {
  background: #a60e24;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(200, 16, 46, 0.3);
}

.btn-submit:disabled {
  background: #cccccc;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .page-header h1 {
    font-size: 2rem;
  }
  
  .create-form-container {
    padding: 1.5rem;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-cancel,
  .btn-submit {
    width: 100%;
  }
}
</style>