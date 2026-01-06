<template>
  <div class="payments-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>收款管理</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="起始日期">
          <el-date-picker
            v-model="searchForm.start_date"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker
            v-model="searchForm.end_date"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
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
          新增收款
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData" border stripe v-loading="loading">
        <el-table-column prop="payment_date" label="收款日期" width="110">
          <template #default="{ row }">
            {{ formatDate(row.payment_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="customer.name" label="客户名称" width="200" />
        <el-table-column prop="agreement.agreement_number" label="协议编号" width="150" />
        <el-table-column prop="amount" label="收款金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount.toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="payment_method" label="收款方式" width="100" />
        <el-table-column prop="period" label="费用期间" width="100" />
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
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
      :title="isEdit ? '编辑收款' : '新增收款'"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="关联客户" prop="customer_id">
          <el-select
            v-model="form.customer_id"
            filterable
            remote
            reserve-keyword
            placeholder="请输入客户名称或税号搜索"
            :remote-method="searchCustomers"
            :loading="customerLoading"
            style="width: 100%"
          >
            <el-option
              v-for="item in customerOptions"
              :key="item.id"
              :label="`${item.name} - ${item.tax_number}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="收款金额" prop="amount">
          <el-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="收款日期" prop="payment_date">
          <el-date-picker
            v-model="form.payment_date"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="收款方式" prop="payment_method">
          <el-select v-model="form.payment_method" placeholder="请选择收款方式" style="width: 100%">
            <el-option label="转账" value="转账" />
            <el-option label="现金" value="现金" />
            <el-option label="支票" value="支票" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="费用期间" prop="period">
          <el-input v-model="form.period" placeholder="如: 2024-01" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" placeholder="请输入备注" />
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
import { getPayments, createPayment, updatePayment, deletePayment } from '@/api/payments'
import { getCustomers } from '@/api/customers'
import { debounce } from '@/utils/clipboard'
import type { Payment } from '@/api/payments'
import type { Customer } from '@/api/customers'

const loading = ref(false)
const tableData = ref<Payment[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

// 客户搜索相关
const customerOptions = ref<Customer[]>([])
const customerLoading = ref(false)

const searchForm = reactive({
  start_date: '',
  end_date: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive<Partial<Payment>>({
  customer_id: undefined,
  agreement_id: undefined,
  amount: 0,
  payment_date: '',
  payment_method: '转账',
  period: '',
  remark: ''
})

const rules = {
  customer_id: [{ required: true, message: '请输入客户ID', trigger: 'blur' }],
  amount: [{ required: true, message: '请输入收款金额', trigger: 'blur' }],
  payment_date: [{ required: true, message: '请选择收款日期', trigger: 'change' }],
  payment_method: [{ required: true, message: '请选择收款方式', trigger: 'change' }]
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

// 搜索客户 - 带防抖
const searchCustomersDebounced = debounce(async (query: string) => {
  if (!query) {
    customerOptions.value = []
    return
  }
  customerLoading.value = true
  try {
    const res = await getCustomers({ keyword: query })
    customerOptions.value = res.items
  } finally {
    customerLoading.value = false
  }
}, 300)

const searchCustomers = (query: string) => {
  searchCustomersDebounced(query)
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getPayments({
      start_date: searchForm.start_date || undefined,
      end_date: searchForm.end_date || undefined
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
  searchForm.start_date = ''
  searchForm.end_date = ''
  handleSearch()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, {
    customer_id: undefined,
    agreement_id: undefined,
    amount: 0,
    payment_date: '',
    payment_method: '转账',
    period: '',
    remark: ''
  })
  customerOptions.value = []
  dialogVisible.value = true
}

const handleEdit = (row: Payment) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = (row: Payment) => {
  ElMessageBox.confirm('确定要删除该收款记录吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    await deletePayment(row.id)
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
      await updatePayment(form.id!, form)
      ElMessage.success('更新成功')
    } else {
      await createPayment(form)
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
.payments-page {
  height: 100%;
}

.search-form {
  margin-bottom: 20px;
}

.toolbar {
  margin-bottom: 20px;
}
</style>
