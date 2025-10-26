import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
    // Load environment var
    const env = loadEnv(mode, process.cwd(), '')

    return {
        plugins: [vue()],
        server: {
            port: 3000
        }
    }
})
