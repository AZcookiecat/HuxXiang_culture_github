<template>
  <div class="unity-webgl-container">
    <div class="unity-content">
      <!-- Unity WebGL 内容将通过iframe加载 -->
      <iframe 
        :src="unityWebGLUrl" 
        frameborder="0" 
        class="unity-iframe"
        allow="accelerometer; autoplay; gyroscope; xr-spatial-tracking"
        xr-spatial-tracking="inline"
        execution-while-out-of-viewport="true"
        execution-while-not-rendered="true"
        loading="lazy"
      ></iframe>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Unity WebGL文件的URL路径
const unityWebGLUrl = ref('')

// 返回上一页
const goBack = () => {
  router.back()
}

// 组件挂载时设置Unity WebGL URL
onMounted(() => {
  // 在开发环境中，我们需要设置一个合适的URL来访问Unity WebGL内容
  // 由于Unity WebGL内容位于HX_project/static/webgl目录中
  // 我们需要通过合适的路径来访问它
  
  // 这里我们使用绝对路径，假设项目在本地运行
  // 在实际部署时，可能需要调整这个路径
  unityWebGLUrl.value = 'http://localhost:5173/unity-webgl/index.html'
})
</script>

<style scoped>
.unity-webgl-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: white;
}

.back-button {
  background-color: #C8102E; /* 湘红色 */
  color: white;
  border: none;
  padding: 0.7rem 1.5rem;
  border-radius: 4px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.back-button:hover {
  background-color: #a60e24;
}

.unity-content {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  position: relative;
}

.unity-iframe {
  width: 100%;
  height: 100%;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .unity-header {
    padding: 1rem;
  }
  
  .unity-header h2 {
    font-size: 1.5rem;
  }
  
  .back-button {
    padding: 0.6rem 1.2rem;
    font-size: 0.9rem;
  }
}

@media (max-width: 480px) {
  .unity-header h2 {
    font-size: 1.3rem;
  }
  
  .back-button {
    padding: 0.5rem 1rem;
    font-size: 0.85rem;
  }
}
</style>