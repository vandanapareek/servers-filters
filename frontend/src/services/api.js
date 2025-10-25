import axios from 'axios'

// Get API base URL from environment variables with fallback
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081'

const api = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
})

// Request interceptor
api.interceptors.request.use(
    (config) => {
        console.log(`Making ${config.method?.toUpperCase()} request to ${config.url}`)
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response interceptor
api.interceptors.response.use(
    (response) => {
        return response
    },
    (error) => {
        console.error('API Error:', error.response?.data || error.message)
        return Promise.reject(error)
    }
)

export const serverService = {
    // Get servers with filters
    getServers(params = {}) {
        return api.get('/servers', { params })
    },


    // Get locations
    getLocations() {
        return api.get('/locations')
    },

    // Get metrics
    getMetrics() {
        return api.get('/metrics')
    },

}

export default api
