import axios from 'axios'

// 1. Buat instance Axios terpusat
const api = axios.create({
  // Ambil base URL dari .env, gunakan localhost sebagai fallback (cadangan)
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000 // 10 detik
})

// 2. Request Interceptor: Otomatis menyuntikkan token ke setiap request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
}, (error) => {
  return Promise.reject(error)
})

// 3. Response Interceptor: Menangani error secara global (seperti 401 Unauthorized)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Jika token expired / tidak valid, otomatis tendang ke halaman login
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
