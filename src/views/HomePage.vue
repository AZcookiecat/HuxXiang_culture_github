<template>
  <div class="home-page">
    <!-- 英雄区域 -->
    <section class="hero">
      <div class="hero-overlay"></div>
      <div class="hero-content">
          <h1>探索千年湖湘文化</h1>
          <p>在这里，您可以深入了解湖湘大地丰富的文化遗产，从历史遗迹到传统艺术，感受千年文脉的魅力。</p>
          <div class="hero-buttons">
            <a href="#" class="btn" @click.prevent="goToResources()">开始探索</a>
            <a href="#" class="btn secondary" @click.prevent="goToResources()">查看更多资源</a>
          </div>
        </div>
    </section>

    <!-- 特色区域 -->
    <section class="features">
      <div class="feature">
        <a href="#" class="feature-link" @click.prevent="navigateTo('cultural-resources')">
          <i class="fas fa-book-reader fa-3x"></i>
          <h3>丰富的文化资源</h3>
          <p>汇集湖湘地区的历史文献、文物古迹、传统艺术等各类文化资源，为您提供全面的文化信息。</p>
        </a>
      </div>
      <div class="feature">
        <a href="#" class="feature-link" @click.prevent="navigateTo('digital-showcase')">
          <i class="fas fa-laptop-code fa-3x"></i>
          <h3>数字化展示</h3>
          <p>运用现代信息技术，以数字化方式展示湖湘文化，让您获得更生动的体验。</p>
        </a>
      </div>
      <div class="feature">
        <a href="#" class="feature-link" @click.prevent="navigateTo('community')">
          <i class="fas fa-users fa-3x"></i>
          <h3>互动社区</h3>
          <p>与志同道合的文化爱好者交流，参与文化活动，共同传承和弘扬湖湘文化。</p>
        </a>
      </div>
    </section>

    <!-- 文化资源展示 -->
    <section class="resources">
      <h2>精选文化资源</h2>
      
      <!-- 轮播图容器 -->
      <div class="carousel-container">
        <!-- 轮播控制按钮 -->
        <button class="carousel-btn carousel-btn-prev" @click="prevSlide">
          <i class="fas fa-chevron-left"></i>
        </button>
        
        <!-- 轮播轨道 -->
        <div class="carousel-track" :style="{ transform: `translateX(-${currentSlide * 100}%)` }">
          <!-- 轮播项（每组3个卡片） -->
          <div class="carousel-slide" v-for="(slide, index) in carouselSlides" :key="index">
            <div class="resource-card" v-for="resource in slide" :key="resource.id">
              <div class="resource-image">
                <img :src="resource.image" :alt="resource.title" loading="lazy" />
              </div>
              <div class="resource-info">
                <h3>{{ resource.title }}</h3>
                <p>{{ resource.description }}</p>
                <div style="margin-top: auto; text-align: right; padding: 0 1rem 1rem;">
                  <button class="btn" @click="viewResourceDetails(resource.id)">查看详情</button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 轮播控制按钮 -->
        <button class="carousel-btn carousel-btn-next" @click="nextSlide">
          <i class="fas fa-chevron-right"></i>
        </button>
      </div>
      
      <!-- 轮播指示器 -->
      <div class="carousel-indicators">
        <button 
          v-for="(slide, index) in carouselSlides" 
          :key="index"
          class="carousel-indicator"
          :class="{ active: currentSlide === index }"
          @click="goToSlide(index)"
        ></button>
      </div>
    </section>

    <!-- 精品数字化展示 -->
    <section class="digital-showcase-home">
      <h2>精品数字化展示</h2>
      <div class="digital-showcase-container">
        <!-- 文化知识图谱 -->
        <div class="digital-card">
          <router-link to="/knowledge-graph" class="digital-link">
            <div class="digital-icon">
              <i class="fas fa-project-diagram fa-3x"></i>
            </div>
            <h3>文化知识图谱</h3>
            <p>通过可视化图谱展示湖湘文化各元素之间的关联关系</p>
            <button class="start-btn">开始查看</button>
          </router-link>
        </div>

        <!-- 虚拟现实体验 -->
        <div class="digital-card">
          <router-link to="/unity-webgl" class="digital-link">
            <div class="digital-icon">
              <i class="fas fa-vr-cardboard fa-3x"></i>
            </div>
            <h3>虚拟现实体验</h3>
            <p>沉浸式VR体验湖湘文化名胜古迹</p>
            <button class="start-btn">开始体验</button>
          </router-link>
        </div>

        <!-- AI文化助手 -->
        <div class="digital-card">
          <router-link to="/ai-assistant" class="digital-link">
            <div class="digital-icon">
              <i class="fas fa-robot fa-3x"></i>
            </div>
            <h3>AI文化助手</h3>
            <p>智能问答系统，解答您的湖湘文化问题</p>
            <button class="start-btn">开始对话</button>
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'HomePage',
  props: {
    showAlert: {
      type: Function,
      default: null
    },
    goToResources: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    // 初始化router
    const router = useRouter()
    
    // 文化资源数据（12个）
    const culturalResources = ref([
      {
        id: '1',
        title: '岳麓书院',
        description: '中国历史上赫赫闻名的四大书院之一',
        image: 'https://img2.baidu.com/it/u=1868825345,3526076103&fm=253&app=138&f=JPEG?w=500&h=582',
        tags: ['岳麓书院', '历史建筑', '四大书院']
      },
      {
        id: '2',
        title: '湘绣',
        description: '中国四大名绣之一，以其精湛的技艺和独特的艺术风格闻名',
        image: 'https://pic.rmb.bdstatic.com/bjh/news/89c508f7a400f1751758f343caae8801.png',
        tags: ['湘绣', '传统工艺', '四大名绣']
      },
      {
        id: '3',
        title: '岳阳楼',
        description: '江南三大名楼之一，范仲淹《岳阳楼记》千古传颂。',
        image: 'https://img0.baidu.com/it/u=1464457526,3199376473&fm=253&app=138&f=JPEG?w=800&h=2997',
        tags: ['岳阳楼', '文化遗产', '江南三大名楼']
      },
       {
        id: '4',
        title: '《桃花源记》',
        description: '陶渊明的代表作，描绘了一个与世隔绝的理想社会。',
        categoryId: 'literature',
        image: 'https://shufa.pku.edu.cn/images/2024-09/90f030ce-1d59-4f14-bfce-00f3534cb9f0.jpeg',
        createdAt: '2025-3-24'
      },
      {
        id: '5',
        title: '湘菜',
        description: '中国八大菜系之一，以辣、香、鲜、嫩为特色。',
        categoryId: 'food',
        image: 'https://qcloud.dpfile.com/pc/y41poElenLBv7nqvkBnoBORgilWDRCfeVQ0djkNKJwMID0bsu-SAtCe9N8iBmkY3.jpg',
        createdAt: '2025-3-24'
      },
      {
        id: '6',
        title: '汨罗端午节',
        description: '纪念屈原的传统节日，有赛龙舟、吃粽子等习俗。',
        categoryId: 'folk',
        image: 'https://pic.rmb.bdstatic.com/bjh/bc1190b5ccbb/250505/b762f2bcc5ec3a93bf3b092a3eef6351.jpeg',
        createdAt: '2025-3-24'
      },
      {
        id: '7',
        title: '衡山',
        description: '中国五岳之一，被誉为"南岳"，是著名的道教和佛教圣地。',
        categoryId: 'history',
        image: 'https://img0.baidu.com/it/u=4250670926,3034496007&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=667',
        createdAt: '2025-3-24'
      },
      {
        id: '8',
        title: '花鼓戏',
        description: '湖南省特有的地方戏曲，具有浓郁的乡土气息。',
        categoryId: 'art',
        image: 'https://p7.itc.cn/q_70/images03/20220411/4af6d947f76f48289ab31aeba6936b91.jpeg',
        createdAt: '2025-3-24'
      },
      {
        id: '9',
        title: '曾国藩故居',
        description: '清代名臣曾国藩的出生地，展示了湘军文化的历史。',
        categoryId: 'history',
        image: 'https://img0.baidu.com/it/u=21139302,2048799699&fm=253&app=138&f=JPEG?w=500&h=652',
        createdAt: '2025-3-24'
      },
      {
        id: '10',
        title: '马王堆汉墓',
        description: '出土了大量珍贵文物的汉代墓葬，为研究汉代历史提供了重要资料。',
        categoryId: 'history',
        image: 'https://gips0.baidu.com/it/u=1600549162,356730080&fm=3074&app=3074&f=JPEG',
        createdAt: '2025-3-24'
      },
      {
        id: '11',
        title: '火宫殿',
        description: '长沙著名的传统小吃聚集地，有"中华老字号"之称。',
        categoryId: 'food',
        image: 'https://img2.baidu.com/it/u=3401007738,1338386073&fm=253&app=138&f=JPEG?w=800&h=1199',
        createdAt: '2025-3-24'
      },
      {
        id: '12',
        title: '铜官窑',
        description: '唐代著名的青瓷窑址，是海上丝绸之路的重要起点之一。',
        categoryId: 'art',
        image: 'https://img2.baidu.com/it/u=3735106663,128794723&fm=253&app=138&f=JPEG?w=800&h=1067',
        createdAt: '2025-3-24'
      }
    ])
    
    // 轮播图状态
    const currentSlide = ref(0)
    let autoplayTimer = null
    
    // 将文化资源数据分成每组3个
    const carouselSlides = computed(() => {
      const slides = []
      for (let i = 0; i < culturalResources.value.length; i += 3) {
        slides.push(culturalResources.value.slice(i, i + 3))
      }
      return slides
    })
    
    // 上一张幻灯片
    const prevSlide = () => {
      if (currentSlide.value > 0) {
        currentSlide.value--
      } else {
        currentSlide.value = carouselSlides.value.length - 1
      }
    }
    
    // 下一张幻灯片
    const nextSlide = () => {
      if (currentSlide.value < carouselSlides.value.length - 1) {
        currentSlide.value++
      } else {
        currentSlide.value = 0
      }
    }
    
    // 跳转到指定幻灯片
    const goToSlide = (index) => {
      currentSlide.value = index
    }
    
    // 自动轮播
    const startAutoplay = () => {
      autoplayTimer = setInterval(() => {
        nextSlide()
      }, 5000) // 每5秒自动切换
    }
    
    // 停止自动轮播
    const stopAutoplay = () => {
      if (autoplayTimer) {
        clearInterval(autoplayTimer)
        autoplayTimer = null
      }
    }

    // 导航到资源页
    const goToResources = () => {
      router.push({ name: 'cultural-resources' })
    }
    
    // 导航到其他页面
    const navigateTo = (page) => {
      router.push({ name: page })
    }

    // 浏览更多资源
    const browseMoreResources = () => {
      router.push({ name: 'cultural-resources' })
    }
    
    // 查看资源详情
    const viewResourceDetails = (resourceId) => {
      router.push({ name: 'resource-detail', params: { id: resourceId } })
    }

    // 页面加载时显示欢迎提示
    onMounted(() => {
      if (props.showAlert) {
        props.showAlert('欢迎来到湖湘文化数字资源平台！', 'success')
      }
      startAutoplay()
    })
    
    // 组件卸载时停止自动轮播
    onUnmounted(() => {
      stopAutoplay()
    })

    return {
        culturalResources,
        browseMoreResources,
        goToResources,
        navigateTo,
        viewResourceDetails,
        currentSlide,
        carouselSlides,
        prevSlide,
        nextSlide,
        goToSlide
      }
  }
}
</script>

<style scoped>
/* 英雄区域按钮样式 */
.hero-buttons {
  display: flex;
  justify-content: center;
  gap: 15px;
  flex-wrap: wrap;
}

/* 资源卡片样式优化 */
.resources-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.resource-card {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.resource-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.resource-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.resource-card:hover .resource-image img {
  transform: scale(1.05);
}

.resource-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 加载更多样式 */
.load-more {
  text-align: center;
  margin-top: 30px;
}

.loading {
  text-align: center;
  margin-top: 30px;
  color: var(--primary-color);
  font-size: 1.1rem;
}

/* 轮播图样式 */
.carousel-container {
  position: relative;
  max-width: 1200px;
  margin: 0 auto 3rem;
  overflow: hidden;
  padding: 0 4rem;
}

.carousel-track {
  display: flex;
  transition: transform 0.5s ease-in-out;
}

.carousel-slide {
  flex: 0 0 100%;
  display: flex;
  gap: 2.5rem;
  padding: 1rem 0;
  transition: all 0.5s ease-in-out;
}

.carousel-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(200, 16, 46, 0.9);
  border: none;
  border-radius: 50%;
  width: 45px;
  height: 45px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 1.1rem;
  color: white;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
  z-index: 10;
}

.carousel-btn:hover {
  background: rgba(200, 16, 46, 1);
  transform: translateY(-50%) scale(1.05);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.4);
}

.carousel-btn-prev {
  left: 0.5rem;
}

.carousel-btn-next {
  right: 0.5rem;
}

.carousel-indicators {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 2rem;
}

.carousel-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 2px solid var(--secondary-color);
  background: transparent;
  cursor: pointer;
  transition: all 0.3s ease;
}

.carousel-indicator.active {
  background: var(--primary-color);
  border-color: var(--primary-color);
  transform: scale(1.2);
}

.carousel-indicator:hover {
  background: rgba(200, 16, 46, 0.2);
}

/* 资源卡片在轮播中的样式调整 */
.carousel-slide .resource-card {
  flex: 1;
  max-width: calc(33.333% - 2rem);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  border-radius: 15px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.carousel-slide .resource-card:hover {
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.15);
  transform: translateY(-5px);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .carousel-container {
    max-width: 90%;
  }
  
  .carousel-slide {
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .carousel-slide .resource-card {
    max-width: 100%;
  }
  
  .carousel-btn {
    width: 40px;
    height: 40px;
    font-size: 1rem;
  }
}

/* 精品数字化展示样式 */
.digital-showcase-home {
  background-color: #f9f9f9;
  padding: 4rem 2rem;
  text-align: center;
}

.digital-showcase-home h2 {
  text-align: center;
  margin-bottom: 2rem;
  font-size: 2rem;
  color: var(--secondary-color);
  position: relative;
  display: inline-block;
  padding-bottom: 0.5rem;
}

.digital-showcase-home h2::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100px;
  height: 3px;
  background-color: var(--primary-color);
  border-radius: 2px;
}

.digital-showcase-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.digital-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.digital-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.digital-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem;
  text-decoration: none;
  color: inherit;
  height: 100%;
}

.digital-icon {
  width: 100px;
  height: 100px;
  background: linear-gradient(135deg, #C8102E, #e53935);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1.5rem;
  color: white;
  box-shadow: 0 4px 15px rgba(200, 16, 46, 0.3);
  transition: all 0.3s ease;
}

.digital-card:hover .digital-icon {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(200, 16, 46, 0.4);
}

.digital-link h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: #333;
}

.digital-link p {
  color: #666;
  margin-bottom: 2rem;
  line-height: 1.6;
  flex-grow: 1;
}

.start-btn {
  background: #C8102E;
  color: white;
  border: none;
  padding: 0.8rem 1.8rem;
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(200, 16, 46, 0.3);
}

.start-btn:hover {
  background: #a60e24;
  transform: translateY(-3px);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .hero h1 {
    font-size: 2rem;
  }
  
  .hero p {
    font-size: 1rem;
  }
  
  .hero-buttons {
    flex-direction: column;
    align-items: center;
  }
  
  .hero-buttons .btn {
    width: 200px;
    margin: 5px 0;
  }
  
  .digital-showcase-home {
    padding: 3rem 1.5rem;
  }
  
  .digital-showcase-home h2 {
    font-size: 2rem;
  }
  
  .digital-showcase-container {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
  
  .digital-link {
    padding: 1.5rem;
  }
}

@media (max-width: 768px) {
  .carousel-slide {
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .carousel-slide .resource-card {
    max-width: 100%;
    min-width: auto;
  }
  
  .carousel-btn {
    width: 40px;
    height: 40px;
    font-size: 1rem;
  }
  
  .digital-showcase-home {
    padding: 3rem 1.5rem;
  }
  
  .digital-showcase-home h2 {
    font-size: 2rem;
  }
  
  .digital-showcase-container {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
  
  .digital-link {
    padding: 1.5rem;
  }
}

@media (max-width: 480px) {
  .carousel-btn {
    width: 35px;
    height: 35px;
    font-size: 0.9rem;
  }
  
  .carousel-indicators {
    gap: 0.3rem;
  }
  
  .carousel-indicator {
    width: 10px;
    height: 10px;
  }
  
  .digital-showcase-home {
    padding: 2rem 1rem;
  }
  
  .digital-showcase-home h2 {
    font-size: 1.8rem;
  }
  
  .digital-icon {
    width: 80px;
    height: 80px;
  }
  
  .digital-icon i {
    font-size: 2.5rem;
  }
  
  .digital-link h3 {
    font-size: 1.3rem;
  }
  
  .start-btn {
    padding: 0.7rem 1.5rem;
    font-size: 0.9rem;
  }
}
</style>