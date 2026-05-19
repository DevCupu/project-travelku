<script setup>
import { ref, computed, onMounted } from 'vue'
import { usePaket } from '../composables/usePaket'
import TopNavbar from '../components/TopNavbar.vue'
import BottomNavbar from '../components/BottomNavbar.vue'
import Sidebar from '../components/Sidebar.vue'
import { useSidebar } from '../composables/useSidebar'

const { pakets, isLoading, error, fetchPakets, savePaket, deletePaket } = usePaket()
const { isCollapsed } = useSidebar()

onMounted(() => fetchPakets(false))

// FITUR PENCARIAN REAL-TIME
const searchQuery = ref('')
const filteredPakets = computed(() => {
  if (!searchQuery.value) return pakets.value
  return pakets.value.filter(p => p.nama_paket.toLowerCase().includes(searchQuery.value.toLowerCase()))
})

const showModal = ref(false)
const isEditMode = ref(false)
const currentId = ref(null)
const currentPaketData = ref(null)

const form = ref({
  nama_paket: '', deskripsi: '', harga: 0, kuota_total: 20, tanggal_berangkat: '', is_active: true
})

// VALIDASI REAL-TIME MENGGUNAKAN COMPUTED PROPERTIES VUE 3
const isNamaInvalid = computed(() => !form.value.nama_paket || form.value.nama_paket.trim().length < 3)
const isHargaInvalid = computed(() => form.value.harga < 0)
const isKuotaInvalid = computed(() => form.value.kuota_total < 1)

const isTanggalInvalid = computed(() => {
  if (!form.value.tanggal_berangkat) return true
  const selected = new Date(form.value.tanggal_berangkat)
  selected.setHours(0, 0, 0, 0)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  // Jika dalam mode edit dan tanggal tidak diubah dari aslinya, izinkan
  if (isEditMode.value && currentPaketData.value?.tanggal_berangkat?.split('T')[0] === form.value.tanggal_berangkat) {
    return false
  }
  return selected < today
})

const isFormInvalid = computed(() => {
  return isNamaInvalid.value || isHargaInvalid.value || isKuotaInvalid.value || isTanggalInvalid.value
})

const resetForm = () => {
  form.value = {
    nama_paket: '', deskripsi: '', harga: 0, kuota_total: 20, tanggal_berangkat: '', is_active: true
  }
  currentPaketData.value = null
}

const handleHargaInput = (e) => {
  const val = e.target.value.replace(/\D/g, '')
  form.value.harga = val ? parseInt(val, 10) : 0
}

const formatRupiahInput = (val) => {
  if (!val) return ''
  return val.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".")
}

const openCreateModal = () => {
  isEditMode.value = false
  currentId.value = null
  resetForm()
  showModal.value = true
}

const openEditModal = (p) => {
  isEditMode.value = true
  currentId.value = p.id
  currentPaketData.value = { ...p }
  form.value = { 
    ...p, 
    tanggal_berangkat: p.tanggal_berangkat ? p.tanggal_berangkat.split('T')[0] : '' 
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

const handleSubmit = async () => {
  if (isFormInvalid.value) {
    showToast('Mohon perbaiki error pada form sebelum menyimpan!', 'error')
    return
  }
  try {
    await savePaket(currentId.value, form.value)
    showModal.value = false
    showToast(isEditMode.value ? 'Paket wisata berhasil diperbarui!' : 'Paket wisata baru berhasil dibuat!', 'success')
    fetchPakets(false)
  } catch (err) {
    showToast(err.response?.data?.message || 'Gagal menyimpan paket wisata', 'error')
  }
}

const handleDelete = async (id) => {
  if (!confirm('Hapus paket wisata ini? Pastikan belum ada peserta yang terdaftar.')) return
  try {
    await deletePaket(id)
    showToast('Paket wisata berhasil dihapus!', 'success')
    fetchPakets(false)
  } catch (err) {
    showToast(err.response?.data?.message || 'Gagal menghapus paket wisata', 'error')
  }
}

const formatCurrency = (amt) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amt)
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
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
          <h1 class="text-3xl font-bold text-gray-900">Manajemen <span class="text-accent-blue">Paket Wisata</span></h1>
          <button class="h-12 px-6 bg-primary hover:bg-primary-active text-white font-semibold rounded-lg transition text-[16px] flex items-center justify-center whitespace-nowrap cursor-pointer" @click="openCreateModal">
            + Buat Paket Baru
          </button>
        </header>

        <div v-if="error" class="bg-red-50 text-red-600 p-4 rounded-xl mb-6 font-medium border border-red-100">{{ error }}</div>

        <!-- Search Filter -->
        <div class="bg-white p-5 rounded-2xl border border-gray-200 shadow-sm mb-8 flex items-center gap-4">
          <div class="flex-1 flex flex-col gap-2">
            <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">Cari Paket Wisata</label>
            <input type="text" v-model="searchQuery" placeholder="Ketik nama paket wisata..." class="px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-slate-50" />
          </div>
        </div>

        <!-- Table Container -->
        <div class="bg-white border border-gray-200 rounded-2xl overflow-hidden shadow-sm">
          <div v-if="isLoading" class="p-12 text-center text-gray-500 font-medium">Memuat data paket wisata...</div>
          
          <div v-else-if="filteredPakets.length > 0" class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-slate-50 border-b border-gray-200">
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Nama Paket</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Deskripsi</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Harga per Pax</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Tgl. Berangkat</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Kuota Total</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Sisa Kuota</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Status</th>
                  <th class="px-6 py-4 text-sm font-semibold text-gray-600">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <tr v-for="p in filteredPakets" :key="p.id" class="hover:bg-slate-50/50 transition">
                  <td class="px-6 py-4">
                    <div class="font-bold text-gray-900">{{ p.nama_paket }}</div>
                  </td>
                  <td class="px-6 py-4">
                    <div class="max-w-[200px] truncate text-gray-600 text-sm">{{ p.deskripsi || '-' }}</div>
                  </td>
                  <td class="px-6 py-4 font-medium text-gray-800">{{ formatCurrency(p.harga) }}</td>
                  <td class="px-6 py-4 text-sm text-gray-600">{{ formatDate(p.tanggal_berangkat) }}</td>
                  <td class="px-6 py-4 text-sm text-gray-600">{{ p.kuota_total }} Pax</td>
                  <td class="px-6 py-4">
                    <span class="px-3 py-1 text-xs font-bold rounded-full" 
                          :class="p.sisa_kuota === 0 ? 'bg-red-100 text-red-700' : p.sisa_kuota <= 5 ? 'bg-orange-100 text-orange-700' : 'bg-accent-soft text-accent-blue'">
                      {{ p.sisa_kuota }} Pax
                    </span>
                  </td>
                  <td class="px-6 py-4">
                    <span class="px-3 py-1 text-xs font-bold rounded-full" 
                          :class="p.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'">
                      {{ p.is_active ? 'Aktif' : 'Nonaktif' }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right">
                    <div class="inline-flex items-center justify-end gap-1.5">
                      <!-- Tombol Edit (Ikon Saja - Sangat Rapi) -->
                      <button @click="openEditModal(p)" class="p-2 text-gray-500 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-all duration-150 cursor-pointer" title="Edit paket wisata">
                        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                        </svg>
                      </button>
                      <span class="w-px h-5 bg-gray-200 self-center"></span>
                      <!-- Tombol Delete (Ikon Saja - Sangat Rapi) -->
                      <button @click="handleDelete(p.id)" class="p-2 text-gray-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-all duration-150 cursor-pointer" title="Hapus paket wisata">
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

          <div v-else class="p-16 text-center text-gray-500 font-medium">
            <p>Belum ada paket wisata yang terdaftar.</p>
          </div>
        </div>
      </main>

      <!-- Modal (Create/Edit) -->
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal-content">
          <h2 class="text-2xl font-bold mb-6 text-accent-blue border-b-2 border-accent-soft pb-3">{{ isEditMode ? 'Edit Paket Wisata' : 'Buat Paket Wisata Baru' }}</h2>
          <form @submit.prevent="handleSubmit" class="space-y-5">
            
            <div>
              <label class="block text-sm font-semibold mb-1 text-gray-700">Nama Paket Wisata <span class="text-primary">*</span></label>
              <input type="text" v-model="form.nama_paket" required placeholder="e.g. Umrah Plus Turki 12 Hari" class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white" />
              <span v-if="form.nama_paket && isNamaInvalid" class="text-xs text-primary mt-1 block font-medium">⚠️ Nama paket minimal 3 karakter</span>
            </div>

            <div>
              <label class="block text-sm font-semibold mb-1 text-gray-700">Deskripsi Singkat</label>
              <textarea v-model="form.deskripsi" rows="3" placeholder="Fasilitas, hotel, maskapai..." class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white"></textarea>
            </div>

            <div class="flex flex-col md:flex-row gap-4">
              <div class="flex-1">
                <label class="block text-sm font-semibold mb-1 text-gray-700">Harga per Pax (IDR) <span class="text-primary">*</span></label>
                <div class="relative">
                  <span class="absolute left-4 top-2.5 text-gray-500 font-medium">Rp</span>
                  <input type="text" :value="formatRupiahInput(form.harga)" @input="handleHargaInput" required class="w-full pl-11 pr-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white" />
                </div>
                <span v-if="isHargaInvalid" class="text-xs text-primary mt-1 block font-medium">⚠️ Harga tidak boleh bernilai negatif</span>
              </div>
              <div class="flex-1">
                <label class="block text-sm font-semibold mb-1 text-gray-700">Tanggal Keberangkatan <span class="text-primary">*</span></label>
                <input type="date" v-model="form.tanggal_berangkat" required class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white" />
                <span v-if="form.tanggal_berangkat && isTanggalInvalid" class="text-xs text-primary mt-1 block font-medium">⚠️ Tanggal keberangkatan tidak boleh di masa lalu</span>
              </div>
            </div>

            <div class="flex flex-col md:flex-row gap-4">
              <div class="flex-1">
                <label class="block text-sm font-semibold mb-1 text-gray-700">Kuota Total (Pax) <span class="text-primary">*</span></label>
                <input type="number" v-model="form.kuota_total" min="1" required class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white" />
                <span v-if="isKuotaInvalid" class="text-xs text-primary mt-1 block font-medium">⚠️ Kuota total minimal 1 pax</span>
              </div>
              <div class="flex-1" v-if="isEditMode">
                <label class="block text-sm font-semibold mb-1 text-gray-700">Status Aktif</label>
                <select v-model="form.is_active" class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent-blue text-sm bg-white">
                  <option :value="true">Aktif</option>
                  <option :value="false">Nonaktif</option>
                </select>
              </div>
            </div>

            <div class="flex justify-end gap-3 mt-8 pt-6 border-t border-gray-100">
              <button type="button" class="h-12 px-6 border border-gray-300 rounded-lg font-semibold text-gray-900 bg-white hover:bg-gray-50 transition cursor-pointer" @click="showModal = false">Batal</button>
              <button type="submit" class="h-12 px-6 bg-primary hover:bg-primary-active text-white font-semibold rounded-lg transition disabled:bg-primary-disabled disabled:cursor-not-allowed cursor-pointer" :disabled="isFormInvalid">Simpan Paket</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Mobile Bottom Navbar -->
      <BottomNavbar class="md:hidden" />
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
