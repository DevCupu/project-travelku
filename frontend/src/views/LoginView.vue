<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api' // <-- Menggunakan instance API sentral

const router = useRouter()
const email = ref('')
const password = ref('')
const isSubmitting = ref(false)
const errorMessage = ref('')

const handleLogin = async () => {
  isSubmitting.value = true
  errorMessage.value = ''
  
  try {
    // URL dasar sudah diatur di services/api.js, jadi cukup '/auth/login'
    const response = await api.post('/auth/login', {
      email: email.value,
      password: password.value
    })
    
    const payload = response.data.data || response.data
    localStorage.setItem('auth_token', payload.token)
    if (payload.user) {
      localStorage.setItem('user', JSON.stringify(payload.user))
    }
    
    router.push('/dashboard')
  } catch (error) {
    if (error.response && error.response.data && error.response.data.message) {
      errorMessage.value = error.response.data.message
    } else {
      errorMessage.value = 'Failed to connect to server. Ensure backend is running.'
    }
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="card-header">
        <div class="brand-logo">
          <!-- 24 Visa Makassar Travel Logo -->
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true" role="presentation" focusable="false" style="display: block; height: 32px; width: 32px; fill: var(--primary);">
            <path d="M21 16v-2l-8-5V3.5c0-.83-.67-1.5-1.5-1.5S10 2.67 10 3.5V9l-8 5v2l8-2.5V19l-2 1.5V22l3.5-1 3.5 1v-1.5L13 19v-5.5l8 2.5z"/>
          </svg>
        </div>
        <h1>Log in to 24 Visa Makassar</h1>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="welcome-text">Sistem Manajemen 24 Visa</div>

        <div v-if="errorMessage" class="error-alert">
          {{ errorMessage }}
        </div>
        
        <div class="input-group">
          <div class="input-wrapper">
            <label for="email">Email</label>
            <input 
              id="email" 
              type="email" 
              v-model="email" 
              placeholder="Email address" 
              required
            >
          </div>
          <div class="input-wrapper input-wrapper-bottom">
            <label for="password">Password</label>
            <input 
              id="password" 
              type="password" 
              v-model="password" 
              placeholder="Password" 
              required
            >
          </div>
        </div>

        <p class="terms">
          Dengan masuk ke sistem, Anda menyetujui Kebijakan Privasi & SOP Internal PT 24 Visa Makassar. <a href="#">Kebijakan Privasi</a>
        </p>

        <button type="submit" class="btn-primary" :disabled="isSubmitting">
          <span v-if="!isSubmitting">Continue</span>
          <span v-else class="loader"></span>
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: var(--bg-page);
}

.login-card {
  width: 100%;
  max-width: 568px;
  border: 1px solid var(--border-color-light);
  border-radius: var(--radius-lg);
  padding: 0 0 32px 0;
  background: var(--bg-page);
  box-shadow: 0 8px 28px rgba(0,0,0,0.05);
}

.card-header {
  text-align: center;
  padding: 24px;
  border-bottom: 1px solid var(--border-color-light);
  margin-bottom: 24px;
  position: relative;
}

.brand-logo {
  display: none;
}

h1 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-main);
  margin: 0;
}

.login-form {
  padding: 0 24px;
}

.welcome-text {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-main);
  margin-bottom: 24px;
}

.error-alert {
  background-color: #fff0f3;
  color: #c13515;
  padding: 12px;
  border-radius: var(--radius-md);
  margin-bottom: 16px;
  font-size: 14px;
  font-weight: 500;
  border: 1px solid #ffb3c1;
}

.input-group {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  margin-bottom: 16px;
  overflow: hidden;
  background: var(--bg-page);
}

.input-wrapper {
  padding: 8px 12px;
  position: relative;
  transition: box-shadow 0.2s ease;
}

.input-wrapper:focus-within {
  box-shadow: inset 0 0 0 2px var(--border-color-focus);
  z-index: 1;
}

.input-wrapper-bottom {
  border-top: 1px solid var(--border-color);
}

label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 2px;
}

input {
  width: 100%;
  border: none;
  outline: none;
  font-size: 16px;
  color: var(--text-main);
  background: transparent;
  padding: 2px 0;
}

.terms {
  font-size: 12px;
  color: var(--text-main);
  line-height: 1.5;
  margin-bottom: 24px;
}

.terms a {
  color: var(--text-main);
  font-weight: 600;
  text-decoration: underline;
}

.btn-primary {
  width: 100%;
  background-color: var(--primary);
  color: white;
  font-size: 16px;
  font-weight: 700;
  padding: 16px 32px;
  border-radius: var(--radius-lg);
  transition: background-color 0.2s ease;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 52px;
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--primary-hover);
}

.btn-primary:disabled {
  background-color: var(--primary-disabled);
  cursor: not-allowed;
}

.loader {
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
