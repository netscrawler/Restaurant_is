<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from './stores/auth';
import { useCartStore } from './stores/cart';
import LoginForm from './components/LoginForm.vue';
import CartModal from './components/CartModal.vue';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const cartStore = useCartStore();
const showLogin = ref(false);
const showCart = ref(false);
const router = useRouter();

const handleLogin = (user: any) => {
  authStore.login(user);
  showLogin.value = false;
  router.push('/profile');
};

const logout = () => {
  authStore.logout();
};

const handleCheckout = () => {
  if (!authStore.isAuthenticated) {
    alert('Для оформления заказа необходимо войти в систему');
    showCart.value = false;
    showLogin.value = true;
    return;
  }
  
  // TODO: Реализовать оформление заказа
  alert('Функция оформления заказа будет добавлена позже');
  showCart.value = false;
};

onMounted(() => {
  authStore.initAuth();
  cartStore.initCart();
});
</script>

<template>
  <div class="min-h-screen bg-primary-50 text-gray-900">
    <header class="bg-primary-500 text-white p-4 shadow">
      <div class="container mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-bold">Меню ресторана</h1>
        <nav class="flex items-center space-x-4">
          <!-- Кнопка корзины -->
          <button 
            @click="showCart = true"
            class="relative bg-white text-primary-600 px-3 py-1 rounded text-sm hover:bg-gray-100 transition-colors"
          >
            🛒 Корзина
            <span v-if="cartStore.totalItems > 0" class="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
              {{ cartStore.totalItems }}
            </span>
          </button>
          <router-link
            v-if="authStore.isAuthenticated"
            to="/profile"
            class="bg-white text-primary-600 px-3 py-1 rounded text-sm hover:bg-gray-100 transition-colors"
          >
            Профиль
          </router-link>
          <div v-if="authStore.isAuthenticated" class="flex items-center space-x-4">
            <span class="text-sm">
              Привет, {{ authStore.user?.name || authStore.user?.phone || 'Пользователь' }}!
            </span>
            <button 
              @click="logout"
              class="bg-white text-primary-600 px-3 py-1 rounded text-sm hover:bg-gray-100 transition-colors"
            >
              Выйти
            </button>
          </div>
          <div v-else>
            <button 
              @click="showLogin = true"
              class="bg-white text-primary-600 px-3 py-1 rounded text-sm hover:bg-gray-100 transition-colors"
            >
              Войти
            </button>
          </div>
        </nav>
      </div>
    </header>
    <main class="container mx-auto py-8">
      <router-view />
    </main>

    <!-- Модальное окно входа -->
    <div v-if="showLogin" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold">Вход в систему</h2>
          <button @click="showLogin = false" class="text-gray-500 hover:text-gray-700">
            ✕
          </button>
        </div>
        <LoginForm @login="handleLogin" />
      </div>
    </div>

    <!-- Модальное окно корзины -->
    <CartModal 
      :is-open="showCart" 
      @close="showCart = false"
      @checkout="handleCheckout"
    />
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
