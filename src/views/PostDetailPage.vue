<template>
  <div class="post-detail-page">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <button class="back-button" @click="goBack">
        <i class="fas fa-arrow-left"></i> 返回社区
      </button>
    </div>

    <!-- 主要内容 -->
    <div class="container">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>正在加载帖子...</p>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="error-message">
        <p>{{ error }}</p>
        <button @click="fetchPostDetail" class="retry-btn">重试</button>
      </div>

      <!-- 帖子详情 -->
      <div class="post-detail" v-else-if="currentPost">
        <!-- 帖子头部 -->
        <div class="post-header">
          <div class="post-author">
            <img :src="currentPost.author?.avatar || 'https://via.placeholder.com/40x40' " alt="Author avatar" class="author-avatar" loading="lazy" />
            <div class="author-info">
              <span class="author-name">{{ currentPost.author?.username }}</span>
              <span class="post-date">{{ formatDate(currentPost.created_at) }}</span>
            </div>
          </div>
          <div class="post-stats">
            <span class="stat-item"><i class="far fa-eye"></i> {{ currentPost.view_count || 0 }} 浏览</span>
            <span class="stat-item"><i class="far fa-comment"></i> <i class="comment-count-loading" v-if="commentsLoading">加载中...</i><span class="comment-count" v-else>{{ commentCount }}</span></span>
            <button @click="toggleLike" class="like-btn" :class="{ liked: isLiked }">
              <i class="far fa-thumbs-up" :class="{ 'fas': isLiked }"></i>
              <span>{{ currentPost.like_count || 0 }} 点赞</span>
            </button>
          </div>
        </div>
        
        <!-- 帖子内容 -->
        <div class="post-content">
          <div class="post-meta">
            <span class="post-category">{{ getCategoryLabel(currentPost.category || '文化讨论') }}</span>
          </div>
          <h1 class="post-title">{{ currentPost.title }}</h1>
          <div class="post-body" v-html="formatPostContent(currentPost.content)"></div>
          <div class="post-actions" v-if="canEditDelete">
            <button class="action-btn" @click="editPost">
              <i class="far fa-edit"></i>
              <span>编辑</span>
            </button>
            <button class="action-btn delete-btn" @click="deletePost">
              <i class="far fa-trash-alt"></i>
              <span>删除</span>
            </button>
          </div>
        </div>
        
        <!-- 相关帖子推荐 -->
        <div v-if="relatedPosts.length > 0" class="related-posts-section">
          <h3 class="section-title">相关推荐</h3>
          <div class="related-posts-container">
            <div 
              v-for="post in relatedPosts" 
              :key="post.id" 
              class="related-post-card"
              @click="goToPost(post.id)"
            >
              <div class="related-post-content">
                <h4 class="related-post-title">{{ post.title }}</h4>
                <p class="related-post-excerpt">{{ getExcerpt(post.content, 100) }}</p>
                <div class="related-post-meta">
                  <span class="related-post-author">{{ post.author?.username || '匿名用户' }}</span>
                  <span class="related-post-date">{{ formatDate(post.created_at) }}</span>
                  <span class="related-post-category">{{ getCategoryLabel(post.category) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 评论区 -->
        <CommentsSection 
          :post-id="postId" 
          :show-alert="showAlert" 
          @comment-added="updateCommentCount"
          @comment-deleted="updateCommentCount"
        />
      </div>
      
      <!-- 帖子不存在 -->
      <div v-else class="error-message">
        <p>帖子不存在或已被删除</p>
        <button @click="goBack" class="back-btn">返回社区</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { request } from '@/services/api.js'  // 导入真实API请求方法
import CommentsSection from '@/components/CommentsSection.vue'  // 导入评论区组件

export default {
  name: 'PostDetailPage',
  components: {
    CommentsSection
  },
  props: {
    showAlert: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    const router = useRouter()
    const route = useRoute()
    
    const postId = route.params.id
    const currentPost = ref(null)
    const loading = ref(true)
    const error = ref(null)
    const commentCount = ref(0)
    const commentsLoading = ref(false)
    const relatedPosts = ref([])  // 相关帖子
    
    // 获取当前用户信息
    const currentUser = computed(() => {
      const userStr = localStorage.getItem('user')
      return userStr ? JSON.parse(userStr) : null
    })
    
    // 检查用户是否有编辑/删除权限
    const canEditDelete = computed(() => {
      if (!currentUser.value || !currentPost.value) return false
      return currentUser.value.id === currentPost.value.author?.id || 
             currentUser.value.role === 'admin'
    })
    
    // 检查当前用户是否已点赞
    const isLiked = computed(() => {
      // 直接使用API返回的点赞状态
      if (currentPost.value && currentPost.value.hasOwnProperty('liked_by_current_user')) {
        return currentPost.value.liked_by_current_user;
      }
      return false; // 默认未点赞
    })
    
    // 更新评论数
    const updateCommentCount = async () => {
      // 重新获取帖子详情，以确保所有数据（包括评论数）都是最新的
      await fetchPostDetail();
    };

    // 获取帖子详情
    const fetchPostDetail = async () => {
      try {
        loading.value = true
        error.value = null
        
        const response = await request(`/community/posts/${postId}`, 'GET')
        
        if (response.success) {
          currentPost.value = response.data
          // 更新评论数
          commentCount.value = response.data?.comment_count || 0
          // 获取相关帖子
          await fetchRelatedPosts(response.data.category)
        } else {
          if (response.error && (response.error.includes('认证') || response.error.includes('令牌') || response.error.includes('token'))) {
            error.value = '登录已过期，请重新登录'
            // 清除过期的令牌
            localStorage.removeItem('access_token')
            localStorage.removeItem('user')
          } else {
            throw new Error(response.message || '获取帖子详情失败')
          }
        }
      } catch (err) {
        console.error('获取帖子详情错误:', err)
        // 检查错误是否与认证有关
        if (err.message.includes('401') || err.message.toLowerCase().includes('unauthorized') || err.message.includes('认证') || err.message.includes('令牌')) {
          error.value = '登录已过期，请重新登录'
          // 清除过期的令牌
          localStorage.removeItem('access_token')
          localStorage.removeItem('user')
        } else {
          error.value = err.message || '获取帖子详情失败'
        }
      } finally {
        loading.value = false
      }
    }
    
    // 获取相关帖子
    const fetchRelatedPosts = async (category) => {
      try {
        const response = await request(`/community/posts/related/${postId}?category=${encodeURIComponent(category)}&limit=2`, 'GET')
        
        if (response.success) {
          relatedPosts.value = response.data || []
        } else {
          console.warn('获取相关帖子失败:', response.message)
          // 如果分类相关的获取失败，尝试获取其他热门帖子
          const fallbackResponse = await request('/community/posts?sortBy=popular&limit=2', 'GET')
          if (fallbackResponse.success) {
            // 过滤掉当前帖子
            relatedPosts.value = (fallbackResponse.data || []).filter(post => post.id !== parseInt(postId))
          }
        }
      } catch (err) {
        console.error('获取相关帖子错误:', err)
        // 发生错误时，获取最新的帖子作为备选
        try {
          const fallbackResponse = await request('/community/posts?sortBy=latest&limit=2', 'GET')
          if (fallbackResponse.success) {
            relatedPosts.value = (fallbackResponse.data || []).filter(post => post.id !== parseInt(postId))
          }
        } catch (fallbackErr) {
          console.error('获取备选相关帖子错误:', fallbackErr)
          relatedPosts.value = []  // 设置为空数组
        }
      }
    }
    
    // 获取评论数量
    const fetchCommentCount = async () => {
      try {
        commentsLoading.value = true
        // 由于帖子详情中已经包含了评论信息，我们可以直接使用
        if (currentPost.value && currentPost.value.comments) {
          commentCount.value = currentPost.value.comments.length
        } else {
          // 如果没有在帖子详情中获取到评论，单独请求评论数
          commentCount.value = currentPost.value?.comment_count || 0
        }
      } catch (err) {
        console.error('获取评论数错误:', err)
        commentCount.value = currentPost.value?.comment_count || 0
      } finally {
        commentsLoading.value = false
      }
    }
    
    // 编辑帖子
    const editPost = () => {
      router.push(`/edit-post/${postId}`)
    }
    
    // 删除帖子
    const deletePost = async () => {
      if (!confirm('确定要删除这个帖子吗？此操作不可撤销。')) {
        return
      }
      
      try {
        const token = localStorage.getItem('access_token')
        if (!token) {
          error.value = '请先登录'
          return
        }
        
        const response = await request(`/community/posts/${postId}`, 'DELETE')
        
        if (response.success) {
          if (props.showAlert) {
            props.showAlert('帖子删除成功', 'success')
          }
          router.push('/community')
        } else {
          throw new Error(response.message || '删除帖子失败')
        }
      } catch (err) {
        console.error('删除帖子错误:', err)
        if (props.showAlert) {
          props.showAlert(err.message || '删除帖子失败', 'error')
        }
      }
    }
    
    // 切换点赞状态
    const toggleLike = async () => {
      console.log('开始点赞操作...');
      const token = localStorage.getItem('access_token');
      if (!token) {
        console.log('未检测到访问令牌');
        if (props.showAlert) {
          props.showAlert('请先登录', 'info')
        }
        router.push('/login')
        return
      }
      
      try {
        console.log('发送点赞请求到:', `/community/posts/${postId}/like`);
        const response = await request(`/community/posts/${postId}/like`, 'POST')
        console.log('点赞API响应:', response);
        
        if (response.success) {
          // 直接更新当前帖子的点赞状态和数量，而不是重新获取整个帖子
          if (currentPost.value) {
            currentPost.value.liked_by_current_user = response.liked;
            currentPost.value.like_count = response.like_count;
          }

          // 更新本地存储的点赞状态
          const likedPosts = JSON.parse(localStorage.getItem('likedPosts') || '[]')
          const postIdInt = parseInt(postId)
          if (response.liked) {
            if (!likedPosts.includes(postIdInt)) {
              likedPosts.push(postIdInt)
            }
            if (props.showAlert) {
              props.showAlert('点赞成功', 'success')
            }
          } else {
            const index = likedPosts.indexOf(postIdInt)
            if (index > -1) {
              likedPosts.splice(index, 1)
            }
            if (props.showAlert) {
              props.showAlert('已取消点赞', 'info')
            }
          }
          localStorage.setItem('likedPosts', JSON.stringify(likedPosts))
          
          console.log('点赞状态更新完成');
        } else {
          throw new Error(response.message || '操作失败')
        }
      } catch (err) {
        console.error('点赞操作错误:', err)
        if (err.message.includes('认证') || err.message.includes('令牌') || err.message.includes('token')) {
          // Token可能过期，清除本地存储
          localStorage.removeItem('access_token')
          localStorage.removeItem('user')
          if (props.showAlert) {
            props.showAlert('登录已过期，请重新登录', 'info')
          }
          router.push('/login')
        } else {
          if (props.showAlert) {
            props.showAlert(err.message || '点赞操作失败', 'error')
          }
        }
      }
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }
    
    // 格式化帖子内容
    const formatPostContent = (content) => {
      if (!content) return ''
      // 简单的内容格式化，换行转为<br>
      return content.replace(/\n/g, '<br>')
    }
    
    // 获取文本摘要
    const getExcerpt = (content, length) => {
      if (!content) return ''
      return content.length > length ? content.substring(0, length) + '...' : content
    }
    
    // 跳转到相关帖子
    const goToPost = (postId) => {
      if (!postId) {
        console.error('Invalid postId:', postId)
        return
      }
      
      // 直接使用 window.location 实现跳转，绕过 Vue Router 的缓存机制
      window.location.href = `/post-detail/${postId}`;
    }
    
    // 返回上一页
    const goBack = () => {
      router.push('/community')
    }
    
    // 分类映射函数
    const getCategoryLabel = (category) => {
      const categoryMap = {
        '文化讨论': '文化讨论',
        '历史研究': '历史研究',
        '传统艺术': '传统艺术',
        '饮食文化': '饮食文化',
        'discussion': '文化讨论',
        'question': '文化讨论',
        'sharing': '文化讨论',
        'activity': '文化讨论',
        'resource': '文化讨论',
        'history': '历史研究',
        'art': '传统艺术',
        'custom': '文化讨论',
        'default': '文化讨论'
      };
      return categoryMap[category] || '文化讨论';
    }
    
    onMounted(() => {
      fetchPostDetail()
    })
    
    return {
      currentPost,
      loading,
      error,
      postId,
      commentCount,
      commentsLoading,
      canEditDelete,
      isLiked,
      editPost,
      deletePost,
      toggleLike,
      formatDate,
      formatPostContent,
      getExcerpt,
      goBack,
      getCategoryLabel,
      relatedPosts,
      goToPost
    }
  }
}
</script>

<style scoped>
.post-detail-page {
  padding: 2rem 0;
  min-height: calc(100vh - 200px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.loading {
  text-align: center;
  padding: 2rem;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid var(--primary-color);
  border-radius: 50%;
  animation: spin 2s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  text-align: center;
  padding: 2rem;
  color: #d32f2f;
}

.error-message p {
  margin-bottom: 1rem;
}

.retry-btn, .back-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.back-button-container {
  margin-bottom: 1rem;
}

.back-button {
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.back-button:hover {
  background: #e9ecef;
}

.post-detail {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.author-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: 600;
  color: #333;
}

.post-date {
  font-size: 0.85rem;
  color: #666;
}

.post-stats {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  color: #666;
  font-size: 0.9rem;
}

.like-btn {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  background: none;
  border: 1px solid #ddd;
  padding: 0.4rem 0.8rem;
  border-radius: 4px;
  cursor: pointer;
  color: #666;
  transition: all 0.3s;
}

.like-btn:hover {
  background: #f8f9fa;
}

.like-btn.liked {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.post-content {
  margin-bottom: 2rem;
}

.post-meta {
  margin-bottom: 1rem;
}

.post-category {
  background: var(--primary-color);
  color: white;
  padding: 0.3rem 0.8rem;
  border-radius: 20px;
  font-size: 0.85rem;
}

.post-title {
  font-size: 1.8rem;
  margin: 0.5rem 0 1.5rem;
  color: #2c3e50;
  line-height: 1.3;
}

.post-body {
  line-height: 1.8;
  color: #34495e;
  margin-bottom: 2rem;
  white-space: pre-line;
}

.post-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  color: #495057;
  transition: all 0.3s;
}

.action-btn:hover {
  background: #e9ecef;
}

.delete-btn {
  background: #ffebee;
  color: #c62828;
  border-color: #ffcdd2;
}

.delete-btn:hover {
  background: #ffcdd2;
}

/* 相关帖子推荐样式 */
.related-posts-section {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #eee;
}

.section-title {
  font-size: 1.4rem;
  color: #2c3e50;
  margin-bottom: 1.5rem;
}

.related-posts-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.related-post-card {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 1.2rem;
  background-color: white; /* 改为白色背景，与帖子详情页一致 */
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.related-post-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  border-color: var(--primary-color);
}

.related-post-content {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.related-post-title {
  font-size: 1.1rem;
  color: #2c3e50;
  margin: 0;
  font-weight: 600;
}

.related-post-excerpt {
  color: #666;
  line-height: 1.5;
  margin: 0;
  font-size: 0.95rem;
}

.related-post-meta {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  margin-top: 0.8rem;
}

.related-post-author {
  font-size: 0.85rem;
  color: #555;
}

.related-post-date {
  font-size: 0.85rem;
  color: #888;
}

.related-post-category {
  background: #e3f2fd;
  color: #1976d2;
  padding: 0.2rem 0.6rem;
  border-radius: 12px;
  font-size: 0.8rem;
}

@media (max-width: 768px) {
  .post-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .post-stats {
    justify-content: space-around;
  }
  
  .post-title {
    font-size: 1.5rem;
  }
  
  .post-actions {
    flex-direction: column;
  }
  
  
}
</style>