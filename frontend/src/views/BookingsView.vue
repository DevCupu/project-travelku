<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useBookings } from '../composables/useBookings'
import { usePaket } from '../composables/usePaket'
import TopNavbar from '../components/TopNavbar.vue'
import BottomNavbar from '../components/BottomNavbar.vue'
import Sidebar from '../components/Sidebar.vue'
import { useSidebar } from '../composables/useSidebar'

const router = useRouter()

// Tarik fungsi dan state dari composable
const { 
  bookings, summary, isLoading, 
  totalBookings, currentPage, itemsPerPage,
  fetchBookings, saveBooking, deleteBooking, updateStatus 
} = useBookings()

const { pakets, fetchPakets } = usePaket()
const { isCollapsed } = useSidebar()

const filters = ref({ status: '', paket_wisata: '', date_from: '', date_to: '' })

// Lifecycle
onMounted(() => {
  fetchBookings({ ...filters.value, page: 1 })
  fetchPakets(true) // Ambil hanya paket yang aktif
})

// Pagenasi
const totalPages = computed(() => Math.ceil(totalBookings.value / itemsPerPage.value) || 1)
const changePage = (p) => {
  if (p < 1 || p > totalPages.value) return
  fetchBookings({ ...filters.value, page: p })
}

// REAKTIVITAS REAL-TIME UNTUK FILTER PENCARIAN
watch(filters, (newFilters) => {
  fetchBookings({ ...newFilters, page: 1 })
}, { deep: true })

// Modal & Form State
const showModal = ref(false)
const isEditMode = ref(false)
const currentId = ref(null)
const currentBookingData = ref(null)

const form = ref({
  nama_pemesan: '',
  kontak: '',
  paket_id: '',
  jumlah_peserta: 1,
  catatan: ''
})

// BEST PRACTICE VUE 3: Gunakan Computed Properties untuk kalkulasi & validasi reaktif
const selectedPaket = computed(() => {
  if (!form.value.paket_id) return null
  return pakets.value.find(p => p.id === form.value.paket_id) || null
})

const calculatedTotal = computed(() => {
  if (!selectedPaket.value || !form.value.jumlah_peserta) return 0
  return selectedPaket.value.harga * form.value.jumlah_peserta
})

const sisaKuota = computed(() => {
  if (!selectedPaket.value) return 0
  // Jika dalam mode edit dan paket tidak diubah, kuota tersedia adalah sisa kuota + jumlah peserta lama
  if (isEditMode.value && currentBookingData.value?.paket_id === form.value.paket_id) {
    return selectedPaket.value.sisa_kuota + currentBookingData.value.jumlah_peserta
  }
  return selectedPaket.value.sisa_kuota
})

const isQuotaExceeded = computed(() => {
  if (!selectedPaket.value) return false
  return form.value.jumlah_peserta > sisaKuota.value
})

const isNamaInvalid = computed(() => !form.value.nama_pemesan || form.value.nama_pemesan.trim().length < 3)
const isKontakInvalid = computed(() => !form.value.kontak || form.value.kontak.trim().length < 3)
const isJumlahPesertaInvalid = computed(() => form.value.jumlah_peserta < 1)

const isFormInvalid = computed(() => {
  return !form.value.paket_id || isNamaInvalid.value || isKontakInvalid.value || isJumlahPesertaInvalid.value || isQuotaExceeded.value
})

const resetForm = () => {
  form.value = {
    nama_pemesan: '',
    kontak: '',
    paket_id: '',
    jumlah_peserta: 1,
    catatan: ''
  }
  currentBookingData.value = null
}

const openCreateModal = () => {
  isEditMode.value = false
  currentId.value = null
  resetForm()
  showModal.value = true
}

const openEditModal = (b) => {
  isEditMode.value = true
  currentId.value = b.id
  currentBookingData.value = { ...b }
  form.value = { 
    nama_pemesan: b.nama_pemesan,
    kontak: b.kontak,
    paket_id: b.paket_id,
    jumlah_peserta: b.jumlah_peserta,
    catatan: b.catatan || ''
  }
  showModal.value = true
}

// State Toast Notification Premium
const toast = ref({ show: false, message: '', type: 'success' })
const showToast = (msg, type = 'success') => {
  toast.value = { show: true, message: msg, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

// UI Actions
const handleSubmit = async () => {
  if (isFormInvalid.value) {
    showToast('Mohon perbaiki error pada form sebelum menyimpan!', 'error')
    return
  }
  try {
    const payload = {
      nama_pemesan: form.value.nama_pemesan,
      kontak: form.value.kontak,
      paket_id: form.value.paket_id,
      jumlah_peserta: form.value.jumlah_peserta,
      catatan: form.value.catatan ? form.value.catatan : null
    }
    await saveBooking(currentId.value, payload)
    showModal.value = false
    showToast(isEditMode.value ? 'Pemesanan berhasil diperbarui!' : 'Pemesanan baru berhasil dibuat!', 'success')
    fetchBookings(filters.value)
  } catch (err) {
    showToast(err.response?.data?.message || 'Gagal menyimpan pemesanan', 'error')
  }
}

// State for Custom Confirmation Modal
const showConfirmModal = ref(false)
const confirmCallback = ref(null)
const confirmMessage = ref('')
const confirmTitle = ref('')

const openConfirm = (title, message, callback) => {
  confirmTitle.value = title
  confirmMessage.value = message
  confirmCallback.value = callback
  showConfirmModal.value = true
}

const handleConfirmAction = async () => {
  showConfirmModal.value = false
  if (confirmCallback.value) {
    await confirmCallback.value()
  }
}

const handleDelete = (id) => {
  openConfirm(
    'Hapus Transaksi Pemesanan?',
    'Apakah Anda yakin ingin menghapus transaksi booking jemaah ini? Tindakan ini bersifat permanen dan hanya bisa dilakukan jika status transaksi belum bersifat final (Selesai/Batal).',
    async () => {
      try {
        await deleteBooking(id)
        showToast('Pemesanan berhasil dihapus!', 'success')
        fetchBookings(filters.value)
      } catch (err) {
        showToast('Gagal menghapus pemesanan', 'error')
      }
    }
  )
}

const handleStatus = async (id, status) => {
  try {
    await updateStatus(id, status)
    showToast(`Status pemesanan berhasil diperbarui menjadi ${status}!`, 'success')
    fetchBookings(filters.value)
  } catch (err) {
    showToast('Gagal mengubah status pemesanan', 'error')
  }
}

// Helpers
const formatCurrency = (amt) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amt)
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}
const getNextStatus = (current) => current === 'MENUNGGU' ? 'DIKONFIRMASI' : current === 'DIKONFIRMASI' ? 'SELESAI' : null
const statusColor = (status) => {
  const colors = { 'MENUNGGU': '#f59e0b', 'DIKONFIRMASI': '#3b82f6', 'SELESAI': '#10b981', 'DIBATALKAN': '#ef4444' }
  return colors[status] || '#717171'
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 flex">
    <!-- Toast Notification Premium -->
    <Transition name="toast">
      <div v-if="toast.show" class="fixed top-6 right-6 z-[999] flex items-center gap-3 px-5 py-3.5 rounded-xl shadow-xl border backdrop-blur-md transition-all duration-300"
           :class="toast.type === 'error' ? 'bg-red-50/95 border-red-200 text-red-800' : 'bg-emerald-50/95 border-emerald-200 text-emerald-800'">
        <svg v-if="toast.type === 'success'" class="w-5 h-5 text-emerald-600 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <svg v-else class="w-5 h-5 text-red-600 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <span class="text-sm font-bold">{{ toast.message }}</span>
      </div>
    </Transition>

    <!-- Desktop Sidebar -->
    <Sidebar />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-h-screen pb-20 md:pb-0 overflow-x-hidden transition-all duration-300" :class="isCollapsed ? 'md:ml-20' : 'md:ml-64'">
      <!-- Mobile Top Navbar -->
      <TopNavbar class="md:hidden" />

      <main class="max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-1">
        <header class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4 border-b border-gray-200 pb-5">
          <h1 class="text-3xl font-bold text-gray-900">Visa & Umrah <span class="text-accent-blue">Bookings</span></h1>
          <button class="h-12 px-6 bg-primary hover:bg-primary-active text-white font-semibold rounded-lg transition text-[16px] flex items-center justify-center whitespace-nowrap cursor-pointer" @click="openCreateModal">
            + Buat Pemesanan
          </button>
        </header>

        <!-- Summaries -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div class="bg-white border border-gray-200 rounded-2xl p-6 shadow-sm flex flex-col justify-center">
            <div class="text-sm font-semibold text-gray-500 mb-2">Total Bookings</div>
            <div class="text-4xl font-bold text-gray-900">{{ summary.jumlah_booking }}</div>
          </div>
          <div class="bg-gradient-to-br from-accent-blue to-blue-800 border border-blue-900 rounded-2xl p-6 shadow-sm flex flex-col justify-center text-white">
            <div class="text-sm font-semibold text-blue-200 mb-2">Est. Revenue (Confirmed/Done)</div>
            <div class="text-2xl sm:text-3xl font-bold tracking-tight truncate w-full block" :title="formatCurrency(summary.total_estimasi_pendapatan)">{{ formatCurrency(summary.total_estimasi_pendapatan) }}</div>
          </div>
        </div>

        <!-- Filters -->
        <div class="bg-white p-5 rounded-2xl border border-gray-200 shadow-sm mb-8 flex flex-col md:flex-row gap-4">
          <div class="flex-1 flex flex-col gap-2">
            <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">Status</label>
            <select v-model="filters.status" class="px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-slate-50">
              <option value="">All Status</option>
              <option value="MENUNGGU">Menunggu</option>
              <option value="DIKONFIRMASI">Dikonfirmasi</option>
              <option value="SELESAI">Selesai</option>
              <option value="DIBATALKAN">Dibatalkan</option>
            </select>
          </div>
          <div class="flex-1 flex flex-col gap-2">
            <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">Paket Wisata</label>
            <input type="text" v-model="filters.paket_wisata" placeholder="Search paket..." class="px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-slate-50" />
          </div>
          <div class="flex-1 flex flex-col gap-2">
            <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">Date From</label>
            <input type="date" v-model="filters.date_from" class="px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-slate-50" />
          </div>
          <div class="flex-1 flex flex-col gap-2">
            <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">Date To</label>
            <input type="date" v-model="filters.date_to" class="px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-slate-50" />
          </div>
        </div>

        <!-- Table Container -->
        <div class="bg-white border border-gray-200 rounded-2xl overflow-hidden shadow-sm">
          <div v-if="isLoading" class="p-12 text-center text-gray-500 font-medium">Loading bookings data...</div>
          
          <div v-else-if="bookings.length > 0" class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-slate-50 border-b border-gray-200">
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Jemaah / Pelanggan</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Kontak</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Layanan (Paket/Visa)</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Tgl. Berangkat</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Pax</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Total Biaya</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Status</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <tr v-for="b in bookings" :key="b.id" class="hover:bg-slate-50/50 transition">
                  <td class="px-6 py-4">
                    <div class="font-bold text-gray-900">{{ b.nama_pemesan }}</div>
                  </td>
                  <td class="px-6 py-4 text-sm text-gray-600">{{ b.kontak }}</td>
                  <td class="px-6 py-4 text-sm font-medium text-gray-800">{{ b.paket_wisata }}</td>
                  <td class="px-6 py-4 text-sm text-gray-600">{{ b.tanggal_berangkat }}</td>
                  <td class="px-6 py-4 text-sm text-gray-900 font-bold">{{ b.jumlah_peserta }}</td>
                  <td class="px-6 py-4 text-sm font-medium text-gray-900">{{ formatCurrency(b.jumlah_peserta * b.harga_per_orang) }}</td>
                  <td class="px-6 py-4">
                    <span class="px-3 py-1 text-xs font-bold rounded-full text-white tracking-wide" 
                          :style="{ backgroundColor: statusColor(b.status) }">
                      {{ b.status }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right">
                    <div class="inline-flex items-center justify-end gap-2.5">
                      <!-- Tombol Status Transisi Utama -->
                      <button 
                        v-if="getNextStatus(b.status)" 
                        @click="handleStatus(b.id, getNextStatus(b.status))"
                        class="inline-flex items-center gap-1 px-3 py-1.5 text-xs font-bold bg-blue-600 hover:bg-blue-700 text-white rounded-lg shadow-sm transition-all duration-150 active:scale-95 cursor-pointer"
                        :title="'Ubah status menjadi ' + getNextStatus(b.status)"
                      >
                        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                          <polyline points="22 4 12 14.01 9 11.01"></polyline>
                        </svg>
                        <span>{{ getNextStatus(b.status) === 'DIKONFIRMASI' ? 'Konfirmasi' : 'Selesai' }}</span>
                      </button>

                      <!-- Tombol Batal -->
                      <button 
                        v-if="b.status === 'MENUNGGU' || b.status === 'DIKONFIRMASI'" 
                        @click="handleStatus(b.id, 'DIBATALKAN')"
                        class="inline-flex items-center gap-1 px-3 py-1.5 text-xs font-bold border border-rose-200 bg-rose-50 hover:bg-rose-100 text-rose-700 rounded-lg transition-all duration-150 active:scale-95 cursor-pointer"
                        title="Batalkan Booking"
                      >
                        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                          <circle cx="12" cy="12" r="10"></circle>
                          <line x1="15" y1="9" x2="9" y2="15"></line>
                          <line x1="9" y1="9" x2="15" y2="15"></line>
                        </svg>
                        <span>Batal</span>
                      </button>
                      
                      <!-- Pembatas vertikal halus -->
                      <span class="w-px h-5 bg-gray-200 self-center" v-if="b.status !== 'SELESAI' && b.status !== 'DIBATALKAN'"></span>
                      
                      <!-- Tombol Edit (Ikon Saja - Sangat Rapi) -->
                      <button 
                        v-if="b.status !== 'SELESAI' && b.status !== 'DIBATALKAN'"
                        @click="openEditModal(b)" 
                        class="p-2 text-gray-500 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-all duration-150 cursor-pointer"
                        title="Edit data pemesanan"
                      >
                        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                        </svg>
                      </button>

                      <!-- Tombol Hapus (Ikon Saja - Sangat Rapi) -->
                      <button 
                        v-if="b.status !== 'SELESAI' && b.status !== 'DIBATALKAN'"
                        @click="handleDelete(b.id)" 
                        class="p-2 text-gray-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-all duration-150 cursor-pointer"
                        title="Hapus pemesanan"
                      >
                        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <polyline points="3 6 5 6 21 6"></polyline>
                          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Pagenasi Premium -->
          <div class="px-6 py-4 flex items-center justify-between border-t border-slate-100 bg-white" v-if="bookings.length > 0">
            <div class="flex items-center gap-2">
              <span class="text-xs text-slate-500 font-medium">Menampilkan</span>
              <span class="text-xs font-bold text-slate-800 bg-slate-100 px-2 py-1 rounded-md">
                {{ bookings.length }} dari {{ totalBookings }} data
              </span>
            </div>
            <div class="flex items-center gap-1">
              <button 
                @click="changePage(currentPage - 1)" 
                :disabled="currentPage === 1"
                class="p-2 border border-slate-200 rounded-lg hover:bg-slate-50 disabled:opacity-40 disabled:hover:bg-transparent transition-all duration-150 cursor-pointer"
                title="Halaman Sebelumnya"
              >
                <svg class="w-4 h-4 text-slate-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="15 18 9 12 15 6"></polyline>
                </svg>
              </button>
              
              <template v-for="p in totalPages" :key="p">
                <button 
                  @click="changePage(p)"
                  class="w-8 h-8 rounded-lg text-xs font-bold transition-all duration-150 cursor-pointer"
                  :class="currentPage === p ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-slate-600 hover:bg-slate-100 border border-transparent hover:border-slate-200'"
                >
                  {{ p }}
                </button>
              </template>
              
              <button 
                @click="changePage(currentPage + 1)" 
                :disabled="currentPage === totalPages"
                class="p-2 border border-slate-200 rounded-lg hover:bg-slate-50 disabled:opacity-40 disabled:hover:bg-transparent transition-all duration-150 cursor-pointer"
                title="Halaman Berikutnya"
              >
                <svg class="w-4 h-4 text-slate-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="9 18 15 12 9 6"></polyline>
                </svg>
              </button>
            </div>
          </div>

          <div v-else class="p-16 text-center text-gray-500 font-medium">
            <p>No bookings found based on your filters.</p>
          </div>
        </div>
      </main>

      <!-- Modal (Create/Edit) -->
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal-content">
          <h2 class="text-2xl font-bold mb-6 text-accent-blue border-b-2 border-accent-soft pb-3">{{ isEditMode ? 'Edit Booking' : 'Create New Booking' }}</h2>
          <form @submit.prevent="handleSubmit" class="space-y-5">
            
            <div class="flex flex-col md:flex-row gap-4">
              <div class="flex-1">
                <label class="block text-sm font-semibold mb-1 text-gray-700">Jemaah / Pelanggan <span class="text-primary">*</span></label>
                <input type="text" v-model="form.nama_pemesan" required placeholder="Nama lengkap jemaah" class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white" />
                <span v-if="form.nama_pemesan && isNamaInvalid" class="text-xs text-primary mt-1 block font-medium">⚠️ Nama pemesan minimal 3 karakter</span>
              </div>
              <div class="flex-1">
                <label class="block text-sm font-semibold mb-1 text-gray-700">Kontak (No. HP / Email) <span class="text-primary">*</span></label>
                <input type="text" v-model="form.kontak" required placeholder="08123456789" class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white" />
                <span v-if="form.kontak && isKontakInvalid" class="text-xs text-primary mt-1 block font-medium">⚠️ Kontak minimal 3 karakter</span>
              </div>
            </div>

            <div>
              <label class="block text-sm font-semibold mb-1 text-gray-700">Paket Wisata / Layanan <span class="text-primary">*</span></label>
              <select v-model="form.paket_id" required class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white">
                <option value="" disabled>Pilih Paket Wisata...</option>
                <option v-for="p in pakets" :key="p.id" :value="p.id">
                  {{ p.nama_paket }} (Sisa: {{ p.sisa_kuota }} Pax - {{ formatCurrency(p.harga) }})
                </option>
              </select>
            </div>

            <div v-if="selectedPaket" class="bg-accent-soft border border-accent-blue/20 rounded-xl p-4 my-2">
              <div class="flex flex-col md:flex-row justify-between items-start md:items-center border-b border-accent-blue/10 pb-3 mb-3 gap-1">
                <span class="font-bold text-gray-900 text-base">{{ selectedPaket.nama_paket }}</span>
                <span class="text-sm text-accent-blue font-bold px-3 py-1 bg-white rounded-full shadow-sm">🛫 Berangkat: {{ formatDate(selectedPaket.tanggal_berangkat) }}</span>
              </div>
              <div class="flex flex-col md:flex-row justify-between gap-4">
                <div class="flex flex-col gap-1">
                  <span class="text-xs text-gray-600">Harga per Pax:</span>
                  <span class="text-sm font-bold text-gray-900">{{ formatCurrency(selectedPaket.harga) }}</span>
                </div>
                <div class="flex flex-col gap-1">
                  <span class="text-xs text-gray-600">Sisa Kuota Tersedia:</span>
                  <span class="text-sm font-bold" :class="sisaKuota <= 5 ? 'text-primary' : 'text-accent-blue'">{{ sisaKuota }} Pax</span>
                </div>
              </div>
            </div>

            <div class="flex flex-col md:flex-row gap-4 items-start md:items-center">
              <div class="flex-1 w-full">
                <label class="block text-sm font-semibold mb-1 text-gray-700">Jumlah Peserta (Pax) <span class="text-primary">*</span></label>
                <input type="number" v-model="form.jumlah_peserta" min="1" required class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white" />
                <span v-if="isJumlahPesertaInvalid" class="text-xs text-primary mt-1 block font-medium">⚠️ Jumlah peserta minimal 1 pax</span>
                <span v-else-if="isQuotaExceeded" class="text-xs text-primary mt-1 block font-medium">⚠️ Melebihi sisa kuota yang tersedia (Sisa: {{ sisaKuota }} Pax)</span>
              </div>
              <div class="flex-1 w-full">
                <label class="block text-sm font-semibold mb-1 text-gray-700">Total Estimasi Biaya</label>
                <div class="px-4 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-base font-bold text-gray-800 flex items-center h-11">
                  {{ formatCurrency(calculatedTotal) }}
                </div>
              </div>
            </div>

            <div>
              <label class="block text-sm font-semibold mb-1 text-gray-700">Catatan Tambahan (Opsional)</label>
              <textarea v-model="form.catatan" rows="3" placeholder="Kebutuhan khusus, preferensi kamar, dll..." class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white"></textarea>
            </div>

            <div class="flex justify-end gap-3 mt-8 pt-6 border-t border-gray-100">
              <button type="button" class="h-12 px-6 border border-gray-300 rounded-lg font-semibold text-gray-900 bg-white hover:bg-gray-50 transition cursor-pointer" @click="showModal = false">Batal</button>
              <button type="submit" class="h-12 px-6 bg-primary hover:bg-primary-active text-white font-semibold rounded-lg transition disabled:bg-primary-disabled disabled:cursor-not-allowed cursor-pointer" :disabled="isFormInvalid">Simpan Pemesanan</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Mobile Bottom Navbar -->
      <BottomNavbar class="md:hidden" />

      <!-- Beautiful Custom Confirmation Modal -->
      <Transition name="confirm-modal">
        <div v-if="showConfirmModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
          <!-- Backdrop -->
          <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm transition-opacity animate-fade-in" @click="showConfirmModal = false"></div>
          
          <!-- Modal Card -->
          <div class="bg-white rounded-2xl shadow-2xl border border-gray-100 max-w-md w-full overflow-hidden transform transition-all z-10 scale-100">
            <div class="p-6">
              <!-- Icon and Header -->
              <div class="flex items-center gap-4 mb-4">
                <div class="w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0 bg-red-50 text-red-600">
                  <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </div>
                <div>
                  <h3 class="text-[17px] font-bold text-gray-900">{{ confirmTitle }}</h3>
                  <p class="text-[10px] text-red-500 uppercase tracking-wider font-bold mt-0.5">Konfirmasi Penghapusan</p>
                </div>
              </div>
              
              <!-- Message -->
              <p class="text-sm text-gray-600 leading-relaxed">{{ confirmMessage }}</p>
            </div>
            
            <!-- Actions Footer -->
            <div class="bg-slate-50 px-6 py-4 flex justify-end gap-3 border-t border-slate-100">
              <button @click="showConfirmModal = false" 
                      class="h-10 px-4 text-sm font-semibold text-gray-600 hover:text-gray-800 bg-white hover:bg-slate-100 rounded-lg border border-gray-200 transition cursor-pointer">
                Batal
              </button>
              <button @click="handleConfirmAction" 
                      class="h-10 px-5 text-sm font-semibold text-white bg-red-600 hover:bg-red-700 active:bg-red-800 rounded-lg transition cursor-pointer">
                Ya, Hapus
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
/* Toast Transition */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-enter-from,
.toast-leave-to {
  transform: translateY(-20px) scale(0.95);
  opacity: 0;
}

/* Custom Confirm Modal Transition */
.confirm-modal-enter-active,
.confirm-modal-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.confirm-modal-enter-from,
.confirm-modal-leave-to {
  opacity: 0;
}

.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal-content {
  background: white;
  width: 100%;
  max-width: 620px;
  border-radius: var(--radius-lg);
  padding: 32px;
  box-shadow: 0 10px 40px rgba(0,0,0,0.15);
  max-height: 90vh;
  overflow-y: auto;
}

@media (max-width: 768px) {
  .modal-content {
    padding: 24px 16px;
    margin: 16px;
    max-height: 90vh;
    overflow-y: auto;
  }
}
</style>
