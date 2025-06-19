<template>
  <div class="admin-panel">
    <div class="admin-header">
      <nav class="admin-nav">
        <button v-for="tab in tabs" :key="tab" :class="{active: currentTab === tab}" @click="currentTab = tab">
          {{ tabNames[tab] }}
        </button>
      </nav>
      <button class="logout-btn" @click="logout">Выйти</button>
    </div>
    <div class="admin-content">
      <div v-if="currentTab === 'staff'">
        <h2>Сотрудники</h2>
        <button class="add-btn" @click="showAddStaff = !showAddStaff">{{ showAddStaff ? 'Отмена' : 'Добавить сотрудника' }}</button>
        <form v-if="showAddStaff" class="add-form" @submit.prevent="addStaff">
          <input v-model="newStaff.work_email" type="email" placeholder="Email" required />
          <input v-model="newStaff.position" type="text" placeholder="Должность" required />
          <input v-model="newStaff.roles" type="text" placeholder="Роли (через запятую)" required />
          <button type="submit">Сохранить</button>
        </form>
        <table class="admin-table">
          <thead>
            <tr>
              <th>ID</th><th>Email</th><th>Должность</th><th>Роли</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="staff in staffList" :key="staff.id">
              <td>{{ staff.id }}</td>
              <td>{{ staff.work_email }}</td>
              <td>{{ staff.position }}</td>
              <td>{{ staff.roles.join(', ') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else-if="currentTab === 'dishes'">
        <h2>Блюда</h2>
        <button class="add-btn" @click="showAddDish = !showAddDish">{{ showAddDish ? 'Отмена' : 'Добавить блюдо' }}</button>
        <form v-if="showAddDish" class="add-form" @submit.prevent="addDish">
          <input v-model="newDish.name" type="text" placeholder="Название" required />
          <input v-model="newDish.description" type="text" placeholder="Описание" required />
          <input v-model.number="newDish.price" type="number" placeholder="Цена" required />
          <input v-model="newDish.category_name" type="text" placeholder="Категория" required />
          <select v-model="newDish.available">
            <option :value="true">Доступно</option>
            <option :value="false">Не доступно</option>
          </select>
          <button type="submit">Сохранить</button>
        </form>
        <table class="admin-table">
          <thead>
            <tr>
              <th>ID</th><th>Название</th><th>Описание</th><th>Цена</th><th>Категория</th><th>Доступно</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="dish in dishes" :key="dish.id">
              <td>{{ dish.id }}</td>
              <td>{{ dish.name }}</td>
              <td>{{ dish.description }}</td>
              <td>{{ dish.price }}</td>
              <td>{{ dish.category_name }}</td>
              <td>{{ dish.available ? 'Да' : 'Нет' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else-if="currentTab === 'orders'">
        <h2>Заказы</h2>
        <button class="add-btn" @click="showAddOrder = !showAddOrder">{{ showAddOrder ? 'Отмена' : 'Добавить заказ' }}</button>
        <form v-if="showAddOrder" class="add-form" @submit.prevent="addOrder">
          <input v-model="newOrder.user_id" type="text" placeholder="ID пользователя" required />
          <input v-model.number="newOrder.total_amount" type="number" placeholder="Сумма" required />
          <select v-model="newOrder.status">
            <option v-for="s in ORDER_STATUSES" :key="s.value" :value="s.value">{{ s.label }}</option>
          </select>
          <input v-model="newOrder.delivery_address" type="text" placeholder="Адрес" required />
          <button type="submit">Сохранить</button>
        </form>
        <table class="admin-table">
          <thead>
            <tr>
              <th>ID</th><th>Пользователь</th><th>Сумма</th><th>Статус</th><th>Адрес</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="order in orders" :key="order.id">
              <td>{{ order.id }}</td>
              <td>{{ order.user_id }}</td>
              <td>{{ order.total_amount }}</td>
              <td>
                <select v-model="order.status" @change="updateOrderStatus(order, order.status)">
                  <option v-for="s in ORDER_STATUSES" :key="s.value" :value="s.value">{{ s.label }}</option>
                </select>
              </td>
              <td>{{ order.delivery_address }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const staffList = ref([
  { id: '1', work_email: 'admin@site.com', position: 'Администратор', roles: ['admin'] },
  { id: '2', work_email: 'staff@site.com', position: 'Официант', roles: ['staff'] }
])
const showAddStaff = ref(false)
const newStaff = ref({ work_email: '', position: '', roles: '' })
function addStaff() {
  staffList.value.push({
    id: String(Date.now()),
    work_email: newStaff.value.work_email,
    position: newStaff.value.position,
    roles: newStaff.value.roles.split(',').map(r => r.trim())
  })
  newStaff.value = { work_email: '', position: '', roles: '' }
  showAddStaff.value = false
}

onMounted(() => {
  const token = localStorage.getItem('access_token')
  if (!token) {
    router.replace('/login')
  }
})

function logout() {
  localStorage.removeItem('access_token')
  router.replace('/login')
}

const tabs = ['staff', 'dishes', 'orders']
const tabNames = {
  staff: 'Сотрудники',
  dishes: 'Блюда',
  orders: 'Заказы'
}
const currentTab = ref('staff')

const dishes = ref([
  { id: '1', name: 'Борщ', description: 'Традиционный суп', price: 250, category_name: 'Супы', available: true },
  { id: '2', name: 'Пицца', description: 'Сырная', price: 500, category_name: 'Пицца', available: false }
])
const showAddDish = ref(false)
const newDish = ref({ name: '', description: '', price: 0, category_name: '', available: true })
function addDish() {
  dishes.value.push({
    id: String(Date.now()),
    name: newDish.value.name,
    description: newDish.value.description,
    price: newDish.value.price,
    category_name: newDish.value.category_name,
    available: newDish.value.available
  })
  newDish.value = { name: '', description: '', price: 0, category_name: '', available: true }
  showAddDish.value = false
}

const orders = ref([
  { id: '1', user_id: '2', total_amount: 750, status: 'ORDER_STATUS_READY', delivery_address: 'ул. Ленина, 1' },
  { id: '2', user_id: '1', total_amount: 250, status: 'ORDER_STATUS_COOKING', delivery_address: 'ул. Мира, 5' }
])
const showAddOrder = ref(false)
const newOrder = ref({ user_id: '', total_amount: 0, status: '', delivery_address: '' })
function addOrder() {
  orders.value.push({
    id: String(Date.now()),
    user_id: newOrder.value.user_id,
    total_amount: newOrder.value.total_amount,
    status: newOrder.value.status,
    delivery_address: newOrder.value.delivery_address
  })
  newOrder.value = { user_id: '', total_amount: 0, status: '', delivery_address: '' }
  showAddOrder.value = false
}

/** СТАТУСЫ ЗАКАЗА **/
const ORDER_STATUSES = [
  { value: 'ORDER_STATUS_CREATED', label: 'Получен' },
  { value: 'ORDER_STATUS_COOKING', label: 'Готовится' },
  { value: 'ORDER_STATUS_READY', label: 'Готов' },
  { value: 'ORDER_STATUS_DELIVERED', label: 'Доставляется' },
  { value: 'ORDER_STATUS_CANCELLED', label: 'Отменён' }
]

function getStatusLabel(status) {
  const found = ORDER_STATUSES.find(s => s.value === status)
  return found ? found.label : status
}

function updateOrderStatus(order, newStatus) {
  order.status = newStatus
}
</script>

<style>
body {
  background: #f3f6fb !important;
}
</style>

<style scoped>
.admin-panel {
  width: 90vw;
  min-width: 1200px;
  max-width: 1800px;
  min-height: 900px;
  margin: 40px auto;
  padding: 60px 0;
  background: #f3f6fb;
  border-radius: 18px;
  box-shadow: 0 12px 40px rgba(44, 62, 80, 0.18);
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-x: auto;
}
.admin-content, .admin-header {
  width: 100%;
  min-width: 1000px;
  max-width: 1600px;
  margin: 0 auto;
}

/* Шапка */
.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  margin-bottom: 2rem;
}

/* Навигация */
.admin-nav {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}

.admin-nav button {
  background: none;
  border: none;
  font-size: 1.6rem;
  font-weight: 700;
  color: #3358e6;
  padding: 1.3rem 2.8rem;
  border-radius: 10px;
  cursor: pointer;
  transition: 0.2s;
}

.admin-nav button.active,
.admin-nav button:hover {
  background: #3358e6;
  color: #fff;
}

/* Контент */
.admin-content {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center; /* центрируем */
  text-align: center;
}

.admin-content h2 {
  font-size: 2.8rem;
  font-weight: 800;
  color: #223;
  margin-bottom: 2.5rem;
}

/* Таблица */
.admin-table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 4px 18px rgba(44, 62, 80, 0.10);
  margin-top: 1rem;
}

.admin-table th,
.admin-table td {
  padding: 2rem 2rem;
  font-size: 1.6rem;
  text-align: center; /* центрируем содержимое таблицы */
}

.admin-table th {
  background: #e3edff;
  font-weight: 700;
  color: #223;
  border-bottom: 2px solid #dbeafe;
}

.admin-table tr:not(:last-child) td {
  border-bottom: 1px solid #e3edff;
}

.admin-table tr:hover td {
  background: #f9fbff;
}

/* Кнопка выхода */
.logout-btn {
  background: #e74c3c;
  color: white;
  border: none;
  font-weight: 700;
  border-radius: 8px;
  padding: 1.2rem 2.8rem;
  cursor: pointer;
  transition: 0.2s;
  font-size: 1.4rem;
  box-shadow: 0 2px 8px rgba(44, 62, 80, 0.10);
}

.logout-btn:hover {
  background: #c0392b;
}

.admin-panel input,
.admin-panel textarea,
.admin-panel select {
  color: #222;
  background: #fff;
}

.admin-table th, .admin-table td {
  color: #222 !important;
}

.add-btn {
  background: #3358e6;
  color: #fff;
  font-size: 1.1rem;
  font-weight: 700;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 2rem;
  margin-bottom: 1.5rem;
  cursor: pointer;
  transition: background 0.2s;
}
.add-btn:hover {
  background: #223;
}
.add-form {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-start;
}
.add-form input, .add-form select {
  padding: 0.7rem 1.2rem;
  font-size: 1.1rem;
  border: 1.5px solid #dbeafe;
  border-radius: 8px;
  color: #222;
  background: #fff;
}
.add-form button[type="submit"] {
  background: #42b883;
  color: #fff;
  font-weight: 700;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 2rem;
  cursor: pointer;
  transition: background 0.2s;
}
.add-form button[type="submit"]:hover {
  background: #2c8c5a;
}

.admin-table select {
  padding: 0.5rem 1rem;
  font-size: 1.1rem;
  border-radius: 8px;
  border: 1.5px solid #dbeafe;
  color: #222;
  background: #fff;
}
</style> 