import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

export function setupRouterGuards(router: Router) {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    // Logged in user visiting login page, redirect to admin
    if (to.meta.guestOnly && authStore.isAuthenticated) {
      return next('/admin')
    }

    // Requires auth
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
      return next('/login')
    }

    // Requires admin
    if (to.meta.requiresAdmin && !authStore.isAdmin) {
      return next('/admin')
    }

    next()
  })
}
