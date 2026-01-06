<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>代理记账ERP系统</h2>
        </div>
      </template>

      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-autocomplete
            v-model="loginForm.username"
            :fetch-suggestions="querySearch"
            placeholder="请输入用户名"
            style="width: 100%"
            @keyup.enter="handleLogin"
            @select="handleSelect"
          >
            <template #default="{ item }">
              <div class="user-item">
                <span>{{ item.value }}</span>
                <span class="user-role">{{ getRoleName(item.role) }}</span>
              </div>
            </template>
          </el-autocomplete>
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleLogin" style="width: 100%" :loading="loading">
            登录
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 提示信息 -->
      <el-alert
        title="超级用户: admin/admin，其他用户默认密码: 123456"
        type="info"
        :closable="false"
        style="margin-top: 20px"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '@/api/auth'
import { setToken, setStoredUser } from '@/api/auth'

const router = useRouter()
const loginFormRef = ref()
const loading = ref(false)
const userSuggestions = ref<any[]>([])

// 登录表单
const loginForm = reactive({
  username: '',
  password: ''
})

// 验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 加载用户列表
const loadUsers = async () => {
  try {
    const response = await fetch('/api/auth/users')
    const result = await response.json()
    if (result.code === 0 && result.data) {
      userSuggestions.value = result.data.map((user: any) => ({
        value: user.username,
        role: user.role
      }))
    }
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

// 搜索建议
const querySearch = (queryString: string, cb: any) => {
  const results = queryString
    ? userSuggestions.value.filter(user => user.value.toLowerCase().includes(queryString.toLowerCase()))
    : userSuggestions.value
  cb(results)
}

// 选择用户
const handleSelect = (item: any) => {
  loginForm.username = item.value
}

// 获取角色名称
const getRoleName = (role: string) => {
  const roleMap: Record<string, string> = {
    'super_admin': '超级管理员',
    'manager': '管理员',
    'service_person': '服务人员'
  }
  return roleMap[role] || role
}

// 登录
const handleLogin = async () => {
  const valid = await loginFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true

  try {
    const res = await login({
      username: loginForm.username,
      password: loginForm.password
    })

    // 保存 token
    setToken(res.token)

    // 获取用户信息
    await fetchUserInfo()

    // 检查是否需要修改密码
    if (res.must_change_password) {
      ElMessage.warning('首次登录请修改密码')
      router.push('/change-password')
      return
    }

    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (error: any) {
    ElMessage.error(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const token = localStorage.getItem('erp_token')
    const user = await (await fetch('/api/user/me', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })).json()
    if (user.code === 0) {
      setStoredUser(user.data)
      return user.data
    }
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 450px;
}

.card-header {
  text-align: center;
}

.card-header h2 {
  margin: 0;
  color: #409EFF;
}

.user-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.user-role {
  font-size: 12px;
  color: #909399;
}
</style>
