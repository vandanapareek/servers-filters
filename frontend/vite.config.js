import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
    // Load environment variables
    const env = loadEnv(mode, process.cwd(), '')

    return {
        plugins: [vue()],
        server: {
            port: 3000,
            proxy: {
                '/api': {
                    target: env.VITE_API_BASE_URL || 'http://localhost:8081',
                    changeOrigin: true,
                    rewrite: (path) => path.replace(/^\/api/, '')
                }
            }
        }
    }
})
