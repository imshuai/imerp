<template>
  <el-container class="layout-container">
    <el-aside width="200px" class="sidebar">
      <div class="logo">
        <h2>代理记账ERP</h2>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      >
        <el-menu-item
          v-for="route in filteredMenuRoutes"
          :key="route.path"
          :index="route.path"
        >
          <el-icon v-if="route.meta?.icon">
            <component :is="route.meta.icon" />
          </el-icon>
          <span>{{ route.meta?.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="main-container">
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item>{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-icon><User /></el-icon>
              <span>{{ userName }}</span>
              <el-tag v-if="userRole" :type="getRoleTagType(userRole)" size="small" style="margin-left: 8px">
                {{ getRoleLabel(userRole) }}
              </el-tag>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item v-if="canChangePassword" @click="handleChangePassword">
                  修改密码
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { removeToken, getStoredUser, removeStoredUser } from '@/api/auth'

const router = useRouter()
const route = useRoute()

// 用户信息
const userRole = computed(() => {
  const user = getStoredUser()
  return user?.role || ''
})

const userName = computed(() => {
  const user = getStoredUser()
  if (user?.role === 'service_person') {
    return user?.name || '服务人员'
  }
  return user?.username || '管理员'
})

// 是否可以修改密码（只有 admin 和 manager 可以）
const canChangePassword = computed(() => {
  const role = userRole.value
  return role === 'super_admin' || role === 'manager'
})

// 根据角色过滤菜单
const filteredMenuRoutes = computed(() => {
  const role = userRole.value
  const routes = router.getRoutes()

  return routes.filter(r => {
    // 不显示登录页和修改密码页
    if (r.path === '/login' || r.path === '/change-password') {
      return false
    }

    // 检查是否有标题
    if (!r.meta?.title) {
      return false
    }

    // 检查权限
    if (r.meta?.requiresSuperAdmin && role !== 'super_admin') {
      return false
    }
    if (r.meta?.requiresManager && role !== 'super_admin' && role !== 'manager') {
      return false
    }

    return r.path.startsWith('/') && r.path !== '/login'
  })
})

const activeMenu = computed(() => {
  return '/' + route.path.split('/')[1]
})

const currentTitle = computed(() => {
  return route.meta?.title || '首页'
})

// 获取角色标签类型
const getRoleTagType = (role: string) => {
  if (role === 'super_admin') return 'danger'
  if (role === 'manager') return 'warning'
  return 'info'
}

// 获取角色标签文本
const getRoleLabel = (role: string) => {
  if (role === 'super_admin') return '超级用户'
  if (role === 'manager') return '管理员'
  if (role === 'service_person') return '服务人员'
  return ''
}

// 修改密码
const handleChangePassword = () => {
  router.push('/change-password')
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确认退出登录？', '提示', {
      type: 'warning'
    })

    // 清除 token 和用户信息
    removeToken()
    removeStoredUser()

    ElMessage.success('已退出登录')
    router.push('/login')
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  // 检查是否已登录
  const user = getStoredUser()
  if (!user && route.path !== '/login') {
    router.push('/login')
  }
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  color: #fff;
}

.logo {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #1f2d3d;
}

.logo h2 {
  font-size: 18px;
  color: #fff;
  margin: 0;
}

.main-container {
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.main-content {
  background-color: #f0f2f5;
  padding: 20px;
  overflow: auto;
}

:deep(.el-menu) {
  border-right: none;
}
</style>
