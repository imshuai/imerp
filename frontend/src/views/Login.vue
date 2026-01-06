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
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            @keyup.enter="handleLogin"
          />
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '@/api/auth'
import { setToken, setStoredUser, getToken } from '@/api/auth'

const router = useRouter()
const loginFormRef = ref()
const loading = ref(false)

// 登录表单
const loginForm = reactive({
  username: 'admin',
  password: 'admin'
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
</style>
