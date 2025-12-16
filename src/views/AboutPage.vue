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
        <div class="partners-grid">
          <div v-for="partner in partners" :key="partner.id" class="partner-item">
            <img :src="partner.logoUrl" :alt="partner.name" loading="lazy" />
            <p>{{ partner.name }}</p>
          </div>
        </div>
      </section>

      <!-- 数据统计 -->
      <section class="stats-section">
        <div class="stats-grid">
          <div v-for="stat in statistics" :key="stat.id" class="stat-item">
            <div class="stat-icon"><i :class="stat.icon"></i></div>
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </section>

      <!-- 文化理念 -->
      <section class="philosophy-section">
        <div class="philosophy-content">
          <div class="philosophy-item">
            <div class="philosophy-icon"><i class="fas fa-book-reader"></i></div>
            <h3>传承创新</h3>
            <p>在传承湖湘文化精髓的基础上，通过数字化手段实现创新表达和传播，让传统文化焕发时代活力。</p>
          </div>
          <div class="philosophy-item">
            <div class="philosophy-icon"><i class="fas fa-users"></i></div>
            <h3>开放共享</h3>
            <p>秉持开放共享的理念，为公众提供免费、便捷的湖湘文化资源获取渠道，促进文化的广泛传播。</p>
          </div>
          <div class="philosophy-item">
            <div class="philosophy-icon"><i class="fas fa-lightbulb"></i></div>
            <h3>融合发展</h3>
            <p>推动文化与科技、教育、旅游等产业的深度融合，实现湖湘文化的多元化发展和价值提升。</p>
          </div>
        </div>
      </section>

      <!-- 用户评价 -->
      <section class="testimonials-section">
        <h2 class="section-title">用户评价</h2>
        <div class="testimonials-slider">
          <div v-for="testimonial in testimonials" :key="testimonial.id" class="testimonial-item">
            <div class="testimonial-quote">
              <i class="fas fa-quote-left"></i>
              <p>{{ testimonial.quote }}</p>
            </div>
            <div class="testimonial-author">
              <img :src="testimonial.authorAvatar" :alt="testimonial.authorName" loading="lazy" />
              <div class="author-details">
                <h4>{{ testimonial.authorName }}</h4>
                <p>{{ testimonial.authorTitle }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
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
    
    // 统计数据
    const statistics = ref([
      {
        id: '1',
        icon: 'fas fa-database',
        value: '10,000+',
        label: '数字化资源'
      },
      {
        id: '2',
        icon: 'fas fa-users',
        value: '100,000+',
        label: '注册用户'
      },
      {
        id: '3',
        icon: 'fas fa-globe-asia',
        value: '20+',
        label: '合作机构'
      },
      {
        id: '4',
        icon: 'fas fa-eye',
        value: '1,000,000+',
        label: '月访问量'
      }
    ])
    
    // 用户评价数据
    const testimonials = ref([
      {
        id: '1',
        quote: '作为一名湖湘文化研究者，这个平台为我提供了丰富的数字化资源，极大地方便了我的研究工作。界面友好，操作简单，内容丰富，是了解湖湘文化的重要窗口。',
        authorName: '周教授',
        authorTitle: '湖南大学历史系教授',
        authorAvatar: 'https://picsum.photos/seed/testimonial1/100/100'
      },
      {
        id: '2',
        quote: '通过这个平台，我第一次如此直观地了解到湖湘文化的博大精深。3D文物展示和虚拟展厅功能非常震撼，让我仿佛身临其境，感受传统文化的魅力。',
        authorName: '王同学',
        authorTitle: '湖南师范大学学生',
        authorAvatar: 'https://picsum.photos/seed/testimonial2/100/100'
      },
      {
        id: '3',
        quote: '平台的数字地图功能非常实用，帮助我规划了一次精彩的湖湘文化之旅。每个文化地标的详细介绍和图片让我在出发前就对目的地有了深入的了解。',
        authorName: '张先生',
        authorTitle: '文化旅游爱好者',
        authorAvatar: 'https://picsum.photos/seed/testimonial3/100/100'
      }
    ])
    
    return {
      teamMembers,
      partners,
      statistics,
      testimonials
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

.partners-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.partner-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1.5rem;
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
  max-width: 150px;
  max-height: 80px;
  margin-bottom: 1rem;
  object-fit: contain;
}

.partner-item p {
  color: #333;
  font-weight: bold;
  text-align: center;
}

/* 数据统计样式 */
.stats-section {
  margin: 4rem 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.stat-item {
  text-align: center;
  padding: 2rem;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 2.5rem;
  color: var(--primary-color);
  margin-bottom: 1rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
  color: #333;
  margin-bottom: 0.5rem;
}

.stat-label {
  color: #666;
  font-size: 1rem;
}

/* 文化理念样式 */
.philosophy-section {
  margin: 4rem 0;
  background-color: #f8f9fa;
  padding: 3rem 0;
}

.philosophy-content {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.philosophy-item {
  text-align: center;
  padding: 2rem;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
}

.philosophy-item:hover {
  transform: translateY(-5px);
}

.philosophy-icon {
  font-size: 2.5rem;
  color: var(--primary-color);
  margin-bottom: 1rem;
}

.philosophy-item h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.philosophy-item p {
  color: #666;
  line-height: 1.5;
}

/* 用户评价样式 */
.testimonials-section {
  margin: 4rem 0;
}

.testimonials-slider {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.testimonial-item {
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.testimonial-quote {
  position: relative;
  margin-bottom: 2rem;
}

.testimonial-quote i {
  position: absolute;
  top: -10px;
  left: -10px;
  font-size: 2rem;
  color: var(--primary-color);
  opacity: 0.2;
}

.testimonial-quote p {
  color: #666;
  line-height: 1.6;
  padding-left: 2rem;
  font-style: italic;
}

.testimonial-author {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.testimonial-author img {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.author-details h4 {
  margin: 0 0 0.25rem 0;
  color: #333;
}

.author-details p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
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
  
  .team-grid,
  .partners-grid,
  .stats-grid,
  .philosophy-content,
  .testimonials-slider {
    grid-template-columns: 1fr;
  }
}
</style>