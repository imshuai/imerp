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
        <!-- 法定代表人姓名（可点击查看详情） -->
        <el-table-column label="法定代表人" width="120">
          <template #default="{ row }">
            <el-link v-if="row.representative" type="primary" @click="handleViewPerson(row.representative)">
              {{ row.representative.name }}
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <!-- 法定代表人电话（可点击复制） -->
        <el-table-column label="法人电话" width="130">
          <template #default="{ row }">
            <el-link v-if="row.representative" type="primary" @click="handleCopy(row.representative.phone)">
              {{ row.representative.phone }}
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <!-- 法定代表人密码（明文显示，可点击复制） -->
        <el-table-column label="法人密码" width="120">
          <template #default="{ row }">
            <el-link v-if="row.representative" type="primary" @click="handleCopy(row.representative.password)">
              {{ row.representative.password }}
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>
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
        <!-- 法定代表人 -->
        <el-form-item label="法定代表人">
          <el-select
            v-model="form.representative_id"
            filterable
            remote
            reserve-keyword
            placeholder="请输入姓名搜索"
            :remote-method="searchRepresentatives"
            :loading="representativeLoading"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="item in representativeOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <!-- 投资人（多选） -->
        <el-form-item label="投资人">
          <el-select
            v-model="investorIds"
            filterable
            remote
            multiple
            reserve-keyword
            placeholder="请输入姓名搜索"
            :remote-method="searchInvestors"
            :loading="investorLoading"
            style="width: 100%"
          >
            <el-option
              v-for="item in investorOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <!-- 服务人员（多选） -->
        <el-form-item label="服务人员">
          <el-select
            v-model="servicePersonIds"
            filterable
            remote
            multiple
            reserve-keyword
            placeholder="请输入姓名搜索"
            :remote-method="searchServicePersons"
            :loading="servicePersonLoading"
            style="width: 100%"
          >
            <el-option
              v-for="item in servicePersonOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 人员详情弹窗 -->
    <el-dialog v-model="personDialogVisible" title="人员详情" width="500px" @close="handlePersonDialogClose">
      <el-form :model="personForm" label-width="100px">
        <el-form-item label="姓名">
          <el-input v-model="personForm.name" />
        </el-form-item>
        <el-form-item label="类型">
          <el-input v-model="personForm.type" disabled />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="personForm.phone" />
        </el-form-item>
        <el-form-item label="身份证号">
          <el-input v-model="personForm.id_card" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="personForm.password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="personDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePerson">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { getCustomers, createCustomer, updateCustomer, deleteCustomer } from '@/api/customers'
import { getPeople } from '@/api/people'
import type { Customer } from '@/api/customers'
import type { Person } from '@/api/people'
import { smartCopy } from '@/utils/clipboard'

const loading = ref(false)
const tableData = ref<Customer[]>([])
const dialogVisible = ref(false)
const personDialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

// 人员搜索相关
const representativeOptions = ref<Person[]>([])
const investorOptions = ref<Person[]>([])
const servicePersonOptions = ref<Person[]>([])
const representativeLoading = ref(false)
const investorLoading = ref(false)
const servicePersonLoading = ref(false)
const investorIds = ref<number[]>([])
const servicePersonIds = ref<number[]>([])
const personForm = reactive<Partial<Person>>({})

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
  registered_capital: 0,
  representative_id: undefined
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

// 搜索法定代表人
const searchRepresentatives = async (query: string) => {
  if (!query) {
    representativeOptions.value = []
    return
  }
  representativeLoading.value = true
  try {
    const res = await getPeople({
      keyword: query,
      type: '法定代表人'
    })
    representativeOptions.value = res.items
  } finally {
    representativeLoading.value = false
  }
}

// 搜索投资人
const searchInvestors = async (query: string) => {
  if (!query) {
    investorOptions.value = []
    return
  }
  investorLoading.value = true
  try {
    const res = await getPeople({
      keyword: query,
      type: '投资人'
    })
    investorOptions.value = res.items
  } finally {
    investorLoading.value = false
  }
}

// 搜索服务人员
const searchServicePersons = async (query: string) => {
  if (!query) {
    servicePersonOptions.value = []
    return
  }
  servicePersonLoading.value = true
  try {
    const res = await getPeople({
      keyword: query,
      type: '服务人员'
    })
    servicePersonOptions.value = res.items
  } finally {
    servicePersonLoading.value = false
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
    registered_capital: 0,
    representative_id: undefined
  })
  investorIds.value = []
  servicePersonIds.value = []
  representativeOptions.value = []
  investorOptions.value = []
  servicePersonOptions.value = []
  dialogVisible.value = true
}

const handleEdit = (row: Customer) => {
  isEdit.value = true
  Object.assign(form, row)
  // 设置已选中的投资人和服务人员
  if (row.investor_list) {
    investorIds.value = row.investor_list.map((p: Person) => p.id)
  }
  if (row.service_persons) {
    servicePersonIds.value = row.service_persons.map((p: Person) => p.id)
  }
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
    // 构建投资人JSON
    const investors = investorIds.value.map(id => ({
      person_id: id,
      share_ratio: 0
    }))

    const submitData = {
      ...form,
      investors: JSON.stringify(investors),
      service_person_ids: servicePersonIds.value.join(',')
    }

    if (isEdit.value) {
      await updateCustomer(form.id!, submitData)
      ElMessage.success('更新成功')
    } else {
      await createCustomer(submitData)
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

// 查看人员详情
const handleViewPerson = (person: Person) => {
  Object.assign(personForm, person)
  personDialogVisible.value = true
}

// 保存人员信息
const handleSavePerson = async () => {
  // TODO: 实现人员信息更新API
  ElMessage.success('人员信息已保存')
  personDialogVisible.value = false
  loadData()
}

const handlePersonDialogClose = () => {
  Object.assign(personForm, {})
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
