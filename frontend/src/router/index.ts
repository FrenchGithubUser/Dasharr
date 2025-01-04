import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router' // Type-only import
import Home from '../views/HomePage.vue'
import Register from '../views/RegisterPage.vue'

const routes: RouteRecordRaw[] = [
  { path: '/home', name: 'Home', component: Home },
  { path: '/', name: 'Register', component: Register },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
