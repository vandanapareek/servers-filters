<template>
    <div class="logo-container">
        <img src="../assets/leaseweb.png" alt="Leaseweb Logo" class="logo" />
    </div>
</template>

<script>
export default {
    name: 'Logo'
}
</script>

<style scoped>
.logo-container {
    position: absolute;
    top: 20px;
    left: 20px;
    z-index: 10;
}

.logo {
    height: 50px;
    width: auto;
    object-fit: contain;
    display: block;
}

/* Responsive design */
@media (max-width: 768px) {
    .logo-container {
        top: 15px;
        left: 15px;
    }
    
    .logo {
        height: 40px;
    }
}

@media (max-width: 480px) {
    .logo-container {
        top: 10px;
        left: 10px;
    }
    
    .logo {
        height: 35px;
    }
}
</style>
