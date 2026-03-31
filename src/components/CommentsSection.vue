<template>
  <section class="comments-section">
    <div class="section-header">
      <h3>评论区</h3>
      <span class="comment-count">{{ comments.length }} 条评论</span>
    </div>

    <div class="comment-form">
      <textarea
        v-model="newCommentContent"
        class="comment-input"
        rows="4"
        placeholder="写下你的评论..."
      ></textarea>
      <div class="comment-form-actions">
        <span class="form-hint">登录后可以发表评论</span>
        <button
          class="submit-comment-btn"
          :disabled="submitting || !newCommentContent.trim()"
          @click="submitComment"
        >
          {{ submitting ? '提交中...' : '发表评论' }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="comments-state">正在加载评论...</div>
    <div v-else-if="comments.length === 0" class="comments-state">暂无评论，欢迎发表第一条评论。</div>

    <div v-else class="comments-list">
      <article v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="comment-header">
          <div class="comment-author-wrap">
            <img
              :src="comment.author?.avatar || 'https://via.placeholder.com/36x36'"
              alt="avatar"
              class="comment-avatar"
            />
            <div>
              <div class="comment-author">{{ comment.author?.username || '匿名用户' }}</div>
              <div class="comment-date">{{ formatDate(comment.created_at) }}</div>
            </div>
          </div>
          <button
            v-if="canDeleteComment(comment)"
            class="delete-comment-btn"
            @click="deleteComment(comment.id)"
          >
            删除
          </button>
        </div>
        <p class="comment-content">{{ comment.content }}</p>
      </article>
    </div>
  </section>
</template>

<script>
import { onMounted, ref, watch } from 'vue'
import { request } from '@/services/api.js'

export default {
  name: 'CommentsSection',
  props: {
    postId: {
      type: [String, Number],
      required: true
    },
    showAlert: {
      type: Function,
      default: null
    }
  },
  emits: ['comment-added', 'comment-deleted'],
  setup(props, { emit }) {
    const comments = ref([])
    const loading = ref(false)
    const submitting = ref(false)
    const newCommentContent = ref('')

    const getCurrentUser = () => {
      const userStr = localStorage.getItem('user')
      return userStr ? JSON.parse(userStr) : null
    }

    const fetchComments = async () => {
      loading.value = true

      try {
        const response = await request(`/community/posts/${props.postId}/comments`, 'GET')
        comments.value = response.success ? response.data || [] : []
      } catch (error) {
        comments.value = []
        props.showAlert?.(error.message || '获取评论失败', 'error')
      } finally {
        loading.value = false
      }
    }

    const submitComment = async () => {
      const content = newCommentContent.value.trim()
      if (!content) {
        return
      }

      if (!localStorage.getItem('access_token')) {
        props.showAlert?.('请先登录后再评论', 'info')
        return
      }

      submitting.value = true

      try {
        const response = await request(`/community/posts/${props.postId}/comments`, 'POST', {
          content
        })

        if (!response.success) {
          throw new Error(response.message || '发表评论失败')
        }

        newCommentContent.value = ''
        await fetchComments()
        emit('comment-added', comments.value.length)
        props.showAlert?.('评论已发布', 'success')
      } catch (error) {
        props.showAlert?.(error.message || '发表评论失败', 'error')
      } finally {
        submitting.value = false
      }
    }

    const deleteComment = async (commentId) => {
      if (!window.confirm('确定删除这条评论吗？')) {
        return
      }

      try {
        const response = await request(`/community/comments/${commentId}`, 'DELETE')
        if (!response.success) {
          throw new Error(response.message || '删除评论失败')
        }

        await fetchComments()
        emit('comment-deleted', comments.value.length)
        props.showAlert?.('评论已删除', 'success')
      } catch (error) {
        props.showAlert?.(error.message || '删除评论失败', 'error')
      }
    }

    const canDeleteComment = (comment) => {
      const currentUser = getCurrentUser()
      if (!currentUser || !comment.author) {
        return false
      }

      return currentUser.id === comment.author.id || currentUser.role === 'admin'
    }

    const formatDate = (dateString) => {
      if (!dateString) {
        return ''
      }

      return new Date(dateString).toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    onMounted(fetchComments)
    watch(() => props.postId, fetchComments)

    return {
      canDeleteComment,
      comments,
      deleteComment,
      formatDate,
      loading,
      newCommentContent,
      submitComment,
      submitting
    }
  }
}
</script>

<style scoped>
.comments-section {
  margin-top: 2.5rem;
  padding-top: 2rem;
  border-top: 1px solid #e8ecf2;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.section-header h3 {
  margin: 0;
  color: #243447;
}

.comment-count {
  color: #6b7280;
  font-size: 0.92rem;
}

.comment-form {
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 1rem;
  margin-bottom: 1.5rem;
}

.comment-input {
  width: 100%;
  resize: vertical;
  min-height: 120px;
  padding: 0.9rem 1rem;
  border: 1px solid #d6dbe3;
  border-radius: 10px;
  font-size: 1rem;
  box-sizing: border-box;
}

.comment-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(200, 16, 46, 0.12);
}

.comment-form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  margin-top: 0.75rem;
}

.form-hint {
  color: #6b7280;
  font-size: 0.9rem;
}

.submit-comment-btn {
  background: var(--primary-color);
  color: #fff;
  border: none;
  border-radius: 999px;
  padding: 0.65rem 1.2rem;
  cursor: pointer;
}

.submit-comment-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.comments-state {
  padding: 1.25rem;
  text-align: center;
  color: #6b7280;
  background: #f8fafc;
  border-radius: 10px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.comment-item {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 1rem;
  background: #fff;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.comment-author-wrap {
  display: flex;
  align-items: center;
  gap: 0.85rem;
}

.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
}

.comment-author {
  font-weight: 600;
  color: #243447;
}

.comment-date {
  color: #6b7280;
  font-size: 0.85rem;
  margin-top: 0.15rem;
}

.comment-content {
  margin: 0.85rem 0 0;
  color: #374151;
  line-height: 1.7;
  white-space: pre-wrap;
}

.delete-comment-btn {
  border: 1px solid #fecaca;
  background: #fff1f2;
  color: #b91c1c;
  border-radius: 999px;
  padding: 0.35rem 0.8rem;
  cursor: pointer;
}

@media (max-width: 768px) {
  .section-header,
  .comment-form-actions,
  .comment-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
