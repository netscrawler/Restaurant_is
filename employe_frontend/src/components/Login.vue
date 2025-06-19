<template>
  <div class="login-bg">
    <div class="login-container">
      <div class="login-header">
        <img src="https://cdn-icons-png.flaticon.com/512/3135/3135715.png" alt="Admin" class="admin-icon" />
        <h2>Вход в админ-панель</h2>
        <p class="login-desc">Пожалуйста, введите рабочий email и пароль администратора</p>
      </div>
      <div class="login-hint">
        <b>Внимание!</b> Первый сотрудник (админ) должен быть зарегистрирован через API (например, через Postman или Swagger). После этого вы сможете добавлять новых сотрудников через эту панель.
      </div>
      <form @submit.prevent="login" class="login-form">
        <div class="form-group">
          <label for="email">Email</label>
          <input v-model="email" id="email" type="email" placeholder="admin@company.com" required autocomplete="username" />
        </div>
        <div class="form-group">
          <label for="password">Пароль</label>
          <input v-model="password" id="password" type="password" placeholder="Введите пароль" required autocomplete="current-password" />
        </div>
        <button type="submit" :disabled="loading">{{ loading ? 'Вход...' : 'Войти' }}</button>
        <div v-if="error" class="error">{{ error }}</div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const router = useRouter()

async function login() {
  loading.value = true
  error.value = ''
  try {
    const response = await fetch('http://localhost:8080/api/v1/auth/staff/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ staff: { work_email: email.value }, password: password.value })
    })
    const data = await response.json()
    console.log('Ответ авторизации:', data)
    if (response.ok && data.access_token) {
      localStorage.setItem('access_token', data.access_token)
      if (data.user && data.user.id) {
        localStorage.setItem('user_id', data.user.id)
      }
      router.push('/admin')
    } else {
      error.value = data.error || data.message || 'Ошибка авторизации'
    }
  } catch (e) {
    error.value = 'Ошибка сети'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-bg {
  min-height: 100vh;
  background: #f4f6fa;
  display: flex;
  align-items: center;
  justify-content: center;
}
.login-container {
  background: #fff;
  width: calc(100vw + 300px);
  min-width: unset;
  max-width: unset;
  min-height: 600px;
  margin: 0;
  border-radius: 0;
  box-shadow: 0 12px 40px rgba(44, 62, 80, 0.13);
  padding: 60px 0 60px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.login-header {
  width: 100%;
  text-align: center;
  margin-bottom: 2.8rem;
}
.admin-icon {
  width: 90px;
  height: 90px;
  margin-bottom: 1rem;
  display: block;
}
.login-header h2 {
  font-size: 3.2rem;
  margin: 0.2rem 0 0.3rem 0;
  color: #223;
}
.login-desc {
  color: #888;
  font-size: 1.7rem;
  margin-top: 0.1rem;
}
.login-form {
  width: 100%;
  max-width: 520px;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.form-group {
  margin-bottom: 2.2rem;
  display: flex;
  flex-direction: column;
  width: 100%;
  align-items: center;
}
label {
  font-weight: 600;
  margin-bottom: 0.6rem;
  color: #223;
  font-size: 1.35rem;
  text-align: center;
}
input[type="email"],
input[type="password"] {
  padding: 1.4rem 1.5rem;
  border: 2px solid #dbeafe;
  border-radius: 12px;
  font-size: 1.6rem;
  outline: none;
  transition: border 0.2s, box-shadow 0.2s;
  background: #f7faff;
  color: #222;
  box-sizing: border-box;
  width: 100%;
  max-width: 420px;
  text-align: center;
}
input[type="email"]::placeholder,
input[type="password"]::placeholder {
  color: #222;
  opacity: 0.7;
}
input[type="email"]:focus,
input[type="password"]:focus {
  border-color: #4f8cff;
  background: #fff;
  box-shadow: 0 0 0 2px #e3edff;
}
button[type="submit"] {
  width: 100%;
  max-width: 420px;
  padding: 1.4rem 0;
  background: linear-gradient(90deg, #4f8cff 0%, #3358e6 100%);
  color: #fff;
  font-weight: 800;
  border: none;
  border-radius: 12px;
  font-size: 1.7rem;
  cursor: pointer;
  transition: background 0.2s, box-shadow 0.2s;
  margin-top: 0.2rem;
  box-shadow: 0 4px 18px rgba(44, 62, 80, 0.13);
  text-align: center;
}
button[type="submit"]:hover:not(:disabled) {
  background: linear-gradient(90deg, #3358e6 0%, #4f8cff 100%);
  box-shadow: 0 8px 32px rgba(44, 62, 80, 0.18);
}
button[disabled] {
  opacity: 0.7;
  cursor: not-allowed;
}
.error {
  color: #e74c3c;
  margin-top: 1.5rem;
  text-align: center;
  font-size: 1.3rem;
}
.login-hint {
  background: #e3edff;
  color: #223;
  border-radius: 8px;
  padding: 1.1rem 1.5rem;
  margin-bottom: 2.2rem;
  font-size: 1.1rem;
  text-align: center;
  box-shadow: 0 2px 8px rgba(44, 62, 80, 0.07);
}
</style> 