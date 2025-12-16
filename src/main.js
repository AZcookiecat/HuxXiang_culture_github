import { createApp } from 'vue'
import './assets/css/style.css'
import '@fortawesome/fontawesome-free/css/all.css'
import App from './App.vue'
import router from './router/index.js'

// 创建Vue应用实例
const app = createApp(App)

// 使用路由
app.use(router)

// 全局事件总线 - 用于组件间通信
app.config.globalProperties.$eventBus = app

// 挂载应用
app.mount('#app')
