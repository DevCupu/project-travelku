import { ref } from 'vue'
import api from '../services/api'

// Composable function untuk membungkus seluruh logika bisnis User/Staff
export function useUsers() {
  const users = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  // Fetch daftar semua user/staf
  const fetchUsers = async () => {
    isLoading.value = true
    error.value = null
    try {
      const res = await api.get('/users')
      users.value = res.data.data?.data || []
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal mengambil data staf'
      console.error(err)
    } finally {
      isLoading.value = false
    }
  }

  // Update profil user (nama, email, phone)
  const updateProfile = async (id, payload) => {
    await api.put(`/users/${id}`, payload)
  }

  // Ganti password user
  const changePassword = async (id, payload) => {
    await api.post(`/users/${id}/change-password`, payload)
  }

  // Hapus user
  const deleteUser = async (id) => {
    await api.delete(`/users/${id}`)
  }

  return {
    users,
    isLoading,
    error,
    fetchUsers,
    updateProfile,
    changePassword,
    deleteUser
  }
}
