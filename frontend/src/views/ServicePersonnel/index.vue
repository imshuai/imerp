<template>
  <div class="service-personnel-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>服务人员管理</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="姓名/电话/身份证" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 操作栏 -->
      <div class="toolbar">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增服务人员
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData" border stripe v-loading="loading">
        <el-table-column prop="name" label="姓名" min-width="120" />
        <el-table-column prop="phone" label="电话" min-width="130">
          <template #default="{ row }">
            <el-link type="primary" @click="handleCopy(row.phone)">
              {{ row.phone }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="password" label="密码" min-width="120">
          <template #default="{ row }">
            <el-link type="primary" @click="handleCopy(row.password)">
              {{ row.password }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="id_card" label="身份证号" min-width="180" />
        <el-table-column prop="customer_count" label="关联客户数量" min-width="120" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑服务人员' : '新增服务人员'"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="身份证号" prop="id_card">
          <el-input v-model="form.id_card" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" placeholder="请输入密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { getPeople, createPerson, updatePerson, deletePerson } from '@/api/people'
import type { Person } from '@/api/people'
import { smartCopy } from '@/utils/clipboard'

interface ServicePersonnelWithCount extends Person {
  customer_count: number
}

const loading = ref(false)
const tableData = ref<ServicePersonnelWithCount[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const currentPerson = ref<Person | null>(null)

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive<Partial<Person>>({
  name: '',
  phone: '',
  id_card: '',
  password: ''
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入电话', trigger: 'blur' }],
  id_card: [{ required: true, message: '请输入身份证号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getPeople({
      keyword: searchForm.keyword || undefined,
      is_service_person: true
    })
    // 计算每个服务人员的客户数量
    tableData.value = res.items.map((person: Person) => ({
      ...person,
      customer_count: person.service_customer_ids
        ? person.service_customer_ids.split(',').filter(Boolean).length
        : 0
    }))
    pagination.total = res.total
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  searchForm.keyword = ''
  handleSearch()
}

const handleAdd = () => {
  isEdit.value = false
  currentPerson.value = null
  Object.assign(form, {
    name: '',
    phone: '',
    id_card: '',
    password: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row: ServicePersonnelWithCount) => {
  isEdit.value = true
  currentPerson.value = row
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = (row: ServicePersonnelWithCount) => {
  ElMessageBox.confirm('确定要删除该服务人员吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    await deletePerson(row.id)
    ElMessage.success('删除成功')
    loadData()
  }).catch(() => {})
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      // 检查是否有任何改变
      let hasChanges = false

      // 检查各个字段是否改变
      const fieldsToCheck = ['name', 'phone', 'id_card', 'password']
      for (const field of fieldsToCheck) {
        const origVal = (currentPerson.value as any)[field]
        const newVal = (form as any)[field]
        if (origVal !== newVal) {
          // 对于 undefined 和 空值的特殊处理
          if ((origVal === undefined || origVal === null || origVal === '') &&
              (newVal === undefined || newVal === null || newVal === '')) {
            continue
          }
          hasChanges = true
          break
        }
      }

      // 如果有任何改变，才调用更新接口
      if (hasChanges) {
        await updatePerson(form.id!, form)
        ElMessage.success('更新成功')
        dialogVisible.value = false
        loadData()
      } else {
        ElMessage.info('没有修改任何内容')
        dialogVisible.value = false
      }
    } else {
      await createPerson({ ...form, is_service_person: true })
      ElMessage.success('创建成功')
      dialogVisible.value = false
      loadData()
    }
  } finally {
    submitting.value = false
  }
}

const handleDialogClose = () => {
  formRef.value?.resetFields()
  currentPerson.value = null
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  loadData()
}

const handleCopy = async (text: string) => {
  await smartCopy(text)
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.service-personnel-page {
  height: 100%;
}

.search-form {
  margin-bottom: 20px;
}

.toolbar {
  margin-bottom: 20px;
}
</style>
