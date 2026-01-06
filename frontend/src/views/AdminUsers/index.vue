<template>
  <div class="admin-users-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>管理员列表</span>
          <el-button type="primary" size="small" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            新增管理员
          </el-button>
        </div>
      </template>

      <el-table :data="adminUsers" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="role" label="角色" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.role === 'super_admin'" type="danger">超级用户</el-tag>
            <el-tag v-else type="warning">管理员</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="关联人员" width="150">
          <template #default="{ row }">
            {{ row.person?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="must_change_password" label="需修改密码" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.must_change_password" type="warning">是</el-tag>
            <el-tag v-else type="info">否</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_password_change" label="最后修改密码时间" width="180">
          <template #default="{ row }">
            {{ row.last_password_change ? formatDateTime(row.last_password_change) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button
              v-if="row.role !== 'super_admin'"
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card style="margin-top: 20px">
      <template #header>
        <div class="card-header">
          <span>服务人员管理（设置管理员）</span>
          <el-button type="primary" size="small" @click="loadServicePeople" :loading="loadingPeople">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table :data="servicePeople" v-loading="loadingPeople" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" width="150" />
        <el-table-column prop="phone" label="电话" width="150" />
        <el-table-column prop="is_manager" label="是否管理员" width="120">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_manager"
              :loading="row.updating"
              @change="handleToggleManager(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建管理员对话框 -->
    <el-dialog v-model="createDialogVisible" title="新增管理员" width="500px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="createForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="createForm.password" type="password" show-password placeholder="请输入密码（至少6位）" />
        </el-form-item>
        <el-form-item label="选择人员" prop="person_id">
          <el-select
            v-model="createForm.person_id"
            placeholder="请选择服务人员"
            style="width: 100%"
            filterable
            :loading="loadingPeople"
          >
            <el-option
              v-for="person in availableServicePeople"
              :key="person.id"
              :label="person.name"
              :value="person.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="loading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import {
  getAdminUsers,
  getServicePeople,
  createAdminUser,
  deleteAdminUser,
  setManager
} from '@/api/admin'
import type { AdminUser } from '@/api/admin'

const loading = ref(false)
const loadingPeople = ref(false)
const adminUsers = ref<AdminUser[]>([])
const servicePeople = ref<any[]>([])
const createDialogVisible = ref(false)
const createFormRef = ref()

const createForm = reactive({
  username: '',
  password: '',
  person_id: undefined as number | undefined
})

const createRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  person_id: [{ required: true, message: '请选择服务人员', trigger: 'change' }]
}

// 可用的服务人员（未设置为管理员的）
const availableServicePeople = computed(() => {
  return servicePeople.value.filter(p => !p.is_manager)
})

// 加载管理员列表
const loadAdminUsers = async () => {
  loading.value = true
  try {
    adminUsers.value = await getAdminUsers()
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

// 加载服务人员列表
const loadServicePeople = async () => {
  loadingPeople.value = true
  try {
    const data = await getServicePeople()
    servicePeople.value = data.map((p: any) => ({ ...p, updating: false }))
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loadingPeople.value = false
  }
}

// 显示创建对话框
const showCreateDialog = () => {
  createForm.username = ''
  createForm.password = ''
  createForm.person_id = undefined
  createDialogVisible.value = true
}

// 创建管理员
const handleCreate = async () => {
  const valid = await createFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await createAdminUser({
      username: createForm.username,
      password: createForm.password,
      person_id: createForm.person_id!
    })
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadAdminUsers()
    loadServicePeople()
  } catch (error: any) {
    ElMessage.error(error.message || '创建失败')
  } finally {
    loading.value = false
  }
}

// 删除管理员
const handleDelete = async (row: AdminUser) => {
  try {
    await ElMessageBox.confirm('确认删除该管理员？', '提示', {
      type: 'warning'
    })

    loading.value = true
    await deleteAdminUser(row.id)
    ElMessage.success('删除成功')
    loadAdminUsers()
    loadServicePeople()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  } finally {
    loading.value = false
  }
}

// 切换管理员状态
const handleToggleManager = async (row: any) => {
  row.updating = true
  try {
    await setManager({
      person_id: row.id,
      is_manager: row.is_manager
    })
    ElMessage.success(row.is_manager ? '已设置为管理员' : '已取消管理员')
    loadAdminUsers()
    loadServicePeople()
  } catch (error: any) {
    // 失败时恢复原状态
    row.is_manager = !row.is_manager
    ElMessage.error(error.message || '操作失败')
  } finally {
    row.updating = false
  }
}

// 格式化日期时间
const formatDateTime = (str: string) => {
  if (!str) return '-'
  return new Date(str).toLocaleString('zh-CN')
}

onMounted(() => {
  loadAdminUsers()
  loadServicePeople()
})
</script>

<style scoped>
.admin-users-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
