import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import BookingsView from '../views/BookingsView.vue'
import UsersView from '../views/UsersView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard' // Mengarah langsung ke dashboard
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { requiresGuest: true } // Hanya boleh diakses tamu (belum login)
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/AnalyticsView.vue'),
      meta: { requiresAuth: true } // Harus login
    },
    {
      path: '/bookings',
      name: 'bookings',
      component: BookingsView,
      meta: { requiresAuth: true } // Harus login
    },
    {
      path: '/users',
      name: 'users',
      component: UsersView,
      meta: { requiresAuth: true } // Harus login
    },
    {
      path: '/pakets',
      name: 'pakets',
      component: () => import('../views/PaketsView.vue'),
      meta: { requiresAuth: true } // Harus login
    }
  ]
})

// GLOBAL ROUTE GUARD (Best Practice)
// Memusatkan pengecekan login di satu tempat
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('auth_token')

  if (to.meta.requiresAuth && !isAuthenticated) {
    // Jika rute butuh login tapi user tidak punya token, tendang ke login
    next({ name: 'login' })
  } else if (to.meta.requiresGuest && isAuthenticated) {
    // Jika user sudah login mencoba buka halaman login, tendang ke dashboard
    next({ name: 'dashboard' })
  } else {
    // Lanjutkan navigasi normal
    next()
  }
})

export default router
