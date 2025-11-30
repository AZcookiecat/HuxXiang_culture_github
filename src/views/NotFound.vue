<template>
  <div class="not-found">
    <div class="container">
      <div class="not-found-content">
        <div class="not-found-icon">
          <i class="fas fa-map-marker-alt"></i>
        </div>
        <h1>404</h1>
        <h2>页面未找到</h2>
        <p>抱歉，您访问的页面不存在或已被移除。</p>
        <button class="btn btn-primary" @click="goHome">返回首页</button>
      </div>
    </div>
  </div>
</template>

<script>
import { onMounted } from 'vue'

export default {
  name: 'NotFound',
  props: {
    showAlert: {
      type: Function,
      required: true
    }
  },
  setup(props) {
    // 返回首页
    const goHome = () => {
      // 通过事件向上传递导航请求
      window.dispatchEvent(new CustomEvent('navigate', { detail: { page: 'home' } }))
    }

    // 组件挂载时的逻辑
    onMounted(() => {
      // 显示错误提示
      props.showAlert('您访问的页面不存在', 'error')
    })

    return {
      goHome
    }
  }
}
</script>

<style scoped>
/* 404页面样式 */
.not-found {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding-top: 100px;
  background-color: var(--background-color);
}

.not-found-content {
  text-align: center;
  padding: 40px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  max-width: 500px;
  width: 100%;
}

.not-found-icon {
  font-size: 5rem;
  color: var(--danger-color);
  margin-bottom: 20px;
  opacity: 0.8;
}

.not-found h1 {
  font-size: 4rem;
  font-weight: bold;
  color: var(--danger-color);
  margin-bottom: 10px;
  line-height: 1;
}

.not-found h2 {
  font-size: 1.8rem;
  color: var(--dark-color);
  margin-bottom: 15px;
}

.not-found p {
  font-size: 1.1rem;
  color: var(--text-color);
  margin-bottom: 30px;
  line-height: 1.5;
}

.not-found .btn {
  min-width: 150px;
  padding: 12px 24px;
  font-size: 1.1rem;
  font-weight: 500;
}

/* 动画效果 */
@keyframes pulse {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
  100% {
    transform: scale(1);
  }
}

.not-found-icon {
  animation: pulse 2s infinite;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .not-found {
    padding: 20px;
  }
  
  .not-found-content {
    padding: 30px 20px;
  }
  
  .not-found-icon {
    font-size: 4rem;
  }
  
  .not-found h1 {
    font-size: 3rem;
  }
  
  .not-found h2 {
    font-size: 1.5rem;
  }
  
  .not-found p {
    font-size: 1rem;
  }
}

@media (max-width: 480px) {
  .not-found-icon {
    font-size: 3rem;
  }
  
  .not-found h1 {
    font-size: 2.5rem;
  }
}
</style>