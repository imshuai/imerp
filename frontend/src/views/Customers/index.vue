<template>
  <div class="customers-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>客户管理</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="公司名称/税号/电话" clearable />
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
          新增客户
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData" border stripe v-loading="loading">
        <el-table-column prop="name" label="公司名称" width="200" />
        <el-table-column prop="tax_number" label="税号" width="180" />
        <el-table-column prop="type" label="客户类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="address" label="地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="registered_capital" label="注册资本" width="120">
          <template #default="{ row }">
            {{ row.registered_capital ? row.registered_capital.toLocaleString() : '-' }}
          </template>
        </el-table-column>
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
      :title="isEdit ? '编辑客户' : '新增客户'"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="公司名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入公司名称" />
        </el-form-item>
        <el-form-item label="税号" prop="tax_number">
          <el-input v-model="form.tax_number" placeholder="请输入税号" />
        </el-form-item>
        <el-form-item label="客户类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择客户类型" style="width: 100%">
            <el-option label="有限公司" value="有限公司" />
            <el-option label="个人独资企业" value="个人独资企业" />
            <el-option label="合伙企业" value="合伙企业" />
            <el-option label="个体工商户" value="个体工商户" />
          </el-select>
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item label="注册资本">
          <el-input-number v-model="form.registered_capital" :min="0" style="width: 100%" />
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
import { getCustomers, createCustomer, updateCustomer, deleteCustomer } from '@/api/customers'
import type { Customer, CustomerType } from '@/api/customers'

const loading = ref(false)
const tableData = ref<Customer[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive<Partial<Customer>>({
  name: '',
  phone: '',
  address: '',
  tax_number: '',
  type: '有限公司',
  registered_capital: 0
})

const rules = {
  name: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
  tax_number: [{ required: true, message: '请输入税号', trigger: 'blur' }],
  type: [{ required: true, message: '请选择客户类型', trigger: 'change' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getCustomers({
      keyword: searchForm.keyword || undefined
    })
    tableData.value = res.items
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
  Object.assign(form, {
    name: '',
    phone: '',
    address: '',
    tax_number: '',
    type: '有限公司',
    registered_capital: 0
  })
  dialogVisible.value = true
}

const handleEdit = (row: Customer) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = (row: Customer) => {
  ElMessageBox.confirm('确定要删除该客户吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    await deleteCustomer(row.id)
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
      await updateCustomer(form.id!, form)
      ElMessage.success('更新成功')
    } else {
      await createCustomer(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } finally {
    submitting.value = false
  }
}

const handleDialogClose = () => {
  formRef.value?.resetFields()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.customers-page {
  height: 100%;
}

.search-form {
  margin-bottom: 20px;
}

.toolbar {
  margin-bottom: 20px;
}
</style>
