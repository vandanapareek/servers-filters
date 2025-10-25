<template>
    <div class="stats" v-if="metrics">
        <div class="stat-card">
            <div class="stat-value">{{ metrics.total_servers }}</div>
            <div class="stat-label">Total Servers</div>
        </div>
        <div class="stat-card">
            <div class="stat-value">{{ metrics.ram_range || 'N/A' }}</div>
            <div class="stat-label">RAM Range</div>
        </div>
        <div class="stat-card">
            <div class="stat-value">{{ metrics.locations_count }}</div>
            <div class="stat-label">Locations</div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Stats',
    props: {
        metrics: {
            type: Object,
            default: null
        }
    }
}
</script>

<style scoped>
@import '../assets/app.css';
</style>
