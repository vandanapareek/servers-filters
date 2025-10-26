import apiClient from '../config/axios'

export const serverService = {
    // get servers
    getServers(params = {}) {
        return apiClient.get('/servers', { params })
    },

    // get locations
    getLocations() {
        return apiClient.get('/locations')
    },

    // get statistics
    getMetrics() {
        return apiClient.get('/metrics')
    },
}

export default apiClient
