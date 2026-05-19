<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

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
