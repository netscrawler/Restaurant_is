import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User } from '../types/api';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const isAuthenticated = computed(() => !!user.value);

  // Инициализация из localStorage
  const initAuth = () => {
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser);
      } catch (e) {
        console.error('Ошибка при парсинге пользователя:', e);
        logout();
      }
    }
  };

  // Вход
  const login = (userData: User) => {
    user.value = userData;
  };

  // Выход
  const logout = () => {
    user.value = null;
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
  };

  return {
    user,
    isAuthenticated,
    initAuth,
    login,
    logout
  };
}); 