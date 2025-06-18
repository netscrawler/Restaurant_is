import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import MenuList from './components/MenuList.vue'

const pinia = createPinia()

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: MenuList },
  ],
})

createApp(App)
  .use(pinia)
  .use(router)
  .mount('#app')
