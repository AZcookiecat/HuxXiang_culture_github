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
      <!-- 帖子详情 -->
      <div class="post-detail">
        <!-- 帖子头部 -->
        <div class="post-header">
          <div class="post-author">
            <img :src="currentPost?.authorAvatar" alt="Author avatar" class="author-avatar" loading="lazy" />
            <div class="author-info">
              <span class="author-name">{{ currentPost?.authorName }}</span>
              <span class="post-date">{{ formatDate(currentPost?.createdAt) }}</span>
            </div>
          </div>
          <div class="post-stats">
            <span class="stat-item"><i class="far fa-eye"></i> {{ currentPost?.views }} 浏览</span>
            <span class="stat-item"><i class="far fa-comment"></i> {{ currentPost?.comments }} 评论</span>
            <span class="stat-item"><i class="far fa-thumbs-up"></i> {{ currentPost?.likes }} 点赞</span>
          </div>
        </div>
        
        <!-- 帖子内容 -->
        <div class="post-content">
          <h1 class="post-title">{{ currentPost?.title }}</h1>
          <div class="post-body" v-html="formatPostContent(currentPost?.content)"></div>
          <div class="post-tags">
            <span class="post-category">{{ currentPost?.category }}</span>
          </div>
        </div>
        
        <!-- 帖子操作 -->
        <div class="post-actions">
          <button class="action-btn" @click="toggleLike">
            <i :class="['far', liked ? 'fa-thumbs-up' : 'fa-thumbs-up', liked ? 'liked' : '']"></i>
            <span>{{ liked ? '已点赞' : '点赞' }}</span>
          </button>
          <button class="action-btn">
            <i class="far fa-share-square"></i>
            <span>分享</span>
          </button>
          <button class="action-btn" v-if="isUserAuthor">
            <i class="far fa-edit"></i>
            <span>编辑</span>
          </button>
        </div>
      </div>
      
      <!-- 相关推荐 -->
      <div class="related-posts">
        <h3>相关推荐</h3>
        <div class="related-posts-list">
          <div v-for="post in relatedPosts" :key="post.id" class="related-post-item" @click="goToPostDetail(post.id)">
            <h4 class="related-post-title">{{ post.title }}</h4>
            <div class="related-post-meta">
              <span>{{ post.authorName }}</span>
              <span>{{ formatDate(post.createdAt) }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 评论区 -->
      <div class="comments-section">
        <h3>评论区</h3>
        
        <!-- 发表评论 -->
        <div class="comment-form" v-if="isLoggedIn">
          <textarea v-model="newComment" placeholder="写下您的评论..." rows="3"></textarea>
          <button class="submit-comment-btn" @click="submitComment">发表评论</button>
        </div>
        <div v-else class="login-prompt">
          请先 <router-link to="/login">登录</router-link> 后发表评论
        </div>
        
        <!-- 评论列表 -->
        <div class="comments-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <img :src="comment.authorAvatar" alt="Comment author" class="comment-avatar" loading="lazy" />
            <div class="comment-content">
              <div class="comment-header">
                <span class="comment-author">{{ comment.authorName }}</span>
                <span class="comment-date">{{ formatDate(comment.createdAt) }}</span>
              </div>
              <p class="comment-text">{{ comment.content }}</p>
            </div>
          </div>
        </div>
        
        <!-- 加载更多评论 -->
        <div class="load-more">
          <button class="load-more-btn" @click="loadMoreComments" v-if="hasMoreComments">
            加载更多评论
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

// 模拟帖子数据服务
const mockPostService = {
  getPostById(id) {
    // 这里模拟从API获取帖子详情
    // 在实际应用中，应该替换为真实的API调用
    const posts = [
      {
        id: '1',
        title: '探讨湖湘文化的现代传承与发展',
        content: `
          <p>湖湘文化作为中国传统文化的重要组成部分，在现代社会如何更好地传承和发展是一个值得思考的问题。</p>
          
          <h3>一、湖湘文化的历史渊源</h3>
          <p>湖湘文化是指湖南省境内的地域文化，具有悠久的历史和丰富的内涵。它起源于先秦时期，经过秦汉、唐宋、明清等各个历史时期的发展，形成了独特的文化特色。</p>
          
          <h3>二、湖湘文化的现代价值</h3>
          <p>在现代社会，湖湘文化仍然具有重要的价值。它不仅是湖南人民的精神财富，也是中华民族文化宝库中的重要组成部分。湖湘文化中的"经世致用"、"敢为人先"等精神，对于现代社会的发展仍然具有重要的启示意义。</p>
          
          <h3>三、湖湘文化的传承与创新</h3>
          <p>传承和创新是湖湘文化发展的两大主题。我们需要在传承传统文化精髓的同时，不断进行创新，使湖湘文化适应现代社会的发展需求。</p>
          
          <p>以上是我对湖湘文化现代传承与发展的一些思考，欢迎大家一起探讨。</p>
        `,
        excerpt: '湖湘文化作为中国传统文化的重要组成部分，在现代社会如何更好地传承和发展是一个值得思考的问题...',
        authorName: '文化探索者',
        authorAvatar: 'https://picsum.photos/seed/user1/100/100',
        createdAt: '2023-06-15',
        views: 528,
        comments: 42,
        likes: 89,
        category: '文化讨论'
      },
      {
        id: '2',
        title: '岳麓书院的历史与文化价值',
        content: `
          <p>岳麓书院是中国古代四大书院之一，位于湖南省长沙市岳麓山下，是湖湘文化的重要载体。</p>
          
          <h3>一、岳麓书院的历史沿革</h3>
          <p>岳麓书院始建于北宋开宝九年（公元976年），历经宋、元、明、清各代，至清末光绪二十九年（公元1903年）改为湖南高等学堂，1926年正式定名为湖南大学。千余年来，岳麓书院弦歌不绝、办学不已，故有"千年学府"的美称。</p>
          
          <h3>二、岳麓书院的文化价值</h3>
          <p>岳麓书院不仅是中国古代教育的重要场所，也是中国传统文化的重要载体。它培养了大批杰出人才，如王夫之、魏源、曾国藩、左宗棠、杨昌济等，对中国历史和文化的发展产生了深远的影响。</p>
          
          <h3>三、岳麓书院的建筑特色</h3>
          <p>岳麓书院的建筑风格独特，体现了中国传统文化的精神内涵。它坐落在岳麓山下，依山而建，布局严谨，环境优美，是中国古代书院建筑的典范。</p>
        `,
        excerpt: '岳麓书院是中国古代四大书院之一，位于湖南省长沙市岳麓山下，是湖湘文化的重要载体...',
        authorName: '历史学者',
        authorAvatar: 'https://picsum.photos/seed/user2/100/100',
        createdAt: '2023-06-10',
        views: 389,
        comments: 27,
        likes: 65,
        category: '历史研究'
      },
      {
        id: '3',
        title: '湘绣艺术的魅力与传承',
        content: `
          <p>湘绣是中国四大名绣之一，起源于湖南省长沙、湘潭一带，具有悠久的历史和独特的艺术风格。</p>
          
          <h3>一、湘绣的历史与发展</h3>
          <p>湘绣的历史可以追溯到两千多年前的春秋战国时期。经过汉、唐、宋、元、明、清等各个历史时期的发展，湘绣逐渐形成了自己独特的艺术风格。</p>
          
          <h3>二、湘绣的艺术特色</h3>
          <p>湘绣以其精湛的技艺、丰富的色彩和生动的形象而著称于世。它的主要特点是：针法多样、色彩鲜艳、形象生动、题材广泛。湘绣的代表作品有《百鸟朝凤》、《芙蓉鲤鱼》等。</p>
          
          <h3>三、湘绣的传承与发展</h3>
          <p>在现代社会，湘绣面临着传承和发展的挑战。一方面，我们需要保护和传承传统的湘绣技艺；另一方面，我们也需要创新湘绣的表现形式和题材内容，使湘绣适应现代社会的审美需求。</p>
        `,
        excerpt: '湘绣是中国四大名绣之一，起源于湖南省长沙、湘潭一带，具有悠久的历史和独特的艺术风格...',
        authorName: '艺术爱好者',
        authorAvatar: 'https://picsum.photos/seed/user3/100/100',
        createdAt: '2023-06-05',
        views: 412,
        comments: 33,
        likes: 76,
        category: '传统艺术'
      }
    ]
    
    return posts.find(post => post.id === id)
  },
  
  getRelatedPosts(id) {
    // 模拟获取相关帖子
    const posts = [
      {
        id: '2',
        title: '岳麓书院的历史与文化价值',
        authorName: '历史学者',
        createdAt: '2023-06-10'
      },
      {
        id: '3',
        title: '湘绣艺术的魅力与传承',
        authorName: '艺术爱好者',
        createdAt: '2023-06-05'
      }
    ]
    
    return posts.filter(post => post.id !== id)
  },
  
  getCommentsByPostId(postId) {
    // 模拟获取评论数据
    const comments = [
      {
        id: '1',
        content: '非常赞同您的观点！湖湘文化确实需要在传承中创新，才能更好地适应现代社会的发展需求。',
        authorName: '文化爱好者',
        authorAvatar: 'https://picsum.photos/seed/comment1/100/100',
        createdAt: '2023-06-16'
      },
      {
        id: '2',
        content: '我认为湖湘文化中的"经世致用"精神对现代社会仍然具有重要的指导意义。我们应该继承和发扬这种精神。',
        authorName: '传统文化研究者',
        authorAvatar: 'https://picsum.photos/seed/comment2/100/100',
        createdAt: '2023-06-17'
      },
      {
        id: '3',
        content: '期待看到更多关于湖湘文化的讨论和研究！',
        authorName: '学生',
        authorAvatar: 'https://picsum.photos/seed/comment3/100/100',
        createdAt: '2023-06-18'
      }
    ]
    
    return comments
  }
}

export default {
  name: 'PostDetailPage',
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
    const relatedPosts = ref([])
    const comments = ref([])
    const newComment = ref('')
    const isLoggedIn = ref(false)
    const user = ref(null)
    const liked = ref(false)
    const isUserAuthor = ref(false)
    const hasMoreComments = ref(false)
    
    // 获取帖子详情
    const fetchPostDetail = () => {
      const post = mockPostService.getPostById(postId)
      if (post) {
        currentPost.value = post
        
        // 更新浏览量
        currentPost.value.views += 1
        
        // 获取相关帖子
        relatedPosts.value = mockPostService.getRelatedPosts(postId)
        
        // 获取评论
        comments.value = mockPostService.getCommentsByPostId(postId)
        
        // 检查是否已登录
        const userInfo = localStorage.getItem('user')
        if (userInfo) {
          isLoggedIn.value = true
          user.value = JSON.parse(userInfo)
          
          // 检查是否是作者
          isUserAuthor.value = user.value.username === currentPost.value.authorName
          
          // 检查是否已点赞（模拟）
          const likedPosts = JSON.parse(localStorage.getItem('likedPosts') || '[]')
          liked.value = likedPosts.includes(postId)
        }
      } else {
        // 如果帖子不存在，显示提示
        if (props.showAlert) {
          props.showAlert('帖子不存在或已被删除', 'error')
        }
        // 重定向回社区页面
        setTimeout(() => {
          router.push('/community')
        }, 1500)
      }
    }
    
    // 格式化帖子内容
    const formatPostContent = (content) => {
      if (!content) return ''
      return content.trim()
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }
    
    // 发表评论
    const submitComment = () => {
      if (!newComment.value.trim()) {
        if (props.showAlert) {
          props.showAlert('评论内容不能为空', 'warning')
        }
        return
      }
      
      // 创建新评论
      const newCommentObj = {
        id: Date.now().toString(),
        content: newComment.value.trim(),
        authorName: user.value.username,
        authorAvatar: `https://picsum.photos/seed/${user.value.username}/100/100`,
        createdAt: new Date().toISOString().split('T')[0]
      }
      
      // 添加到评论列表
      comments.value.unshift(newCommentObj)
      
      // 更新评论数
      if (currentPost.value) {
        currentPost.value.comments += 1
      }
      
      // 清空评论框
      newComment.value = ''
      
      // 显示成功提示
      if (props.showAlert) {
        props.showAlert('评论发表成功', 'success')
      }
    }
    
    // 切换点赞状态
    const toggleLike = () => {
      if (!isLoggedIn.value) {
        if (props.showAlert) {
          props.showAlert('请先登录后点赞', 'info')
        }
        return
      }
      
      liked.value = !liked.value
      
      // 更新点赞数
      if (currentPost.value) {
        currentPost.value.likes += liked.value ? 1 : -1
      }
      
      // 保存点赞状态（模拟）
      const likedPosts = JSON.parse(localStorage.getItem('likedPosts') || '[]')
      if (liked.value) {
        likedPosts.push(postId)
      } else {
        const index = likedPosts.indexOf(postId)
        if (index > -1) {
          likedPosts.splice(index, 1)
        }
      }
      localStorage.setItem('likedPosts', JSON.stringify(likedPosts))
    }
    
    // 跳转到其他帖子详情
    const goToPostDetail = (id) => {
      router.push(`/post-detail/${id}`)
    }
    
    // 返回社区页面
    const goBack = () => {
      router.push('/community')
    }
    
    // 加载更多评论
    const loadMoreComments = () => {
      // 模拟加载更多评论
      if (props.showAlert) {
        props.showAlert('已加载全部评论', 'info')
      }
      hasMoreComments.value = false
    }
    
    // 页面加载时获取数据
    onMounted(() => {
      fetchPostDetail()
    })
    
    return {
      currentPost,
      relatedPosts,
      comments,
      newComment,
      isLoggedIn,
      liked,
      isUserAuthor,
      hasMoreComments,
      formatPostContent,
      formatDate,
      submitComment,
      toggleLike,
      goToPostDetail,
      goBack,
      loadMoreComments
    }
  }
}
</script>

<style scoped>
.post-detail-page {
  padding: 2rem 0;
}

.back-button-container {
  padding: 0 1rem;
  margin-bottom: 1rem;
}

.back-button {
  background-color: transparent;
  border: 1px solid #ddd;
  padding: 0.5rem 1rem;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s;
}

.back-button:hover {
  background-color: #f8f9fa;
  border-color: #ccc;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.post-detail {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 2rem;
  margin-bottom: 2rem;
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.author-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: bold;
  color: #333;
  font-size: 1.1rem;
}

.post-date {
  font-size: 0.9rem;
  color: #999;
}

.post-stats {
  display: flex;
  gap: 1.5rem;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.9rem;
  color: #666;
}

.post-title {
  font-size: 2rem;
  color: #333;
  margin: 0 0 1.5rem 0;
  line-height: 1.3;
}

.post-body {
  color: #333;
  line-height: 1.8;
  margin-bottom: 2rem;
}

.post-body p {
  margin-bottom: 1rem;
}

.post-body h3 {
  font-size: 1.5rem;
  color: #333;
  margin: 1.5rem 0 1rem 0;
}

.post-tags {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.post-category {
  background-color: #f0f0f0;
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
  font-size: 0.9rem;
  color: #666;
}

.post-actions {
  display: flex;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.action-btn {
  background-color: transparent;
  border: 1px solid #ddd;
  padding: 0.5rem 1rem;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s;
}

.action-btn:hover {
  background-color: #f8f9fa;
  border-color: #ccc;
}

.action-btn i.liked {
  color: var(--primary-color);
}

.related-posts {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.related-posts h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #333;
  font-size: 1.3rem;
}

.related-posts-list {
  display: grid;
  gap: 1rem;
}

.related-post-item {
  padding: 1rem;
  border: 1px solid #eee;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.related-post-item:hover {
  background-color: #f8f9fa;
  border-color: #ddd;
}

.related-post-title {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  color: #333;
}

.related-post-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.85rem;
  color: #999;
}

.comments-section {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
}

.comments-section h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #333;
  font-size: 1.3rem;
}

.comment-form {
  margin-bottom: 2rem;
}

.comment-form textarea {
  width: 100%;
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: vertical;
  font-size: 1rem;
  margin-bottom: 1rem;
}

.submit-comment-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 4px;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.submit-comment-btn:hover {
  background-color: var(--primary-dark);
}

.login-prompt {
  padding: 1.5rem;
  background-color: #f8f9fa;
  border-radius: 4px;
  text-align: center;
  margin-bottom: 2rem;
  color: #666;
}

.login-prompt a {
  color: var(--primary-color);
  text-decoration: none;
}

.login-prompt a:hover {
  text-decoration: underline;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.comment-item {
  display: flex;
  gap: 1rem;
}

.comment-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.comment-author {
  font-weight: bold;
  color: #333;
}

.comment-date {
  font-size: 0.85rem;
  color: #999;
}

.comment-text {
  color: #333;
  line-height: 1.6;
  margin: 0;
}

.load-more {
  text-align: center;
  margin-top: 2rem;
}

.load-more-btn {
  background-color: transparent;
  border: 1px solid #ddd;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.load-more-btn:hover {
  background-color: #f8f9fa;
  border-color: #ccc;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .post-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .post-title {
    font-size: 1.5rem;
  }
  
  .post-actions {
    flex-wrap: wrap;
  }
  
  .comment-item {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .comment-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.25rem;
  }
}
</style>