import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, __dirname, '')
  const flaskApiTarget = env.VITE_FLASK_API_URL || 'http://127.0.0.1:5000'
  const communityApiTarget = env.VITE_COMMUNITY_API_URL || 'http://127.0.0.1:8080'

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src')
      }
    },
    server: {
      proxy: {
        '/api/community': {
          target: communityApiTarget,
          changeOrigin: true,
          secure: false
        },
        '/api': {
          target: flaskApiTarget,
          changeOrigin: true,
          secure: false
        }
      }
    }
  }
})
