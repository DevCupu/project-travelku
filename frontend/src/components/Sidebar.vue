<script setup>
import { useRouter, useRoute } from 'vue-router'
import { useSidebar } from '../composables/useSidebar'

const router = useRouter()
const route = useRoute()
const { isCollapsed, toggleSidebar } = useSidebar()

const handleLogout = () => {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<template>
  <aside class="hidden md:flex flex-col fixed inset-y-0 left-0 bg-white border-r border-gray-200 z-50 shadow-[4px_0_24px_rgba(0,0,0,0.02)] transition-all duration-300" :class="isCollapsed ? 'w-20' : 'w-64'">
    <!-- Brand & Toggle Button -->
    <div class="h-20 flex items-center justify-between px-6 border-b border-gray-100">
      <div v-if="!isCollapsed" class="flex items-center gap-1 cursor-pointer overflow-hidden" @click="router.push('/dashboard')">
        <span class="text-2xl font-bold tracking-tight text-primary">24visa</span>
        <span class="text-xs font-bold bg-accent-soft text-accent-blue px-2 py-0.5 rounded-full mt-1">Makassar</span>
      </div>
      <div v-else class="flex items-center justify-center w-full cursor-pointer" @click="router.push('/dashboard')">
        <span class="text-xl font-bold tracking-tight text-primary">24</span>
      </div>

      <button @click="toggleSidebar" class="p-1.5 rounded-lg border border-gray-200 text-gray-500 hover:text-primary hover:bg-red-50 hover:border-red-200 transition-colors cursor-pointer bg-white flex items-center justify-center" :class="isCollapsed ? 'mx-auto' : ''" title="Toggle Sidebar">
        <svg class="w-4 h-4 transition-transform duration-300" :class="isCollapsed ? 'rotate-180' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="15 18 9 12 15 6"></polyline>
        </svg>
      </button>
    </div>

    <!-- Navigation Links -->
    <nav class="flex-1 px-4 py-8 space-y-2 overflow-y-auto">
      <div v-if="!isCollapsed" class="px-4 mb-4 text-xs font-bold text-gray-400 uppercase tracking-wider">Main Menu</div>
      <div v-else class="mb-4 text-center text-[10px] font-bold text-gray-400 uppercase">Menu</div>
      
      <router-link to="/dashboard" class="flex items-center gap-3 py-3 rounded-xl font-semibold transition-all duration-200" :class="[isCollapsed ? 'justify-center px-0' : 'px-4', route.path === '/dashboard' ? 'bg-red-50 text-primary shadow-sm' : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900']" title="Dashboard Analytics">
        <svg class="w-5 h-5 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 3v18h18"></path>
          <path d="M18.7 8l-5.1 5.2-2.8-2.7L7 14.3"></path>
        </svg>
        <span v-if="!isCollapsed">Dashboard</span>
      </router-link>

      <router-link to="/bookings" class="flex items-center gap-3 py-3 rounded-xl font-semibold transition-all duration-200" :class="[isCollapsed ? 'justify-center px-0' : 'px-4', route.path === '/bookings' ? 'bg-red-50 text-primary shadow-sm' : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900']" title="Bookings">
        <svg class="w-5 h-5 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
          <line x1="16" y1="2" x2="16" y2="6"></line>
          <line x1="8" y1="2" x2="8" y2="6"></line>
          <line x1="3" y1="10" x2="21" y2="10"></line>
        </svg>
        <span v-if="!isCollapsed">Bookings</span>
      </router-link>

      <router-link to="/pakets" class="flex items-center gap-3 py-3 rounded-xl font-semibold transition-all duration-200" :class="[isCollapsed ? 'justify-center px-0' : 'px-4', route.path === '/pakets' ? 'bg-red-50 text-primary shadow-sm' : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900']" title="Paket Wisata">
        <svg class="w-5 h-5 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
          <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
          <line x1="12" y1="22.08" x2="12" y2="12"></line>
        </svg>
        <span v-if="!isCollapsed">Paket Wisata</span>
      </router-link>

      <router-link to="/users" class="flex items-center gap-3 py-3 rounded-xl font-semibold transition-all duration-200" :class="[isCollapsed ? 'justify-center px-0' : 'px-4', route.path === '/users' ? 'bg-red-50 text-primary shadow-sm' : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900']" title="Profil Admin">
        <svg class="w-5 h-5 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
          <circle cx="12" cy="7" r="4"></circle>
        </svg>
        <span v-if="!isCollapsed">Profil Admin</span>
      </router-link>
    </nav>

    <!-- User & Logout -->
    <div class="p-4 border-t border-gray-100 bg-gray-50/50">
      <button @click="handleLogout" class="flex w-full items-center justify-center gap-2 py-3 rounded-xl font-semibold text-gray-600 border border-gray-200 bg-white hover:border-red-200 hover:bg-red-50 hover:text-primary transition-all duration-200 cursor-pointer shadow-sm" :class="isCollapsed ? 'px-0' : 'px-4'" title="Log out">
        <svg class="w-4 h-4 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
        <span v-if="!isCollapsed">Log out</span>
      </button>
    </div>
  </aside>
</template>
