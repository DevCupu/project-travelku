<script setup>
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const handleLogout = () => {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<template>
  <div class="md:hidden sticky top-0 z-50 bg-white border-b border-gray-200 h-16 px-4 flex items-center justify-between shadow-sm">
    <!-- Brand (Left) -->
    <div class="flex items-center gap-1 cursor-pointer" @click="router.push('/dashboard')">
      <span class="text-xl font-bold tracking-tight text-primary">24visa</span>
      <span class="text-[10px] font-bold bg-accent-soft text-accent-blue px-2 py-0.5 rounded-full mt-1">Makassar</span>
    </div>

    <!-- User Menu (Right) -->
    <button @click="handleLogout" class="flex items-center justify-center p-2 rounded-full border border-gray-200 text-gray-500 hover:text-primary hover:bg-red-50 hover:border-red-200 transition-colors cursor-pointer bg-white">
      <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
        <polyline points="16 17 21 12 16 7"></polyline>
        <line x1="21" y1="12" x2="9" y2="12"></line>
      </svg>
    </button>
  </div>
</template>

<style scoped>
/* Scoped styles dihapus karena seluruhnya sudah dikelola oleh Tailwind CSS utility classes */
</style>
