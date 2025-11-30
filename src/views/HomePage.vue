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
        <a href="#" class="feature-link" @click.prevent="navigateTo('community')">
          <i class="fas fa-users fa-3x"></i>
          <h3>互动社区</h3>
          <p>与志同道合的文化爱好者交流，参与文化活动，共同传承和弘扬湖湘文化。</p>
        </a>
      </div>
      <div class="feature">
        <a href="#" class="feature-link" @click.prevent="navigateTo('digital-showcase')">
          <i class="fas fa-laptop-code fa-3x"></i>
          <h3>数字化展示</h3>
          <p>运用现代信息技术，以数字化方式展示湖湘文化，让您获得更生动的体验。</p>
        </a>
      </div>
    </section>

    <!-- 文化资源展示 -->
    <section class="resources">
      <h2>精选文化资源</h2>
      <div class="resources-container">
        <div class="resource-card" v-for="resource in culturalResources" :key="resource.id">
          <div class="resource-image">
            <img :src="resource.image" :alt="resource.title" loading="lazy" />
          </div>
          <div class="resource-info">
            <h3>{{ resource.title }}</h3>
            <p>{{ resource.description }}</p>
            <div style="margin-top: auto; text-align: right; padding: 0 1rem 1rem;">
              <button class="btn" @click="browseMoreResources()">查看详情</button>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
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
    
    // 文化资源数据
    const culturalResources = ref([
      {
        id: 1,
        title: '岳麓书院',
        description: '中国历史上赫赫闻名的四大书院之一',
        image: 'https://img2.baidu.com/it/u=1868825345,3526076103&fm=253&app=138&f=JPEG?w=500&h=582',
        tags: ['岳麓书院', '历史建筑', '四大书院']
      },
      {
        id: 2,
        title: '湘绣',
        description: '中国四大名绣之一，以其精湛的技艺和独特的艺术风格闻名',
        image: 'https://pic.rmb.bdstatic.com/bjh/news/89c508f7a400f1751758f343caae8801.png',
        tags: ['湘绣', '传统工艺', '四大名绣']
      },
      {
        id: 3,
        title: '《岳阳楼记》',
        description: '江南三大名楼之一，范仲淹《岳阳楼记》千古传颂。',
        image: 'https://img0.baidu.com/it/u=1464457526,3199376473&fm=253&app=138&f=JPEG?w=800&h=2997',
        tags: ['岳阳楼', '文化遗产', '江南三大名楼']
      }
    ])

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

    // 页面加载时显示欢迎提示
    onMounted(() => {
      if (props.showAlert) {
        props.showAlert('欢迎来到湖湘文化数字资源平台！', 'success')
      }
    })

    return {
        culturalResources,
        browseMoreResources,
        goToResources,
        navigateTo
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
}
</style>