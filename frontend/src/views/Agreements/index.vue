<template>
  <div class="agreements-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>协议管理</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="协议编号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="有效" value="有效" />
            <el-option label="已过期" value="已过期" />
            <el-option label="已取消" value="已取消" />
          </el-select>
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
          新增协议
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData" border stripe v-loading="loading">
        <el-table-column prop="agreement_number" label="协议编号" width="150" />
        <el-table-column prop="customer.name" label="客户名称" width="200" />
        <el-table-column prop="fee_type" label="收费类型" width="100" />
        <el-table-column prop="amount" label="服务费" width="100">
          <template #default="{ row }">
            ¥{{ row.amount.toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="start_date" label="开始日期" width="110">
          <template #default="{ row }">
            {{ formatDate(row.start_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="end_date" label="结束日期" width="110">
          <template #default="{ row }">
            {{ formatDate(row.end_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">{{ row.status }}</el-tag>
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
      :title="isEdit ? '编辑协议' : '新增协议'"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="协议编号" prop="agreement_number">
          <el-input v-model="form.agreement_number" placeholder="请输入协议编号" />
        </el-form-item>
        <el-form-item label="关联客户" prop="customer_id">
          <el-input-number v-model="form.customer_id" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="收费类型" prop="fee_type">
          <el-select v-model="form.fee_type" placeholder="请选择收费类型" style="width: 100%">
            <el-option label="月度" value="月度" />
            <el-option label="季度" value="季度" />
            <el-option label="年度" value="年度" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务费金额" prop="amount">
          <el-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="开始日期" prop="start_date">
          <el-date-picker
            v-model="form.start_date"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束日期" prop="end_date">
          <el-date-picker
            v-model="form.end_date"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
            <el-option label="有效" value="有效" />
            <el-option label="已过期" value="已过期" />
            <el-option label="已取消" value="已取消" />
          </el-select>
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
import { getAgreements, createAgreement, updateAgreement, deleteAgreement } from '@/api/agreements'
import type { Agreement, AgreementStatus, FeeType } from '@/api/agreements'

const loading = ref(false)
const tableData = ref<Agreement[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const searchForm = reactive({
  keyword: '',
  status: '' as AgreementStatus | ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive<Partial<Agreement>>({
  customer_id: 0,
  agreement_number: '',
  start_date: '',
  end_date: '',
  fee_type: '月度',
  amount: 0,
  status: '有效'
})

const rules = {
  agreement_number: [{ required: true, message: '请输入协议编号', trigger: 'blur' }],
  customer_id: [{ required: true, message: '请输入客户ID', trigger: 'blur' }],
  fee_type: [{ required: true, message: '请选择收费类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入服务费金额', trigger: 'blur' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end_date: [{ required: true, message: '请选择结束日期', trigger: 'change' }]
}

const getStatusTagType = (status: string) => {
  const map: Record<string, string> = {
    '有效': 'success',
    '已过期': 'info',
    '已取消': 'danger'
  }
  return map[status] || ''
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getAgreements({
      keyword: searchForm.keyword || undefined,
      status: searchForm.status || undefined
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
  searchForm.status = ''
  handleSearch()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, {
    customer_id: 0,
    agreement_number: '',
    start_date: '',
    end_date: '',
    fee_type: '月度',
    amount: 0,
    status: '有效'
  })
  dialogVisible.value = true
}

const handleEdit = (row: Agreement) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = (row: Agreement) => {
  ElMessageBox.confirm('确定要删除该协议吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    await deleteAgreement(row.id)
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
      await updateAgreement(form.id!, form)
      ElMessage.success('更新成功')
    } else {
      await createAgreement(form)
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
.agreements-page {
  height: 100%;
}

.search-form {
  margin-bottom: 20px;
}

.toolbar {
  margin-bottom: 20px;
}
</style>
