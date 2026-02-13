<template>
  <div class="community-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="container">
        <h1>互动社区</h1>
        <p>加入我们的社区，分享您对湖湘文化的见解与热情</p>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="container">
      <!-- 社区导航 -->
      <div class="community-nav">
        <button 
          class="nav-tab" 
          :class="{ active: activeTab === 'forum' }" 
          @click="switchTab('forum')"
        >
          <i class="fas fa-comments"></i> 文化论坛
        </button>
        <button 
          class="nav-tab" 
          :class="{ active: activeTab === 'activities' }" 
          @click="switchTab('activities')"
        >
          <i class="fas fa-calendar-alt"></i> 文化活动
        </button>
        <button 
          class="nav-tab" 
          :class="{ active: activeTab === 'contributions' }" 
          @click="switchTab('contributions')"
        >
          <i class="fas fa-hand-holding-heart"></i> 内容贡献
        </button>
        <button 
          class="nav-tab" 
          :class="{ active: activeTab === 'feedback' }" 
          @click="switchTab('feedback')"
        >
          <i class="fas fa-comment-dots"></i> 意见反馈
        </button>
      </div>

      <!-- 论坛内容 -->
      <div v-if="activeTab === 'forum'" class="forum-content">
        <!-- 发帖按钮 -->
        <div class="forum-actions">
          <button class="create-post-btn" @click="createNewPost">
            <i class="fas fa-plus-circle"></i> 发布新主题
          </button>
          
          <div class="forum-filters">
            <select v-model="forumSortBy" @change="sortForumPosts">
              <option value="latest">最新发布</option>
              <option value="popular">热门讨论</option>
              <option value="comments">评论最多</option>
            </select>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading">
          <p>正在加载帖子...</p>
        </div>

        <!-- 论坛帖子列表 -->
        <div v-else class="forum-posts">
          <div v-for="post in currentPagePosts" :key="post.id" class="forum-post">
            <div class="post-header">
              <div class="post-author">
                <img :src="post.author?.avatar || 'https://via.placeholder.com/40x40' " alt="Author avatar" class="author-avatar" loading="lazy" />
                <div class="author-info">
                  <span class="author-name">{{ post.author?.username }}</span>
                  <span class="post-date">{{ formatDate(post.created_at) }}</span>
                </div>
              </div>
              <div class="post-stats">
                <span class="stat-item"><i class="far fa-eye"></i> {{ post.view_count }} 浏览</span>
                <span class="stat-item"><i class="far fa-comment"></i> {{ post.comment_count }} 评论</span>
                <span class="stat-item"><i class="far fa-thumbs-up"></i> {{ post.like_count }} 点赞</span>
              </div>
            </div>
            <div class="post-content">
              <h3 class="post-title">{{ post.title }}</h3>
              <p class="post-excerpt">{{ post.summary }}</p>
            </div>
            <div class="post-footer">
              <span class="post-category">{{ getCategoryLabel(post.category) }}</span>
              <button class="view-post-btn" @click="viewPostDetails(post.id)">查看详情</button>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="!loading && forumPosts.length > 0" class="pagination">
          <button class="pagination-btn" :disabled="currentPage === 1" @click="goToPage(currentPage - 1)">
            上一页
          </button>
          <button class="pagination-btn" v-for="page in forumPages" :key="page" :class="{ active: currentPage === page }" @click="goToPage(page)">
            {{ page }}
          </button>
          <button class="pagination-btn" :disabled="currentPage === forumPages" @click="goToPage(currentPage + 1)">
            下一页
          </button>
        </div>
      </div>

      <!-- 文化活动 -->
      <div v-if="activeTab === 'activities'" class="activities-content">
        <div class="activities-grid">
          <div v-for="activity in activities" :key="activity.id" class="activity-card">
            <div class="activity-image">
              <img :src="activity.imageUrl" :alt="activity.title" />
              <div class="activity-date">
                {{ formatActivityDate(activity.startDate) }}
              </div>
            </div>
            <div class="activity-info">
              <h3 class="activity-title">{{ activity.title }}</h3>
              <p class="activity-location"><i class="fas fa-map-marker-alt"></i> {{ activity.location }}</p>
              <p class="activity-description">{{ activity.description }}</p>
              <div class="activity-footer">
                <button class="activity-btn" @click="joinActivity(activity.id)">
                  {{ activity.isJoined ? '已报名' : '我要报名' }}
                </button>
                <span class="participants-count">{{ activity.participants }} 人已报名</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 内容贡献 -->
      <div v-if="activeTab === 'contributions'" class="contributions-content">
        <div class="contribution-form">
          <h2>贡献您的文化资源</h2>
          <form @submit.prevent="submitContribution">
            <div class="form-group">
              <label for="contribution-title">标题</label>
              <input type="text" id="contribution-title" v-model="contribution.title" required placeholder="请输入资源标题" />
            </div>
            <div class="form-group">
              <label for="contribution-category">分类</label>
              <select id="contribution-category" v-model="contribution.category" required>
                <option value="">请选择分类</option>
                <option value="历史遗迹">历史遗迹</option>
                <option value="传统艺术">传统艺术</option>
                <option value="文学作品">文学作品</option>
                <option value="民俗风情">民俗风情</option>
                <option value="饮食文化">饮食文化</option>
                <option value="建筑风格">建筑风格</option>
              </select>
            </div>
            <div class="form-group">
              <label for="contribution-description">描述</label>
              <textarea id="contribution-description" v-model="contribution.description" required rows="4" placeholder="请详细描述您要分享的文化资源"></textarea>
            </div>
            <div class="form-group">
              <label for="contribution-image">上传图片 (可选)</label>
              <input type="file" id="contribution-image" accept="image/*" @change="handleImageUpload" />
            </div>
            <button type="submit" class="submit-btn">提交贡献</button>
          </form>
        </div>
        
        <div class="contribution-guidelines">
          <h3>贡献指南</h3>
          <ul>
            <li>请确保您分享的内容与湖湘文化相关</li>
            <li>请勿上传侵犯他人版权的内容</li>
            <li>内容将经过审核后才会发布到平台上</li>
            <li>优质贡献者将获得平台颁发的荣誉证书</li>
          </ul>
        </div>
      </div>

      <!-- 意见反馈 -->
      <div v-if="activeTab === 'feedback'" class="feedback-content">
        <div class="feedback-form">
          <h2>给我们留言</h2>
          <form @submit.prevent="submitFeedback">
            <div class="form-group">
              <label for="feedback-type">反馈类型</label>
              <select id="feedback-type" v-model="feedback.type" required>
                <option value="">请选择反馈类型</option>
                <option value="建议">功能建议</option>
                <option value="bug">问题反馈</option>
                <option value="compliment">表扬鼓励</option>
                <option value="other">其他</option>
              </select>
            </div>
            <div class="form-group">
              <label for="feedback-content">反馈内容</label>
              <textarea id="feedback-content" v-model="feedback.content" required rows="5" placeholder="请详细描述您的反馈内容"></textarea>
            </div>
            <div class="form-group">
              <label for="feedback-contact">联系方式 (可选)</label>
              <input type="text" id="feedback-contact" v-model="feedback.contact" placeholder="邮箱或电话，方便我们回复您" />
            </div>
            <button type="submit" class="submit-btn">提交反馈</button>
          </form>
        </div>
        
        <div class="recent-feedback" v-if="recentFeedback.length > 0">
          <h3>近期反馈回复</h3>
          <div v-for="item in recentFeedback" :key="item.id" class="feedback-item">
            <div class="feedback-header">
              <span class="feedback-type">{{ item.type }}</span>
              <span class="feedback-date">{{ item.date }}</span>
            </div>
            <p class="feedback-question">{{ item.question }}</p>
            <div class="feedback-reply">
              <p>{{ item.reply }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { request } from '@/services/api.js'

export default {
  name: 'CommunityPage',
  props: {
    showAlert: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    const router = useRouter()
    
    // 状态管理
    const activeTab = ref('forum')
    const currentPage = ref(1)
    const forumSortBy = ref('latest')
    const forumPages = ref(1)
    const loading = ref(true)
    
    // 论坛帖子数据 - 初始为空数组，后续从API获取
    const forumPosts = ref([])
    
    // 文化活动数据
    const activities = ref([
      {
        id: '1',
        title: '湖湘文化艺术节',
        description: '一场集音乐、舞蹈、戏剧、美术于一体的综合性文化艺术盛宴，展示湖湘文化的独特魅力。',
        imageUrl: '/src/assets/Imgs/湖湘文化.png',
        startDate: '2023-07-15',
        endDate: '2023-07-20',
        location: '长沙市湖南大剧院',
        participants: 356,
        isJoined: false
      },
      {
        id: '2',
        title: '岳麓书院文化讲座系列',
        description: '邀请知名学者解读湖湘文化经典，探讨传统文化在当代的价值和意义。',
        imageUrl: '/src/assets/Imgs/岳麓书院1.jpg',
        startDate: '2023-07-08',
        endDate: '2023-07-08',
        location: '岳麓书院明伦堂',
        participants: 128,
        isJoined: true
      },
      {
        id: '3',
        title: '湘绣技艺体验工作坊',
        description: '由资深湘绣艺人亲自指导，让参与者亲身体验湘绣的制作过程，感受传统工艺的魅力。',
        imageUrl: '/src/assets/Imgs/湘绣1.jpg',
        startDate: '2023-07-22',
        endDate: '2023-07-22',
        location: '湖南省博物馆',
        participants: 85,
        isJoined: false
      }
    ])
    
    // 内容贡献表单数据
    const contribution = ref({
      title: '',
      category: '',
      description: '',
      image: null
    })
    
    // 意见反馈表单数据
    const feedback = ref({
      type: '',
      content: '',
      contact: ''
    })
    
    // 近期反馈数据
    const recentFeedback = ref([
      {
        id: '1',
        type: '建议',
        question: '希望平台能够增加更多湖湘地方方言的学习资源',
        reply: '感谢您的建议！我们正在策划湖湘方言保护与传承项目，敬请期待。',
        date: '2023-06-12'
      },
      {
        id: '2',
        type: '问题反馈',
        question: '移动端浏览时图片加载较慢，希望能够优化',
        reply: '您好，我们已经注意到这个问题，技术团队正在进行图片加载优化，预计下周完成。',
        date: '2023-06-08'
      }
    ])
    
    // 从API获取帖子列表
    const fetchForumPosts = async () => {
      try {
        loading.value = true
        // 添加分页参数：每页3条记录
        const params = new URLSearchParams({
          page: currentPage.value,
          limit: 3,
          sortBy: forumSortBy.value
        });
        const response = await request(`/community/posts?${params.toString()}`, 'GET')
        
        if (response.success) {
          forumPosts.value = response.data
          // 根据后端返回的分页信息更新总页数
          forumPages.value = response.pagination?.pages || 1
        } else {
          throw new Error(response.message || '获取帖子列表失败')
        }
      } catch (err) {
        console.error('获取帖子列表错误:', err)
        if (props.showAlert) {
          props.showAlert(err.message || '获取帖子列表失败', 'error')
        }
      } finally {
        loading.value = false
      }
    }
    
    // 方法：切换标签页
    const switchTab = (tab) => {
      activeTab.value = tab
      // 当切换到论坛标签时，刷新帖子列表
      if (tab === 'forum') {
        fetchForumPosts()
      }
    }
    
    // 方法：创建新帖子
    const createNewPost = () => {
      router.push('/create-post')
    }
    
    // 方法：查看帖子详情
    const viewPostDetails = (postId) => {
      router.push(`/post-detail/${postId}`)
    }
    
    // 方法：论坛帖子排序
    const sortForumPosts = () => {
      // 调用获取帖子列表函数，这会使用当前的排序方式
      fetchForumPosts()
    }
    
    // 计算属性：当前页显示的帖子
    const currentPagePosts = computed(() => {
      // 直接返回从后端获取的当前页数据，不需要前端再进行切片
      return forumPosts.value;
    })
    
    // 方法：分页导航
    const goToPage = (page) => {
      if (page >= 1 && page <= forumPages.value) {
        currentPage.value = page
        // 切换页面时重新获取数据
        fetchForumPosts()
      }
    }
    
    // 方法：报名活动
    const joinActivity = (activityId) => {
      const activity = activities.value.find(a => a.id === activityId)
      if (activity) {
        if (activity.isJoined) {
          if (props.showAlert) {
            props.showAlert('您已取消报名该活动', 'success')
          }
          activity.isJoined = false
          activity.participants--
        } else {
          if (props.showAlert) {
            props.showAlert('报名成功！我们将通过站内信通知您活动详情', 'success')
          }
          activity.isJoined = true
          activity.participants++
        }
      }
    }
    
    // 方法：提交内容贡献
    const submitContribution = () => {
      if (props.showAlert) {
        props.showAlert('感谢您的贡献！我们将尽快审核您提交的内容', 'success')
      }
      // 重置表单
      contribution.value = {
        title: '',
        category: '',
        description: '',
        image: null
      }
    }
    
    // 方法：处理图片上传
    const handleImageUpload = (event) => {
      const file = event.target.files[0]
      if (file) {
        contribution.value.image = file
      }
    }
    
    // 方法：提交反馈
    const submitFeedback = () => {
      if (props.showAlert) {
        props.showAlert('感谢您的反馈！我们会认真对待每一条建议', 'success')
      }
      // 重置表单
      feedback.value = {
        type: '',
        content: '',
        contact: ''
      }
    }
    
    // 方法：格式化日期
    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }
    
    // 方法：格式化活动日期
    const formatActivityDate = (dateString) => {
      const date = new Date(dateString)
      const month = date.getMonth() + 1
      const day = date.getDate()
      return `${month}月${day}日`
    }

    // 分类映射函数 - 只保留四种主要分类
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

    // 组件挂载时获取帖子列表
    onMounted(() => {
      fetchForumPosts()
    })
    
    return {
      activeTab,
      currentPage,
      forumSortBy,
      forumPages,
      forumPosts,
      currentPagePosts,
      loading,
      activities,
      contribution,
      feedback,
      recentFeedback,
      switchTab,
      createNewPost,
      viewPostDetails,
      sortForumPosts,
      goToPage,
      joinActivity,
      submitContribution,
      handleImageUpload,
      submitFeedback,
      formatDate,
      formatActivityDate,
      fetchForumPosts,
      getCategoryLabel
    }
  }
}
</script>

<style scoped>
.community-page {
  padding: 2rem 0;
}

.page-header {
  background-size: cover;
  background-position: center;
  color: white;
  padding: 2rem 0;
  text-align: center;
  margin-top: 20px;
}

.page-header h1 {
  color: black;
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.page-header p {
  color: black;
  font-size: 1.2rem;
  max-width: 800px;
  margin: 0 auto;
}

.community-nav {
  display: flex;
  margin: 2rem 0;
  background-color: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
}

.nav-tab {
  flex: 1;
  background-color: transparent;
  border: none;
  padding: 1rem;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  gap: 0.5rem;
}

.nav-tab:hover {
  background-color: #e9ecef;
}

.nav-tab.active {
  background-color: var(--primary-color);
  color: white;
}

/* 论坛样式 */
.forum-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.create-post-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: background-color 0.3s;
}

.create-post-btn:hover {
  background-color: var(--primary-dark);
}

.forum-filters select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
}

.forum-post {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1rem;
  transition: box-shadow 0.3s;
}

.forum-post:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.author-avatar {
  width: 40px;
  height: 40px;
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
}

.post-date {
  font-size: 0.85rem;
  color: #999;
}

.post-stats {
  display: flex;
  gap: 1rem;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.85rem;
  color: #666;
}

.post-title {
  margin: 0 0 0.5rem 0;
  font-size: 1.3rem;
  color: #333;
}

.post-excerpt {
  margin: 0 0 1rem 0;
  color: #666;
  line-height: 1.5;
}

.post-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.post-category {
  background-color: #f0f0f0;
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
  font-size: 0.85rem;
  color: #666;
}

.view-post-btn {
  background-color: var(--secondary-color);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  cursor: pointer;
  border-radius: 4px;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.view-post-btn:hover {
  background-color: var(--secondary-dark);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.5rem;
  margin: 2rem 0;
}

.pagination-btn {
  background-color: white;
  border: 1px solid #ddd;
  padding: 0.5rem 1rem;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.pagination-btn:hover:not(:disabled) {
  background-color: #f8f9fa;
}

.pagination-btn.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.pagination-btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.loading {
  text-align: center;
  padding: 2rem;
  font-size: 1.2rem;
  color: #666;
}

/* 活动样式 */
.activities-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.activity-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
}

.activity-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.activity-image {
  position: relative;
  height: 200px;
  overflow: hidden;
}

.activity-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.activity-card:hover .activity-image img {
  transform: scale(1.05);
}

.activity-date {
  position: absolute;
  top: 10px;
  left: 10px;
  background-color: var(--primary-color);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  font-weight: bold;
}

.activity-info {
  padding: 1.5rem;
}

.activity-title {
  margin: 0 0 0.5rem 0;
  font-size: 1.3rem;
  color: #333;
}

.activity-location {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #666;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.activity-description {
  color: #666;
  line-height: 1.5;
  margin-bottom: 1.5rem;
}

.activity-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.activity-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.activity-btn:hover {
  background-color: var(--primary-dark);
}

.activity-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.participants-count {
  color: #666;
  font-size: 0.9rem;
}

/* 内容贡献样式 */
.contributions-content {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 2rem;
}

.contribution-form h2 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
  color: #333;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group textarea {
  resize: vertical;
}

.submit-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 4px;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.submit-btn:hover {
  background-color: var(--primary-dark);
}

.contribution-guidelines {
  background-color: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
}

.contribution-guidelines h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #333;
}

.contribution-guidelines ul {
  margin: 0;
  padding-left: 1.5rem;
  color: #666;
}

.contribution-guidelines li {
  margin-bottom: 0.5rem;
  line-height: 1.5;
}

/* 反馈样式 */
.feedback-content {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 2rem;
}

.feedback-form h2 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #333;
}

.recent-feedback h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #333;
}

.feedback-item {
  background-color: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.feedback-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.feedback-type {
  background-color: var(--secondary-color);
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
  font-size: 0.85rem;
}

.feedback-date {
  font-size: 0.85rem;
  color: #999;
}

.feedback-question {
  margin: 0 0 0.5rem 0;
  color: #666;
  line-height: 1.5;
}

.feedback-reply {
  background-color: white;
  padding: 0.75rem;
  border-radius: 4px;
  border-left: 3px solid var(--primary-color);
}

.feedback-reply p {
  margin: 0;
  color: #333;
}

/* 响应式设计 */
@media (max-width: 992px) {
  .contributions-content,
  .feedback-content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .page-header h1 {
    font-size: 2rem;
  }
  
  .community-nav {
    flex-direction: column;
  }
  
  .nav-tab {
    justify-content: flex-start;
    padding: 1rem 1.5rem;
  }
  
  .forum-actions {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .activities-grid {
    grid-template-columns: 1fr;
  }
  
  .post-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .pagination {
    flex-wrap: wrap;
  }
}
</style>