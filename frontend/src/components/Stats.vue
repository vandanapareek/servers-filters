<template>
    <div class="stats" v-if="metrics">
        <div class="stat-card">
            <div class="stat-value">{{ metrics.total_servers }}</div>
            <div class="stat-label">Total Servers</div>
        </div>
        <div class="stat-card">
            <div class="stat-value">{{ CURRENCY_SYMBOL }}{{ formatPrice(metrics.min_price) }} - {{ CURRENCY_SYMBOL }}{{ formatPrice(metrics.max_price) }}</div>
            <div class="stat-label">Price Range</div>
        </div>
        <div class="stat-card">
            <div class="stat-value">{{ metrics.locations_count }}</div>
            <div class="stat-label">Locations</div>
        </div>
    </div>
</template>

<script>
import { CURRENCY_SYMBOL, PRICE_DECIMAL_PLACES } from '../constants/index.js'

export default {
    name: 'Stats',
    props: {
        metrics: {
            type: Object,
            default: null
        }
    },
    setup() {
        const formatPrice = (price) => {
            return price?.toFixed(PRICE_DECIMAL_PLACES) || '0.00'
        }

        return {
            formatPrice,
            CURRENCY_SYMBOL
        }
    }
}
</script>

<style scoped>
@import '../assets/app.css';
</style>
