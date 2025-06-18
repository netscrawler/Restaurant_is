<template>
  <div>
    <h2 class="text-xl font-semibold mb-4">Меню ресторана</h2>
    <div v-if="loading" class="text-center py-8 text-lg">Загрузка меню...</div>
    <div v-else-if="error" class="text-center py-8 text-red-600">{{ error }}</div>
    <div v-else>
      <div v-if="dishes.length === 0" class="text-center text-gray-500 py-8">
        Нет доступных блюд
      </div>
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        <div v-for="dish in dishes" :key="dish.id" class="bg-white rounded-lg shadow-md p-4 flex flex-col">
          <div class="mb-3 h-40 bg-gray-200 rounded flex items-center justify-center">
            <span class="text-gray-500">Фото блюда</span>
          </div>
          <div class="flex-1">
            <h3 class="font-bold text-lg mb-2">{{ dish.name }}</h3>
            <p class="text-gray-600 text-sm mb-3">{{ dish.description }}</p>
            <div class="text-sm text-gray-500 mb-2">Категория: {{ dish.category_name }}</div>
          </div>
          <div class="flex justify-between items-center mt-auto pt-3 border-t">
            <span class="font-semibold text-lg text-primary-600">{{ dish.price }} ₽</span>
            <button 
              class="bg-primary-500 hover:bg-primary-600 text-white px-4 py-2 rounded transition-colors"
              @click="addToCart(dish)"
            >
              В корзину
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useCartStore } from '../stores/cart';
import { apiService } from '../services/api';

interface Dish {
  id: string;
  name: string;
  description: string;
  price: number;
  category_id: number;
  category_name: string;
  image_url?: string;
  available: boolean;
  created_at: string;
  updated_at: string;
}

const authStore = useAuthStore();
const cartStore = useCartStore();

const dishes = ref<Dish[]>([]);
const loading = ref(true);
const error = ref('');

const fetchDishes = async () => {
  loading.value = true;
  error.value = '';
  try {
    const res = await apiService.getDishes({ only_available: true });
    dishes.value = res.dishes || [];
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Ошибка при загрузке меню';
  } finally {
    loading.value = false;
  }
};

const addToCart = (dish: Dish) => {
  if (!authStore.isAuthenticated) {
    alert('Для добавления в корзину необходимо войти в систему');
    return;
  }
  cartStore.addItem({
    id: dish.id,
    name: dish.name,
    price: dish.price,
    category: dish.category_name,
    image_url: dish.image_url
  });
  alert(`Добавлено в корзину: ${dish.name}`);
};

onMounted(fetchDishes);
</script> 