import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
  },
  {
    path: '/cases',
    name: 'Cases',
    component: () => import('@/views/CasesView.vue'),
  },
  {
    path: '/cases/jwt-auth',
    name: 'JwtAuth',
    component: () => import('@/views/cases/JwtAuthView.vue'),
  },
  {
    path: '/cases/redis-lock',
    name: 'RedisLock',
    component: () => import('@/views/cases/RedisLockView.vue'),
  },
  {
    path: '/roadmap',
    name: 'Roadmap',
    component: () => import('@/views/RoadmapView.vue'),
  },
  {
    path: '/docs-link',
    name: 'DocsLink',
    component: () => import('@/views/DocsLinkView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    return savedPosition || { top: 0 }
  },
})

let initialized = false
router.beforeEach(async () => {
  if (!initialized) {
    const authStore = useAuthStore()
    await authStore.initAuth()
    initialized = true
  }
})

export default router
