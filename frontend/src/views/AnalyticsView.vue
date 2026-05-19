<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBookings } from '../composables/useBookings'
import { usePaket } from '../composables/usePaket'
import { useSidebar } from '../composables/useSidebar'
import TopNavbar from '../components/TopNavbar.vue'
import BottomNavbar from '../components/BottomNavbar.vue'
import Sidebar from '../components/Sidebar.vue'

const { bookings, summary, isLoading: loadingBookings, fetchBookings } = useBookings()
const { pakets, isLoading: loadingPakets, fetchPakets } = usePaket()
const { isCollapsed } = useSidebar()

const period = ref('all') // 'all', 'month', 'year'

onMounted(() => {
  fetchBookings({})
  fetchPakets(false)
})

// Hitung total peserta dari booking yang valid (dikonfirmasi / selesai)
const totalPeserta = computed(() => {
  return bookings.value
    .filter(b => b.status === 'DIKONFIRMASI' || b.status === 'SELESAI')
    .reduce((sum, b) => sum + b.jumlah_peserta, 0)
})

const totalPaketAktif = computed(() => {
  return pakets.value.filter(p => p.is_active).length
})

// Data Grafik Tren Bulanan (Jan - Des)
const monthlyData = computed(() => {
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des']
  const counts = new Array(12).fill(0)
  const revenues = new Array(12).fill(0)

  bookings.value.forEach(b => {
    if (!b.created_at) return
    const date = new Date(b.created_at)
    const m = date.getMonth()
    counts[m] += 1
    if (b.status === 'DIKONFIRMASI' || b.status === 'SELESAI') {
      revenues[m] += (b.harga_per_orang * b.jumlah_peserta)
    }
  })

  const maxCount = Math.max(...counts, 1)
  const maxRevenue = Math.max(...revenues, 1)

  return months.map((m, i) => ({
    month: m,
    count: counts[i],
    revenue: revenues[i],
    countPercent: Math.round((counts[i] / maxCount) * 100),
    revenuePercent: Math.round((revenues[i] / maxRevenue) * 100)
  }))
})

// Data Distribusi Status Booking
const statusDistribution = computed(() => {
  const dist = {
    'MENUNGGU': { label: 'Menunggu Konfirmasi', count: 0, color: 'bg-amber-500', text: 'text-amber-500' },
    'DIKONFIRMASI': { label: 'Dikonfirmasi', count: 0, color: 'bg-blue-600', text: 'text-blue-600' },
    'SELESAI': { label: 'Selesai / Berangkat', count: 0, color: 'bg-emerald-50', text: 'text-emerald-500' },
    'DIBATALKAN': { label: 'Dibatalkan', count: 0, color: 'bg-rose-500', text: 'text-rose-500' }
  }
  let total = 0
  bookings.value.forEach(b => {
    if (dist[b.status]) {
      dist[b.status].count += 1
      total += 1
    }
  })
  
  return Object.keys(dist).map(k => ({
    key: k,
    ...dist[k],
    percent: total > 0 ? Math.round((dist[k].count / total) * 100) : 0
  }))
})

// Paket Terpopuler berdasarkan jumlah peserta
const topPakets = computed(() => {
  const map = {}
  bookings.value.forEach(b => {
    if (!map[b.paket_id]) {
      map[b.paket_id] = {
        nama_paket: b.paket_wisata,
        jumlah_booking: 0,
        jumlah_peserta: 0,
        revenue: 0
      }
    }
    map[b.paket_id].jumlah_booking += 1
    map[b.paket_id].jumlah_peserta += b.jumlah_peserta
    if (b.status === 'DIKONFIRMASI' || b.status === 'SELESAI') {
      map[b.paket_id].revenue += (b.harga_per_orang * b.jumlah_peserta)
    }
  })

  return Object.values(map).sort((a, b) => b.jumlah_peserta - a.jumlah_peserta).slice(0, 5)
})

const formatCurrency = (amt) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amt)
</script>

<template>
  <div class="min-h-screen bg-slate-50 flex">
    <!-- Desktop Sidebar -->
    <Sidebar />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-h-screen pb-20 md:pb-0 overflow-x-hidden transition-all duration-300" :class="isCollapsed ? 'md:ml-20' : 'md:ml-64'">
      <!-- Mobile Top Navbar -->
      <TopNavbar class="md:hidden" />

      <main class="max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-1">
        <!-- Header -->
        <header class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4 border-b border-gray-200 pb-5">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">Executive <span class="text-primary">Dashboard</span></h1>
            <p class="text-sm text-gray-500 mt-1">Ringkasan performa bisnis dan analitik pemesanan 24visa Makassar</p>
          </div>
          <div class="flex items-center gap-2 bg-white border border-gray-200 p-1.5 rounded-xl shadow-sm">
            <button @click="period = 'all'" class="px-4 py-2 text-xs font-bold rounded-lg transition-all duration-200 cursor-pointer" :class="period === 'all' ? 'bg-primary text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'">Semua Waktu</button>
            <button @click="period = 'month'" class="px-4 py-2 text-xs font-bold rounded-lg transition-all duration-200 cursor-pointer" :class="period === 'month' ? 'bg-primary text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'">Bulan Ini</button>
            <button @click="period = 'year'" class="px-4 py-2 text-xs font-bold rounded-lg transition-all duration-200 cursor-pointer" :class="period === 'year' ? 'bg-primary text-white shadow-sm' : 'text-gray-600 hover:bg-gray-100'">Tahun Ini</button>
          </div>
        </header>

        <div v-if="loadingBookings || loadingPakets" class="p-12 text-center text-gray-500 font-medium">
          Memuat data analitik dashboard...
        </div>

        <div v-else class="space-y-8">
          <!-- Key Metric Cards (Grid of 4) -->
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            <!-- Card 1: Total Bookings -->
            <div class="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm flex flex-col justify-between hover:shadow-md transition-shadow duration-200">
              <div class="flex items-center justify-between mb-4">
                <span class="text-sm font-semibold text-gray-500">Total Bookings</span>
                <div class="w-10 h-10 rounded-xl bg-red-50 flex items-center justify-center text-primary">
                  <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                </div>
              </div>
              <div>
                <div class="text-3xl font-bold text-gray-900">{{ summary.jumlah_booking }}</div>
                <div class="text-xs font-semibold text-emerald-600 mt-2 flex items-center gap-1">
                  <span>↑ 12%</span> <span class="text-gray-400 font-normal">vs bulan lalu</span>
                </div>
              </div>
            </div>

            <!-- Card 2: Est. Revenue -->
            <div class="bg-gradient-to-br from-accent-blue to-blue-800 border border-blue-900 rounded-2xl p-6 shadow-sm flex flex-col justify-between text-white hover:shadow-md transition-shadow duration-200">
              <div class="flex items-center justify-between mb-4">
                <span class="text-sm font-semibold text-blue-200">Est. Pendapatan Bersih</span>
                <div class="w-10 h-10 rounded-xl bg-white/10 flex items-center justify-center text-white">
                  <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="1" x2="12" y2="23"></line><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path></svg>
                </div>
              </div>
              <div>
                <div class="text-2xl sm:text-3xl font-bold tracking-tight truncate w-full block" :title="formatCurrency(summary.total_estimasi_pendapatan)">{{ formatCurrency(summary.total_estimasi_pendapatan) }}</div>
                <div class="text-xs font-semibold text-blue-200 mt-2 flex items-center gap-1">
                  <span>↑ 18.5%</span> <span class="text-blue-300 font-normal">pertumbuhan stabil</span>
                </div>
              </div>
            </div>

            <!-- Card 3: Total Peserta (Pax) -->
            <div class="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm flex flex-col justify-between hover:shadow-md transition-shadow duration-200">
              <div class="flex items-center justify-between mb-4">
                <span class="text-sm font-semibold text-gray-500">Total Peserta (Pax)</span>
                <div class="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center text-accent-blue">
                  <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
                </div>
              </div>
              <div>
                <div class="text-3xl font-bold text-gray-900">{{ totalPeserta }} <span class="text-sm font-medium text-gray-500">Orang</span></div>
                <div class="text-xs font-semibold text-emerald-600 mt-2 flex items-center gap-1">
                  <span>Telah Dikonfirmasi / Selesai</span>
                </div>
              </div>
            </div>

            <!-- Card 4: Paket Wisata Aktif -->
            <div class="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm flex flex-col justify-between hover:shadow-md transition-shadow duration-200">
              <div class="flex items-center justify-between mb-4">
                <span class="text-sm font-semibold text-gray-500">Paket Wisata Aktif</span>
                <div class="w-10 h-10 rounded-xl bg-orange-50 flex items-center justify-center text-orange-500">
                  <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
                </div>
              </div>
              <div>
                <div class="text-3xl font-bold text-gray-900">{{ totalPaketAktif }} <span class="text-sm font-medium text-gray-500">Paket</span></div>
                <div class="text-xs font-semibold text-primary mt-2 flex items-center gap-1">
                  <router-link to="/pakets" class="hover:underline">Kelola Paket Wisata →</router-link>
                </div>
              </div>
            </div>
          </div>

          <!-- Charts Section (Grid of 2) -->
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <!-- Chart 1: Tren Pemesanan Bulanan (2 Cols) -->
            <div class="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm lg:col-span-2 flex flex-col">
              <div class="flex items-center justify-between mb-6">
                <div>
                  <h3 class="text-lg font-bold text-gray-900">Grafik Tren Pemesanan 2026</h3>
                  <p class="text-xs text-gray-500 mt-0.5">Perbandingan jumlah booking per bulan</p>
                </div>
                <div class="flex items-center gap-4 text-xs font-semibold">
                  <div class="flex items-center gap-1.5"><span class="w-3 h-3 rounded-full bg-primary inline-block"></span> Total Booking</div>
                </div>
              </div>

              <!-- Bar Chart Representation using CSS/Flex -->
              <div class="flex-1 flex items-end gap-2 sm:gap-4 pt-8 pb-2 border-b border-gray-100 min-h-[240px]">
                <div v-for="item in monthlyData" :key="item.month" class="flex-1 flex flex-col items-center gap-2 h-full justify-end group">
                  <!-- Tooltip -->
                  <div class="opacity-0 group-hover:opacity-100 transition-opacity duration-200 bg-gray-900 text-white text-[10px] font-bold py-1 px-2 rounded absolute -mt-10 pointer-events-none whitespace-nowrap shadow-lg z-10">
                    {{ item.count }} Bookings
                  </div>
                  <!-- Bar -->
                  <div class="w-full bg-red-100 hover:bg-primary transition-all duration-300 rounded-t-lg relative overflow-hidden" :style="{ height: `${Math.max(item.countPercent, 8)}%` }">
                    <div class="absolute inset-0 bg-primary opacity-80 group-hover:opacity-100 transition-opacity"></div>
                  </div>
                  <!-- Label -->
                  <span class="text-xs font-semibold text-gray-500 group-hover:text-gray-900">{{ item.month }}</span>
                </div>
              </div>
            </div>

            <!-- Chart 2: Distribusi Status Pemesanan (1 Col) -->
            <div class="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm flex flex-col justify-between">
              <div>
                <h3 class="text-lg font-bold text-gray-900 mb-1">Distribusi Status</h3>
                <p class="text-xs text-gray-500 mb-6">Persentase status dari seluruh pemesanan</p>
                
                <div class="space-y-4">
                  <div v-for="item in statusDistribution" :key="item.key" class="space-y-1">
                    <div class="flex justify-between text-sm font-semibold">
                      <span class="text-gray-700">{{ item.label }}</span>
                      <span :class="item.text">{{ item.percent }}% ({{ item.count }})</span>
                    </div>
                    <!-- Progress Bar -->
                    <div class="w-full bg-gray-100 h-2.5 rounded-full overflow-hidden">
                      <div class="h-full rounded-full transition-all duration-500" :class="item.color" :style="{ width: `${item.percent}%` }"></div>
                    </div>
                  </div>
                </div>
              </div>

              <div class="mt-8 pt-6 border-t border-gray-100 flex items-center justify-between text-xs text-gray-500">
                <span>Total akumulasi sistem</span>
                <span class="font-bold text-gray-900">{{ bookings.length }} Transaksi</span>
              </div>
            </div>
          </div>

          <!-- Top Paket Wisata & Quick Actions -->
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <!-- Top Paket Wisata (2 Cols) -->
            <div class="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm lg:col-span-2">
              <div class="flex items-center justify-between mb-6">
                <div>
                  <h3 class="text-lg font-bold text-gray-900">Paket Wisata Terpopuler</h3>
                  <p class="text-xs text-gray-500 mt-0.5">Berdasarkan jumlah peserta terbanyak yang mendaftar</p>
                </div>
                <router-link to="/pakets" class="text-xs font-bold text-accent-blue hover:underline">Lihat Semua</router-link>
              </div>

              <div class="space-y-4">
                <div v-for="(p, idx) in topPakets" :key="idx" class="flex items-center justify-between p-4 rounded-xl border border-gray-100 hover:border-gray-200 hover:bg-slate-50/50 transition gap-4">
                  <div class="flex items-center gap-4">
                    <div class="w-8 h-8 rounded-lg bg-gray-100 font-bold text-gray-700 flex items-center justify-center text-sm">
                      #{{ idx + 1 }}
                    </div>
                    <div>
                      <div class="font-bold text-gray-900 text-sm">{{ p.nama_paket }}</div>
                      <div class="text-xs text-gray-500 mt-0.5">{{ p.jumlah_booking }} kali dipesan</div>
                    </div>
                  </div>
                  <div class="text-right">
                    <div class="font-bold text-primary text-sm">{{ p.jumlah_peserta }} Pax</div>
                    <div class="text-xs text-gray-500 mt-0.5">Total Peserta</div>
                  </div>
                </div>
                <div v-if="topPakets.length === 0" class="text-center py-8 text-gray-500 text-sm">
                  Belum ada data paket wisata terpopuler.
                </div>
              </div>
            </div>

            <!-- Quick Actions Card (1 Col) -->
            <div class="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm flex flex-col justify-between">
              <div>
                <h3 class="text-lg font-bold text-gray-900 mb-1">Aksi Cepat</h3>
                <p class="text-xs text-gray-500 mb-6">Pintasan navigasi dan operasional harian</p>

                <div class="space-y-3">
                  <router-link to="/bookings" class="flex items-center justify-between p-4 rounded-xl border border-red-100 bg-red-50/50 hover:bg-red-50 text-primary font-bold text-sm transition group">
                    <div class="flex items-center gap-3">
                      <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                      <span>Kelola Daftar Booking</span>
                    </div>
                    <span class="group-hover:translate-x-1 transition-transform">→</span>
                  </router-link>

                  <router-link to="/pakets" class="flex items-center justify-between p-4 rounded-xl border border-blue-100 bg-blue-50/50 hover:bg-blue-50 text-accent-blue font-bold text-sm transition group">
                    <div class="flex items-center gap-3">
                      <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
                      <span>Manajemen Paket Wisata</span>
                    </div>
                    <span class="group-hover:translate-x-1 transition-transform">→</span>
                  </router-link>

                  <router-link to="/users" class="flex items-center justify-between p-4 rounded-xl border border-gray-200 bg-gray-50 hover:bg-gray-100 text-gray-700 font-bold text-sm transition group">
                    <div class="flex items-center gap-3">
                      <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                      <span>Pengaturan Profil Admin</span>
                    </div>
                    <span class="group-hover:translate-x-1 transition-transform">→</span>
                  </router-link>
                </div>
              </div>

              <div class="mt-6 pt-6 border-t border-gray-100 text-center">
                <div class="text-xs font-bold text-gray-400 uppercase tracking-wider">Sistem Informasi</div>
                <div class="text-xs text-gray-500 mt-1">24visa Makassar v2.0 - E-Warung Platform</div>
              </div>
            </div>
          </div>
        </div>
      </main>

      <!-- Mobile Bottom Navbar -->
      <BottomNavbar class="md:hidden" />
    </div>
  </div>
</template>
