<template>
  <div class="contact-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="container">
        <h1>联系我们</h1>
        <p>如有任何问题、建议或合作意向，请随时与我们联系</p>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="container">
      <div class="contact-content">
        <!-- 联系信息 -->
        <div class="contact-info">
          <h2 class="section-title">联系方式</h2>
          <div class="info-items">
            <div class="info-item">
              <div class="info-icon"><i class="fas fa-map-marker-alt"></i></div>
              <div class="info-text">
                <h3>办公地址</h3>
                <p>湖南省长沙市岳麓区岳麓大道567号文化大厦15层</p>
              </div>
            </div>
            <div class="info-item">
              <div class="info-icon"><i class="fas fa-phone"></i></div>
              <div class="info-text">
                <h3>联系电话</h3>
                <p>0731-88888888</p>
              </div>
            </div>
            <div class="info-item">
              <div class="info-icon"><i class="fas fa-envelope"></i></div>
              <div class="info-text">
                <h3>电子邮箱</h3>
                <p>contact@huxiangwenhua.com</p>
              </div>
            </div>
            <div class="info-item">
              <div class="info-icon"><i class="fas fa-clock"></i></div>
              <div class="info-text">
                <h3>工作时间</h3>
                <p>周一至周五: 9:00 - 17:30</p>
                <p>周六至周日: 10:00 - 16:00</p>
              </div>
            </div>
          </div>

          <!-- 社交媒体 -->
          <div class="social-media">
            <h3>关注我们</h3>
            <div class="social-links">
              <a href="#" class="social-link"><i class="fab fa-weixin"></i><span>微信公众号</span></a>
              <a href="#" class="social-link"><i class="fab fa-weibo"></i><span>新浪微博</span></a>
              <a href="#" class="social-link"><i class="fab fa-youtube"></i><span>YouTube</span></a>
              <a href="#" class="social-link"><i class="fab fa-tiktok"></i><span>抖音</span></a>
            </div>
          </div>

          <!-- 常见问题 -->
          <div class="faq-section">
            <h3>常见问题</h3>
            <div class="faq-items">
              <div class="faq-item" v-for="faq in faqs" :key="faq.id">
                <div class="faq-question" @click="toggleFaq(faq.id)">
                  <span>{{ faq.question }}</span>
                  <i :class="openFaqId === faq.id ? 'fas fa-chevron-down' : 'fas fa-chevron-right' "></i>
                </div>
                <div class="faq-answer" v-if="openFaqId === faq.id">
                  <p>{{ faq.answer }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 联系表单 -->
        <div class="contact-form-container">
          <h2 class="section-title">发送消息</h2>
          <form @submit.prevent="submitForm" class="contact-form">
            <div class="form-group">
              <label for="name">姓名</label>
              <input type="text" id="name" v-model="formData.name" required placeholder="请输入您的姓名" />
            </div>
            <div class="form-group">
              <label for="email">电子邮箱</label>
              <input type="email" id="email" v-model="formData.email" required placeholder="请输入您的电子邮箱" />
            </div>
            <div class="form-group">
              <label for="phone">联系电话</label>
              <input type="tel" id="phone" v-model="formData.phone" placeholder="请输入您的联系电话（选填）" />
            </div>
            <div class="form-group">
              <label for="subject">主题</label>
              <select id="subject" v-model="formData.subject" required>
                <option value="">请选择消息主题</option>
                <option value="suggestion">建议反馈</option>
                <option value="cooperation">合作咨询</option>
                <option value="resource">资源提交</option>
                <option value="question">问题咨询</option>
                <option value="other">其他</option>
              </select>
            </div>
            <div class="form-group">
              <label for="message">消息内容</label>
              <textarea id="message" v-model="formData.message" required placeholder="请输入您的消息内容" rows="6"></textarea>
            </div>
            <div class="form-group checkbox-group">
              <input type="checkbox" id="agreement" v-model="formData.agreement" required />
              <label for="agreement">我同意平台的<a href="#">隐私政策</a>和<a href="#">服务条款</a></label>
            </div>
            <button type="submit" class="submit-button" :disabled="isSubmitting">
              <span v-if="isSubmitting">发送中...</span>
              <span v-else>发送消息</span>
            </button>
          </form>
        </div>
      </div>

      <!-- 地图 -->
      <div class="map-section">
        <h2 class="section-title">位置地图</h2>
        <div class="map-container">
          <!-- 这里使用占位图片，实际项目中应该使用真实地图API -->
          <img src="https://picsum.photos/seed/map/1200/500" alt="位置地图" class="map-image" />
          <div class="map-overlay">
            <div class="map-marker"></div>
          </div>
        </div>
      </div>

      <!-- 合作咨询 -->
      <div class="cooperation-section">
        <div class="cooperation-content">
          <div class="cooperation-text">
            <h2>有合作意向？</h2>
            <p>如果您有湖湘文化相关的合作项目或资源想与我们分享，请联系我们的合作部门，我们期待与您携手共进，共同推动湖湘文化的传承与发展。</p>
          </div>
          <div class="cooperation-button">
            <a href="#" @click.prevent="goToCooperation" class="btn btn-primary">合作咨询</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'

export default {
  name: 'ContactPage',
  props: {
    showAlert: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    // 表单数据
    const formData = reactive({
      name: '',
      email: '',
      phone: '',
      subject: '',
      message: '',
      agreement: false
    })
    
    // 提交状态
    const isSubmitting = ref(false)
    
    // FAQ相关
    const openFaqId = ref(null)
    const faqs = ref([
      {
        id: '1',
        question: '如何在平台上提交湖湘文化资源？',
        answer: '您可以在"资源中心"页面点击"提交资源"按钮，按照提示填写资源信息并上传相关文件。我们的管理员会在3-5个工作日内审核您提交的资源，审核通过后会在平台上展示。'
      },
      {
        id: '2',
        question: '平台上的文化资源可以用于商业用途吗？',
        answer: '平台上的部分资源可能受版权保护，具体使用权限请查看每个资源的版权说明。如需将资源用于商业用途，建议您联系资源提供者或我们的版权部门获取授权。'
      },
      {
        id: '3',
        question: '如何加入平台的志愿者团队？',
        answer: '您可以通过联系页面的表单或直接发送邮件至volunteer@huxiangwenhua.com申请加入志愿者团队。我们会定期组织志愿者招募活动，欢迎对湖湘文化有热情的人士参与。'
      },
      {
        id: '4',
        question: '平台提供哪些类型的湖湘文化资源？',
        answer: '我们提供多种类型的湖湘文化资源，包括历史文物、传统技艺、民俗文化、地方戏曲、名人故事、历史建筑、特色小吃等。您可以通过资源分类或搜索功能快速找到您感兴趣的内容。'
      },
      {
        id: '5',
        question: '如何获取平台的最新活动信息？',
        answer: '您可以关注我们的官方微信公众号、微博或定期访问平台首页的"活动公告"栏目获取最新活动信息。我们也会通过邮件向注册用户发送重要活动通知。'
      }
    ])
    
    // 切换FAQ展开/收起
    const toggleFaq = (id) => {
      if (openFaqId.value === id) {
        openFaqId.value = null
      } else {
        openFaqId.value = id
      }
    }
    
    // 提交表单
    const submitForm = async () => {
      isSubmitting.value = true
      
      try {
        // 模拟表单提交
        await new Promise(resolve => setTimeout(resolve, 1500))
        
        // 显示成功提示
        if (props.showAlert) {
          props.showAlert('success', '您的消息已成功发送，我们将尽快回复您！')
        }
        
        // 重置表单
        formData.name = ''
        formData.email = ''
        formData.phone = ''
        formData.subject = ''
        formData.message = ''
        formData.agreement = false
      } catch (error) {
        // 显示错误提示
        if (props.showAlert) {
          props.showAlert('error', '消息发送失败，请稍后再试或通过其他方式联系我们。')
        }
        console.error('表单提交失败:', error)
      } finally {
        isSubmitting.value = false
      }
    }
    
    // 前往合作咨询
    const goToCooperation = () => {
      // 滚动到表单区域
      const formElement = document.querySelector('.contact-form-container')
      if (formElement) {
        formElement.scrollIntoView({ behavior: 'smooth' })
      }
      
      // 自动选择合作咨询主题
      formData.subject = 'cooperation'
    }
    
    return {
      formData,
      isSubmitting,
      openFaqId,
      faqs,
      toggleFaq,
      submitForm,
      goToCooperation
    }
  }
}
</script>

<style scoped>
.contact-page {
  padding: 2rem 0;
}

.page-header {
  background-image: linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5)), url('https://picsum.photos/seed/contact/1920/400');
  background-size: cover;
  background-position: center;
  color: white;
  padding: 4rem 0;
  text-align: center;
}

.page-header h1 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.page-header p {
  font-size: 1.2rem;
  max-width: 800px;
  margin: 0 auto;
}

.contact-content {
  display: flex;
  flex-wrap: wrap;
  gap: 2rem;
  margin: 3rem 0;
}

.contact-info,
.contact-form-container {
  flex: 1;
  min-width: 300px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 2rem;
}

.section-title {
  font-size: 1.8rem;
  margin-bottom: 1.5rem;
  color: #333;
  position: relative;
  padding-bottom: 0.5rem;
}

.section-title::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 60px;
  height: 3px;
  background-color: var(--primary-color);
}

/* 联系信息样式 */
.info-items {
  margin-bottom: 2rem;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.info-icon {
  font-size: 1.5rem;
  color: var(--primary-color);
  background-color: #f8f9fa;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.info-text h3 {
  margin: 0 0 0.5rem 0;
  color: #333;
  font-size: 1.1rem;
}

.info-text p {
  margin: 0;
  color: #666;
  line-height: 1.5;
}

/* 社交媒体样式 */
.social-media {
  margin-bottom: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #eee;
}

.social-media h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.social-links {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.social-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background-color: #f8f9fa;
  color: #333;
  border-radius: 4px;
  transition: all 0.3s;
}

.social-link:hover {
  background-color: var(--primary-color);
  color: white;
}

/* 常见问题样式 */
.faq-section {
  padding-top: 2rem;
  border-top: 1px solid #eee;
}

.faq-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.faq-item {
  border-bottom: 1px solid #eee;
}

.faq-item:last-child {
  border-bottom: none;
}

.faq-question {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  cursor: pointer;
  transition: color 0.3s;
}

.faq-question:hover {
  color: var(--primary-color);
}

.faq-question span {
  font-weight: 500;
  color: #333;
}

.faq-question:hover span {
  color: var(--primary-color);
}

.faq-answer {
  padding: 0 0 1rem 0;
}

.faq-answer p {
  color: #666;
  line-height: 1.5;
  margin: 0;
}

/* 联系表单样式 */
.contact-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  margin-bottom: 0.5rem;
  color: #333;
  font-weight: 500;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

.checkbox-group {
  flex-direction: row;
  align-items: center;
  gap: 0.5rem;
}

.checkbox-group input {
  width: auto;
}

.checkbox-group label {
  margin: 0;
  font-weight: normal;
  color: #666;
}

.checkbox-group a {
  color: var(--primary-color);
  text-decoration: none;
}

.checkbox-group a:hover {
  text-decoration: underline;
}

.submit-button {
  padding: 0.75rem 1.5rem;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s;
}

.submit-button:hover:not(:disabled) {
  background-color: #c62828;
}

.submit-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

/* 地图样式 */
.map-section {
  margin: 3rem 0;
}

.map-container {
  position: relative;
  width: 100%;
  height: 500px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.map-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.map-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.map-marker {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 30px;
  height: 30px;
  background-color: var(--primary-color);
  border-radius: 50%;
  box-shadow: 0 0 0 5px rgba(220, 38, 38, 0.3);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(220, 38, 38, 0.7);
  }
  70% {
    box-shadow: 0 0 0 15px rgba(220, 38, 38, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(220, 38, 38, 0);
  }
}

/* 合作咨询样式 */
.cooperation-section {
  margin: 3rem 0;
  background-image: linear-gradient(rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.7)), url('https://picsum.photos/seed/cooperation/1920/400');
  background-size: cover;
  background-position: center;
  color: white;
  padding: 3rem 0;
  border-radius: 8px;
}

.cooperation-content {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
}

.cooperation-text {
  flex: 1;
  min-width: 300px;
}

.cooperation-text h2 {
  margin: 0 0 1rem 0;
  font-size: 1.8rem;
}

.cooperation-text p {
  margin: 0;
  line-height: 1.6;
}

.cooperation-button {
  flex-shrink: 0;
}

.btn {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background-color: #c62828;
  transform: translateY(-2px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header h1 {
    font-size: 2rem;
  }
  
  .contact-content {
    flex-direction: column;
  }
  
  .contact-info,
  .contact-form-container {
    width: 100%;
  }
  
  .social-links {
    flex-direction: column;
  }
  
  .cooperation-content {
    text-align: center;
    justify-content: center;
  }
  
  .map-container {
    height: 300px;
  }
}
</style>