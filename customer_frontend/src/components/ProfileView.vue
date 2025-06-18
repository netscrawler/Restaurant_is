<template>
  <div class="max-w-lg mx-auto bg-white rounded shadow p-6">
    <h2 class="text-2xl font-bold mb-4">Профиль пользователя</h2>
    <form @submit.prevent="saveProfile" class="space-y-4">
      <div>
        <label class="block text-sm font-medium mb-1">Имя</label>
        <input v-model="form.name" class="w-full border rounded px-3 py-2" />
      </div>
      <div>
        <label class="block text-sm font-medium mb-1">Телефон</label>
        <input v-model="form.phone" class="w-full border rounded px-3 py-2" />
      </div>
      <div>
        <label class="block text-sm font-medium mb-1">Email</label>
        <input v-model="form.email" class="w-full border rounded px-3 py-2" />
      </div>
      <div class="flex space-x-2">
        <button type="submit" class="bg-primary-600 text-white px-4 py-2 rounded hover:bg-primary-700">Сохранить</button>
        <span v-if="success" class="text-green-600 self-center">Сохранено!</span>
        <span v-if="error" class="text-red-600 self-center">{{ error }}</span>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { apiService } from '../services/api';

const form = ref({ name: '', phone: '', email: '' });
const error = ref('');
const success = ref(false);

const fetchProfile = async () => {
  error.value = '';
  try {
    const user = await apiService.getProfile();
    form.value = {
      name: user.name || '',
      phone: user.phone || '',
      email: user.email || ''
    };
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Ошибка загрузки профиля';
  }
};

const saveProfile = async () => {
  error.value = '';
  success.value = false;
  try {
    await apiService.updateProfile(form.value);
    success.value = true;
    setTimeout(() => (success.value = false), 2000);
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Ошибка сохранения профиля';
  }
};

onMounted(fetchProfile);
</script> 