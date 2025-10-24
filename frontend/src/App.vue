<template>
    <div class="app">
        <!-- Header -->
        <header class="header">
            <div class="header-content">
                <div class="header-text">
                    <h1>Server Marketplace</h1>
                </div>
            </div>
        </header>

        <!-- Stats -->
        <div class="stats" v-if="metrics">
            <div class="stat-card">
                <div class="stat-value">{{ metrics.total_servers }}</div>
                <div class="stat-label">Total Servers</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">€{{ metrics.min_price?.toFixed(2) || '0.00' }} - €{{ metrics.max_price?.toFixed(2) || '0.00' }}</div>
                <div class="stat-label">Price Range</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{ metrics.locations_count }}</div>
                <div class="stat-label">Locations</div>
            </div>
        </div>

        <!-- Filters -->
        <div class="filters">
            <h3>Filters</h3>

            <!-- Search -->
            <div class="filter-row">
                <div class="filter-group">
                    <label for="query">Search</label>
                    <input id="query" v-model="filters.query" type="text" placeholder="Search by model or CPU..."
                        @input="debouncedSearch" />
                </div>
            </div>

        <!-- Custom Storage Slider -->
        <div class="filter-row">
            <div class="filter-group storage-slider">
                <label class="slider-label">
                    💾 Storage Capacity
                    <span class="current-value">{{ formatStorageValue(storageOptions[storageSliderIndex]) }}</span>
                </label>
                <div class="custom-slider">
                    <div class="slider-track" @click="handleSliderClick" ref="sliderTrack">
                        <div class="slider-progress" :style="{ width: (storageSliderIndex / (storageOptions.length - 1)) * 100 + '%' }"></div>
                        <div class="slider-thumb" 
                             :style="{ left: (storageSliderIndex / (storageOptions.length - 1)) * 100 + '%' }"
                             @mousedown="startDrag"
                             @touchstart="startDrag">
                        </div>
                    </div>
                    <div class="slider-labels">
                        <div v-for="(option, index) in storageOptions" :key="index" 
                             class="slider-label-item"
                             :class="{ active: index <= storageSliderIndex }"
                             :style="{ left: (index / (storageOptions.length - 1)) * 100 + '%' }"
                             @click="setSliderValue(index)">
                            <span class="label-text">{{ formatStorageValue(option) }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

            <!-- RAM Checkboxes -->
            <div class="filter-row">
                <div class="filter-group ram-checkboxes">
                    <label class="ram-label">💾 RAM Memory</label>
                    <div class="ram-grid">
                        <div v-for="ram in ramOptions" :key="ram" class="ram-option" 
                             :class="{ active: filters.ramValues.includes(ram) }"
                             @click="toggleRam(ram, $event)">
                            <div class="ram-checkbox">
                                <input type="checkbox" :id="`ram-${ram}`" :checked="filters.ramValues.includes(ram)" />
                                <div class="checkmark"></div>
                            </div>
                            <label :for="`ram-${ram}`" class="ram-text">
                                {{ ram }}
                            </label>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Harddisk Type Dropdown -->
            <div class="filter-row">
                <div class="filter-group">
                    <label for="hddType">Harddisk Type</label>
                    <select id="hddType" v-model="filters.hddType" @change="applyFilters">
                        <option value="">All Types</option>
                        <option v-for="type in hddTypes" :key="type" :value="type">{{ type }}</option>
                    </select>
                </div>

                <!-- Location Dropdown -->
                <div class="filter-group">
                    <label for="location">Location</label>
                    <select id="location" v-model="filters.location" @change="applyFilters">
                        <option value="">All Locations</option>
                        <option v-for="loc in locations" :key="loc" :value="loc">{{ loc }}</option>
                    </select>
                </div>
            </div>

            <!-- Sort and Per Page -->
            <div class="filter-row">
                <div class="filter-group">
                    <label for="sort">Sort By</label>
                    <select id="sort" v-model="filters.sort" @change="applyFilters">
                        <option value="id.asc">ID (Ascending)</option>
                        <option value="id.desc">ID (Descending)</option>
                        <option value="price_eur.asc">Price (Low to High)</option>
                        <option value="price_eur.desc">Price (High to Low)</option>
                        <option value="ram_gb.asc">RAM (Low to High)</option>
                        <option value="ram_gb.desc">RAM (High to Low)</option>
                        <option value="storage_gb.asc">Storage (Low to High)</option>
                        <option value="storage_gb.desc">Storage (High to Low)</option>
                        <option value="location_city.asc">Location (A-Z)</option>
                        <option value="location_city.desc">Location (Z-A)</option>
                    </select>
                </div>

                <div class="filter-group">
                    <label for="perPage">Per Page</label>
                    <select id="perPage" v-model="filters.perPage" @change="handlePerPageChange">
                        <option value="10">10</option>
                        <option value="20">20</option>
                        <option value="50">50</option>
                        <option value="100">100</option>
                    </select>
                </div>
            </div>

            <div class="filter-row">
                <button class="btn btn-secondary" @click="clearFilters">Clear Filters</button>
            </div>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="error">
            {{ error }}
        </div>

        <!-- Servers Table -->
        <div class="servers-table">
            <div v-if="loading" class="loading">
                Loading servers...
            </div>

            <div v-else-if="servers.length === 0" class="loading">
                No servers found matching your criteria.
            </div>

            <div v-else class="table-container">
                <table class="table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Model</th>
                        <th>CPU</th>
                        <th>RAM</th>
                        <th>Storage</th>
                        <th>HDD Type</th>
                        <th>Location</th>
                        <th>Price</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="server in servers" :key="server.id">
                        <td>{{ server.id }}</td>
                        <td>{{ server.model }}</td>
                        <td>{{ server.cpu || '-' }}</td>
                        <td>{{ server.ram_gb ? `${server.ram_gb}GB` : '-' }}</td>
                        <td>
                            <div class="storage-capacity">{{ server.storage_display || formatStorageFromGB(server.storage_gb) }}</div>
                        </td>
                        <td>
                            <div class="hdd-type">{{ server.hdd_type || '-' }}</div>
                        </td>
                        <td>
                            <div>{{ server.location_city || '-' }}</div>
                        </td>
                        <td>
                            <div class="price">€{{ server.price_eur || '-' }}</div>
                        </td>
                    </tr>
                </tbody>
                </table>
            </div>
        </div>

        <!-- Enhanced Server-Side Pagination -->
        <div v-if="pagination && pagination.total_pages > 1" class="pagination">
            <button @click="goToFirstPage" :disabled="pagination.page <= 1" class="pagination-btn">
                ⏮️ First
            </button>
            
            <button @click="goToPage(pagination.page - 1)" :disabled="pagination.page <= 1" class="pagination-btn">
                ⏪ Previous
            </button>

            <button v-for="page in visiblePages" :key="page" @click="goToPage(page)"
                :class="{ active: page === pagination.page }" class="pagination-btn">
                {{ page }}
            </button>

            <button @click="goToPage(pagination.page + 1)" :disabled="pagination.page >= pagination.total_pages" class="pagination-btn">
                Next ⏩
            </button>
            
            <button @click="goToLastPage" :disabled="pagination.page >= pagination.total_pages" class="pagination-btn">
                Last ⏭️
            </button>
        </div>
    </div>
</template>
  
<script>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { serverService } from './services/api.js'

export default {
    name: 'App',
    setup() {
        const servers = ref([])
        const locations = ref([])
        const metrics = ref(null)
        const loading = ref(false)
        const error = ref('')
        const pagination = ref(null)


        const filters = reactive({
            query: '',
            location: '',
            hddType: '',
            ramValues: [],
            storageMin: null,
            storageMax: null,
            sort: 'id.asc',
            page: 1,
            perPage: 20
        })

        // Filter options
        const storageOptions = [0, 0.25, 0.5, 1, 2, 3, 4, 8, 12, 24, 48, 72] // in TB
        const ramOptions = ['2GB', '4GB', '8GB', '12GB', '16GB', '24GB', '32GB', '48GB', '64GB', '96GB']
        const hddTypes = ['SAS', 'SATA', 'SSD']

        // Storage slider
        const storageSliderIndex = ref(0)

        // Debounced search
        let searchTimeout = null
        const debouncedSearch = () => {
            clearTimeout(searchTimeout)
            searchTimeout = setTimeout(() => {
                filters.page = 1
                applyFilters()
            }, 500)
        }

        // Computed properties
        const visiblePages = computed(() => {
            if (!pagination.value) return []

            const current = pagination.value.page
            const total = pagination.value.total_pages
            const pages = []

            // Show up to 5 pages around current page
            const start = Math.max(1, current - 2)
            const end = Math.min(total, current + 2)

            for (let i = start; i <= end; i++) {
                pages.push(i)
            }

            return pages
        })


        // Methods
        const formatStorage = (gb) => {
            if (gb >= 1024) {
                return `${(gb / 1024).toFixed(1)}TB`
            }
            return `${gb}GB`
        }

        const buildParams = () => {
            const params = {}

            if (filters.query) params.q = filters.query
            if (filters.location) params.location = filters.location
            if (filters.hddType) params.hdd = filters.hddType
            if (filters.ramValues.length > 0) {
                // Convert RAM values to GB for API - send as comma-separated list
                const ramGBs = filters.ramValues.map(ram => {
                    const match = ram.match(/(\d+)GB/)
                    return match ? parseInt(match[1]) : 0
                }).filter(gb => gb > 0)

                if (ramGBs.length > 0) {
                    // Send as comma-separated values for exact matching
                    params.ram_values = ramGBs.join(',')
                }
            }
            if (filters.storageMin !== null) params.storage_min = filters.storageMin
            if (filters.storageMax !== null) params.storage_max = filters.storageMax
            if (filters.sort) params.sort = filters.sort
            if (filters.page) params.page = filters.page
            if (filters.perPage) params.per_page = filters.perPage

            return params
        }

        const updateStorageFromSlider = () => {
            const selectedStorage = storageOptions[storageSliderIndex.value]
            filters.storageMax = selectedStorage
            filters.storageMin = null // Allow all storage up to selected value
            applyFilters()
        }

        const setSliderValue = (index) => {
            // Only update if the value actually changed
            if (storageSliderIndex.value !== index) {
                storageSliderIndex.value = index
                updateStorageFromSlider()
            }
        }

        const handleSliderClick = (event) => {
            const rect = event.currentTarget.getBoundingClientRect()
            const x = event.clientX - rect.left
            const percentage = x / rect.width
            const index = Math.round(percentage * (storageOptions.length - 1))
            setSliderValue(Math.max(0, Math.min(index, storageOptions.length - 1)))
        }

        const startDrag = (event) => {
            event.preventDefault()
            const handleMouseMove = (e) => {
                const rect = event.currentTarget.getBoundingClientRect()
                const x = e.clientX - rect.left
                const percentage = x / rect.width
                const index = Math.round(percentage * (storageOptions.length - 1))
                setSliderValue(Math.max(0, Math.min(index, storageOptions.length - 1)))
            }
            const handleMouseUp = () => {
                document.removeEventListener('mousemove', handleMouseMove)
                document.removeEventListener('mouseup', handleMouseUp)
            }
            document.addEventListener('mousemove', handleMouseMove)
            document.addEventListener('mouseup', handleMouseUp)
        }

        const toggleRam = (ram, event) => {
            event.preventDefault()
            event.stopPropagation()
            
            const index = filters.ramValues.indexOf(ram)
            if (index > -1) {
                // Remove from array
                filters.ramValues = filters.ramValues.filter(item => item !== ram)
            } else {
                // Add to array
                filters.ramValues = [...filters.ramValues, ram]
            }
            applyFilters()
        }

const formatStorageValue = (tb) => {
    if (tb === 0) return '0'
    if (tb < 1) {
        return `${(tb * 1000).toFixed(0)}GB`
    }
    return `${tb}TB`
}

const formatStorageFromGB = (gb) => {
    if (!gb) return '-'
    if (gb >= 1024) {
        const tb = gb / 1024
        if (tb === Math.floor(tb)) {
            return `${tb}TB`
        } else {
            return `${tb.toFixed(1)}TB`
        }
    }
    return `${gb}GB`
}


        const loadServers = async () => {
            loading.value = true
            error.value = ''

            try {
                const params = buildParams()
                console.log('Loading servers with params:', params)
                const response = await serverService.getServers(params)
                console.log('API response:', response.data)

                servers.value = response.data.data
                pagination.value = response.data.pagination
                console.log('Servers loaded:', servers.value.length, 'Pagination:', pagination.value)
            } catch (err) {
                error.value = err.response?.data?.message || 'Failed to load servers'
                console.error('Error loading servers:', err)
            } finally {
                loading.value = false
            }
        }

        const loadLocations = async () => {
            try {
                const response = await serverService.getLocations()
                locations.value = response.data.data
            } catch (err) {
                console.error('Error loading locations:', err)
            }
        }

        const loadMetrics = async () => {
            try {
                const response = await serverService.getMetrics()
                metrics.value = response.data
            } catch (err) {
                console.error('Error loading metrics:', err)
            }
        }

        const applyFilters = () => {
            loadServers()
        }

        const clearFilters = () => {
            Object.assign(filters, {
                query: '',
                location: '',
                hddType: '',
                ramValues: [],
                storageMin: null,
                storageMax: null,
                sort: 'id.asc',
                page: 1,
                perPage: 20
            })
            storageSliderIndex.value = 0
            applyFilters()
        }


        const goToPage = (page) => {
            if (page >= 1 && page <= pagination.value?.total_pages) {
                filters.page = page
                applyFilters()
            }
        }

        const goToFirstPage = () => {
            goToPage(1)
        }

        const goToLastPage = () => {
            goToPage(pagination.value?.total_pages || 1)
        }

        const handlePerPageChange = () => {
            // Reset to page 1 when changing per page limit
            filters.page = 1
            applyFilters()
        }


        // Lifecycle
        onMounted(() => {
            loadServers()
            loadLocations()
            loadMetrics()
        })

        return {
            servers,
            locations,
            metrics,
            loading,
            error,
            pagination,
            filters,
            visiblePages,
            storageOptions,
            ramOptions,
            hddTypes,
            storageSliderIndex,
            debouncedSearch,
            formatStorage,
            formatStorageValue,
            formatStorageFromGB,
            updateStorageFromSlider,
            setSliderValue,
            handleSliderClick,
            startDrag,
            toggleRam,
            applyFilters,
            clearFilters,
            goToPage,
            goToFirstPage,
            goToLastPage,
            handlePerPageChange,
            buildParams
        }
    }
}
</script>
  
<style scoped>
/* Global App Styles */
.app {
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    padding: 1rem;
}

/* Header Styles */
.header {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    padding: 2rem 0;
    margin: 0 0 2rem 0;
    border-radius: 20px;
}

.header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
    display: flex;
    justify-content: center;
    align-items: center;
}

.header-text h1 {
    font-size: 2.5rem;
    font-weight: 700;
    background: linear-gradient(135deg, #667eea, #764ba2);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    margin: 0 0 0.5rem 0;
}

.header-text p {
    color: #666;
    font-size: 1.1rem;
    margin: 0;
}

.header-stats {
    display: flex;
    gap: 2rem;
}

.quick-stat {
    text-align: center;
}

.stat-number {
    display: block;
    font-size: 2rem;
    font-weight: 700;
    color: #667eea;
}

.stat-label {
    font-size: 0.9rem;
    color: #666;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

/* Container */
.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
}


/* Stats Cards */
.stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
}

.stat-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    padding: 2rem;
    text-align: center;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    transition: transform 0.3s ease;
}

.stat-card:hover {
    transform: translateY(-4px);
}

.stat-value {
    font-size: 2.5rem;
    font-weight: 700;
    color: #667eea;
    margin-bottom: 0.5rem;
}

.stat-label {
    color: #666;
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

/* Filters */
.filters {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    padding: 2rem;
    margin-bottom: 2rem;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
}

.filters h3 {
    margin: 0 0 1.5rem 0;
    color: #333;
    font-size: 1.5rem;
    font-weight: 600;
}

.filter-row {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
    flex-wrap: wrap;
}

.filter-group {
    flex: 1;
    min-width: 200px;
}

.filter-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
    color: #333;
    font-size: 0.9rem;
}

.filter-group input,
.filter-group select {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 2px solid #e1e5e9;
    border-radius: 12px;
    font-size: 1rem;
    transition: all 0.3s ease;
    background: rgba(255, 255, 255, 0.8);
}

.filter-group select {
    padding-right: 2.5rem;
    background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='m6 8 4 4 4-4'/%3e%3c/svg%3e");
    background-position: right 0.75rem center;
    background-repeat: no-repeat;
    background-size: 1.5em 1.5em;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
}

.filter-group input:focus,
.filter-group select:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    background: white;
}

/* Custom Storage Slider */
.storage-slider {
    flex: 2;
}

.slider-label {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    font-weight: 600;
    color: #333;
    font-size: 0.9rem;
}

.current-value {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-weight: 600;
    font-size: 0.9rem;
}

.custom-slider {
    position: relative;
    padding: 1rem 0;
}

.slider-track {
    position: relative;
    height: 6px;
    background: #e1e5e9;
    border-radius: 3px;
    cursor: pointer;
    margin-bottom: 1rem;
}

.slider-progress {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    background: linear-gradient(135deg, #667eea, #764ba2);
    border-radius: 3px;
    transition: width 0.2s ease;
}

.slider-thumb {
    position: absolute;
    top: 50%;
    width: 20px;
    height: 20px;
    background: white;
    border: 3px solid #667eea;
    border-radius: 50%;
    cursor: pointer;
    transform: translate(-50%, -50%);
    box-shadow: 0 2px 6px rgba(102, 126, 234, 0.3);
    transition: all 0.2s ease;
    z-index: 2;
}

.slider-thumb:hover {
    transform: translate(-50%, -50%) scale(1.1);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.slider-labels {
    position: relative;
    height: 30px;
}

.slider-label-item {
    position: absolute;
    top: 0;
    transform: translateX(-50%);
    cursor: pointer;
    display: flex;
    flex-direction: column;
    align-items: center;
}

.slider-label-item:first-child {
    left: 0 !important;
    transform: translateX(0);
}

.slider-label-item:last-child {
    left: 100% !important;
    transform: translateX(-100%);
}

.label-text {
    font-size: 0.65rem;
    color: #667eea;
    font-weight: 600;
    white-space: nowrap;
    background: rgba(102, 126, 234, 0.1);
    padding: 0.15rem 0.3rem;
    border-radius: 3px;
    transition: all 0.2s ease;
}

.slider-label-item.active .label-text {
    background: rgba(102, 126, 234, 0.2);
    font-weight: 700;
}


.text-muted {
    color: #666;
    font-size: 0.9em;
}

/* RAM Checkboxes Styles */
.ram-checkboxes {
    flex: 2;
}

.ram-label {
    font-weight: 600;
    color: #333;
    font-size: 0.9rem;
    margin-bottom: 1rem;
    display: block;
}

.ram-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    gap: 0.75rem;
    margin-top: 0.5rem;
}

.ram-option {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 1rem;
    border: 2px solid #e1e5e9;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    background: rgba(255, 255, 255, 0.8);
    position: relative;
    overflow: hidden;
    min-height: 60px;
}

.ram-option:hover {
    border-color: #667eea;
    background: white;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

.ram-option.active {
    border-color: #667eea;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.ram-checkbox {
    position: relative;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.ram-checkbox input[type="checkbox"] {
    position: absolute;
    opacity: 0;
    cursor: pointer;
    width: 20px;
    height: 20px;
}

.checkmark {
    position: absolute;
    top: 0;
    left: 0;
    width: 20px;
    height: 20px;
    background: white;
    border: 2px solid #e1e5e9;
    border-radius: 4px;
    transition: all 0.3s ease;
}

.ram-checkbox input:checked + .checkmark {
    background: linear-gradient(135deg, #667eea, #764ba2);
    border-color: #667eea;
}

.ram-checkbox input:checked + .checkmark::after {
    content: '';
    position: absolute;
    left: 6px;
    top: 3px;
    width: 4px;
    height: 8px;
    border: solid white;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg);
}

.ram-text {
    font-weight: 600;
    color: #333;
    cursor: pointer;
    user-select: none;
    transition: color 0.3s ease;
    display: flex;
    align-items: center;
    height: 10px;
    line-height: 15px;
}

.ram-option.active .ram-text {
    color: #667eea;
    font-weight: 700;
}

/* Buttons */
.btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 12px;
    cursor: pointer;
    font-size: 1rem;
    font-weight: 600;
    transition: all 0.3s ease;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.btn-primary {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
    background: rgba(255, 255, 255, 0.8);
    color: #666;
    border: 2px solid #e1e5e9;
}

.btn-secondary:hover {
    background: white;
    border-color: #667eea;
    color: #667eea;
}

/* Table Styles */
.servers-table {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    margin-bottom: 2rem;
}

.table {
    width: 100%;
    border-collapse: collapse;
}

.table th {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    padding: 1.5rem 1rem;
    text-align: left;
    font-weight: 600;
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.table td {
    padding: 1.5rem 1rem;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
    transition: background-color 0.3s ease;
}

.table tbody tr:hover {
    background: rgba(102, 126, 234, 0.05);
}

.table tbody tr:last-child td {
    border-bottom: none;
}

/* Loading and Error States */
.loading {
    text-align: center;
    padding: 4rem 2rem;
    color: #666;
    font-size: 1.1rem;
}

.error {
    background: linear-gradient(135deg, #ff6b6b, #ee5a52);
    color: white;
    padding: 1rem 1.5rem;
    border-radius: 12px;
    margin-bottom: 1.5rem;
    font-weight: 500;
}

/* Pagination */
.pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0.5rem;
    padding: 2rem;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
}

.pagination button {
    padding: 0.75rem 1rem;
    border: 2px solid #e1e5e9;
    background: white;
    cursor: pointer;
    border-radius: 12px;
    font-weight: 600;
    transition: all 0.3s ease;
    min-width: 44px;
}

.pagination button:hover:not(:disabled) {
    background: #667eea;
    color: white;
    border-color: #667eea;
    transform: translateY(-2px);
}

.pagination button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
}

.pagination .active {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    border-color: #667eea;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

/* Responsive Design */
@media (max-width: 768px) {
    .header-content {
        flex-direction: column;
        gap: 1rem;
        text-align: center;
    }

    .header-text h1 {
        font-size: 2rem;
    }

    .filter-row {
        flex-direction: column;
    }

    .filter-group {
        min-width: auto;
    }

    .checkbox-grid {
        grid-template-columns: repeat(auto-fit, minmax(80px, 1fr));
    }

    .storage-slider {
        flex: 1;
    }

    .ram-checkboxes {
        flex: 1;
    }

    /* Responsive Slider */
    .slider-container {
        padding: 0 20px;
    }

    .slider-mark {
        width: 30px;
        margin-left: -15px;
    }

    .mark-label {
        font-size: 0.75rem;
    }

    .table {
        font-size: 0.9rem;
    }

    .table th,
    .table td {
        padding: 1rem 0.5rem;
    }

    .pagination {
        flex-wrap: wrap;
        gap: 0.25rem;
    }

    .pagination button {
        padding: 0.5rem 0.75rem;
        min-width: 40px;
    }
}

@media (max-width: 480px) {
    .container {
        padding: 0 1rem;
    }

    .filters {
        padding: 1.5rem;
    }

    .header {
        padding: 1.5rem 0;
    }

    .header-text h1 {
        font-size: 1.75rem;
    }

    .stat-card {
        padding: 1.5rem;
    }

    .table th,
    .table td {
        padding: 0.75rem 0.25rem;
        font-size: 0.8rem;
    }

    /* Mobile Slider */
    .slider-container {
        padding: 0 15px;
    }

    .slider-mark {
        width: 25px;
        margin-left: -12.5px;
    }

    .mark-label {
        font-size: 0.7rem;
    }
}

/* Data Table Styles */
.data-table-container {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    padding: 2rem;
    margin: 2rem 0;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.table-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    flex-wrap: wrap;
    gap: 1rem;
}

.search-box {
    position: relative;
    flex: 1;
    max-width: 400px;
}

.search-input {
    width: 100%;
    padding: 0.75rem 3rem 0.75rem 1rem;
    border: 2px solid #e1e5e9;
    border-radius: 12px;
    font-size: 1rem;
    background: white;
    transition: all 0.3s ease;
}

.search-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.search-icon {
    position: absolute;
    right: 1rem;
    top: 50%;
    transform: translateY(-50%);
    font-size: 1.2rem;
    color: #666;
}

.table-info {
    color: #666;
    font-size: 0.9rem;
    font-weight: 500;
}

.data-table {
    overflow-x: auto;
    border-radius: 12px;
    border: 1px solid #e1e5e9;
}

.sortable {
    cursor: pointer;
    user-select: none;
    position: relative;
    transition: background-color 0.2s ease;
}

.sortable:hover {
    background-color: #f8f9fa;
}

.sort-icon {
    margin-left: 0.5rem;
    font-size: 0.9rem;
    opacity: 0.7;
}

.pagination-btn {
    padding: 0.5rem 1rem;
    border: 2px solid #e1e5e9;
    background: white;
    color: #333;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 500;
    transition: all 0.3s ease;
    min-width: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.pagination-btn:hover:not(:disabled) {
    background: #667eea;
    color: white;
    border-color: #667eea;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.pagination-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.pagination-btn.active {
    background: #667eea;
    color: white;
    border-color: #667eea;
}

/* Table Container for Mobile Responsiveness */
.table-container {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    border-radius: 12px;
    border: 1px solid #e1e5e9;
    background: white;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.storage-capacity {
    font-weight: 600;
    color: #2d3748;
}

.hdd-type {
    font-weight: 500;
    color: #4a5568;
    background: #f7fafc;
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
    display: inline-block;
    font-size: 0.875rem;
}

.price {
    font-weight: 600;
    color: #2d3748;
    font-size: 1.1rem;
}

/* Mobile Table Responsiveness */
@media (max-width: 768px) {
    .table-container {
        margin: 0 -1rem;
        border-radius: 0;
        border-left: none;
        border-right: none;
    }
    
    .table {
        min-width: 700px; /* Ensure table doesn't get too cramped with new HDD column */
    }
    
    .table th,
    .table td {
        padding: 0.5rem 0.75rem;
        font-size: 0.875rem;
    }
    
    .table th:first-child,
    .table td:first-child {
        position: sticky;
        left: 0;
        background: white;
        z-index: 10;
        border-right: 1px solid #e1e5e9;
    }
}

@media (max-width: 480px) {
    .table {
        min-width: 600px;
    }
    
    .table th,
    .table td {
        padding: 0.375rem 0.5rem;
        font-size: 0.8rem;
    }
}
</style>
  