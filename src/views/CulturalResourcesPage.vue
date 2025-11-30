<template>
  <div class="cultural-resources-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="container">
        <h1>文化资源库</h1>
        <p>探索湖湘大地丰富的文化遗产和历史底蕴</p>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="container">
      <!-- 搜索和筛选 -->
      <div class="search-filter-section">
        <div class="search-box">
          <input type="text" v-model="searchQuery" placeholder="搜索文化资源..." @input="searchResources" />
          <button class="search-btn"><i class="fas fa-search"></i></button>
        </div>
        
        <div class="resource-categories">
          <button 
            v-for="category in categories" 
            :key="category.id"
            :class="['category-btn', { active: currentCategory === category.id }]"
            @click="filterByCategory(category.id)"
          >
            {{ category.name }}
          </button>
        </div>
      </div>

      <!-- 资源网格 -->
      <div class="resources-grid">
        <div 
          v-for="resource in filteredResources" 
          :key="resource.id"
          class="resource-card"
          @click="viewResourceDetails(resource.id)"
        >
          <div class="resource-image">
            <img :src="resource.imageUrl" :alt="resource.title" loading="lazy" />
          </div>
          <div class="resource-info">
            <h3>{{ resource.title }}</h3>
            <p class="resource-description">{{ resource.description }}</p>
            <div class="resource-meta">
              <span class="resource-category">{{ getCategoryName(resource.categoryId) }}</span>
              <button class="btn view-detail-btn" @click.stop="viewResourceDetails(resource.id)">查看详情</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination" v-if="totalPages > 1">
        <button 
          class="pagination-btn"
          :disabled="currentPage === 1"
          @click="goToPage(currentPage - 1)"
        >
          上一页
        </button>
        
        <button 
          v-for="page in totalPages" 
          :key="page"
          :class="['pagination-btn', { active: currentPage === page }]"
          @click="goToPage(page)"
        >
          {{ page }}
        </button>
        
        <button 
          class="pagination-btn"
          :disabled="currentPage === totalPages"
          @click="goToPage(currentPage + 1)"
        >
          下一页
        </button>
      </div>

      <!-- 空数据提示 -->
      <div class="empty-state" v-if="filteredResources.length === 0">
        <i class="fas fa-search"></i>
        <p>未找到匹配的文化资源，请尝试其他关键词或分类。</p>
        <button class="reset-btn" @click="resetFilters">重置筛选条件</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'CulturalResourcesPage',
  props: {
    showAlert: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    const router = useRouter()
    
    // 状态管理
    const searchQuery = ref('')
    const currentCategory = ref('all')
    const currentPage = ref(1)
    const itemsPerPage = ref(12)
    
    // 分类数据
    const categories = ref([
      { id: 'all', name: '全部资源' },
      { id: 'history', name: '历史遗迹' },
      { id: 'art', name: '传统艺术' },
      { id: 'literature', name: '文学作品' },
      { id: 'folk', name: '民俗风情' },
      { id: 'food', name: '饮食文化' },
      { id: 'architecture', name: '建筑风格' },
      { id: 'religion', name: '宗教文化' },
      { id: 'greatmen', name: '历史伟人' }
    ])
    
    // 模拟资源数据
    const resources = ref([
      {
        id: '1',
        title: '岳麓书院',
        description: '中国历史上著名的四大书院之一，始建于北宋开宝九年（976年）。',
        categoryId: 'history',
        imageUrl: 'https://img2.baidu.com/it/u=1868825345,3526076103&fm=253&app=138&f=JPEG?w=500&h=582',
        createdAt: '2025-3-24'
      },
      {
        id: '2',
        title: '湘绣',
        description: '中国四大名绣之一，以其精湛的技艺和独特的风格闻名于世。',
        categoryId: 'art',
        imageUrl: 'https://pic.rmb.bdstatic.com/bjh/news/89c508f7a400f1751758f343caae8801.png',
        createdAt: '2025-3-24'
      },
      {
        id: '3',
        title: '《岳阳楼记》',
        description: '江南三大名楼之一，因范仲淹的《岳阳楼记》而名扬天下。',
        categoryId: 'architecture',
        imageUrl: 'https://img0.baidu.com/it/u=1464457526,3199376473&fm=253&app=138&f=JPEG?w=800&h=2997',
        createdAt: '2025-3-24'
      },
      {
        id: '4',
        title: '《桃花源记》',
        description: '陶渊明的代表作，描绘了一个与世隔绝的理想社会。',
        categoryId: 'literature',
        imageUrl: 'https://shufa.pku.edu.cn/images/2024-09/90f030ce-1d59-4f14-bfce-00f3534cb9f0.jpeg',
        createdAt: '2025-3-24'
      },
      {
        id: '5',
        title: '湘菜',
        description: '中国八大菜系之一，以辣、香、鲜、嫩为特色。',
        categoryId: 'food',
        imageUrl: 'https://qcloud.dpfile.com/pc/y41poElenLBv7nqvkBnoBORgilWDRCfeVQ0djkNKJwMID0bsu-SAtCe9N8iBmkY3.jpg',
        createdAt: '2025-3-24'
      },
      {
        id: '6',
        title: '汨罗端午节',
        description: '纪念屈原的传统节日，有赛龙舟、吃粽子等习俗。',
        categoryId: 'folk',
        imageUrl: 'https://pic.rmb.bdstatic.com/bjh/bc1190b5ccbb/250505/b762f2bcc5ec3a93bf3b092a3eef6351.jpeg',
        createdAt: '2025-3-24'
      },
      {
        id: '7',
        title: '衡山',
        description: '中国五岳之一，被誉为"南岳"，是著名的道教和佛教圣地。',
        categoryId: 'history',
        imageUrl: 'https://img0.baidu.com/it/u=4250670926,3034496007&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=667',
        createdAt: '2025-3-24'
      },
      {
        id: '8',
        title: '花鼓戏',
        description: '湖南省特有的地方戏曲，具有浓郁的乡土气息。',
        categoryId: 'art',
        imageUrl: 'https://p7.itc.cn/q_70/images03/20220411/4af6d947f76f48289ab31aeba6936b91.jpeg',
        createdAt: '2025-3-24'
      },
      {
        id: '9',
        title: '曾国藩故居',
        description: '清代名臣曾国藩的出生地，展示了湘军文化的历史。',
        categoryId: 'history',
        imageUrl: 'https://img0.baidu.com/it/u=21139302,2048799699&fm=253&app=138&f=JPEG?w=500&h=652',
        createdAt: '2025-3-24'
      },
      {
        id: '10',
        title: '马王堆汉墓',
        description: '出土了大量珍贵文物的汉代墓葬，为研究汉代历史提供了重要资料。',
        categoryId: 'history',
        imageUrl: 'https://gips0.baidu.com/it/u=1600549162,356730080&fm=3074&app=3074&f=JPEG',
        createdAt: '2025-3-24'
      },
      {
        id: '11',
        title: '火宫殿',
        description: '长沙著名的传统小吃聚集地，有"中华老字号"之称。',
        categoryId: 'food',
        imageUrl: 'https://img2.baidu.com/it/u=3401007738,1338386073&fm=253&app=138&f=JPEG?w=800&h=1199',
        createdAt: '2025-3-24'
      },
      {
        id: '12',
        title: '铜官窑',
        description: '唐代著名的青瓷窑址，是海上丝绸之路的重要起点之一。',
        categoryId: 'art',
        imageUrl: 'https://img2.baidu.com/it/u=3735106663,128794723&fm=253&app=138&f=JPEG?w=800&h=1067',
        createdAt: '2025-3-24'
      },
      {
        id: '13',
        title: '曾国藩',
        description: '晚清时期政治家、战略家、理学家、文学家，湘军的创立者和统帅。',
        categoryId: 'greatmen',
        imageUrl: 'https://pic.rmb.bdstatic.com/bjh/3f19cdcfc55/241008/0b198489d085fc5234e8d11e0f8edc2c.png',
        createdAt: '2025-3-24'
      }
    ])
    
    // 计算属性：筛选后的资源
    const filteredResources = computed(() => {
      let filtered = [...resources.value]
      
      // 按分类筛选
      if (currentCategory.value !== 'all') {
        filtered = filtered.filter(resource => resource.categoryId === currentCategory.value)
      }
      
      // 按关键词搜索
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(resource => 
          resource.title.toLowerCase().includes(query) || 
          resource.description.toLowerCase().includes(query) ||
          getCategoryName(resource.categoryId).toLowerCase().includes(query)
        )
      }
      
      // 分页处理
      const startIndex = (currentPage.value - 1) * itemsPerPage.value
      const endIndex = startIndex + itemsPerPage.value
      
      return filtered.slice(startIndex, endIndex)
    })
    
    // 计算属性：总页数
    const totalPages = computed(() => {
      let filtered = [...resources.value]
      
      // 按分类筛选
      if (currentCategory.value !== 'all') {
        filtered = filtered.filter(resource => resource.categoryId === currentCategory.value)
      }
      
      // 按关键词搜索
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(resource => 
          resource.title.toLowerCase().includes(query) || 
          resource.description.toLowerCase().includes(query) ||
          getCategoryName(resource.categoryId).toLowerCase().includes(query)
        )
      }
      
      return Math.ceil(filtered.length / itemsPerPage.value)
    })
    
    // 方法：获取分类名称
    const getCategoryName = (categoryId) => {
      const category = categories.value.find(cat => cat.id === categoryId)
      return category ? category.name : '未知分类'
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
    
    // 方法：按分类筛选
    const filterByCategory = (categoryId) => {
      currentCategory.value = categoryId
      currentPage.value = 1 // 重置到第一页
    }
    
    // 方法：搜索资源
    const searchResources = () => {
      currentPage.value = 1 // 重置到第一页
    }
    
    // 方法：分页导航
    const goToPage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
      }
    }
    
    // 方法：重置筛选条件
    const resetFilters = () => {
      searchQuery.value = ''
      currentCategory.value = 'all'
      currentPage.value = 1
    }
    
    // 方法：查看资源详情
    const viewResourceDetails = (resourceId) => {
      router.push({ name: 'resource-detail', params: { id: resourceId } })
    }
    
    return {
      searchQuery,
      currentCategory,
      currentPage,
      categories,
      resources,
      filteredResources,
      totalPages,
      getCategoryName,
      formatDate,
      filterByCategory,
      searchResources,
      goToPage,
      resetFilters,
      viewResourceDetails
    }
  }
}
</script>

<style scoped>
.cultural-resources-page {
  padding: 20px 0;
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
  margin-bottom: 0.5rem;
}

.page-header p {
  color: black;
  font-size: 1.2rem;
  max-width: 800px;
  margin: 0 auto;
}

.search-filter-section {
  margin: 1rem 0;
}

.search-box {
  display: flex;
  margin: 0 auto 1rem auto;
  max-width: 600px;
}

.search-box input {
  flex: 1;
  padding: 0.75rem;
  font-size: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px 0 0 4px;
}

.search-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  cursor: pointer;
  border-radius: 0 4px 4px 0;
  transition: background-color 0.3s;
}

.search-btn:hover {
  background-color: var(--primary-dark);
}

.resource-categories {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  justify-content: center;
}

.category-btn {
  background-color: white;
  border: 1px solid #ddd;
  padding: 0.5rem 1rem;
  cursor: pointer;
  border-radius: 20px;
  transition: all 0.3s;
  font-size: 0.9rem;
}

.category-btn:hover {
  background-color: #f8f9fa;
}

.category-btn.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.resources-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  margin: 1.5rem 0;
}

.resource-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.resource-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.resource-image {
  height: 200px;
  overflow: hidden;
}

.resource-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.resource-card:hover .resource-image img {
  transform: scale(1.05);
}

.resource-info {
  padding: 1rem;
}

.resource-info h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.2rem;
  color: #333;
}

.resource-description {
  margin: 0 0 1rem 0;
  color: #666;
  font-size: 0.9rem;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.resource-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: #999;
  margin-top: 0.5rem;
}

.btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: var(--primary-dark);
}

.view-detail-btn {
  padding: 0.4rem 0.8rem;
  font-size: 0.8rem;
}

.resource-category {
  background-color: #f0f0f0;
  padding: 0.25rem 0.5rem;
  border-radius: 10px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.5rem;
  margin: 1.5rem 0;
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

.empty-state {
  text-align: center;
  padding: 4rem 0;
  color: #666;
}

.empty-state i {
  font-size: 3rem;
  margin-bottom: 1rem;
  color: #ccc;
}

.reset-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 4px;
  margin-top: 1rem;
  transition: background-color 0.3s;
}

.reset-btn:hover {
  background-color: var(--primary-dark);
}

@media (max-width: 768px) {
  .page-header h1 {
    font-size: 2rem;
  }
  
  .resources-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  }
  
  .resource-categories {
    justify-content: center;
  }
  
  .pagination {
    flex-wrap: wrap;
  }
}
</style>