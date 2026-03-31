      <template>
  <div class="digital-showcase-page">
    <!-- 主标题区域 -->
    <div class="section-title">
      <h2>湖湘文化数字化展示</h2>
      <p>通过现代信息技术，以数字化方式展示湖湘文化的各个方面，让您获得更生动的体验。</p>
    </div>

    <!-- 主要展示区域 -->
    <div class="showcase-container">
      <!-- 已有的核心展示项目 -->
      <h3 class="category-title">核心数字化体验</h3>
      <div class="showcase-grid">
        <div class="showcase-item">
          <router-link to="/knowledge-graph" class="showcase-link">
            <div class="showcase-icon">
              <i class="fas fa-project-diagram fa-3x"></i>
            </div>
            <h3>文化知识图谱</h3>
            <p>通过可视化图谱展示湖湘文化各元素之间的关联关系</p>
            <button class="start-btn">开始查看</button>
          </router-link>
        </div>
        
        <div class="showcase-item">
          <router-link to="/unity-webgl" class="showcase-link">
            <div class="showcase-icon">
              <i class="fas fa-vr-cardboard fa-3x"></i>
            </div>
            <h3>虚拟现实体验</h3>
            <p>沉浸式VR体验湖湘文化名胜古迹</p>
            <button class="start-btn">开始体验</button>
          </router-link>
        </div>
        
        <div class="showcase-item">
          <a href="#" @click.prevent="showAlert('360度全景展示功能即将上线！')">
            <div class="showcase-icon">
              <i class="fas fa-street-view fa-3x"></i>
            </div>
            <h3>360度全景展示</h3>
            <p>全方位欣赏湖湘文化景观</p>
          </a>
        </div>
        
        <div class="showcase-item">
          <router-link to="/ai-assistant" class="ai-assistant-link">
            <div class="showcase-icon">
              <i class="fas fa-robot fa-3x"></i>
            </div>
            <h3>AI文化助手</h3>
            <p>智能问答系统，解答您的湖湘文化问题</p>
            <button class="start-ai-btn">开始对话</button>
          </router-link>
        </div>
      </div>

      <!-- 新增博物馆展示项目 -->
      <h3 class="category-title">数字博物馆</h3>
      <div class="showcase-grid">
        <div class="showcase-item">
          <router-link to="/unity-webgl">
            <div class="showcase-icon">
              <i class="fas fa-museum fa-3x"></i>
            </div>
            <h3>数字博物馆展览</h3>
            <p>在线浏览湖湘文化珍贵文物和历史资料</p>
            <button class="start-btn">开始展览</button>
          </router-link>
        </div>
        
        <div class="showcase-item">
          <a href="#" @click.prevent="showAlert('文物AR展示功能即将上线！')">
            <div class="showcase-icon">
              <i class="fas fa-cube fa-3x"></i>
            </div>
            <h3>文物AR展示</h3>
            <p>通过增强现实技术查看文物的3D模型</p>
          </a>
        </div>
        
        <div class="showcase-item">
          <a href="#" @click.prevent="showAlert('非物质文化遗产数字保护项目即将上线！')">
            <div class="showcase-icon">
              <i class="fas fa-hands-helping fa-3x"></i>
            </div>
            <h3>非遗数字保护</h3>
            <p>数字化记录和保护湖湘地区非物质文化遗产</p>
          </a>
        </div>
      </div>

      <!-- 新增互动体验项目 -->
      <h3 class="category-title">互动文化体验</h3>
      <div class="showcase-grid">
        <div class="showcase-item">
          <a href="#" @click.prevent="showAlert('湖湘文化互动地图功能即将上线！')">
            <div class="showcase-icon">
              <i class="fas fa-map-marked-alt fa-3x"></i>
            </div>
            <h3>文化互动地图</h3>
            <p>探索湖湘地区文化景点和历史遗迹分布</p>
          </a>
        </div>
        
        <div class="showcase-item">
          <a href="#" @click.prevent="showAlert('湖湘文化音频可视化功能即将上线！')">
            <div class="showcase-icon">
              <i class="fas fa-music fa-3x"></i>
            </div>
            <h3>音频可视化</h3>
            <p>聆听湖湘传统音乐，体验音频可视化效果</p>
          </a>
        </div>
        
        <div class="showcase-item">
          <a href="#" @click.prevent="showAlert('数字文创产品展示功能即将上线！')">
            <div class="showcase-icon">
              <i class="fas fa-paint-brush fa-3x"></i>
            </div>
            <h3>数字文创产品</h3>
            <p>浏览和了解基于湖湘文化元素的数字文创产品</p>
          </a>
        </div>
      </div>

      <!-- 特色专题区域 -->
      <div class="featured-topics">
        <h3 class="category-title">特色专题</h3>
        <div class="topics-container">
          <div class="topic-card">
            <img src="/src/assets/Imgs/岳阳楼记1.jpg" alt="湖湘诗词数字化" class="topic-image">
            <div class="topic-content">
              <h4>湖湘诗词数字化</h4>
              <p>数字化呈现湖湘地区经典诗词作品，配合注释和朗诵音频</p>
              <button class="topic-btn" @click="goToPage('poetry-digitalization')">了解更多</button>
            </div>
          </div>
          
          <div class="topic-card">
            <img src="/src/assets/Imgs/岳麓书院1.jpg" alt="湖湘建筑3D模型" class="topic-image">
            <div class="topic-content">
              <h4>湖湘建筑3D模型</h4>
              <p>探索湖湘地区特色建筑的精细3D模型，了解建筑结构和文化内涵</p>
              <button class="topic-btn" @click="goToPage('architecture-3d')">了解更多</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { useRouter } from 'vue-router'

// 获取全局事件总线和路由实例
const eventBus = inject('$eventBus')
const router = useRouter()

// 显示提示消息
const showAlert = (message, type = 'info') => {
  if (eventBus) {
    eventBus.emit('show-alert', { message, type })
  } else {
    // 如果没有eventBus，不显示任何弹窗，只在控制台打印消息
    console.log(`[${type}] ${message}`)
  }
}

// 跳转到指定页面
const goToPage = (pageName) => {
  router.push({ name: pageName })
}
</script>

<style scoped>
:root {
  --primary-color: #C8102E; /* 湘红色 */
  --bg-color: #F9FAFB;
  --text-color: #1F2937;
  --light-text: #6B7280;
  --heading-font: 'Microsoft YaHei', 'STXihei', 'SimHei', sans-serif;
  --body-font: 'Microsoft YaHei', 'STXihei', 'SimSun', serif;
}

.digital-showcase-page {
  padding: 2rem 0;
  background-color: #fafafa;
}

.section-title {
  text-align: center;
  margin-bottom: 3rem;
  padding: 0 1rem;
  margin-top: 20px;
}

.section-title h2 {
  font-size: 2.5rem;
  color: var(--text-color);
  margin-bottom: 0.5rem;
  font-family: var(--heading-font);
  font-weight: 700;
  line-height: 1.3;
  position: relative;
  display: inline-block;
}

.section-title h2::after {
  content: '';
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%);
  width: 100px;
  height: 3px;
  background-color: var(--primary-color);
}

.section-title p {
  color: var(--light-text);
  max-width: 600px;
  margin: 1.5rem auto 0;
  line-height: 1.7;
  font-size: 16px;
}

.showcase-container {
  padding: 0 2rem 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

/* 分类标题样式 */
.category-title {
  font-size: 1.8rem;
  color: var(--text-color);
  margin: 3rem 0 1.5rem;
  font-family: var(--heading-font);
  font-weight: 600;
  position: relative;
  padding-left: 1rem;
}

.category-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  background-color: var(--primary-color);
}

/* 网格展示样式 */
.showcase-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-top: 1rem;
}

.showcase-item {
  background: #ffffff;
  border-radius: 8px;
  padding: 2rem;
  text-align: center;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border: 1px solid #e9ecef;
  position: relative;
  overflow: hidden;
}

.showcase-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background-color: transparent;
  transition: background-color 0.3s ease;
}

.showcase-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.1);
}

.showcase-item:hover::before {
  background-color: var(--primary-color);
}

.showcase-item a,
.showcase-item router-link {
  text-decoration: none;
  color: inherit;
  display: block;
  height: 100%;
}

.showcase-icon {
  color: var(--primary-color);
  margin-bottom: 1rem;
  background-color: rgba(200, 16, 46, 0.05);
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
  transition: background-color 0.3s ease;
}

.showcase-item:hover .showcase-icon {
  background-color: rgba(200, 16, 46, 0.1);
}

.showcase-item h3 {
  margin: 1rem 0;
  color: #333;
  font-family: var(--heading-font);
  font-weight: 700;
  line-height: 1.3;
  font-size: 1.4rem;
  transition: color 0.3s ease;
}

.showcase-item:hover h3 {
  color: var(--primary-color);
}

.showcase-item p {
  color: #666;
  font-size: 0.95rem;
  line-height: 1.7;
}

/* 特色专题样式 */
.featured-topics {
  margin-top: 3rem;
}

.topics-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 1rem;
}

.topic-card {
  background: #ffffff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.topic-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.topic-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  display: block;
}

.topic-content {
  padding: 1.5rem;
}

.topic-content h4 {
  font-size: 1.3rem;
  color: #333;
  margin-bottom: 0.8rem;
  font-family: var(--heading-font);
  font-weight: 600;
}

.topic-content p {
  color: #666;
  font-size: 0.95rem;
  line-height: 1.7;
  margin-bottom: 1rem;
}

.topic-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.8rem 1.5rem;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.topic-btn:hover {
  background-color: #a60e24;
  transform: translateY(-1px);
}

.topic-btn:active {
  transform: translateY(0);
}

/* AI文化助手样式 */
.ai-assistant-placeholder {
  cursor: pointer;
  transition: all 0.3s ease;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.ai-assistant-placeholder:hover {
  transform: translateY(-5px);
}

.start-ai-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.8rem 1.8rem;
  border-radius: 25px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 1.5rem;
  box-shadow: 0 4px 12px rgba(200, 16, 46, 0.3);
}

.start-ai-btn:hover {
  background-color: #a60e24;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

.start-ai-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(200, 16, 46, 0.3);
}

.ai-assistant-item {
  grid-column: 1 / -1;
  min-height: 600px;
  padding: 0;
  overflow: hidden;
}

.close-btn {
  background: none;
  border: none;
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 0.2rem;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: rotate(90deg);
}

.menu-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e9ecef;
  background: var(--primary-color);
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.menu-header h4 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 600;
}

.ai-assistant-container {
  display: flex;
  height: 100%;
  width: 100%;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

/* 左侧菜单样式 */
.ai-menu {
  width: 250px;
  background: #f8f9fa;
  border-right: 1px solid #e9ecef;
  display: flex;
  flex-direction: column;
}

.menu-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e9ecef;
  background: var(--primary-color);
  color: white;
}

.menu-header h4 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 600;
}

.menu-items {
  list-style: none;
  padding: 0;
  margin: 0;
  flex: 1;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 1rem 1.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
}

.menu-item:hover {
  background: rgba(200, 16, 46, 0.05);
}

.menu-item.active {
  background: rgba(200, 16, 46, 0.1);
  border-left-color: var(--primary-color);
  color: var(--primary-color);
}

.menu-item i {
  margin-right: 0.8rem;
  font-size: 1.1rem;
}

/* 展示项目链接样式 */
.showcase-link {
  text-decoration: none;
  color: inherit;
  display: block;
  height: 100%;
}

.ai-assistant-link {
  text-decoration: none;
  color: inherit;
  display: block;
  height: 100%;
}

/* 开始按钮样式 */
.start-btn,
.start-ai-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.8rem 1.8rem;
  border-radius: 25px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 1.5rem;
  box-shadow: 0 4px 12px rgba(200, 16, 46, 0.3);
}

.start-btn:hover,
.start-ai-btn:hover {
  background-color: #a60e24;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

.start-btn:active,
.start-ai-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(200, 16, 46, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .section-title h2 {
    font-size: 2rem;
  }
  
  .showcase-container {
    padding: 0 1rem 2rem;
  }
  
  .showcase-grid,
  .topics-container {
    gap: 1.2rem;
  }
  
  .showcase-item {
    padding: 1.5rem;
  }
  
  .category-title {
    font-size: 1.6rem;
  }
  
  .topic-card {
    margin-bottom: 1rem;
  }
}

@media (max-width: 480px) {
  .section-title h2 {
    font-size: 1.8rem;
  }
  
  .showcase-item h3 {
    font-size: 1.25rem;
  }
  
  .category-title {
    font-size: 1.4rem;
  }
  
  .topic-content h4 {
    font-size: 1.2rem;
  }
}
</style>