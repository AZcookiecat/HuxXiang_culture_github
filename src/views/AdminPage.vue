<template>
  <div class="admin-page">
    <div class="admin-container">
      <!-- 管理员页面头部 -->
      <div class="admin-header">
        <h1>管理中心</h1>
        <p>欢迎回来，{{ user?.username }}！您可以在这里管理平台的各项内容。</p>
      </div>

      <!-- 管理员统计卡片 -->
      <div class="admin-stats">
        <div class="stat-card">
          <div class="stat-icon"><i class="fas fa-landmark"></i></div>
          <div class="stat-content">
            <div class="stat-number">{{ resourceCount }}</div>
            <div class="stat-label">文化资源</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon"><i class="fas fa-project-diagram"></i></div>
          <div class="stat-content">
            <div class="stat-number">{{ graphNodeCount }}</div>
            <div class="stat-label">知识图谱节点</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon"><i class="fas fa-comments"></i></div>
          <div class="stat-content">
            <div class="stat-number">{{ postCount }}</div>
            <div class="stat-label">社区帖子</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon"><i class="fas fa-users"></i></div>
          <div class="stat-content">
            <div class="stat-number">{{ userCount }}</div>
            <div class="stat-label">注册用户</div>
          </div>
        </div>
      </div>

      <!-- 管理员选项卡 -->
      <div class="admin-tabs">
        <button 
          class="tab-button" 
          :class="{ active: activeTab === 'resources' }"
          @click="activeTab = 'resources'"
        >
          <i class="fas fa-landmark"></i> 文化资源管理
        </button>
        <button 
          class="tab-button" 
          :class="{ active: activeTab === 'knowledge-graph' }"
          @click="activeTab = 'knowledge-graph'"
        >
          <i class="fas fa-project-diagram"></i> 知识图谱管理
        </button>
        <button 
          class="tab-button" 
          :class="{ active: activeTab === 'community' }"
          @click="activeTab = 'community'"
        >
          <i class="fas fa-comments"></i> 互动社区管理
        </button>
        <button 
          class="tab-button" 
          :class="{ active: activeTab === 'users' }"
          @click="activeTab = 'users'"
        >
          <i class="fas fa-users"></i> 用户管理
        </button>
      </div>

      <!-- 管理内容区域 -->
      <div class="admin-content">
        <!-- 文化资源管理 -->
        <div v-if="activeTab === 'resources'" class="tab-content">
          <div class="content-header">
            <h2>文化资源管理</h2>
            <button class="btn btn-primary" @click="showAddResourceModal = true">
              <i class="fas fa-plus"></i> 添加资源
            </button>
          </div>
          <div class="resource-table-container">
            <table class="admin-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>标题</th>
                  <th>类型</th>
                  <th>创建时间</th>
                  <th>状态</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="resource in resources" :key="resource.id">
                  <td>{{ resource.id }}</td>
                  <td>{{ resource.title }}</td>
                  <td>{{ resource.type }}</td>
                  <td>{{ formatDate(resource.createdAt) }}</td>
                  <td>
                    <span :class="['status-badge', resource.status]">
                      {{ resource.status === 'published' ? '已发布' : '草稿' }}
                    </span>
                  </td>
                  <td>
                    <button class="btn btn-sm btn-info" @click="editResource(resource)">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" @click="deleteResource(resource.id)">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 知识图谱管理 -->
        <div v-if="activeTab === 'knowledge-graph'" class="tab-content">
          <div class="content-header">
            <h2>知识图谱管理</h2>
            <button class="btn btn-primary" @click="showAddNodeModal = true">
              <i class="fas fa-plus"></i> 添加节点
            </button>
          </div>
          <div class="graph-table-container">
            <table class="admin-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>节点名称</th>
                  <th>类型</th>
                  <th>相关资源</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="node in graphNodes" :key="node.id">
                  <td>{{ node.id }}</td>
                  <td>{{ node.name }}</td>
                  <td>{{ node.type }}</td>
                  <td>{{ node.relatedResources || 0 }}</td>
                  <td>
                    <button class="btn btn-sm btn-info" @click="editNode(node)">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" @click="deleteNode(node.id)">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 互动社区管理 -->
        <div v-if="activeTab === 'community'" class="tab-content">
          <div class="content-header">
            <h2>互动社区管理</h2>
            <div class="filter-buttons">
              <button class="btn btn-sm btn-outline" :class="{ active: postFilter === 'all' }" @click="postFilter = 'all'">
                全部
              </button>
              <button class="btn btn-sm btn-outline" :class="{ active: postFilter === 'pending' }" @click="postFilter = 'pending'">
                待审核
              </button>
              <button class="btn btn-sm btn-outline" :class="{ active: postFilter === 'approved' }" @click="postFilter = 'approved'">
                已审核
              </button>
            </div>
          </div>
          <div class="post-table-container">
            <table class="admin-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>标题</th>
                  <th>作者</th>
                  <th>发布时间</th>
                  <th>状态</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="post in filteredPosts" :key="post.id">
                  <td>{{ post.id }}</td>
                  <td>{{ post.title }}</td>
                  <td>{{ post.author }}</td>
                  <td>{{ formatDate(post.createdAt) }}</td>
                  <td>
                    <span :class="['status-badge', post.status]">
                      {{ post.status === 'approved' ? '已审核' : '待审核' }}
                    </span>
                  </td>
                  <td>
                    <button class="btn btn-sm btn-info" @click="viewPost(post.id)">
                      <i class="fas fa-eye"></i>
                    </button>
                    <button class="btn btn-sm btn-warning" v-if="post.status !== 'approved'" @click="approvePost(post.id)">
                      <i class="fas fa-check"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" @click="deletePost(post.id)">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 用户管理 -->
        <div v-if="activeTab === 'users'" class="tab-content">
          <div class="content-header">
            <h2>用户管理</h2>
            <button class="btn btn-primary" @click="showAddUserModal = true">
              <i class="fas fa-plus"></i> 添加用户
            </button>
          </div>
          <div class="user-table-container">
            <table class="admin-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>用户名</th>
                  <th>邮箱</th>
                  <th>注册时间</th>
                  <th>角色</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="user in users" :key="user.id" :class="{ 'admin-row': user.isAdmin }">
                  <td>{{ user.id }}</td>
                  <td>{{ user.username }}</td>
                  <td>{{ user.email }}</td>
                  <td>{{ formatDate(user.registeredAt) }}</td>
                  <td>
                    <span :class="['status-badge', user.isAdmin ? 'admin' : 'user']">
                      {{ user.isAdmin ? '管理员' : '普通用户' }}
                    </span>
                  </td>
                  <td>
                    <button class="btn btn-sm btn-info" @click="editUser(user)">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button 
                      class="btn btn-sm btn-danger" 
                      @click="deleteUser(user.id)"
                      :disabled="user.isAdmin && user.username === currentUser.username"
                    >
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

// 模拟数据服务
const mockAdminService = {
  // 获取统计数据
  getStats: async () => {
    return {
      resourceCount: 42,
      graphNodeCount: 128,
      postCount: 256,
      userCount: 1024
    }
  },
  
  // 获取资源列表
  getResources: async () => {
    return [
      { id: 1, title: '岳麓书院历史', type: '历史建筑', createdAt: '2024-01-15T08:30:00', status: 'published' },
      { id: 2, title: '岳阳楼记赏析', type: '文学作品', createdAt: '2024-01-16T14:20:00', status: 'published' },
      { id: 3, title: '湘菜文化研究', type: '饮食文化', createdAt: '2024-01-17T10:15:00', status: 'draft' }
    ]
  },
  
  // 获取知识图谱节点
  getGraphNodes: async () => {
    return [
      { id: 1, name: '毛泽东', type: '历史人物', relatedResources: 15 },
      { id: 2, name: '曾国藩', type: '历史人物', relatedResources: 8 },
      { id: 3, name: '岳麓书院', type: '文化遗产', relatedResources: 12 },
      { id: 4, name: '湘菜', type: '饮食文化', relatedResources: 6 }
    ]
  },
  
  // 获取社区帖子
  getPosts: async () => {
    return [
      { id: 1, title: '讨论湖湘文化的现代意义', author: 'user1', createdAt: '2024-01-18T09:45:00', status: 'approved' },
      { id: 2, title: '请教关于曾国藩家训的解读', author: 'user2', createdAt: '2024-01-18T11:30:00', status: 'pending' },
      { id: 3, title: '分享一次湖南文化之旅', author: 'user3', createdAt: '2024-01-18T15:20:00', status: 'approved' }
    ]
  },
  
  // 获取用户列表
  getUsers: async () => {
    return [
      { id: 1, username: 'admin', email: '18821959907@qq.com', registeredAt: '2024-01-01T10:00:00', isAdmin: true },
      { id: 2, username: 'user1', email: 'user1@example.com', registeredAt: '2024-01-05T14:30:00', isAdmin: false },
      { id: 3, username: 'user2', email: 'user2@example.com', registeredAt: '2024-01-10T09:15:00', isAdmin: false }
    ]
  }
}

export default {
  name: 'AdminPage',
  setup() {
    const router = useRouter()
    
    // 状态管理
    const activeTab = ref('resources')
    const postFilter = ref('all')
    const resourceCount = ref(0)
    const graphNodeCount = ref(0)
    const postCount = ref(0)
    const userCount = ref(0)
    const resources = ref([])
    const graphNodes = ref([])
    const posts = ref([])
    const users = ref([])
    const user = ref(null)
    const currentUser = ref(null)
    
    // 模态框状态
    const showAddResourceModal = ref(false)
    const showAddNodeModal = ref(false)
    const showAddUserModal = ref(false)
    
    // 计算过滤后的帖子
    const filteredPosts = computed(() => {
      if (postFilter.value === 'all') return posts.value
      return posts.value.filter(post => post.status === postFilter.value)
    })
    
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
    
    // 加载数据
    const loadData = async () => {
      try {
        // 获取统计数据
        const stats = await mockAdminService.getStats()
        resourceCount.value = stats.resourceCount
        graphNodeCount.value = stats.graphNodeCount
        postCount.value = stats.postCount
        userCount.value = stats.userCount
        
        // 获取资源列表
        resources.value = await mockAdminService.getResources()
        
        // 获取知识图谱节点
        graphNodes.value = await mockAdminService.getGraphNodes()
        
        // 获取社区帖子
        posts.value = await mockAdminService.getPosts()
        
        // 获取用户列表
        users.value = await mockAdminService.getUsers()
        
        // 获取当前用户信息
        const storedUser = localStorage.getItem('user')
        if (storedUser) {
          currentUser.value = JSON.parse(storedUser)
          user.value = currentUser.value
        }
      } catch (error) {
        console.error('加载管理数据失败:', error)
      }
    }
    
    // 资源操作
    const editResource = (resource) => {
      alert(`编辑资源: ${resource.title}`)
      // 实际项目中应该打开编辑模态框
    }
    
    const deleteResource = (id) => {
      if (confirm('确定要删除这个资源吗？')) {
        resources.value = resources.value.filter(r => r.id !== id)
        resourceCount.value--
        alert('资源已删除')
      }
    }
    
    // 知识图谱节点操作
    const editNode = (node) => {
      alert(`编辑节点: ${node.name}`)
      // 实际项目中应该打开编辑模态框
    }
    
    const deleteNode = (id) => {
      if (confirm('确定要删除这个节点吗？')) {
        graphNodes.value = graphNodes.value.filter(n => n.id !== id)
        graphNodeCount.value--
        alert('节点已删除')
      }
    }
    
    // 帖子操作
    const viewPost = (id) => {
      alert(`查看帖子ID: ${id}`)
      // 实际项目中应该跳转到帖子详情页或打开详情模态框
    }
    
    const approvePost = (id) => {
      const post = posts.value.find(p => p.id === id)
      if (post) {
        post.status = 'approved'
        alert('帖子已审核通过')
      }
    }
    
    const deletePost = (id) => {
      if (confirm('确定要删除这个帖子吗？')) {
        posts.value = posts.value.filter(p => p.id !== id)
        postCount.value--
        alert('帖子已删除')
      }
    }
    
    // 用户操作
    const editUser = (user) => {
      alert(`编辑用户: ${user.username}`)
      // 实际项目中应该打开编辑模态框
    }
    
    const deleteUser = (id) => {
      const userToDelete = users.value.find(u => u.id === id)
      if (userToDelete && userToDelete.isAdmin && userToDelete.username === currentUser.value?.username) {
        alert('无法删除当前登录的管理员账号')
        return
      }
      
      if (confirm(`确定要删除用户 ${userToDelete?.username || ''} 吗？`)) {
        users.value = users.value.filter(u => u.id !== id)
        userCount.value--
        alert('用户已删除')
      }
    }
    
    // 页面加载时获取数据
    onMounted(() => {
      loadData()
    })
    
    return {
      activeTab,
      postFilter,
      resourceCount,
      graphNodeCount,
      postCount,
      userCount,
      resources,
      graphNodes,
      posts,
      users,
      user,
      currentUser,
      showAddResourceModal,
      showAddNodeModal,
      showAddUserModal,
      filteredPosts,
      formatDate,
      editResource,
      deleteResource,
      editNode,
      deleteNode,
      viewPost,
      approvePost,
      deletePost,
      editUser,
      deleteUser
    }
  }
}
</script>

<style scoped>
:root {
  --primary-color: #C8102E; /* 湘红 */
  --secondary-color: #1E40AF; /* 湘蓝 */
  --accent-color: #D97706; /* 湘金 */
  --bg-color: #F9FAFB;
  --text-color: #1F2937;
  --light-text: #6B7280;
  --border-color: #E5E7EB;
  --success-color: #10B981;
  --warning-color: #F59E0B;
  --danger-color: #EF4444;
  --info-color: #3B82F6;
}

.admin-page {
  padding: 2rem 0;
  min-height: calc(100vh - 100px);
  background-color: var(--bg-color);
}

.admin-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

.admin-header {
  margin-bottom: 2rem;
  text-align: center;
}

.admin-header h1 {
  font-size: 2.5rem;
  color: black;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.admin-header p {
  color: var(--light-text);
  font-size: 1.1rem;
}

/* 统计卡片样式 */
.admin-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background-color: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.stat-icon {
  font-size: 2rem;
  color: var(--primary-color);
  background-color: rgba(200, 16, 46, 0.1);
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
  color: var(--text-color);
  line-height: 1;
}

.stat-label {
  color: var(--light-text);
  font-size: 0.9rem;
  margin-top: 0.25rem;
}

/* 选项卡样式 */
.admin-tabs {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.tab-button {
  padding: 0.75rem 1.5rem;
  border: 1px solid var(--border-color);
  background-color: white;
  color: var(--text-color);
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tab-button:hover {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.tab-button.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

/* 内容区域样式 */
.admin-content {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 2rem;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.content-header h2 {
  color: var(--text-color);
  font-size: 1.5rem;
  margin: 0;
}

.filter-buttons {
  display: flex;
  gap: 0.5rem;
}

/* 表格样式 */
.admin-table {
  width: 100%;
  border-collapse: collapse;
}

.admin-table th {
  background-color: var(--bg-color);
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: var(--text-color);
  border-bottom: 2px solid var(--border-color);
}

.admin-table td {
  padding: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.admin-table tr:hover {
  background-color: var(--bg-color);
}

.admin-table tr.admin-row {
  background-color: rgba(200, 16, 46, 0.05);
}

/* 状态徽章 */
.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.8rem;
  font-weight: 500;
}

.status-badge.published,
.status-badge.approved {
  background-color: rgba(16, 185, 129, 0.1);
  color: var(--success-color);
}

.status-badge.draft,
.status-badge.pending {
  background-color: rgba(245, 158, 11, 0.1);
  color: var(--warning-color);
}

.status-badge.admin {
  background-color: rgba(200, 16, 46, 0.1);
  color: var(--primary-color);
}

.status-badge.user {
  background-color: rgba(59, 130, 246, 0.1);
  color: var(--info-color);
}

/* 按钮样式 */
.btn {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  text-decoration: none;
  font-size: 0.9rem;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background-color: #A80D27;
}

.btn-outline {
  background-color: transparent;
  border: 1px solid var(--border-color);
  color: var(--text-color);
}

.btn-outline:hover {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.btn-outline.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.btn-sm {
  padding: 0.25rem 0.75rem;
  font-size: 0.8rem;
}

.btn-info {
  background-color: var(--info-color);
  color: white;
}

.btn-info:hover {
  background-color: #2563EB;
}

.btn-danger {
  background-color: var(--danger-color);
  color: white;
}

.btn-danger:hover {
  background-color: #DC2626;
}

.btn-danger:disabled {
  background-color: var(--light-text);
  cursor: not-allowed;
}

.btn-warning {
  background-color: var(--warning-color);
  color: white;
}

.btn-warning:hover {
  background-color: #D97706;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .admin-stats {
    grid-template-columns: 1fr;
  }
  
  .admin-tabs {
    justify-content: center;
  }
  
  .content-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .admin-table {
    display: block;
    overflow-x: auto;
  }
}
</style>