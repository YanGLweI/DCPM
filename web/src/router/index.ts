import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/change-password',
      name: 'ChangePassword',
      component: () => import('../views/ChangePassword.vue'),
      meta: { requiresAuth: false } // 密码过期时也可访问
    },
    {
      path: '/account',
      name: 'AccountManage',
      component: () => import('../views/AccountManage.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  
  // 需要认证但没有 token 的路由
  if (to.meta.requiresAuth && !token) {
    next('/login')
    return
  }
  
  // 已登录用户访问登录页，重定向到账号管理
  if (to.path === '/login' && token) {
    next('/account')
    return
  }
  
  next()
})

export default router
