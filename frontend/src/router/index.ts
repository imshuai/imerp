import { createRouter, createWebHistory } from 'vue-router'
import { isLoggedIn } from '@/api/auth'

// 角色类型
type UserRole = 'super_admin' | 'service_person'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { title: '登录', requiresAuth: false }
    },
    {
      path: '/change-password',
      name: 'ChangePassword',
      component: () => import('@/views/ChangePassword.vue'),
      meta: { title: '修改密码', requiresAuth: true }
    },
    {
      path: '/',
      component: () => import('@/components/Layout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue'),
          meta: { title: '首页', icon: 'Odometer' }
        },
        {
          path: 'service-personnel',
          name: 'ServicePersonnel',
          component: () => import('@/views/ServicePersonnel/index.vue'),
          meta: { title: '服务人员管理', icon: 'User' }
        },
        {
          path: 'customers',
          name: 'Customers',
          component: () => import('@/views/Customers/index.vue'),
          meta: { title: '客户管理', icon: 'OfficeBuilding' }
        },
        {
          path: 'tasks',
          name: 'Tasks',
          component: () => import('@/views/Tasks/index.vue'),
          meta: { title: '任务管理', icon: 'List' }
        },
        {
          path: 'agreements',
          name: 'Agreements',
          component: () => import('@/views/Agreements/index.vue'),
          meta: { title: '协议管理', icon: 'Document' }
        },
        {
          path: 'payments',
          name: 'Payments',
          component: () => import('@/views/Payments/index.vue'),
          meta: { title: '收款管理', icon: 'Money' }
        },
        {
          path: 'import-export',
          name: 'ImportExport',
          component: () => import('@/views/ImportExport.vue'),
          meta: { title: '导入导出', icon: 'Download' }
        },
        {
          path: 'audit-logs',
          name: 'AuditLogs',
          component: () => import('@/views/AuditLogs/index.vue'),
          meta: { title: '审计日志', icon: 'Document' }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/dashboard'
    }
  ]
})

// 全局前置守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - ERP系统`
  }

  // 检查是否需要登录
  const requiresAuth = to.meta?.requiresAuth !== false

  if (requiresAuth && !isLoggedIn()) {
    // 需要登录但未登录，跳转到登录页
    next('/login')
    return
  }

  // 已登录用户访问登录页，跳转到首页
  if (to.path === '/login' && isLoggedIn()) {
    next('/dashboard')
    return
  }

  // 检查管理员权限
  const userStr = localStorage.getItem('erp_user')
  if (userStr) {
    const user = JSON.parse(userStr)
    const role = user.role as UserRole

    // 需要超级管理员权限
    if (to.meta?.requiresSuperAdmin && role !== 'super_admin') {
      next('/dashboard')
      return
    }
  }

  next()
})

export default router
