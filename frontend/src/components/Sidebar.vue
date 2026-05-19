<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSidebar } from '../composables/useSidebar'

const router = useRouter()
const route = useRoute()
const { isCollapsed, toggleSidebar } = useSidebar()

const showLogoutModal = ref(false)

const handleLogout = () => {
  showLogoutModal.value = true
}

const confirmLogout = () => {
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

  <!-- Premium Logout Confirmation Modal -->
  <Transition name="fade">
    <div v-if="showLogoutModal" class="fixed inset-0 z-[999] flex items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm transition-opacity" @click="showLogoutModal = false"></div>
      
      <!-- Modal Card -->
      <div class="bg-white rounded-2xl shadow-2xl border border-gray-100 max-w-sm w-full overflow-hidden transform transition-all z-10 p-6 animate-scale-up">
        <!-- Icon and Header -->
        <div class="flex flex-col items-center text-center mb-6">
          <div class="w-12 h-12 rounded-full flex items-center justify-center bg-primary/10 text-primary mb-4 shadow-sm">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900">Yakin Ingin Keluar?</h3>
          <p class="text-xs text-gray-500 mt-2 px-2">Anda akan keluar dari akun TravelKu dan harus login kembali untuk mengakses dashboard.</p>
        </div>
        
        <!-- Actions Footer -->
        <div class="flex gap-3 justify-stretch">
          <button @click="showLogoutModal = false" 
                  class="flex-1 py-2.5 text-xs font-bold text-gray-500 hover:text-gray-700 bg-gray-50 hover:bg-gray-100 rounded-xl border border-gray-200 transition-colors cursor-pointer text-center">
            Batal
          </button>
          <button @click="confirmLogout" 
                  class="flex-1 py-2.5 text-xs font-bold text-white bg-primary hover:bg-primary-hover active:bg-primary-active rounded-xl transition-colors shadow-md shadow-red-100 hover:shadow-none cursor-pointer text-center">
            Ya, Keluar
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes scaleUp {
  from {
    transform: scale(0.95);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.animate-scale-up {
  animation: scaleUp 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
</style>
