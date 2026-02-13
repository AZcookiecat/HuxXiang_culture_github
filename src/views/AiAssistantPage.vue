<template>
  <div class="ai-assistant-page">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <button class="back-button" @click="goBack">
        <i class="fas fa-arrow-left"></i>
        <span>返回数字化展示</span>
      </button>
    </div>

    <!-- 主标题区域 -->
    <div class="section-title">
      <h2>AI文化助手</h2>
      <p>智能问答系统，解答您的湖湘文化问题</p>
    </div>

    <!-- AI助手主体 -->
    <div class="ai-assistant-container">
      <!-- 左侧菜单 -->
      <div class="ai-menu">
        <div class="menu-header">
          <h4>分类导航</h4>
        </div>
        <ul class="menu-items">
          <li class="menu-item" :class="{ active: currentCategory === '文化问答' }" @click="selectCategory('文化问答')">
            <i class="fas fa-question-circle"></i>
            <span>文化问答</span>
          </li>
          <li class="menu-item" :class="{ active: currentCategory === '历史人物' }" @click="selectCategory('历史人物')">
            <i class="fas fa-user-circle"></i>
            <span>历史人物</span>
          </li>
          <li class="menu-item" :class="{ active: currentCategory === '文化遗产' }" @click="selectCategory('文化遗产')">
            <i class="fas fa-university"></i>
            <span>文化遗产</span>
          </li>
          <li class="menu-item" :class="{ active: currentCategory === '传统习俗' }" @click="selectCategory('传统习俗')">
            <i class="fas fa-calendar-alt"></i>
            <span>传统习俗</span>
          </li>
          <li class="menu-item" :class="{ active: currentCategory === '湖湘美食' }" @click="selectCategory('湖湘美食')">
            <i class="fas fa-utensils"></i>
            <span>湖湘美食</span>
          </li>
          <li class="menu-item" :class="{ active: currentCategory === '旅游景点' }" @click="selectCategory('旅游景点')">
            <i class="fas fa-map-marker-alt"></i>
            <span>旅游景点</span>
          </li>
        </ul>
      </div>
      
      <!-- 聊天界面 -->
      <div class="chat-container">
        <div class="chat-header">
          <h4>{{ currentCategory }}</h4>
        </div>
        <div class="chat-messages" ref="chatMessages">
          <!-- 系统消息 -->
          <div class="message system-message">
            <div class="message-bubble">
              <p>您好！我是AI文化助手，很高兴为您解答湖湘文化相关问题。</p>
            </div>
          </div>
          
          <!-- 对话消息 -->
          <div v-for="(msg, index) in messages" :key="index" :class="['message', msg.type === 'user' ? 'user-message' : 'ai-message']">
            <div class="message-bubble">
              <p>{{ msg.content }}</p>
            </div>
          </div>
        </div>
        <div class="chat-input-container">
          <input type="text" v-model="inputMessage" placeholder="请输入您的问题..." @keyup.enter="sendMessage">
          <button class="send-btn" @click="sendMessage">
            <i class="fas fa-paper-plane"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'

// 路由实例
const router = useRouter()

// AI助手相关状态
const currentCategory = ref('文化问答')
const messages = ref([])
const inputMessage = ref('')
const chatMessages = ref(null)

// 返回上一页
const goBack = () => {
  router.back()
}

// 选择分类
const selectCategory = (category) => {
  currentCategory.value = category
  // 清空消息历史
  messages.value = []
  // 滚动到底部
  scrollToBottom()
}

// 发送消息
const sendMessage = () => {
  if (inputMessage.value.trim()) {
    // 添加用户消息
    messages.value.push({
      type: 'user',
      content: inputMessage.value.trim()
    })
    
    // 清空输入框
    inputMessage.value = ''
    
    // 滚动到底部
    scrollToBottom()
    
    // 模拟AI回复
    setTimeout(() => {
      let reply = ''
      switch (currentCategory.value) {
        case '文化问答':
          reply = '湖湘文化是中国传统文化的重要组成部分，以其独特的地域特色和历史底蕴而闻名。'
          break
        case '历史人物':
          reply = '湖湘地区涌现了许多历史名人，如曾国藩、左宗棠、毛泽东等。'
          break
        case '文化遗产':
          reply = '湖湘地区有许多国家级非物质文化遗产，如湘剧、花鼓戏、剪纸等。'
          break
        case '传统习俗':
          reply = '湖湘地区的传统习俗丰富多彩，如春节的贴春联、吃团圆饭，端午节的赛龙舟等。'
          break
        case '湖湘美食':
          reply = '湖湘美食以辣著称，代表菜品有红烧肉、剁椒鱼头、麻辣小龙虾等。'
          break
        case '旅游景点':
          reply = '湖湘地区有许多著名的旅游景点，如张家界、岳阳楼、橘子洲头等。'
          break
        default:
          reply = '感谢您的提问，我会为您提供更多湖湘文化相关的信息。'
      }
      
      messages.value.push({
        type: 'ai',
        content: reply
      })
      
      // 滚动到底部
      scrollToBottom()
    }, 1000)
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (chatMessages.value) {
      chatMessages.value.scrollTop = chatMessages.value.scrollHeight
    }
  })
}

// 监听消息变化，自动滚动到底部
watch(messages, () => {
  scrollToBottom()
}, { deep: true })
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

.ai-assistant-page {
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
  animation: backButtonFloat 1.5s infinite;
}

.back-button i {
  font-size: 1.1rem;
}

@keyframes backButtonFloat {
  0%, 100% {
    transform: translateX(-5px) scale(1.05);
  }
  50% {
    transform: translateX(0px) scale(1.05);
  }
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

/* AI助手容器 */
.ai-assistant-container {
  display: flex;
  max-width: 1200px;
  margin: 0 auto;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  min-height: 700px;
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

.menu-item span {
  font-size: 0.95rem;
  font-weight: 500;
}

/* 聊天界面样式 */
.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fafafa;
}

.chat-header {
  padding: 1.5rem;
  background: white;
  border-bottom: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.chat-header h4 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-color);
}

.chat-messages {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

/* 消息气泡样式 */
.message {
  display: flex;
  max-width: 80%;
  animation: fadeIn 0.3s ease;
}

.system-message {
  align-self: center;
  max-width: 90%;
}

.user-message {
  align-self: flex-end;
  justify-content: flex-end;
}

.ai-message {
  align-self: flex-start;
  justify-content: flex-start;
}

.message-bubble {
  padding: 1rem 1.2rem;
  border-radius: 18px;
  line-height: 1.5;
  font-size: 0.95rem;
}

.system-message .message-bubble {
  background: rgba(200, 16, 46, 0.1);
  color: var(--text-color);
  border-radius: 12px;
  font-style: italic;
}

.user-message .message-bubble {
  background: var(--primary-color);
  color: white;
  border-bottom-right-radius: 4px;
}

.ai-message .message-bubble {
  background: white;
  color: var(--text-color);
  border-bottom-left-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* 输入区域样式 */
.chat-input-container {
  padding: 1rem 1.5rem;
  background: white;
  border-top: 1px solid #e9ecef;
  display: flex;
  gap: 0.8rem;
  align-items: center;
}

.chat-input-container input {
  flex: 1;
  padding: 0.8rem 1.2rem;
  border: 1px solid #e9ecef;
  border-radius: 25px;
  font-size: 0.95rem;
  transition: all 0.3s ease;
}

.chat-input-container input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(200, 16, 46, 0.1);
}

.send-btn {
  background: var(--primary-color);
  color: white;
  border: none;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.send-btn:hover {
  background: #a60e24;
  transform: scale(1.05);
}

.send-btn:active {
  transform: scale(0.95);
}

/* 滚动条样式 */
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: #a1a1a1;
}

/* 动画效果 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .section-title h2 {
    font-size: 2rem;
  }
  
  .ai-assistant-container {
    flex-direction: column;
    margin: 0 1rem;
  }
  
  .ai-menu {
    width: 100%;
    max-height: 200px;
    border-right: none;
    border-bottom: 1px solid #e9ecef;
  }
  
  .menu-items {
    display: flex;
    flex-wrap: wrap;
    overflow-y: auto;
  }
  
  .menu-item {
    flex: 1 1 calc(50% - 2px);
    border-left: none;
    border-bottom: 1px solid #e9ecef;
  }
  
  .menu-item.active {
    border-left: none;
    border-bottom: 3px solid var(--primary-color);
  }
  
  .message {
    max-width: 90%;
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
  
  .chat-messages {
    padding: 1rem;
  }
  
  .chat-input-container {
    padding: 0.8rem 1rem;
  }
  
  .menu-item {
    flex: 1 1 100%;
  }
  
  .ai-assistant-container {
    min-height: 600px;
  }
}
</style>