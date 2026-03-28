<template>
  <div class="post-editor-page">
    <div class="container">
      <button class="back-button" @click="goBack">
        <i class="fas fa-arrow-left"></i>
        <span>返回帖子</span>
      </button>

      <section class="editor-card">
        <div class="editor-header">
          <p class="eyebrow">Edit Post</p>
          <h1>编辑帖子</h1>
          <p>修改帖子内容后重新保存。</p>
        </div>

        <div v-if="loading" class="state-card">正在加载帖子内容...</div>

        <form v-else class="editor-form" @submit.prevent="submitForm">
          <label>
            <span>标题</span>
            <input
              v-model="formData.title"
              type="text"
              maxlength="200"
              placeholder="请输入帖子标题"
              required
            />
          </label>

          <label>
            <span>分类</span>
            <select v-model="formData.category" required>
              <option value="">请选择分类</option>
              <option value="文化讨论">文化讨论</option>
              <option value="历史研究">历史研究</option>
              <option value="传统艺术">传统艺术</option>
              <option value="饮食文化">饮食文化</option>
            </select>
          </label>

          <label>
            <span>内容</span>
            <textarea
              v-model="formData.content"
              rows="14"
              maxlength="5000"
              placeholder="请输入帖子正文"
              required
            ></textarea>
          </label>

          <div class="form-actions">
            <button type="button" class="secondary-btn" @click="goBack">取消</button>
            <button type="submit" class="primary-btn" :disabled="isSubmitting">
              {{ isSubmitting ? '保存中...' : '保存修改' }}
            </button>
          </div>
        </form>
      </section>
    </div>
  </div>
</template>

<script>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { request } from '@/services/api.js'

export default {
  name: 'EditPostPage',
  props: {
    showAlert: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    const postId = route.params.id
    const loading = ref(true)
    const isSubmitting = ref(false)
    const formData = ref({
      title: '',
      category: '',
      content: ''
    })

    const fetchPostDetail = async () => {
      loading.value = true

      try {
        const response = await request(`/community/posts/${postId}`, 'GET')
        if (!response.success) {
          throw new Error(response.message || '获取帖子详情失败')
        }

        formData.value = {
          title: response.data.title || '',
          category: response.data.category || '',
          content: response.data.content || ''
        }
      } catch (error) {
        props.showAlert?.(error.message || '获取帖子详情失败', 'error')
        router.push('/community')
      } finally {
        loading.value = false
      }
    }

    const submitForm = async () => {
      const payload = {
        title: formData.value.title.trim(),
        category: formData.value.category,
        content: formData.value.content.trim()
      }

      if (!payload.title || !payload.category || !payload.content) {
        props.showAlert?.('请完整填写标题、分类和内容。', 'error')
        return
      }

      isSubmitting.value = true

      try {
        const response = await request(`/community/posts/${postId}`, 'PUT', payload)
        if (!response.success) {
          throw new Error(response.message || '更新帖子失败')
        }

        props.showAlert?.('帖子已更新', 'success')
        router.push(`/post-detail/${postId}`)
      } catch (error) {
        props.showAlert?.(error.message || '更新帖子失败', 'error')
      } finally {
        isSubmitting.value = false
      }
    }

    const goBack = () => {
      router.push(`/post-detail/${postId}`)
    }

    onMounted(fetchPostDetail)

    return {
      formData,
      goBack,
      isSubmitting,
      loading,
      submitForm
    }
  }
}
</script>

<style scoped>
.post-editor-page {
  padding: 2rem 0 3rem;
}

.container {
  max-width: 920px;
  margin: 0 auto;
  padding: 0 1rem;
}

.back-button {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
  border: none;
  background: transparent;
  color: #374151;
  cursor: pointer;
  margin-bottom: 1rem;
  font-weight: 600;
}

.editor-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 20px;
  box-shadow: 0 18px 38px rgba(15, 23, 42, 0.08);
  padding: 1.5rem;
}

.editor-header {
  margin-bottom: 1.25rem;
}

.eyebrow {
  margin: 0 0 0.45rem;
  color: #9b5a30;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-size: 0.78rem;
  font-weight: 700;
}

.editor-header h1 {
  margin: 0;
  color: #1f2937;
}

.editor-header p:last-child {
  color: #6b7280;
  margin-top: 0.55rem;
}

.state-card {
  padding: 2rem;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  background: #f8fafc;
  color: #6b7280;
}

.editor-form {
  display: grid;
  gap: 1rem;
}

.editor-form label {
  display: grid;
  gap: 0.5rem;
}

.editor-form span {
  font-weight: 600;
  color: #374151;
}

.editor-form input,
.editor-form select,
.editor-form textarea {
  width: 100%;
  box-sizing: border-box;
  border: 1px solid #d6dbe3;
  border-radius: 12px;
  padding: 0.9rem 1rem;
  font-size: 1rem;
}

.editor-form textarea {
  resize: vertical;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.secondary-btn,
.primary-btn {
  border: none;
  border-radius: 999px;
  padding: 0.78rem 1.2rem;
  cursor: pointer;
  font-weight: 600;
}

.secondary-btn {
  background: #eef2f7;
  color: #374151;
}

.primary-btn {
  background: var(--primary-color);
  color: #fff;
}

.primary-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 640px) {
  .form-actions {
    flex-direction: column;
  }

  .secondary-btn,
  .primary-btn {
    width: 100%;
  }
}
</style>
