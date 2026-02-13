<template>
  <div class="architecture-3d-page">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <button class="back-button" @click="goBack">
        <i class="fas fa-arrow-left"></i>
        <span>返回数字化展示</span>
      </button>
    </div>

    <!-- 主标题区域 -->
    <div class="section-title">
      <h2>湖湘建筑3D模型</h2>
      <p>探索湖湘地区特色建筑的精细3D模型，了解建筑结构和文化内涵</p>
    </div>

    <!-- 内容区域 -->
    <div class="content-container">
      <!-- 专题介绍 -->
      <div class="introduction-section">
        <h3>专题介绍</h3>
        <p>湖湘建筑3D模型项目旨在通过先进的3D建模技术，对湖湘地区的特色建筑进行数字化重建和展示。我们选择了具有代表性的湖湘传统建筑，如岳阳楼、岳麓书院、橘子洲头等，通过高精度3D扫描和建模技术，创建了详细的建筑3D模型。</p>
        <p>该项目不仅展示了湖湘建筑的外观美感，还通过分层展示、结构解析等方式，让用户能够深入了解建筑的内部结构、构造原理和文化内涵。同时，我们还提供了建筑的历史背景、建筑特色、保护现状等相关信息，帮助用户全面认识湖湘建筑文化。</p>
      </div>

      <!-- 建筑模型展示 -->
      <div class="model-showcase">
        <h3>建筑模型展示</h3>
        <div class="model-grid">
          <!-- 建筑模型卡片 1 -->
          <div class="model-card">
            <div class="model-image-container">
              <img src="/src/assets/Imgs/岳阳楼记1.jpg" alt="岳阳楼3D模型" class="model-image">
              <div class="model-overlay">
                <button class="view-model-btn" @click="viewModel('yueyanglou')">
                  <i class="fas fa-cube"></i>
                  <span>查看3D模型</span>
                </button>
              </div>
            </div>
            <div class="model-info">
              <h4>岳阳楼</h4>
              <p>江南三大名楼之一，始建于东汉末年，因范仲淹《岳阳楼记》而闻名天下。</p>
              <div class="model-features">
                <span class="feature-tag">木质结构</span>
                <span class="feature-tag">江南风格</span>
                <span class="feature-tag">历史名楼</span>
              </div>
            </div>
          </div>

          <!-- 建筑模型卡片 2 -->
          <div class="model-card">
            <div class="model-image-container">
              <img src="/src/assets/Imgs/岳麓书院1.jpg" alt="岳麓书院3D模型" class="model-image">
              <div class="model-overlay">
                <button class="view-model-btn" @click="viewModel('yuelushuyuan')">
                  <i class="fas fa-cube"></i>
                  <span>查看3D模型</span>
                </button>
              </div>
            </div>
            <div class="model-info">
              <h4>岳麓书院</h4>
              <p>中国历史上赫赫闻名的四大书院之一，始建于北宋开宝九年（公元976年）。</p>
              <div class="model-features">
                <span class="feature-tag">传统书院</span>
                <span class="feature-tag">木质结构</span>
                <span class="feature-tag">文化遗产</span>
              </div>
            </div>
          </div>

          <!-- 建筑模型卡片 3 -->
          <div class="model-card">
            <div class="model-image-container">
              <img src="/src/assets/Imgs/衡山1.jpeg" alt="橘子洲3D模型" class="model-image">
              <div class="model-overlay">
                <button class="view-model-btn" @click="viewModel('juzizhou')">
                  <i class="fas fa-cube"></i>
                  <span>查看3D模型</span>
                </button>
              </div>
            </div>
            <div class="model-info">
              <h4>橘子洲</h4>
              <p>长沙著名的风景名胜区，位于湘江之中，因毛泽东《沁园春·长沙》而闻名。</p>
              <div class="model-features">
                <span class="feature-tag">园林建筑</span>
                <span class="feature-tag">现代重建</span>
                <span class="feature-tag">风景名胜</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 技术特点 -->
      <div class="technology-section">
        <h3>技术特点</h3>
        <div class="technology-features">
          <div class="tech-card">
            <div class="tech-icon">
              <i class="fas fa-ruler-combined fa-3x"></i>
            </div>
            <h4>高精度建模</h4>
            <p>采用先进的3D扫描技术，确保模型的几何精度和细节还原度，真实呈现建筑的每一个细节。</p>
          </div>
          <div class="tech-card">
            <div class="tech-icon">
              <i class="fas fa-layer-group fa-3x"></i>
            </div>
            <h4>分层展示</h4>
            <p>通过分层展示功能，用户可以查看建筑的不同结构层次，了解建筑的构造原理和空间布局。</p>
          </div>
          <div class="tech-card">
            <div class="tech-icon">
              <i class="fas fa-info-circle fa-3x"></i>
            </div>
            <h4>信息集成</h4>
            <p>集成建筑的历史背景、建筑特色、保护现状等相关信息，帮助用户全面认识湖湘建筑文化。</p>
          </div>
          <div class="tech-card">
            <div class="tech-icon">
              <i class="fas fa-mobile-alt fa-3x"></i>
            </div>
            <h4>多端适配</h4>
            <p>支持在电脑、平板、手机等多种设备上查看3D模型，提供流畅的交互体验。</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { useRouter } from 'vue-router'

// 获取路由实例
const router = useRouter()
// 获取全局事件总线
const eventBus = inject('$eventBus')

// 返回上一页
const goBack = () => {
  router.back()
}

// 查看3D模型
const viewModel = (buildingId) => {
  showAlert('3D模型查看功能即将上线，敬请期待！')
}

// 显示提示消息
const showAlert = (message, type = 'info') => {
  if (eventBus) {
    eventBus.emit('show-alert', { message, type })
  } else {
    console.log(`[${type}] ${message}`)
  }
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

.architecture-3d-page {
  padding: 2rem 0;
  background-color: #fafafa;
  min-height: 100vh;
}

/* 返回按钮样式 */
.back-button-container {
  position: fixed;
  top: 2rem;
  left: 2rem;
  z-index: 1000;
}

.back-button {
  background: var(--primary-color);
  color: white;
  border: none;
  padding: 0.8rem 1.5rem;
  border-radius: 25px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(200, 16, 46, 0.3);
}

.back-button:hover {
  background: #a60e24;
  transform: translateX(-5px) scale(1.05);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

.back-button i {
  font-size: 1.1rem;
}

/* 主标题样式 */
.section-title {
  text-align: center;
  margin-bottom: 3rem;
  padding: 0 1rem;
  margin-top: 60px;
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

/* 内容容器 */
.content-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

/* 介绍部分 */
.introduction-section {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  margin-bottom: 3rem;
}

.introduction-section h3 {
  font-size: 1.8rem;
  color: var(--text-color);
  margin-bottom: 1.5rem;
  font-family: var(--heading-font);
  font-weight: 600;
  position: relative;
  padding-left: 1rem;
}

.introduction-section h3::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  background-color: var(--primary-color);
}

.introduction-section p {
  color: var(--text-color);
  line-height: 1.7;
  margin-bottom: 1.2rem;
}

/* 模型展示 */
.model-showcase {
  margin-bottom: 3rem;
}

.model-showcase h3 {
  font-size: 1.8rem;
  color: var(--text-color);
  margin-bottom: 1.5rem;
  font-family: var(--heading-font);
  font-weight: 600;
  position: relative;
  padding-left: 1rem;
}

.model-showcase h3::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  background-color: var(--primary-color);
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 2rem;
}

.model-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.model-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.model-image-container {
  position: relative;
  overflow: hidden;
  height: 250px;
}

.model-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.model-card:hover .model-image {
  transform: scale(1.05);
}

.model-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.model-card:hover .model-overlay {
  opacity: 1;
}

.view-model-btn {
  background: var(--primary-color);
  color: white;
  border: none;
  padding: 1rem 2rem;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.view-model-btn:hover {
  background: #a60e24;
  transform: scale(1.05);
}

.model-info {
  padding: 1.5rem;
}

.model-info h4 {
  font-size: 1.3rem;
  color: var(--text-color);
  margin-bottom: 1rem;
  font-family: var(--heading-font);
  font-weight: 600;
}

.model-info p {
  color: var(--light-text);
  line-height: 1.7;
  margin-bottom: 1.5rem;
}

.model-features {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.feature-tag {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  padding: 0.3rem 0.8rem;
  border-radius: 20px;
  font-size: 0.8rem;
  color: var(--light-text);
}

/* 技术特点 */
.technology-section {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.technology-section h3 {
  font-size: 1.8rem;
  color: var(--text-color);
  margin-bottom: 1.5rem;
  font-family: var(--heading-font);
  font-weight: 600;
  position: relative;
  padding-left: 1rem;
}

.technology-section h3::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  background-color: var(--primary-color);
}

.technology-features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.tech-card {
  text-align: center;
  padding: 2rem;
  border-radius: 8px;
  background: #f8f9fa;
  transition: transform 0.3s ease;
}

.tech-card:hover {
  transform: translateY(-5px);
}

.tech-icon {
  color: var(--primary-color);
  margin-bottom: 1.5rem;
  background-color: rgba(200, 16, 46, 0.05);
  width: 100px;
  height: 100px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
}

.tech-card h4 {
  font-size: 1.2rem;
  color: var(--text-color);
  margin-bottom: 1rem;
  font-family: var(--heading-font);
  font-weight: 600;
}

.tech-card p {
  color: var(--light-text);
  line-height: 1.7;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .section-title h2 {
    font-size: 2rem;
  }
  
  .content-container {
    padding: 0 1rem;
  }
  
  .introduction-section,
  .technology-section {
    padding: 1.5rem;
  }
  
  .model-grid,
  .technology-features {
    gap: 1.5rem;
  }
  
  .model-card {
    margin-bottom: 1rem;
  }
  
  .back-button {
    padding: 0.6rem 1.2rem;
    font-size: 0.85rem;
  }
}

@media (max-width: 480px) {
  .section-title h2 {
    font-size: 1.8rem;
  }
  
  .model-grid,
  .technology-features {
    grid-template-columns: 1fr;
  }
  
  .tech-card {
    padding: 1.5rem;
  }
}
</style>