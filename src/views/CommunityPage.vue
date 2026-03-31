<template>
  <div class="community-page">
    <header class="page-header">
      <div class="container">
        <p class="eyebrow">Community</p>
        <h1>互动社区</h1>
        <p class="subtitle">围绕湖湘文化交流观点、分享资源，也可以报名参加线下活动。</p>
      </div>
    </header>

    <div class="container">
      <div class="community-nav">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="nav-tab"
          :class="{ active: activeTab === tab.key }"
          @click="switchTab(tab.key)"
        >
          <i :class="tab.icon"></i>
          <span>{{ tab.label }}</span>
        </button>
      </div>

      <section v-if="activeTab === 'forum'" class="forum-content">
        <div class="forum-toolbar">
          <div>
            <h2>文化论坛</h2>
            <p>按照发布时间、热度或评论数浏览帖子。</p>
          </div>
          <div class="toolbar-actions">
            <select v-model="forumSortBy" @change="resetForumPagination">
              <option value="latest">最新发布</option>
              <option value="popular">热门讨论</option>
              <option value="comments">评论最多</option>
            </select>
            <button class="create-post-btn" @click="createNewPost">
              <i class="fas fa-pen"></i>
              <span>发布帖子</span>
            </button>
          </div>
        </div>

        <div v-if="loading" class="state-card">正在加载帖子...</div>
        <div v-else-if="forumPosts.length === 0" class="state-card">暂无帖子，来发布第一篇内容。</div>

        <div v-else class="forum-posts">
          <article v-for="post in forumPosts" :key="post.id" class="forum-post">
            <div class="post-header">
              <div class="post-author">
                <img
                  :src="post.author?.avatar || 'https://via.placeholder.com/48x48'"
                  alt="avatar"
                  class="author-avatar"
                />
                <div>
                  <div class="author-name">{{ post.author?.username || '匿名用户' }}</div>
                  <div class="post-date">{{ formatDate(post.created_at) }}</div>
                </div>
              </div>
              <div class="post-category">{{ getCategoryLabel(post.category) }}</div>
            </div>

            <div class="post-body">
              <h3>{{ post.title }}</h3>
              <p>{{ post.summary }}</p>
            </div>

            <div class="post-footer">
              <div class="post-stats">
                <span><i class="far fa-eye"></i> {{ post.view_count }}</span>
                <span><i class="far fa-comment"></i> {{ post.comment_count }}</span>
                <span><i class="far fa-thumbs-up"></i> {{ post.like_count }}</span>
              </div>
              <button class="view-post-btn" @click="viewPostDetails(post.id)">查看详情</button>
            </div>
          </article>
        </div>

        <div v-if="forumPosts.length > 0" class="pagination">
          <button class="pagination-btn" :disabled="pageIndex === 0 || loading" @click="goToPreviousPage">
            上一页
          </button>
          <span class="page-indicator">第 {{ pageIndex + 1 }} 页</span>
          <button class="pagination-btn" :disabled="!hasMore || loading" @click="goToNextPage">
            下一页
          </button>
        </div>
      </section>

      <section v-else-if="activeTab === 'activities'" class="card-section">
        <div class="section-intro">
          <h2>文化活动</h2>
          <p>关注近期活动，报名参与线下文化体验。</p>
        </div>
        <div class="activities-grid">
          <article v-for="activity in activities" :key="activity.id" class="activity-card">
            <img :src="activity.imageUrl" :alt="activity.title" class="activity-image" />
            <div class="activity-info">
              <div class="activity-date">{{ formatActivityDate(activity.startDate) }}</div>
              <h3>{{ activity.title }}</h3>
              <p class="activity-location"><i class="fas fa-map-marker-alt"></i> {{ activity.location }}</p>
              <p class="activity-description">{{ activity.description }}</p>
              <div class="activity-footer">
                <span>{{ activity.participants }} 人已报名</span>
                <button class="activity-btn" @click="joinActivity(activity.id)">
                  {{ activity.isJoined ? '已报名' : '我要报名' }}
                </button>
              </div>
            </div>
          </article>
        </div>
      </section>

      <section v-else-if="activeTab === 'contributions'" class="card-section two-column">
        <div class="panel">
          <h2>内容贡献</h2>
          <p>提交你整理的湖湘文化素材，后台审核后会进入资源库。</p>
          <form class="stack-form" @submit.prevent="submitContribution">
            <input v-model="contribution.title" type="text" placeholder="资源标题" required />
            <select v-model="contribution.category" required>
              <option value="">选择分类</option>
              <option value="历史遗迹">历史遗迹</option>
              <option value="传统艺术">传统艺术</option>
              <option value="文学作品">文学作品</option>
              <option value="民俗风情">民俗风情</option>
              <option value="饮食文化">饮食文化</option>
            </select>
            <textarea
              v-model="contribution.description"
              rows="5"
              placeholder="请描述资源内容与来源"
              required
            ></textarea>
            <input type="file" accept="image/*" @change="handleImageUpload" />
            <button class="primary-btn" type="submit">提交贡献</button>
          </form>
        </div>

        <div class="panel muted">
          <h3>贡献说明</h3>
          <ul>
            <li>仅提交与湖湘文化相关的原创或可授权内容。</li>
            <li>资源会进入人工审核流程。</li>
            <li>图片、文字请尽量附带来源说明。</li>
            <li>优质贡献会在平台首页推荐展示。</li>
          </ul>
        </div>
      </section>

      <section v-else class="card-section two-column">
        <div class="panel">
          <h2>意见反馈</h2>
          <p>告诉我们你希望新增的内容或遇到的问题。</p>
          <form class="stack-form" @submit.prevent="submitFeedback">
            <select v-model="feedback.type" required>
              <option value="">反馈类型</option>
              <option value="建议">功能建议</option>
              <option value="bug">问题反馈</option>
              <option value="compliment">表扬鼓励</option>
              <option value="other">其他</option>
            </select>
            <textarea
              v-model="feedback.content"
              rows="6"
              placeholder="请尽量详细描述问题或建议"
              required
            ></textarea>
            <input v-model="feedback.contact" type="text" placeholder="邮箱或电话（选填）" />
            <button class="primary-btn" type="submit">提交反馈</button>
          </form>
        </div>

        <div class="panel muted">
          <h3>近期回复</h3>
          <div v-for="item in recentFeedback" :key="item.id" class="feedback-item">
            <div class="feedback-meta">
              <span>{{ item.type }}</span>
              <span>{{ item.date }}</span>
            </div>
            <p class="feedback-question">{{ item.question }}</p>
            <p class="feedback-reply">{{ item.reply }}</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { request } from '@/services/api.js'

const CATEGORY_MAP = {
  discussion: '文化讨论',
  question: '文化讨论',
  sharing: '文化讨论',
  activity: '文化讨论',
  resource: '文化讨论',
  history: '历史研究',
  art: '传统艺术',
  custom: '文化讨论',
  '文化讨论': '文化讨论',
  '历史研究': '历史研究',
  '传统艺术': '传统艺术',
  '饮食文化': '饮食文化'
}

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
    const activeTab = ref('forum')
    const forumSortBy = ref('latest')
    const forumPosts = ref([])
    const loading = ref(false)
    const cursorStack = ref([null])
    const pageIndex = ref(0)
    const nextCursor = ref(null)
    const hasMore = ref(false)
    const limit = 6

    const tabs = [
      { key: 'forum', label: '文化论坛', icon: 'fas fa-comments' },
      { key: 'activities', label: '文化活动', icon: 'fas fa-calendar-alt' },
      { key: 'contributions', label: '内容贡献', icon: 'fas fa-hand-holding-heart' },
      { key: 'feedback', label: '意见反馈', icon: 'fas fa-comment-dots' }
    ]

    const activities = ref([
      {
        id: '1',
        title: '湖湘文化艺术节',
        description: '涵盖戏曲、音乐、非遗手作与青年论坛的综合活动周。',
        imageUrl: 'https://picsum.photos/seed/huxiang-art/640/400',
        startDate: '2026-04-12',
        location: '长沙市梅溪湖艺术中心',
        participants: 356,
        isJoined: false
      },
      {
        id: '2',
        title: '岳麓书院公开讲座',
        description: '围绕湖湘学派、书院文化与当代传播进行专题分享。',
        imageUrl: 'https://picsum.photos/seed/yuelu-lecture/640/400',
        startDate: '2026-04-18',
        location: '岳麓书院',
        participants: 128,
        isJoined: true
      },
      {
        id: '3',
        title: '湘绣体验工作坊',
        description: '邀请非遗传承人现场示范，体验传统刺绣工艺。',
        imageUrl: 'https://picsum.photos/seed/xiangxiu/640/400',
        startDate: '2026-04-26',
        location: '湖南省博物院',
        participants: 85,
        isJoined: false
      }
    ])

    const contribution = ref({
      title: '',
      category: '',
      description: '',
      image: null
    })

    const feedback = ref({
      type: '',
      content: '',
      contact: ''
    })

    const recentFeedback = ref([
      {
        id: '1',
        type: '建议',
        question: '希望增加更多方言与地方戏曲资料。',
        reply: '内容团队已纳入下一批专题选题，后续会逐步补充。',
        date: '2026-03-12'
      },
      {
        id: '2',
        type: '问题反馈',
        question: '移动端图片加载速度偏慢。',
        reply: '已在排查资源压缩与缓存策略，本轮更新会同步优化。',
        date: '2026-03-08'
      }
    ])

    const fetchForumPosts = async () => {
      loading.value = true

      try {
        const params = new URLSearchParams({
          limit: String(limit),
          sortBy: forumSortBy.value
        })

        const currentCursor = cursorStack.value[pageIndex.value]
        if (currentCursor) {
          params.set('cursor', String(currentCursor))
        }

        const response = await request(`/community/posts?${params.toString()}`, 'GET')
        if (!response.success) {
          throw new Error(response.message || '获取帖子失败')
        }

        forumPosts.value = response.data || []
        nextCursor.value = response.pagination?.cursor ?? null
        hasMore.value = Boolean(response.pagination?.has_more)
      } catch (error) {
        forumPosts.value = []
        hasMore.value = false
        nextCursor.value = null
        props.showAlert?.(error.message || '获取帖子失败', 'error')
      } finally {
        loading.value = false
      }
    }

    const resetForumPagination = async () => {
      cursorStack.value = [null]
      pageIndex.value = 0
      await fetchForumPosts()
    }

    const goToNextPage = async () => {
      if (!hasMore.value || !nextCursor.value) {
        return
      }

      const nextIndex = pageIndex.value + 1
      if (cursorStack.value.length === nextIndex) {
        cursorStack.value.push(nextCursor.value)
      } else {
        cursorStack.value[nextIndex] = nextCursor.value
      }

      pageIndex.value = nextIndex
      await fetchForumPosts()
    }

    const goToPreviousPage = async () => {
      if (pageIndex.value === 0) {
        return
      }

      pageIndex.value -= 1
      await fetchForumPosts()
    }

    const switchTab = async (tab) => {
      activeTab.value = tab
      if (tab === 'forum' && forumPosts.value.length === 0) {
        await fetchForumPosts()
      }
    }

    const createNewPost = () => {
      router.push('/create-post')
    }

    const viewPostDetails = (postId) => {
      router.push(`/post-detail/${postId}`)
    }

    const joinActivity = (activityId) => {
      const activity = activities.value.find((item) => item.id === activityId)
      if (!activity) {
        return
      }

      activity.isJoined = !activity.isJoined
      activity.participants += activity.isJoined ? 1 : -1

      props.showAlert?.(
        activity.isJoined ? '报名成功，活动提醒将通过站内消息发送。' : '已取消报名。',
        'success'
      )
    }

    const submitContribution = () => {
      props.showAlert?.('内容贡献已提交，审核通过后会展示到平台中。', 'success')
      contribution.value = {
        title: '',
        category: '',
        description: '',
        image: null
      }
    }

    const handleImageUpload = (event) => {
      contribution.value.image = event.target.files?.[0] || null
    }

    const submitFeedback = () => {
      props.showAlert?.('反馈已收到，感谢你的建议。', 'success')
      feedback.value = {
        type: '',
        content: '',
        contact: ''
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) {
        return ''
      }

      return new Date(dateString).toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    }

    const formatActivityDate = (dateString) => {
      const date = new Date(dateString)
      return `${date.getMonth() + 1} 月 ${date.getDate()} 日`
    }

    const getCategoryLabel = (category) => CATEGORY_MAP[category] || '文化讨论'

    onMounted(fetchForumPosts)

    return {
      activeTab,
      activities,
      contribution,
      createNewPost,
      feedback,
      forumPosts,
      forumSortBy,
      formatActivityDate,
      formatDate,
      getCategoryLabel,
      goToNextPage,
      goToPreviousPage,
      handleImageUpload,
      hasMore,
      joinActivity,
      loading,
      pageIndex,
      recentFeedback,
      resetForumPagination,
      submitContribution,
      submitFeedback,
      switchTab,
      tabs,
      viewPostDetails
    }
  }
}
</script>

<style scoped>
.community-page {
  padding-bottom: 3rem;
}

.page-header {
  padding: 3rem 0 2.5rem;
  margin-top: 20px;
  background:
    radial-gradient(circle at top right, rgba(200, 16, 46, 0.16), transparent 32%),
    linear-gradient(135deg, #f8f1ea 0%, #fff 45%, #f4f7fb 100%);
}

.container {
  max-width: 1180px;
  margin: 0 auto;
  padding: 0 1rem;
}

.eyebrow {
  margin: 0 0 0.5rem;
  color: #9b5a30;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-size: 0.8rem;
  font-weight: 700;
}

.page-header h1 {
  margin: 0;
  font-size: clamp(2.3rem, 5vw, 3.4rem);
  color: #1f2937;
}

.subtitle {
  max-width: 640px;
  margin: 0.9rem 0 0;
  color: #4b5563;
  line-height: 1.7;
}

.community-nav {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.85rem;
  margin: 2rem 0 1.75rem;
}

.nav-tab {
  border: 1px solid #e5e7eb;
  background: #fff;
  border-radius: 16px;
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.6rem;
  cursor: pointer;
  color: #374151;
  font-weight: 600;
  transition: all 0.2s ease;
}

.nav-tab.active {
  background: #1f2937;
  color: #fff;
  border-color: #1f2937;
  transform: translateY(-2px);
}

.forum-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 1rem;
  margin-bottom: 1rem;
}

.forum-toolbar h2,
.section-intro h2,
.panel h2 {
  margin: 0 0 0.3rem;
  color: #1f2937;
}

.forum-toolbar p,
.section-intro p,
.panel p {
  margin: 0;
  color: #6b7280;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.toolbar-actions select,
.stack-form input,
.stack-form select,
.stack-form textarea {
  width: 100%;
  border: 1px solid #d6dbe3;
  border-radius: 12px;
  padding: 0.8rem 0.95rem;
  font-size: 1rem;
  box-sizing: border-box;
}

.create-post-btn,
.primary-btn,
.activity-btn,
.view-post-btn,
.pagination-btn {
  border: none;
  border-radius: 999px;
  padding: 0.75rem 1.15rem;
  cursor: pointer;
  font-weight: 600;
}

.create-post-btn,
.primary-btn,
.activity-btn,
.view-post-btn {
  background: var(--primary-color);
  color: #fff;
}

.state-card,
.panel,
.forum-post {
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  background: #fff;
  box-shadow: 0 14px 36px rgba(15, 23, 42, 0.05);
}

.state-card {
  padding: 2rem;
  text-align: center;
  color: #6b7280;
}

.forum-posts {
  display: grid;
  gap: 1rem;
}

.forum-post {
  padding: 1.35rem;
}

.post-header,
.post-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 0.9rem;
}

.author-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.author-name {
  font-weight: 700;
  color: #1f2937;
}

.post-date {
  color: #6b7280;
  font-size: 0.9rem;
  margin-top: 0.2rem;
}

.post-category {
  background: #fef3c7;
  color: #92400e;
  border-radius: 999px;
  padding: 0.45rem 0.9rem;
  font-size: 0.88rem;
}

.post-body h3 {
  margin: 1rem 0 0.65rem;
  color: #111827;
}

.post-body p {
  margin: 0;
  color: #4b5563;
  line-height: 1.75;
}

.post-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  color: #6b7280;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-top: 1.5rem;
}

.pagination-btn {
  background: #fff;
  border: 1px solid #d6dbe3;
  color: #1f2937;
}

.pagination-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.page-indicator {
  color: #4b5563;
  font-weight: 600;
}

.card-section {
  margin-top: 1rem;
}

.activities-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
}

.activity-card {
  overflow: hidden;
  border-radius: 18px;
  border: 1px solid #e5e7eb;
  background: #fff;
}

.activity-image {
  width: 100%;
  height: 220px;
  object-fit: cover;
}

.activity-info {
  padding: 1.2rem;
}

.activity-date {
  color: #9b5a30;
  font-weight: 700;
  margin-bottom: 0.4rem;
}

.activity-location,
.activity-description {
  color: #4b5563;
}

.activity-footer,
.feedback-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.two-column {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(280px, 1fr);
  gap: 1rem;
}

.panel {
  padding: 1.4rem;
}

.panel.muted {
  background: #f8fafc;
}

.stack-form {
  display: grid;
  gap: 0.9rem;
  margin-top: 1rem;
}

.feedback-item + .feedback-item {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
}

.feedback-meta {
  color: #6b7280;
  font-size: 0.9rem;
}

.feedback-question {
  margin: 0.45rem 0;
  font-weight: 600;
  color: #1f2937;
}

.feedback-reply {
  margin: 0;
  color: #4b5563;
  line-height: 1.7;
}

@media (max-width: 900px) {
  .community-nav,
  .two-column {
    grid-template-columns: 1fr;
  }

  .forum-toolbar,
  .toolbar-actions,
  .post-header,
  .post-footer,
  .activity-footer {
    flex-direction: column;
    align-items: flex-start;
  }

  .toolbar-actions {
    width: 100%;
  }

  .toolbar-actions select,
  .create-post-btn {
    width: 100%;
  }
}
</style>
