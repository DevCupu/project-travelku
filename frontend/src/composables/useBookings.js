import { ref } from 'vue'
import api from '../services/api'

// Composable function untuk membungkus seluruh logika bisnis Booking
export function useBookings() {
  const bookings = ref([])
  const summary = ref({ jumlah_booking: 0, total_estimasi_pendapatan: 0 })
  const isLoading = ref(false)
  const error = ref(null)
  
  // State Pagenasi
  const totalBookings = ref(0)
  const currentPage = ref(1)
  const itemsPerPage = ref(10)

  // Fetch daftar booking & ringkasan
  const fetchBookings = async (filters) => {
    isLoading.value = true
    error.value = null
    try {
      const params = new URLSearchParams()
      if (filters?.status) params.append('status', filters.status)
      if (filters?.paket_wisata) params.append('paket_wisata', filters.paket_wisata)
      if (filters?.date_from) params.append('date_from', filters.date_from)
      if (filters?.date_to) params.append('date_to', filters.date_to)
      
      // Append pagination params
      const pageVal = filters?.page || currentPage.value
      const limitVal = filters?.limit || itemsPerPage.value
      params.append('page', pageVal)
      params.append('limit', limitVal)
      
      const [listRes, summaryRes] = await Promise.all([
        api.get(`/bookings?${params.toString()}`),
        api.get(`/bookings/summary?${params.toString()}`)
      ])
      
      bookings.value = listRes.data.data?.data || []
      totalBookings.value = listRes.data.data?.total || 0
      currentPage.value = parseInt(pageVal)
      itemsPerPage.value = parseInt(limitVal)
      summary.value = summaryRes.data.data || { jumlah_booking: 0, total_estimasi_pendapatan: 0 }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal mengambil data pemesanan'
      console.error(err)
    } finally {
      isLoading.value = false
    }
  }

  // Simpan booking (Create atau Update)
  const saveBooking = async (id, payload) => {
    if (id) {
      await api.put(`/bookings/${id}`, payload)
    } else {
      await api.post('/bookings', payload)
    }
  }

  // Hapus booking
  const deleteBooking = async (id) => {
    await api.delete(`/bookings/${id}`)
  }

  // Ubah status booking
  const updateStatus = async (id, status) => {
    await api.patch(`/bookings/${id}/status`, { status })
  }

  return {
    bookings,
    summary,
    isLoading,
    error,
    totalBookings,
    currentPage,
    itemsPerPage,
    fetchBookings,
    saveBooking,
    deleteBooking,
    updateStatus
  }
}
