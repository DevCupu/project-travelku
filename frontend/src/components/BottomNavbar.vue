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
  <nav class="md:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 flex justify-around items-center pt-3 pb-[calc(env(safe-area-inset-bottom)+12px)] z-50">
    <router-link to="/dashboard" class="flex flex-col items-center gap-1 text-[11px] font-semibold transition-colors flex-1" :class="route.path === '/dashboard' ? 'text-primary' : 'text-gray-500 hover:text-gray-900'">
      <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M3 3v18h18"></path>
        <path d="M18.7 8l-5.1 5.2-2.8-2.7L7 14.3"></path>
      </svg>
      <span>Dashboard</span>
    </router-link>
    <router-link to="/bookings" class="flex flex-col items-center gap-1 text-[11px] font-semibold transition-colors flex-1" :class="route.path === '/bookings' ? 'text-primary' : 'text-gray-500 hover:text-gray-900'">
      <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
        <line x1="16" y1="2" x2="16" y2="6"></line>
        <line x1="8" y1="2" x2="8" y2="6"></line>
        <line x1="3" y1="10" x2="21" y2="10"></line>
      </svg>
      <span>Bookings</span>
    </router-link>
    <router-link to="/pakets" class="flex flex-col items-center gap-1 text-[11px] font-semibold transition-colors flex-1" :class="route.path === '/pakets' ? 'text-primary' : 'text-gray-500 hover:text-gray-900'">
      <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
        <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
        <line x1="12" y1="22.08" x2="12" y2="12"></line>
      </svg>
      <span>Pakets</span>
    </router-link>
    <router-link to="/users" class="flex flex-col items-center gap-1 text-[11px] font-semibold transition-colors flex-1" :class="route.path === '/users' ? 'text-primary' : 'text-gray-500 hover:text-gray-900'">
      <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
        <circle cx="12" cy="7" r="4"></circle>
      </svg>
      <span>Profil</span>
    </router-link>
    <button @click="handleLogout" class="flex flex-col items-center gap-1 text-[11px] font-semibold text-gray-500 hover:text-primary transition-colors flex-1 bg-transparent border-none cursor-pointer">
      <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
        <polyline points="16 17 21 12 16 7"></polyline>
        <line x1="21" y1="12" x2="9" y2="12"></line>
      </svg>
      <span>Log out</span>
    </button>
  </nav>
</template>

<style scoped>
/* Scoped styles dihapus karena seluruhnya sudah dikelola oleh Tailwind CSS utility classes */
</style>
