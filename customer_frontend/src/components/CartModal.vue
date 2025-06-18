<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-full max-w-2xl mx-4 max-h-[80vh] overflow-hidden flex flex-col">
      <!-- Заголовок -->
      <div class="flex justify-between items-center mb-4 pb-4 border-b">
        <h2 class="text-xl font-semibold">Корзина</h2>
        <button @click="$emit('close')" class="text-gray-500 hover:text-gray-700">
          ✕
        </button>
      </div>

      <!-- Содержимое корзины -->
      <div class="flex-1 overflow-y-auto">
        <div v-if="cartStore.isEmpty" class="text-center py-8">
          <div class="text-gray-500 text-lg mb-2">Корзина пуста</div>
          <div class="text-sm text-gray-400">Добавьте товары из меню</div>
        </div>
        
        <div v-else class="space-y-4">
          <div v-for="item in cartStore.items" :key="item.id" class="flex items-center space-x-4 p-4 border rounded-lg">
            <!-- Изображение -->
            <div class="w-16 h-16 bg-gray-200 rounded flex items-center justify-center flex-shrink-0">
              <span v-if="!item.image_url" class="text-gray-500 text-xs">Фото</span>
              <img v-else :src="item.image_url" :alt="item.name" class="w-full h-full object-cover rounded" />
            </div>
            
            <!-- Информация о товаре -->
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-gray-900 truncate">{{ item.name }}</h3>
              <p class="text-sm text-gray-500">{{ item.category }}</p>
              <p class="text-primary-600 font-semibold">{{ item.price }} ₽</p>
            </div>
            
            <!-- Управление количеством -->
            <div class="flex items-center space-x-2">
              <button 
                @click="cartStore.updateQuantity(item.id, item.quantity - 1)"
                class="w-8 h-8 bg-gray-200 rounded flex items-center justify-center hover:bg-gray-300 transition-colors"
              >
                -
              </button>
              <span class="w-8 text-center font-semibold">{{ item.quantity }}</span>
              <button 
                @click="cartStore.updateQuantity(item.id, item.quantity + 1)"
                class="w-8 h-8 bg-gray-200 rounded flex items-center justify-center hover:bg-gray-300 transition-colors"
              >
                +
              </button>
            </div>
            
            <!-- Общая цена за товар -->
            <div class="text-right min-w-[80px]">
              <p class="font-semibold text-gray-900">{{ item.price * item.quantity }} ₽</p>
            </div>
            
            <!-- Кнопка удаления -->
            <button 
              @click="cartStore.removeItem(item.id)"
              class="text-red-500 hover:text-red-700 p-1"
            >
              ✕
            </button>
          </div>
        </div>
      </div>

      <!-- Итого и кнопки -->
      <div v-if="!cartStore.isEmpty" class="border-t pt-4 mt-4">
        <div class="flex justify-between items-center mb-4">
          <div>
            <p class="text-sm text-gray-600">Товаров: {{ cartStore.totalItems }}</p>
            <p class="text-lg font-semibold text-gray-900">Итого: {{ cartStore.totalPrice }} ₽</p>
          </div>
          <div class="flex space-x-2">
            <button 
              @click="cartStore.clearCart()"
              class="px-4 py-2 border border-gray-300 rounded text-gray-700 hover:bg-gray-50 transition-colors"
            >
              Очистить
            </button>
            <button 
              @click="checkout"
              class="px-6 py-2 bg-primary-600 text-white rounded hover:bg-primary-700 transition-colors"
            >
              Оформить заказ
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useCartStore } from '../stores/cart';

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  close: [];
  checkout: [];
}>();

const cartStore = useCartStore();

const checkout = () => {
  emit('checkout');
};
</script> 