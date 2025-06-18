import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export interface CartItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
  image_url?: string;
  category: string;
}

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([]);

  // Инициализация из localStorage
  const initCart = () => {
    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
      try {
        items.value = JSON.parse(savedCart);
      } catch (e) {
        console.error('Ошибка при загрузке корзины:', e);
        items.value = [];
      }
    }
  };

  // Сохранение в localStorage
  const saveCart = () => {
    localStorage.setItem('cart', JSON.stringify(items.value));
  };

  // Добавление товара в корзину
  const addItem = (item: Omit<CartItem, 'quantity'>) => {
    const existingItem = items.value.find(i => i.id === item.id);
    
    if (existingItem) {
      existingItem.quantity += 1;
    } else {
      items.value.push({ ...item, quantity: 1 });
    }
    
    saveCart();
  };

  // Удаление товара из корзины
  const removeItem = (itemId: string) => {
    items.value = items.value.filter(item => item.id !== itemId);
    saveCart();
  };

  // Изменение количества товара
  const updateQuantity = (itemId: string, quantity: number) => {
    const item = items.value.find(i => i.id === itemId);
    if (item) {
      if (quantity <= 0) {
        removeItem(itemId);
      } else {
        item.quantity = quantity;
        saveCart();
      }
    }
  };

  // Очистка корзины
  const clearCart = () => {
    items.value = [];
    saveCart();
  };

  // Вычисляемые свойства
  const totalItems = computed(() => {
    return items.value.reduce((sum, item) => sum + item.quantity, 0);
  });

  const totalPrice = computed(() => {
    return items.value.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  });

  const isEmpty = computed(() => items.value.length === 0);

  return {
    items,
    totalItems,
    totalPrice,
    isEmpty,
    initCart,
    addItem,
    removeItem,
    updateQuantity,
    clearCart
  };
}); 