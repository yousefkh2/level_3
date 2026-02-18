import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import DatabasesView from '../views/DatabasesView.vue'
import { getToken } from '@/services/api'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/databases'
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/databases',
      name: 'databases',
      component: DatabasesView,
      meta: { requiresAuth: true } // This is how we tag a route as protected
    },
  ],
})

// Navigation guard to protect routes
router.beforeEach((to, from, next) => {
  const token = getToken()
  
  // If route requires auth and no token, redirect to login
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } 
  // If logged in and trying to access login, redirect to databases
  else if (to.path === '/login' && token) {
    next('/databases')
  } 
  // Otherwise proceed
  else {
    next()
  }
})

export default router