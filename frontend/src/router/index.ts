import { createRouter, createWebHistory } from 'vue-router'

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
      meta: { title: '登录' }
    },
    {
      path: '/',
      component: () => import('@/components/Layout.vue'),
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue'),
          meta: { title: '首页', icon: 'Odometer' }
        },
        {
          path: 'people',
          name: 'People',
          component: () => import('@/views/People/index.vue'),
          meta: { title: '人员管理', icon: 'User' }
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
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/dashboard'
    }
  ]
})

export default router
