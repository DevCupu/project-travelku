import { ref } from 'vue'
import api from '../services/api'

export function usePaket() {
  const pakets = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  const fetchPakets = async (activeOnly = false) => {
    isLoading.value = true
    error.value = null
    try {
      const res = await api.get(`/pakets?active_only=${activeOnly}`)
      pakets.value = res.data.data || []
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal mengambil data paket wisata'
      console.error(err)
    } finally {
      isLoading.value = false
    }
  }

  const savePaket = async (id, payload) => {
    if (id) {
      await api.put(`/pakets/${id}`, payload)
    } else {
      await api.post('/pakets', payload)
    }
  }

  const deletePaket = async (id) => {
    await api.delete(`/pakets/${id}`)
  }

  return {
    pakets,
    isLoading,
    error,
    fetchPakets,
    savePaket,
    deletePaket
  }
}
