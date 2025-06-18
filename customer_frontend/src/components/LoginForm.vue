<template>
  <div class="min-h-screen flex items-center justify-center bg-primary-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Вход в систему
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Введите номер телефона для получения кода подтверждения
        </p>
      </div>
      
      <!-- Шаг 1: Ввод телефона -->
      <div v-if="step === 1" class="mt-8 space-y-6">
        <div>
          <label for="phone" class="block text-sm font-medium text-gray-700">
            Номер телефона
          </label>
          <input
            id="phone"
            v-model="phone"
            type="tel"
            required
            class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
            placeholder="+7 (999) 123-45-67"
          />
        </div>
        
        <div>
          <button
            @click="sendCode"
            :disabled="loading || !phone"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading">Отправка...</span>
            <span v-else>Отправить код</span>
          </button>
        </div>
        
        <div v-if="error" class="text-red-600 text-sm text-center">
          {{ error }}
        </div>
      </div>
      
      <!-- Шаг 2: Ввод кода -->
      <div v-if="step === 2" class="mt-8 space-y-6">
        <div>
          <label for="code" class="block text-sm font-medium text-gray-700">
            Код подтверждения
          </label>
          <input
            id="code"
            v-model="code"
            type="text"
            required
            class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
            placeholder="123456"
            maxlength="6"
          />
        </div>
        
        <div class="flex space-x-3">
          <button
            @click="step = 1"
            class="flex-1 py-2 px-4 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Назад
          </button>
          <button
            @click="confirmCode"
            :disabled="loading || !code"
            class="flex-1 py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading">Проверка...</span>
            <span v-else>Войти</span>
          </button>
        </div>
        
        <div v-if="error" class="text-red-600 text-sm text-center">
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { apiService } from '../services/api';

const emit = defineEmits<{ login: [user: any] }>();

const step = ref(1);
const phone = ref('');
const code = ref('');
const loading = ref(false);
const error = ref('');

const sendCode = async () => {
  if (!phone.value) return;
  loading.value = true;
  error.value = '';
  try {
    await apiService.loginInit({ phone: phone.value });
    step.value = 2;
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Ошибка при отправке кода';
  } finally {
    loading.value = false;
  }
};

const confirmCode = async () => {
  if (!code.value) return;
  loading.value = true;
  error.value = '';
  try {
    const response = await apiService.loginConfirm({
      phone: phone.value,
      code: code.value
    });
    localStorage.setItem('access_token', response.access_token);
    localStorage.setItem('refresh_token', response.refresh_token);
    localStorage.setItem('user', JSON.stringify(response.user));
    emit('login', response.user);
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Неверный код подтверждения';
  } finally {
    loading.value = false;
  }
};
</script> 