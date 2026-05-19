<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUsers } from '../composables/useUsers'
import TopNavbar from '../components/TopNavbar.vue'
import BottomNavbar from '../components/BottomNavbar.vue'
import Sidebar from '../components/Sidebar.vue'
import { useSidebar } from '../composables/useSidebar'

const { updateProfile, changePassword } = useUsers()
const { isCollapsed } = useSidebar()

// Ambil data user yang sedang login dari localStorage
const currentUser = ref(null)

onMounted(() => {
  const stored = localStorage.getItem('user')
  if (stored) {
    currentUser.value = JSON.parse(stored)
  }
})

// Format tanggal
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit'
  })
}

// Inisial avatar (huruf pertama nama)
const initials = computed(() => {
  if (!currentUser.value?.name) return '?'
  return currentUser.value.name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
})

// State Toast Notification Premium
const toast = ref({ show: false, message: '', type: 'success' })
const showToast = (msg, type = 'success') => {
  toast.value = { show: true, message: msg, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

// Modal Edit Profile
const showEditModal = ref(false)
const editForm = ref({ name: '', email: '', phone: '' })

const openEditModal = () => {
  editForm.value = {
    name: currentUser.value.name,
    email: currentUser.value.email,
    phone: currentUser.value.phone
  }
  showEditModal.value = true
}

const handleEditSubmit = async () => {
  try {
    await updateProfile(currentUser.value.id, editForm.value)
    // Update data lokal juga agar tampilan langsung berubah
    currentUser.value = { ...currentUser.value, ...editForm.value }
    localStorage.setItem('user', JSON.stringify(currentUser.value))
    showEditModal.value = false
    showToast('Profil berhasil diperbarui!', 'success')
  } catch (err) {
    showToast(err.response?.data?.message || 'Gagal mengupdate profil', 'error')
  }
}

// Modal Change Password
const showPasswordModal = ref(false)
const passwordForm = ref({ old_password: '', new_password: '' })
const passwordSuccess = ref(false)

const openPasswordModal = () => {
  passwordForm.value = { old_password: '', new_password: '' }
  passwordSuccess.value = false
  showPasswordModal.value = true
}

const handlePasswordSubmit = async () => {
  try {
    await changePassword(currentUser.value.id, passwordForm.value)
    passwordSuccess.value = true
    showToast('Password berhasil diubah!', 'success')
    setTimeout(() => {
      showPasswordModal.value = false
    }, 1500)
  } catch (err) {
    showToast(err.response?.data?.message || 'Gagal mengubah password', 'error')
  }
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

      <main class="max-w-4xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-1">
        <header class="page-header">
          <h1>Profil Saya</h1>
        </header>

        <div v-if="currentUser" class="profile-container">
          <!-- Profile Card -->
          <div class="profile-card">
            <div class="avatar">{{ initials }}</div>
            <div class="profile-info">
              <h2>{{ currentUser.name }}</h2>
              <span class="status-badge" :class="currentUser.is_active ? 'active' : 'inactive'">
                {{ currentUser.is_active ? 'Aktif' : 'Nonaktif' }}
              </span>
            </div>
          </div>

          <!-- Details -->
          <div class="details-card">
            <h3>Informasi Akun</h3>
            <div class="detail-row">
              <div class="detail-label">Email</div>
              <div class="detail-value">{{ currentUser.email }}</div>
            </div>
            <div class="detail-row">
              <div class="detail-label">Telepon</div>
              <div class="detail-value">{{ currentUser.phone }}</div>
            </div>
            <div class="detail-row">
              <div class="detail-label">Login Terakhir</div>
              <div class="detail-value">{{ formatDate(currentUser.last_login) }}</div>
            </div>
            <div class="detail-row">
              <div class="detail-label">Terdaftar Sejak</div>
              <div class="detail-value">{{ formatDate(currentUser.created_at) }}</div>
            </div>

            <div class="profile-actions">
              <button class="btn-primary" @click="openEditModal">Edit Profil</button>
              <button class="btn-outline" @click="openPasswordModal">Ubah Password</button>
            </div>
          </div>
        </div>

        <div v-else class="empty-state">
          <p>Data profil tidak ditemukan. Silakan login ulang.</p>
        </div>
      </main>

      <!-- Modal Edit Profile -->
      <div v-if="showEditModal" class="modal-overlay" @click.self="showEditModal = false">
        <div class="modal-content">
          <h2>Edit Profil</h2>
          <form @submit.prevent="handleEditSubmit" class="modal-form">
            <div class="form-group">
              <label>Nama Lengkap</label>
              <input type="text" v-model="editForm.name" required minlength="3" />
            </div>
            <div class="form-group">
              <label>Email</label>
              <input type="email" v-model="editForm.email" required />
            </div>
            <div class="form-group">
              <label>Nomor Telepon</label>
              <input type="text" v-model="editForm.phone" required minlength="10" />
            </div>
            <div class="modal-actions">
              <button type="button" class="btn-outline" @click="showEditModal = false">Batal</button>
              <button type="submit" class="btn-primary">Simpan</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Modal Change Password -->
      <div v-if="showPasswordModal" class="modal-overlay" @click.self="showPasswordModal = false">
        <div class="modal-content">
          <h2>Ubah Password</h2>
          <div v-if="passwordSuccess" class="success-message">
            ✅ Password berhasil diubah!
          </div>
          <form v-else @submit.prevent="handlePasswordSubmit" class="modal-form">
            <div class="form-group">
              <label>Password Lama</label>
              <input type="password" v-model="passwordForm.old_password" required minlength="6" />
            </div>
            <div class="form-group">
              <label>Password Baru</label>
              <input type="password" v-model="passwordForm.new_password" required minlength="6" />
            </div>
            <div class="modal-actions">
              <button type="button" class="btn-outline" @click="showPasswordModal = false">Batal</button>
              <button type="submit" class="btn-primary">Ubah Password</button>
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

.dashboard-layout {
  min-height: 100vh;
  background-color: var(--bg-subtle);
}

.content {
  max-width: 720px;
  margin: 0 auto;
  padding: 40px 24px;
}

.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 600;
  margin: 0;
}

/* Profile Card */
.profile-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.profile-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: 32px;
  border: 1px solid var(--border-color-light);
  display: flex;
  align-items: center;
  gap: 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.02);
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), var(--primary-hover));
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 700;
  flex-shrink: 0;
}

.profile-info h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.status-badge {
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: 12px;
  font-weight: 600;
  display: inline-block;
}

.status-badge.active {
  background-color: #d1fae5;
  color: #065f46;
}

.status-badge.inactive {
  background-color: #fee2e2;
  color: #991b1b;
}

/* Details Card */
.details-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: 32px;
  border: 1px solid var(--border-color-light);
  box-shadow: 0 2px 8px rgba(0,0,0,0.02);
}

.details-card h3 {
  margin: 0 0 24px 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-main);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid var(--border-color-light);
}

.detail-row:last-of-type {
  border-bottom: none;
}

.detail-label {
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
}

.detail-value {
  color: var(--text-main);
  font-size: 14px;
  font-weight: 600;
}

.profile-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--border-color-light);
}

/* Buttons */
.btn-primary {
  background-color: var(--primary);
  color: white;
  font-weight: 700;
  padding: 16px 32px;
  border-radius: var(--radius-lg);
  border: none;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary:hover {
  background-color: var(--primary-hover);
}

.btn-outline {
  padding: 12px 24px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  font-weight: 500;
  color: var(--text-main);
  background: white;
  cursor: pointer;
  transition: box-shadow 0.2s;
}

.btn-outline:hover {
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.empty-state {
  color: var(--text-secondary);
  font-size: 16px;
  padding: 48px;
  text-align: center;
}

/* Success Message */
.success-message {
  text-align: center;
  padding: 32px;
  font-size: 18px;
  font-weight: 600;
  color: #065f46;
}

/* Modal */
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
  max-width: 480px;
  border-radius: var(--radius-lg);
  padding: 32px;
  box-shadow: 0 10px 40px rgba(0,0,0,0.15);
}

.modal-content h2 {
  margin-top: 0;
  margin-bottom: 24px;
  font-size: 24px;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 13px;
  font-weight: 600;
}

.form-group input {
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  font-size: 14px;
  outline: none;
  font-family: inherit;
}

.form-group input:focus {
  border-color: var(--border-color-focus);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}

/* Mobile */
@media (max-width: 768px) {
  .content {
    padding: 24px 16px 80px 16px;
  }

  .profile-card {
    flex-direction: column;
    text-align: center;
  }

  .profile-actions {
    flex-direction: column;
  }

  .detail-row {
    flex-direction: column;
    gap: 4px;
  }

  .modal-content {
    padding: 24px 16px;
    margin: 16px;
    max-height: 90vh;
    overflow-y: auto;
  }
}
</style>
