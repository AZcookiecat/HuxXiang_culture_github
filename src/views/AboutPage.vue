<template>
  <div class="about-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="container">
        <h1>关于我们</h1>
        <p>了解湖湘文化数字化开发与传播平台的使命、愿景和团队</p>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="container">
      <!-- 平台介绍 -->
      <section class="about-section">
        <div class="section-content">
          <div class="text-content">
            <h2>平台使命</h2>
            <p>湖湘文化数字化开发与传播平台致力于利用现代数字技术，系统性地整理、保护、展示和传播湖湘文化资源，让更多人了解、热爱和传承湖湘文化的精髓。</p>
            <p>我们相信，文化是一个民族的精神命脉，是人民的精神家园。通过数字化手段，我们能够突破时间和空间的限制，让湖湘文化焕发新的生机与活力。</p>
          </div>
          <div class="image-content">
            <img src="https://picsum.photos/seed/mission/600/400" alt="平台使命" loading="lazy" />
          </div>
        </div>
      </section>

      <!-- 核心团队 -->
      <section class="team-section">
        <h2 class="section-title">核心团队</h2>
        <div class="team-grid">
          <div v-for="member in teamMembers" :key="member.id" class="team-member">
            <div class="member-photo">
              <img :src="member.photoUrl" :alt="member.name" loading="lazy" />
            </div>
            <div class="member-info">
              <h3>{{ member.name }}</h3>
              <p class="member-position">{{ member.position }}</p>
              <p class="member-bio">{{ member.bio }}</p>
              <div class="member-social">
                <a href="#" v-if="member.email"><i class="fas fa-envelope"></i></a>
                <a href="#" v-if="member.linkedin"><i class="fab fa-linkedin"></i></a>
                <a href="#" v-if="member.twitter"><i class="fab fa-twitter"></i></a>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 合作伙伴 -->
      <section class="partners-section">
        <h2 class="section-title">合作伙伴</h2>
        <p class="section-description">我们与众多高校、研究机构、文化单位和企业建立了紧密的合作关系，共同推动湖湘文化的数字化保护与传播。</p>
        
        <!-- 合作伙伴轮播图 -->
        <div class="partners-carousel">
          <!-- 轮播控制按钮 -->
          <button class="carousel-btn carousel-btn-prev" @click="prevSlide">
            <i class="fas fa-chevron-left"></i>
          </button>
          
          <!-- 轮播轨道 -->
          <div class="carousel-track" :style="{ transform: `translateX(-${currentSlide * 100}%)` }">
            <!-- 轮播项（每组3个合作伙伴） -->
            <div v-for="(slide, index) in carouselSlides" :key="index" class="carousel-slide">
              <div v-for="partner in slide" :key="partner.id" class="partner-item">
                <img :src="partner.logoUrl" :alt="partner.name" loading="lazy" />
                <p>{{ partner.name }}</p>
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
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import XiaozhuImage from '../assets/imgs/xiaozhu.jpeg' // 导入图片

export default {
  name: 'AboutPage',
  props: {
    showAlert: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    
    // 团队成员数据
    const teamMembers = ref([
      {
        id: '1',
        name: '朱理婧',
        position: '项目指导老师',
        bio: '本科毕业于西南交通大学轨道交通信号与控制专业，硕士研究生毕业于英国曼彻斯特大学管理与信息系统专业，现任湖南工商大学数字媒体工程与人文学院专任教师。主要研究领域为管理与信息系统、数字经济等。',
        photoUrl: XiaozhuImage, // 使用导入的图片
        email: 'zhangming@example.com',
        linkedin: '#',
        twitter: '#'
      },
      {
        id: '2',
        name: '李华',
        position: '技术总监',
        bio: '计算机科学博士，具有丰富的数字化平台开发经验，负责平台的技术架构设计和实施。',
        photoUrl: 'https://picsum.photos/seed/member2/300/300',
        email: 'lihua@example.com',
        linkedin: '#',
        twitter: '#'
      },
      {
        id: '3',
        name: '王芳',
        position: '内容总监',
        bio: '文物与博物馆学硕士，负责平台文化资源的收集、整理和内容策划工作。',
        photoUrl: 'https://picsum.photos/seed/member3/300/300',
        email: 'wangfang@example.com',
        linkedin: '#',
        twitter: null
      },
      {
        id: '4',
        name: '赵强',
        position: '设计总监',
        bio: '数字媒体艺术硕士，负责平台的视觉设计和用户体验优化工作。',
        photoUrl: 'https://picsum.photos/seed/member4/300/300',
        email: 'zhaoqiang@example.com',
        linkedin: '#',
        twitter: '#'
      },
      {
        id: '5',
        name: '陈静',
        position: '运营总监',
        bio: '传播学硕士，负责平台的运营推广和用户社区建设工作。',
        photoUrl: 'https://picsum.photos/seed/member5/300/300',
        email: 'chenjing@example.com',
        linkedin: '#',
        twitter: '#'
      },
      {
        id: '6',
        name: '刘博',
        position: '研究员',
        bio: '文化遗产学博士，负责平台文化内容的学术审核和研究工作。',
        photoUrl: 'https://picsum.photos/seed/member6/300/300',
        email: 'liubo@example.com',
        linkedin: '#',
        twitter: null
      }
    ])
    
    // 合作伙伴数据
    const partners = ref([
      {
        id: '1',
        name: '湖南省博物馆',
        logoUrl: 'https://picsum.photos/seed/partner1/200/100'
      },
      {
        id: '2',
        name: '岳麓书院',
        logoUrl: 'https://picsum.photos/seed/partner2/200/100'
      },
      {
        id: '3',
        name: '湖南大学',
        logoUrl: 'https://picsum.photos/seed/partner3/200/100'
      },
      {
        id: '4',
        name: '中南大学',
        logoUrl: 'https://picsum.photos/seed/partner4/200/100'
      },
      {
        id: '5',
        name: '湖南省文化和旅游厅',
        logoUrl: 'https://picsum.photos/seed/partner5/200/100'
      },
      {
        id: '6',
        name: '湖南出版集团',
        logoUrl: 'https://picsum.photos/seed/partner6/200/100'
      }
    ])
    
    // 轮播图状态
    const currentSlide = ref(0)
    let autoplayTimer = null
    
    // 将合作伙伴数据分成每组3个
    const carouselSlides = computed(() => {
      const slides = []
      for (let i = 0; i < partners.value.length; i += 3) {
        slides.push(partners.value.slice(i, i + 3))
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
    
    // 组件挂载时启动自动轮播
    onMounted(() => {
      startAutoplay()
    })
    
    // 组件卸载时停止自动轮播
    onUnmounted(() => {
      stopAutoplay()
    })
    
    return {
      teamMembers,
      partners,
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
.about-page {
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

.about-section {
  margin: 3rem 0;
}

.section-content {
  display: flex;
  flex-wrap: wrap;
  gap: 2rem;
  align-items: center;
}

.text-content {
  flex: 1;
  min-width: 300px;
}

.text-content h2 {
  font-size: 1.8rem;
  margin-bottom: 1.5rem;
  color: #333;
}

.text-content p {
  color: #666;
  line-height: 1.6;
  margin-bottom: 1rem;
}

.image-content {
  flex: 1;
  min-width: 300px;
}

.image-content img {
  width: 100%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.section-title {
  text-align: center;
  font-size: 2rem;
  margin-bottom: 2rem;
  color: #333;
}

.section-description {
  text-align: center;
  color: #666;
  max-width: 800px;
  margin: 0 auto 2rem auto;
}

/* 发展历程样式 */
.timeline-section {
  margin: 4rem 0;
}

.timeline {
  position: relative;
  max-width: 800px;
  margin: 0 auto;
}

.timeline::after {
  content: '';
  position: absolute;
  width: 4px;
  background-color: var(--primary-color);
  top: 0;
  bottom: 0;
  left: 50%;
  margin-left: -2px;
}

.timeline-item {
  padding: 10px 40px;
  position: relative;
  width: 50%;
  box-sizing: border-box;
}

.timeline-item:nth-child(odd) {
  left: 0;
}

.timeline-item:nth-child(even) {
  left: 50%;
}

.timeline-date {
  position: absolute;
  top: 15px;
  width: 80px;
  background-color: var(--primary-color);
  color: white;
  text-align: center;
  padding: 5px 10px;
  border-radius: 4px;
  font-weight: bold;
  z-index: 1;
}

.timeline-item:nth-child(odd) .timeline-date {
  right: -40px;
}

.timeline-item:nth-child(even) .timeline-date {
  left: -40px;
}

.timeline-content {
  padding: 20px;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.timeline-content h3 {
  margin-top: 0;
  color: #333;
}

.timeline-content p {
  color: #666;
  line-height: 1.5;
}

/* 团队成员样式 */
.team-section {
  margin: 4rem 0;
}

.team-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}

.team-member {
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
}

.team-member:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.member-photo {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  margin: 0 auto 1rem auto;
}

.member-photo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.member-info h3 {
  margin: 0 0 0.5rem 0;
  color: #333;
}

.member-position {
  color: var(--primary-color);
  font-weight: bold;
  margin-bottom: 1rem;
}

.member-bio {
  color: #666;
  line-height: 1.5;
  margin-bottom: 1rem;
}

.member-social {
  display: flex;
  justify-content: center;
  gap: 1rem;
}

.member-social a {
  color: #666;
  font-size: 1.2rem;
  transition: color 0.3s;
}

.member-social a:hover {
  color: var(--primary-color);
}

/* 合作伙伴样式 */
.partners-section {
  margin: 4rem 0;
  background-color: #f8f9fa;
  padding: 3rem 0;
}

/* 合作伙伴轮播图样式 */
.partners-carousel {
  position: relative;
  max-width: 1000px;
  margin: 0 auto 2rem;
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
  padding: 0 1rem;
}

.partner-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
}

.partner-item:hover {
  transform: translateY(-5px);
}

.partner-item img {
  max-width: 180px;
  max-height: 100px;
  margin-bottom: 1.5rem;
  object-fit: contain;
}

.partner-item p {
  color: #333;
  font-weight: bold;
  text-align: center;
  font-size: 1.1rem;
}

/* 轮播控制按钮 */
.carousel-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: #C8102E;
  border: none;
  border-radius: 50%;
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 1.2rem;
  color: white;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
  z-index: 10;
}

.carousel-btn:hover {
  background: #a60e24;
  transform: translateY(-50%) scale(1.1);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.4);
}

.carousel-btn-prev {
  left: 0.5rem;
}

.carousel-btn-next {
  right: 0.5rem;
}

/* 轮播指示器 */
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
  border: 2px solid #C8102E;
  background: transparent;
  cursor: pointer;
  transition: all 0.3s ease;
}

.carousel-indicator.active {
  background: #C8102E;
  border-color: #C8102E;
  transform: scale(1.2);
}

.carousel-indicator:hover {
  background: rgba(200, 16, 46, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header h1 {
    font-size: 2rem;
  }
  
  .section-content {
    flex-direction: column;
  }
  
  .text-content,
  .image-content {
    width: 100%;
  }
  
  .timeline::after {
    left: 31px;
  }
  
  .timeline-item {
    width: 100%;
    padding-left: 70px;
    padding-right: 25px;
  }
  
  .timeline-item:nth-child(even) {
    left: 0;
  }
  
  .timeline-item:nth-child(odd) .timeline-date,
  .timeline-item:nth-child(even) .timeline-date {
    left: 0;
    right: auto;
  }
  
  .team-grid {
    grid-template-columns: 1fr;
  }
  
  .partners-carousel {
    max-width: 90%;
  }
  
  .carousel-slide {
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .partner-item {
    padding: 1.5rem;
  }
  
  .partner-item img {
    max-width: 150px;
    max-height: 80px;
  }
  
  .partner-item p {
    font-size: 1rem;
  }
  
  .carousel-btn {
    width: 40px;
    height: 40px;
    font-size: 1rem;
  }
  
  .carousel-indicators {
    gap: 0.3rem;
  }
  
  .carousel-indicator {
    width: 8px;
    height: 8px;
  }
}
</style>