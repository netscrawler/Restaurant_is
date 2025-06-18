import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import MenuList from './components/MenuList.vue'
import ProfileView from './components/ProfileView.vue'

const pinia = createPinia()

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: MenuList },
    { path: '/profile', component: ProfileView },
  ],
})

createApp(App)
  .use(pinia)
  .use(router)
  .mount('#app')
