<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>代理记账ERP系统</h2>
        </div>
      </template>

      <el-tabs v-model="activeTab" class="login-tabs">
        <!-- 超级用户/管理员登录 -->
        <el-tab-pane label="管理员登录" name="admin">
          <el-form :model="loginForm" :rules="rules" ref="adminFormRef" label-width="80px">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="loginForm.username" placeholder="请输入用户名" />
            </el-form-item>

            <el-form-item label="密码" prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                show-password
                @keyup.enter="handleAdminLogin"
              />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="handleAdminLogin" style="width: 100%" :loading="loading">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 服务人员登录 -->
        <el-tab-pane label="服务人员登录" name="service">
          <el-form :model="serviceForm" :rules="serviceRules" ref="serviceFormRef" label-width="80px">
            <el-form-item label="选择人员" prop="person_id">
              <el-select
                v-model="serviceForm.person_id"
                placeholder="请选择您的姓名"
                style="width: 100%"
                filterable
                :loading="loadingPeople"
              >
                <el-option
                  v-for="person in servicePeople"
                  :key="person.id"
                  :label="person.name"
                  :value="person.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="handleServiceLogin" style="width: 100%" :loading="loading">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <!-- 提示信息 -->
      <el-alert
        v-if="activeTab === 'admin'"
        title="超级用户: admin/admin"
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
import { getPeople } from '@/api/people'

const router = useRouter()

const activeTab = ref('admin')
const loading = ref(false)
const loadingPeople = ref(false)
const adminFormRef = ref()
const serviceFormRef = ref()

// 管理员登录表单
const loginForm = reactive({
  username: 'admin',
  password: 'admin'
})

// 服务人员登录表单
const serviceForm = reactive({
  person_id: undefined as number | undefined
})

// 服务人员列表
const servicePeople = ref<any[]>([])

// 验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const serviceRules = {
  person_id: [
    { required: true, message: '请选择您的姓名', trigger: 'change' }
  ]
}

// 加载服务人员列表
const loadServicePeople = async () => {
  loadingPeople.value = true
  try {
    const data = await getPeople({ is_service_person: true })
    servicePeople.value = data.items || []
  } catch (error) {
    console.error('Failed to load service people:', error)
  } finally {
    loadingPeople.value = false
  }
}

// 管理员登录
const handleAdminLogin = async () => {
  const valid = await adminFormRef.value?.validate().catch(() => false)
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

// 服务人员登录
const handleServiceLogin = async () => {
  const valid = await serviceFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true

  try {
    const res = await login({
      person_id: serviceForm.person_id
    })

    // 保存 token
    setToken(res.token)

    // 获取用户信息
    await fetchUserInfo()

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
    const user = await (await fetch('/api/user/me', {
      headers: {
        'Authorization': `Bearer ${getToken()}`
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

// 获取 token（用于 fetchUserInfo）
const getToken = () => {
  return localStorage.getItem('erp_token')
}

onMounted(() => {
  loadServicePeople()
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

.login-tabs {
  margin-top: 10px;
}
</style>
