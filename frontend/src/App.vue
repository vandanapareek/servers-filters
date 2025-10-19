<template>
    <div class="container">
        <!-- Header -->
        <div class="header">
            <h1>Servers Listing</h1>
            <p>Browse and filter server configurations</p>
        </div>

        <!-- Stats -->
        <div class="stats" v-if="metrics">
            <div class="stat-card">
                <div class="stat-value">{{ metrics.total_servers }}</div>
                <div class="stat-label">Total Servers</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">€{{ metrics.average_price?.toFixed(2) || '0.00' }}</div>
                <div class="stat-label">Average Price</div>
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

            <!-- Storage Range Slider -->
            <div class="filter-row">
                <div class="filter-group storage-slider">
                    <label for="storage">Storage</label>
                    <div class="slider-container">
                        <input id="storage" type="range" :min="0" :max="storageOptions.length - 1"
                            v-model="storageSliderIndex" @input="updateStorageFromSlider" class="slider" />
                        <div class="slider-labels">
                            <span>{{ formatStorageValue(storageOptions[storageSliderIndex]) }}</span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- RAM Checkboxes -->
            <div class="filter-row">
                <div class="filter-group ram-checkboxes">
                    <label>RAM</label>
                    <div class="checkbox-grid">
                        <label v-for="ram in ramOptions" :key="ram" class="checkbox-item">
                            <input type="checkbox" :value="ram" v-model="filters.ramValues" @change="applyFilters" />
                            <span>{{ ram }}</span>
                        </label>
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
                    <select id="perPage" v-model="filters.perPage" @change="applyFilters">
                        <option value="10">10</option>
                        <option value="20">20</option>
                        <option value="50">50</option>
                        <option value="100">100</option>
                    </select>
                </div>
            </div>

            <div class="filter-row">
                <button class="btn btn-primary" @click="applyFilters">Apply Filters</button>
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

            <table v-else class="table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Model</th>
                        <th>CPU</th>
                        <th>RAM</th>
                        <th>Storage</th>
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
                            <div>{{ server.hdd }}</div>
                            <small v-if="server.storage_gb" class="text-muted">
                                ({{ formatStorage(server.storage_gb) }})
                            </small>
                        </td>
                        <td>
                            <div>{{ server.location_city || '-' }}</div>
                            <small v-if="server.location_code" class="text-muted">
                                {{ server.location_code }}
                            </small>
                        </td>
                        <td>
                            <div>€{{ server.price_eur || '-' }}</div>
                            <small v-if="server.raw_price" class="text-muted">
                                {{ server.raw_price }}
                            </small>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Pagination -->
        <div v-if="pagination && pagination.total_pages > 1" class="pagination">
            <button @click="goToPage(pagination.page - 1)" :disabled="pagination.page <= 1">
                Previous
            </button>

            <button v-for="page in visiblePages" :key="page" @click="goToPage(page)"
                :class="{ active: page === pagination.page }">
                {{ page }}
            </button>

            <button @click="goToPage(pagination.page + 1)" :disabled="pagination.page >= pagination.total_pages">
                Next
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
        const storageOptions = [0, 250, 500, 1000, 2000, 3000, 4000, 8000, 12000, 24000, 48000, 72000] // in GB
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

        const formatStorageValue = (gb) => {
            if (gb === 0) return '0'
            if (gb >= 1000) {
                return `${(gb / 1000).toFixed(gb % 1000 === 0 ? 0 : 1)}TB`
            }
            return `${gb}GB`
        }

        const loadServers = async () => {
            loading.value = true
            error.value = ''

            try {
                const params = buildParams()
                const response = await serverService.getServers(params)

                servers.value = response.data.data
                pagination.value = response.data.pagination
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
            updateStorageFromSlider,
            applyFilters,
            clearFilters,
            goToPage
        }
    }
}
</script>
  
<style scoped>
.text-muted {
    color: #666;
    font-size: 0.9em;
}

/* Storage Slider Styles */
.storage-slider {
    flex: 2;
}

.slider-container {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.slider {
    width: 100%;
    height: 6px;
    border-radius: 3px;
    background: #ddd;
    outline: none;
    -webkit-appearance: none;
}

.slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #007bff;
    cursor: pointer;
}

.slider::-moz-range-thumb {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #007bff;
    cursor: pointer;
    border: none;
}

.slider-labels {
    display: flex;
    justify-content: center;
    font-weight: 500;
    color: #007bff;
}

/* RAM Checkboxes Styles */
.ram-checkboxes {
    flex: 2;
}

.checkbox-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(80px, 1fr));
    gap: 8px;
    margin-top: 8px;
}

.checkbox-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 14px;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
    transition: background-color 0.2s;
}

.checkbox-item:hover {
    background-color: #f8f9fa;
}

.checkbox-item input[type="checkbox"] {
    margin: 0;
    cursor: pointer;
}

.checkbox-item span {
    user-select: none;
}

/* Responsive adjustments */
@media (max-width: 768px) {
    .checkbox-grid {
        grid-template-columns: repeat(auto-fit, minmax(70px, 1fr));
    }

    .storage-slider {
        flex: 1;
    }

    .ram-checkboxes {
        flex: 1;
    }
}
</style>
  