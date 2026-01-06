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
      <el-table :data="tableData" border stripe v-loading="loading" style="width: 100%">
        <el-table-column label="公司名称" min-width="180">
          <template #default="{ row }">
            <el-link type="primary" @click="handleViewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="tax_number" label="税号" min-width="170" />
        <el-table-column label="法定代表人" min-width="100">
          <template #default="{ row }">
            {{ row.representative?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="法人电话" min-width="120">
          <template #default="{ row }">
            <el-link v-if="row.representative?.phone" type="primary" @click="handleCopy(row.representative.phone)">
              {{ row.representative.phone }}
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="法人密码" min-width="110">
          <template #default="{ row }">
            <el-link v-if="row.representative?.password" type="primary" @click="handleCopy(row.representative.password)">
              {{ row.representative.password }}
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="taxpayer_type" label="纳税人类型" min-width="120">
          <template #default="{ row }">
            <el-tag v-if="row.taxpayer_type" :type="row.taxpayer_type === '一般纳税人' ? 'success' : 'warning'" size="small">
              {{ row.taxpayer_type }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="tax_office" label="税务所" min-width="120" show-overflow-tooltip />
        <el-table-column prop="tax_administrator" label="管理员" min-width="100" show-overflow-tooltip />
        <el-table-column prop="tax_administrator_phone" label="管理员电话" min-width="120" />
        <el-table-column label="操作" min-width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="info" size="small" @click="handleViewDetail(row)">详情</el-button>
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
      width="850px"
      @close="handleDialogClose"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="140px">
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
        <el-form-item label="执照登记日">
          <el-date-picker v-model="form.license_registration_date" type="date" placeholder="请选择" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="form.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item label="注册资本">
          <el-input-number v-model="form.registered_capital" :min="0" style="width: 280px" />
          <span style="margin-left: 10px">元</span>
        </el-form-item>

        <!-- 税务信息 -->
        <el-divider content-position="left">税务信息</el-divider>
        <el-form-item label="税务登记日">
          <el-date-picker v-model="form.tax_registration_date" type="date" placeholder="请选择" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="税务所">
          <el-input v-model="form.tax_office" placeholder="请输入税务所" />
        </el-form-item>
        <el-form-item label="税务管理员">
          <el-input v-model="form.tax_administrator" placeholder="请输入税务管理员" />
        </el-form-item>
        <el-form-item label="税务管理员联系电话">
          <el-input v-model="form.tax_administrator_phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="纳税人类型">
          <el-select v-model="form.taxpayer_type" placeholder="请选择纳税人类型" style="width: 100%">
            <el-option label="一般纳税人" value="一般纳税人" />
            <el-option label="小规模纳税人" value="小规模纳税人" />
          </el-select>
        </el-form-item>

        <!-- 法定代表人 -->
        <el-divider content-position="left">法定代表人</el-divider>
        <el-form :model="representativeForm" label-width="100px" style="padding: 0 20px; background: #f5f5f5; padding: 15px; border-radius: 4px;">
          <el-form-item label="姓名">
            <el-autocomplete
              v-model="representativeForm.name"
              :fetch-suggestions="searchRepresentativesAuto"
              placeholder="请输入姓名搜索或直接填写"
              :trigger-on-focus="false"
              clearable
              style="width: 100%"
              @select="(item: any) => handleRepresentativeAutoSelect(item)"
              @clear="handleRepresentativeClear"
            >
              <template #default="{ item }">
                <div>{{ item.name }}</div>
                <div style="font-size: 12px; color: #999;">{{ item.phone }}</div>
              </template>
            </el-autocomplete>
          </el-form-item>
          <el-form-item label="电话">
            <el-input v-model="representativeForm.phone" placeholder="请输入电话" />
          </el-form-item>
          <el-form-item label="身份证号">
            <el-input v-model="representativeForm.id_card" placeholder="请输入身份证号" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="representativeForm.password" placeholder="请输入密码" />
          </el-form-item>
        </el-form>

        <!-- 投资人 -->
        <el-divider content-position="left">投资人</el-divider>
        <el-button type="dashed" style="width: 100%; margin-bottom: 15px" @click="handleAddInvestor">
          <el-icon><Plus /></el-icon> 添加投资人
        </el-button>
        <el-card v-for="(investor, index) in investorsForm" :key="index" style="margin-bottom: 10px">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span>投资人 {{ index + 1 }}</span>
              <el-button type="danger" size="small" text @click="handleRemoveInvestor(index)">删除</el-button>
            </div>
          </template>
          <el-form :model="investor" label-width="100px">
            <el-form-item label="姓名">
              <el-autocomplete
                v-model="investor.name"
                :fetch-suggestions="(q: string) => searchInvestorsForItemAuto(q, index)"
                placeholder="请输入姓名搜索或直接填写"
                :trigger-on-focus="false"
                clearable
                style="width: 100%"
                @select="(item: any) => handleInvestorAutoSelect(item, index)"
                @clear="handleInvestorClear(index)"
              >
                <template #default="{ item }">
                  <div>{{ item.name }}</div>
                  <div style="font-size: 12px; color: #999;">{{ item.phone }}</div>
                </template>
              </el-autocomplete>
            </el-form-item>
            <el-form-item label="电话">
              <el-input v-model="investor.phone" placeholder="请输入电话" />
            </el-form-item>
            <el-form-item label="身份证号">
              <el-input v-model="investor.id_card" placeholder="请输入身份证号" />
            </el-form-item>
            <el-form-item label="投资比例">
              <el-input-number v-model="investor.share_ratio" :min="0" :max="100" :precision="2" style="width: 200px" />
              <span style="margin-left: 10px">%</span>
            </el-form-item>
            <!-- 出资信息 -->
            <el-divider content-position="left" style="margin: 10px 0;">出资记录</el-divider>
            <el-button type="dashed" size="small" style="width: 100%; margin-bottom: 10px" @click="handleAddInvestmentRecord(index)">
              <el-icon><Plus /></el-icon> 添加出资记录
            </el-button>
            <div v-for="(record, rIndex) in investor.investment_records" :key="rIndex" style="display: flex; gap: 10px; margin-bottom: 10px; align-items: center;">
              <el-date-picker
                v-model="record.date"
                type="date"
                placeholder="出资日期"
                value-format="YYYY-MM-DD"
                style="flex: 1"
              />
              <el-input-number v-model="record.amount" :min="0" :precision="2" placeholder="出资金额" style="flex: 1" />
              <span style="flex-shrink: 0">元</span>
              <el-button type="danger" size="small" text @click="handleRemoveInvestmentRecord(index, rIndex)">删除</el-button>
            </div>
          </el-form>
        </el-card>

        <!-- 服务人员 -->
        <el-form-item label="服务人员">
          <el-select
            v-model="servicePersonIds"
            multiple
            placeholder="请选择服务人员"
            style="width: 100%"
            @focus="loadServicePersons"
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

    <!-- 客户详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" title="客户详情" width="700px" @close="handleDetailDialogClose">
      <el-descriptions :column="2" border v-if="currentCustomer">
        <el-descriptions-item label="公司名称">{{ currentCustomer.name }}</el-descriptions-item>
        <el-descriptions-item label="税号">{{ currentCustomer.tax_number }}</el-descriptions-item>
        <el-descriptions-item label="客户类型">{{ currentCustomer.type }}</el-descriptions-item>
        <el-descriptions-item label="执照登记日">{{ currentCustomer.license_registration_date || '-' }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ currentCustomer.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="注册资本">{{ currentCustomer.registered_capital ? currentCustomer.registered_capital.toLocaleString() + ' 元' : '-' }}</el-descriptions-item>
        <el-descriptions-item label="地址" :span="2">{{ currentCustomer.address || '-' }}</el-descriptions-item>

        <!-- 税务信息 -->
        <el-descriptions-item label="税务登记日">{{ currentCustomer.tax_registration_date || '-' }}</el-descriptions-item>
        <el-descriptions-item label="纳税人类型">{{ currentCustomer.taxpayer_type || '-' }}</el-descriptions-item>
        <el-descriptions-item label="税务所">{{ currentCustomer.tax_office || '-' }}</el-descriptions-item>
        <el-descriptions-item label="税务管理员">{{ currentCustomer.tax_administrator || '-' }}</el-descriptions-item>
        <el-descriptions-item label="税务管理员电话" :span="2">{{ currentCustomer.tax_administrator_phone || '-' }}</el-descriptions-item>

        <!-- 法定代表人 -->
        <el-descriptions-item label="法定代表人">{{ currentCustomer.representative?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="法人电话">{{ currentCustomer.representative?.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="法人身份证">{{ currentCustomer.representative?.id_card || '-' }}</el-descriptions-item>
        <el-descriptions-item label="法人密码">{{ currentCustomer.representative?.password || '-' }}</el-descriptions-item>

        <!-- 投资人列表 -->
        <el-descriptions-item label="投资人" :span="2">
          <div v-if="customerInvestors.length > 0">
            <div v-for="(inv, idx) in customerInvestors" :key="idx" style="margin-bottom: 8px;">
              <span style="font-weight: 500;">{{ inv.name }}</span>
              <span style="margin-left: 10px;">持股比例: {{ inv.share_ratio }}%</span>
              <div v-if="inv.investment_records && inv.investment_records.length > 0" style="margin-left: 20px; font-size: 12px; color: #666;">
                出资记录:
                <span v-for="(rec, rIdx) in inv.investment_records" :key="rIdx">
                  {{ rec.date }} {{ rec.amount }}元{{ rIdx < inv.investment_records.length - 1 ? '；' : '' }}
                </span>
              </div>
            </div>
          </div>
          <span v-else>-</span>
        </el-descriptions-item>

        <!-- 服务人员 -->
        <el-descriptions-item label="服务人员" :span="2">
          <div v-if="currentCustomer.service_persons && currentCustomer.service_persons.length > 0">
            <el-tag v-for="person in currentCustomer.service_persons" :key="person.id" style="margin-right: 5px;">
              {{ person.name }}
            </el-tag>
          </div>
          <span v-else>-</span>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleEditFromDetail">编辑</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { getCustomers, createCustomer, updateCustomer, deleteCustomer } from '@/api/customers'
import { getPeople, createPerson, updatePerson } from '@/api/people'
import type { Customer } from '@/api/customers'
import type { Person } from '@/api/people'
import { smartCopy, debounce } from '@/utils/clipboard'

// 出资记录接口
interface InvestmentRecord {
  date: string
  amount: number
}

// 投资人表单接口
interface InvestorForm {
  id?: number
  name: string
  phone: string
  id_card: string
  share_ratio: number
  investment_records: InvestmentRecord[]
}

const loading = ref(false)
const tableData = ref<Customer[]>([])
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
const currentCustomer = ref<Customer | null>(null)
const customerInvestors = ref<any[]>([])

// 人员搜索相关
const servicePersonOptions = ref<Person[]>([])
const servicePersonIds = ref<number[]>([])

// 法定代表人详细信息表单
const representativeForm = reactive<Partial<Person>>({
  id: undefined,
  name: '',
  phone: '',
  id_card: '',
  password: ''
})

// 投资人表单列表
const investorsForm = ref<InvestorForm[]>([])

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
  license_registration_date: undefined,
  tax_registration_date: undefined,
  tax_office: '',
  tax_administrator: '',
  tax_administrator_phone: '',
  taxpayer_type: ''
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

// 搜索法定代表人（autocomplete使用） - 带防抖
const searchRepresentativesDebounced = debounce(async (queryString: string, cb: any) => {
  if (!queryString) {
    cb([])
    return
  }
  try {
    const res = await getPeople({
      keyword: queryString
    })
    cb(res.items)
  } catch (error) {
    cb([])
  }
}, 300)

const searchRepresentativesAuto = (queryString: string, cb: any) => {
  searchRepresentativesDebounced(queryString, cb)
}

// 选择法定代表人（autocomplete）
const handleRepresentativeAutoSelect = (item: any) => {
  Object.assign(representativeForm, {
    id: item.id,
    name: item.name,
    phone: item.phone || '',
    id_card: item.id_card || '',
    password: item.password || ''
  })
}

// 清空法定代表人选择
const handleRepresentativeClear = () => {
  Object.assign(representativeForm, {
    id: undefined,
    name: '',
    phone: '',
    id_card: '',
    password: ''
  })
}

// 为特定投资人索引搜索（autocomplete使用） - 带防抖
const searchInvestorsDebouncedMap = new Map<number, ReturnType<typeof debounce>>()

const searchInvestorsForItemAuto = async (queryString: string, index: number) => {
  if (!queryString) {
    return []
  }

  if (!searchInvestorsDebouncedMap.has(index)) {
    searchInvestorsDebouncedMap.set(index, debounce(async (qs: string) => {
      try {
        const res = await getPeople({
          keyword: qs
        })
        return res.items
      } catch (error) {
        return []
      }
    }, 300))
  }

  return await searchInvestorsDebouncedMap.get(index)!(queryString)
}

// 选择投资人（autocomplete）
const handleInvestorAutoSelect = (item: any, index: number) => {
  const investor = investorsForm.value[index]
  investor.id = item.id
  investor.name = item.name
  investor.phone = item.phone || ''
  investor.id_card = item.id_card || ''
}

// 清空投资人选择
const handleInvestorClear = (index: number) => {
  const investor = investorsForm.value[index]
  investor.id = undefined
}

// 添加投资人
const handleAddInvestor = () => {
  investorsForm.value.push({
    id: undefined,
    name: '',
    phone: '',
    id_card: '',
    share_ratio: 0,
    investment_records: []
  })
}

// 删除投资人
const handleRemoveInvestor = (index: number) => {
  investorsForm.value.splice(index, 1)
}

// 添加出资记录
const handleAddInvestmentRecord = (investorIndex: number) => {
  investorsForm.value[investorIndex].investment_records.push({
    date: '',
    amount: 0
  })
}

// 删除出资记录
const handleRemoveInvestmentRecord = (investorIndex: number, recordIndex: number) => {
  investorsForm.value[investorIndex].investment_records.splice(recordIndex, 1)
}

// 加载所有服务人员（静态下拉）
const loadServicePersons = async () => {
  if (servicePersonOptions.value.length > 0) return
  try {
    const res = await getPeople({ is_service_person: true })
    servicePersonOptions.value = res.items
  } catch (error) {
    console.error('加载服务人员失败:', error)
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

// 查看客户详情
const handleViewDetail = (row: Customer) => {
  currentCustomer.value = row

  // 处理投资人数据
  if (row.investors) {
    try {
      const investorInfos = JSON.parse(row.investors)
      customerInvestors.value = investorInfos.map((info: any) => {
        const person = row.investor_list?.find((p: Person) => p.id === info.person_id)
        return {
          ...info,
          name: person?.name || '',
          investment_records: info.investment_records || []
        }
      })
    } catch (e) {
      console.error('解析投资人数据失败:', e)
      customerInvestors.value = []
    }
  } else {
    customerInvestors.value = []
  }

  detailDialogVisible.value = true
}

// 从详情页打开编辑
const handleEditFromDetail = () => {
  detailDialogVisible.value = false
  if (currentCustomer.value) {
    handleEdit(currentCustomer.value)
  }
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
    license_registration_date: undefined,
    tax_registration_date: undefined,
    tax_office: '',
    tax_administrator: '',
    tax_administrator_phone: '',
    taxpayer_type: ''
  })
  servicePersonIds.value = []
  Object.assign(representativeForm, {
    id: undefined,
    name: '',
    phone: '',
    id_card: '',
    password: ''
  })
  investorsForm.value = []
  dialogVisible.value = true
}

const handleEdit = (row: Customer) => {
  isEdit.value = true
  Object.assign(form, row)

  // 服务人员
  if (row.service_persons) {
    servicePersonIds.value = row.service_persons.map((p: Person) => p.id)
  } else {
    servicePersonIds.value = []
  }

  // 法定代表人
  if (row.representative) {
    Object.assign(representativeForm, {
      id: row.representative.id,
      name: row.representative.name,
      phone: row.representative.phone || '',
      id_card: row.representative.id_card || '',
      password: row.representative.password || ''
    })
  } else {
    Object.assign(representativeForm, {
      id: undefined,
      name: '',
      phone: '',
      id_card: '',
      password: ''
    })
  }

  // 投资人 - 从investors字段解析
  investorsForm.value = []
  if (row.investors) {
    try {
      const investorInfos = JSON.parse(row.investors)
      for (const info of investorInfos) {
        let personData: Partial<Person> = {
          name: '',
          phone: '',
          id_card: ''
        }
        if (row.investor_list && row.investor_list.length > 0) {
          const person = row.investor_list.find((p: Person) => p.id === info.person_id)
          if (person) {
            personData = person
          }
        }
        investorsForm.value.push({
          id: info.person_id,
          name: personData.name || '',
          phone: personData.phone || '',
          id_card: personData.id_card || '',
          share_ratio: info.share_ratio || 0,
          investment_records: info.investment_records || []
        })
      }
    } catch (e) {
      console.error('解析投资人数据失败:', e)
    }
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
    const idCardToPersonId = new Map<string, number>()

    // 处理法定代表人
    let representativeId = form.representative_id
    if (representativeForm.name) {
      if (representativeForm.id) {
        await updatePerson(representativeForm.id, {
          name: representativeForm.name,
          phone: representativeForm.phone,
          id_card: representativeForm.id_card,
          password: representativeForm.password
        })
        representativeId = representativeForm.id
        if (representativeForm.id_card) {
          idCardToPersonId.set(representativeForm.id_card, representativeForm.id)
        }
      } else {
        if (representativeForm.id_card) {
          const existingRes = await getPeople({ keyword: representativeForm.id_card })
          const existing = existingRes.items.find((p: Person) => p.id_card === representativeForm.id_card)
          if (existing) {
            await updatePerson(existing.id, {
              name: representativeForm.name,
              phone: representativeForm.phone,
              password: representativeForm.password
            })
            representativeId = existing.id
            idCardToPersonId.set(representativeForm.id_card, existing.id)
          } else {
            const newPerson = await createPerson({
              name: representativeForm.name,
              phone: representativeForm.phone,
              id_card: representativeForm.id_card,
              password: representativeForm.password
            })
            representativeId = newPerson.id
            idCardToPersonId.set(representativeForm.id_card, newPerson.id)
          }
        } else {
          const newPerson = await createPerson({
            name: representativeForm.name,
            phone: representativeForm.phone,
            id_card: representativeForm.id_card,
            password: representativeForm.password
          })
          representativeId = newPerson.id
        }
      }
    }

    // 处理投资人
    const investors = []

    // 如果填写了法定代表人但投资人列表为空，自动将法定代表人添加为投资人（100%）
    if (representativeId && investorsForm.value.length === 0) {
      investors.push({
        person_id: representativeId,
        share_ratio: 100,
        investment_records: []
      })
    }

    for (const investor of investorsForm.value) {
      if (investor.name) {
        let personId = investor.id

        if (investor.id) {
          await updatePerson(investor.id, {
            name: investor.name,
            phone: investor.phone,
            id_card: investor.id_card
          })
          personId = investor.id
          if (investor.id_card) {
            idCardToPersonId.set(investor.id_card, investor.id)
          }
        } else {
          if (investor.id_card && idCardToPersonId.has(investor.id_card)) {
            personId = idCardToPersonId.get(investor.id_card)!
          } else if (investor.id_card) {
            const existingRes = await getPeople({ keyword: investor.id_card })
            const existing = existingRes.items.find((p: Person) => p.id_card === investor.id_card)
            if (existing) {
              await updatePerson(existing.id, {
                name: investor.name,
                phone: investor.phone
              })
              personId = existing.id
              idCardToPersonId.set(investor.id_card, existing.id)
            } else {
              const newPerson = await createPerson({
                name: investor.name,
                phone: investor.phone,
                id_card: investor.id_card
              })
              personId = newPerson.id
              idCardToPersonId.set(investor.id_card!, newPerson.id)
            }
          } else {
            const newPerson = await createPerson({
              name: investor.name,
              phone: investor.phone,
              id_card: investor.id_card
            })
            personId = newPerson.id
          }
        }

        investors.push({
          person_id: personId,
          share_ratio: investor.share_ratio || 0,
          investment_records: investor.investment_records || []
        })
      }
    }

    const submitData = {
      ...form,
      representative_id: representativeId,
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

const handleDetailDialogClose = () => {
  currentCustomer.value = null
  customerInvestors.value = []
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
