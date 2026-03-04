<template>
  <div class="comments-section">
    <h3 class="comments-title">评论区</h3>
    
    <!-- 发表评论 -->
    <div class="comment-form-container">
      <textarea 
        v-model="newCommentContent" 
        class="comment-input" 
        placeholder="请输入您的评论..." 
        rows="4"
      ></textarea>
      <button 
        @click="submitComment" 
        class="submit-comment-btn" 
        :disabled="!newCommentContent.trim() || submitting"
      >
        {{ submitting ? '提交中...' : '发表评论' }}
      </button>
    </div>
    
    <!-- 评论列表 -->
    <div v-if="comments.length > 0" class="comments-list">
      <div 
        v-for="comment in comments" 
        :key="comment.id" 
        class="comment-item"
      >
        <div class="comment-header">
          <img :src="comment.author?.avatar || 'https://via.placeholder.com/30x30' " alt="Avatar" class="comment-avatar" />
          <div class="comment-author-info">
            <span class="comment-author">{{ comment.author?.username || '匿名用户' }}</span>
            <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
          </div>
          <button 
            v-if="canDeleteComment(comment)" 
            @click="deleteComment(comment.id)" 
            class="delete-comment-btn"
          >
            删除
          </button>
        </div>
        <div class="comment-content">{{ comment.content }}</div>
      </div>
    </div>
    
    <!-- 无评论提示 -->
    <div v-else class="no-comments">
      暂无评论，快来抢沙发吧~
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
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
  emits: ['comment-added', 'comment-deleted'], // 添加事件声明
  setup(props, { emit }) {
    const comments = ref([])
    const newCommentContent = ref('')
    const submitting = ref(false)
    
    // 获取当前用户信息
    const getCurrentUser = () => {
      const userStr = localStorage.getItem('user')
      return userStr ? JSON.parse(userStr) : null
    }
    
    // 获取评论列表
    const fetchComments = async () => {
      try {
        const response = await request(`/community/posts/${props.postId}/comments`, 'GET')
        if (response.success) {
          comments.value = response.data || []
        } else {
          throw new Error(response.message || '获取评论失败')
        }
      } catch (err) {
        console.error('获取评论错误:', err)
        if (props.showAlert) {
          props.showAlert(err.message || '获取评论失败', 'error')
        }
      }
    }
    
    // 提交评论
    const submitComment = async () => {
      if (!newCommentContent.value.trim()) {
        if (props.showAlert) {
          props.showAlert('评论内容不能为空', 'error')
        }
        return
      }
      
      if (!localStorage.getItem('access_token')) {
        if (props.showAlert) {
          props.showAlert('请先登录', 'info')
        }
        return
      }
      
      submitting.value = true
      
      try {
        const response = await request(`/community/posts/${props.postId}/comments`, 'POST', {
          content: newCommentContent.value
        })
        
        if (response.success) {
          if (props.showAlert) {
            props.showAlert('评论发表成功', 'success')
          }
          newCommentContent.value = ''
          await fetchComments() // 刷新评论列表
          emit('comment-added') // 发出评论添加事件
        } else {
          throw new Error(response.message || '发表评论失败')
        }
      } catch (err) {
        console.error('发表评论错误:', err)
        if (props.showAlert) {
          props.showAlert(err.message || '发表评论失败', 'error')
        }
      } finally {
        submitting.value = false
      }
    }
    
    // 删除评论
    const deleteComment = async (commentId) => {
      if (!confirm('确定要删除这条评论吗？')) {
        return
      }
      
      try {
        const response = await request(`/community/comments/${commentId}`, 'DELETE')
        
        if (response.success) {
          if (props.showAlert) {
            props.showAlert('评论删除成功', 'success')
          }
          await fetchComments() // 刷新评论列表
          emit('comment-deleted') // 发出评论删除事件
        } else {
          throw new Error(response.message || '删除评论失败')
        }
      } catch (err) {
        console.error('删除评论错误:', err)
        if (props.showAlert) {
          props.showAlert(err.message || '删除评论失败', 'error')
        }
      }
    }
    
    // 检查是否可以删除评论
    const canDeleteComment = (comment) => {
      const currentUser = getCurrentUser()
      if (!currentUser || !comment.author) return false
      return currentUser.id === comment.author.id || currentUser.role === 'admin'
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }
    
    onMounted(() => {
      fetchComments()
    })
    
    return {
      comments,
      newCommentContent,
      submitting,
      submitComment,
      deleteComment,
      canDeleteComment,
      formatDate
    }
  }
}
</script>

<style scoped>
.comments-section {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #eee;
}

.comments-title {
  font-size: 1.4rem;
  color: #2c3e50;
  margin-bottom: 1.5rem;
}

.comment-form-container {
  margin-bottom: 2rem;
}

.comment-input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  resize: vertical;
  margin-bottom: 0.5rem;
  box-sizing: border-box;
}

.comment-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(200, 16, 46, 0.1);
}

.submit-comment-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.6rem 1.2rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.submit-comment-btn:hover:not(:disabled) {
  background-color: #a60e24;
}

.submit-comment-btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.comment-item {
  padding: 1rem;
  border: 1px solid #eee;
  border-radius: 6px;
  background-color: #fafafa;
}

.comment-header {
  display: flex;
  align-items: center;
  margin-bottom: 0.5rem;
  gap: 0.8rem;
}

.comment-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  object-fit: cover;
}

.comment-author-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.comment-author {
  font-weight: 600;
  color: #333;
  font-size: 0.95rem;
}

.comment-date {
  font-size: 0.8rem;
  color: #666;
}

.delete-comment-btn {
  background-color: #f8f9fa;
  border: 1px solid #dee2e6;
  color: #dc3545;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
}

.delete-comment-btn:hover {
  background-color: #ffebee;
}

.comment-content {
  color: #444;
  line-height: 1.6;
  padding-left: 38px;
}

.no-comments {
  text-align: center;
  color: #666;
  padding: 1.5rem;
  font-style: italic;
}

@media (max-width: 768px) {
  .comment-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .comment-content {
    padding-left: 0;
  }
  
  .delete-comment-btn {
    align-self: flex-end;
  }
}
</style>