<template>
    <div class="app">
        <!-- Header -->
        <Header />

        <!-- Stats -->
        <Stats :metrics="metrics" />

        <!-- Filters -->
        <div class="filters">
            <h3>Filters</h3>

            <!-- Search -->
            <div class="filter-row">
                <div class="filter-group">
                    <label for="query">Search</label>
                    <input id="query" v-model="filters.query" type="text" placeholder="Search by model..."
                        @input="debouncedSearch" />
                </div>
            </div>

        <!-- Custom Storage Slider -->
        <div class="filter-row">
            <div class="filter-group storage-slider">
                <label class="slider-label">
                    Storage Capacity
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
                    <label class="ram-label">RAM Memory</label>
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
                First
            </button>
            
            <button @click="goToPage(pagination.page - 1)" :disabled="pagination.page <= 1" class="pagination-btn">
                Previous
            </button>

            <button v-for="page in visiblePages" :key="page" @click="goToPage(page)"
                :class="{ active: page === pagination.page }" class="pagination-btn">
                {{ page }}
            </button>

            <button @click="goToPage(pagination.page + 1)" :disabled="pagination.page >= pagination.total_pages" class="pagination-btn">
                Next
            </button>
            
            <button @click="goToLastPage" :disabled="pagination.page >= pagination.total_pages" class="pagination-btn">
                Last
            </button>
        </div>
    </div>
</template>
  
<script>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { serverService } from './services/api.js'
import Header from './components/Header.vue'
import Stats from './components/Stats.vue'
import { 
    STORAGE_OPTIONS, 
    RAM_OPTIONS, 
    HDD_TYPES, 
    PER_PAGE_OPTIONS, 
    DEFAULT_PAGE, 
    DEFAULT_PER_PAGE, 
    DEBOUNCE_DELAY,
    STORAGE_UNITS
} from './constants/index.js'

export default {
    name: 'App',
    components: {
        Header,
        Stats
    },
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
            page: DEFAULT_PAGE,
            perPage: DEFAULT_PER_PAGE
        })

        // Filter options (imported from constants)
        const storageOptions = STORAGE_OPTIONS
        const ramOptions = RAM_OPTIONS
        const hddTypes = HDD_TYPES

        // Storage slider
        const storageSliderIndex = ref(0)

        // Debounced search
        let searchTimeout = null
        const debouncedSearch = () => {
            clearTimeout(searchTimeout)
            searchTimeout = setTimeout(() => {
                filters.page = DEFAULT_PAGE
                applyFilters()
            }, DEBOUNCE_DELAY)
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
        return `${(tb * 1000).toFixed(0)}${STORAGE_UNITS.GB}`
    }
    return `${tb}${STORAGE_UNITS.TB}`
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
@import './assets/app.css';
</style>
  