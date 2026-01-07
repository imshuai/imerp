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
        <el-table-column label="税号" min-width="170">
          <template #default="{ row }">
            <span class="copy-cell" @click="handleCopy(row.tax_number)">
              {{ row.tax_number }}
              <el-icon class="copy-icon"><DocumentCopy /></el-icon>
            </span>
          </template>
        </el-table-column>
        <el-table-column label="法定代表人" min-width="100">
          <template #default="{ row }">
            {{ row.representative?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="法人电话" min-width="120">
          <template #default="{ row }">
            <span v-if="row.representative?.phone" class="copy-cell" @click="handleCopy(row.representative.phone)">
              {{ row.representative.phone }}
              <el-icon class="copy-icon"><DocumentCopy /></el-icon>
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="法人密码" min-width="110">
          <template #default="{ row }">
            <span v-if="row.representative?.password" class="copy-cell" @click="handleCopy(row.representative.password)">
              {{ row.representative.password }}
              <el-icon class="copy-icon"><DocumentCopy /></el-icon>
            </span>
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
                <div>{{ item.name }} - {{ item.id_card || '无身份证号' }}</div>
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
                :fetch-suggestions="getSearchInvestorsForItem(index)"
                placeholder="请输入姓名搜索或直接填写"
                :trigger-on-focus="false"
                clearable
                style="width: 100%"
                @select="(item: any) => handleInvestorAutoSelect(item, index)"
                @clear="handleInvestorClear(index)"
              >
                <template #default="{ item }">
                  <div>{{ item.name }} - {{ item.id_card || '无身份证号' }}</div>
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
          >
            <el-option
              v-for="item in servicePersonOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <!-- 办税人 -->
        <el-form-item label="办税人">
          <el-select
            v-model="taxAgentIds"
            multiple
            placeholder="请选择办税人"
            style="width: 100%"
          >
            <el-option
              v-for="item in taxAgentOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <!-- 新增字段 v0.4.0 -->
        <el-divider content-position="left">其他信息</el-divider>
        <el-form-item label="信用等级">
          <el-select v-model="form.credit_rating" placeholder="请选择信用等级" style="width: 100%">
            <el-option v-for="rating in creditRatingOptions" :key="rating" :label="rating" :value="rating" />
          </el-select>
        </el-form-item>
        <el-form-item label="社保号">
          <el-input v-model="form.social_security_number" placeholder="请输入社保号" />
        </el-form-item>
        <el-form-item label="渝快办密码">
          <el-input v-model="form.yukuai_ban_password" type="password" show-password placeholder="请输入渝快办密码" />
        </el-form-item>
        <el-form-item label="经营范围">
          <el-input v-model="form.business_scope" type="textarea" :rows="3" placeholder="请输入经营范围" />
        </el-form-item>

        <!-- 对公账户 -->
        <el-divider content-position="left">对公账户</el-divider>
        <el-button type="dashed" style="width: 100%; margin-bottom: 15px" @click="handleAddBankAccount">
          <el-icon><Plus /></el-icon> 添加对公账户
        </el-button>
        <el-card v-for="(account, index) in bankAccountsForm" :key="index" style="margin-bottom: 10px">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span>账户 {{ index + 1 }}</span>
              <el-button type="danger" size="small" text @click="handleRemoveBankAccount(index)">删除</el-button>
            </div>
          </template>
          <el-form :model="account" label-width="100px">
            <el-form-item label="开户银行">
              <el-input v-model="account.bank_name" placeholder="请输入开户银行" />
            </el-form-item>
            <el-form-item label="账号">
              <el-input v-model="account.account_number" placeholder="请输入账号" />
            </el-form-item>
            <el-form-item label="开户行号">
              <el-input v-model="account.bank_code" placeholder="请输入开户行号" />
            </el-form-item>
            <el-form-item label="联系电话">
              <el-input v-model="account.contact_phone" placeholder="请输入联系电话" />
            </el-form-item>
            <el-form-item label="账户类型">
              <el-select v-model="account.account_type" placeholder="请选择账户类型" style="width: 100%">
                <el-option v-for="type in accountTypeOptions" :key="type" :label="type" :value="type" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-card>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 客户详情弹窗 -->
    <el-dialog v-model="detailDialogVisible" title="客户详情" width="800px" @close="handleDetailDialogClose">
      <el-descriptions :column="2" border v-if="currentCustomer">
        <el-descriptions-item label="公司名称">{{ currentCustomer.name }}</el-descriptions-item>
        <el-descriptions-item label="税号">{{ currentCustomer.tax_number }}</el-descriptions-item>
        <el-descriptions-item label="客户类型">{{ currentCustomer.type }}</el-descriptions-item>
        <el-descriptions-item label="信用等级">
          <el-tag v-if="currentCustomer.credit_rating" :type="getCreditRatingType(currentCustomer.credit_rating)" size="small">
            {{ currentCustomer.credit_rating }}
          </el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
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

        <!-- 新增字段 v0.4.0 -->
        <el-descriptions-item label="社保号" :span="2">
          <span v-if="currentCustomer.social_security_number" class="copy-cell" @click="handleCopy(currentCustomer.social_security_number)">
            {{ currentCustomer.social_security_number }}
            <el-icon class="copy-icon"><DocumentCopy /></el-icon>
          </span>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="经营范围" :span="2">{{ currentCustomer.business_scope || '-' }}</el-descriptions-item>

        <!-- 法定代表人 -->
        <el-descriptions-item label="法定代表人">{{ currentCustomer.representative?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="法人电话">{{ currentCustomer.representative?.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="法人身份证">{{ currentCustomer.representative?.id_card || '-' }}</el-descriptions-item>
        <el-descriptions-item label="法人密码">
          <span v-if="currentCustomer.representative?.password" class="copy-cell" @click="handleCopy(currentCustomer.representative.password)">
            {{ currentCustomer.representative.password }}
            <el-icon class="copy-icon"><DocumentCopy /></el-icon>
          </span>
          <span v-else>-</span>
        </el-descriptions-item>

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

        <!-- 办税人 -->
        <el-descriptions-item label="办税人" :span="2">
          <div v-if="currentCustomer.tax_agents && currentCustomer.tax_agents.length > 0">
            <el-tag v-for="person in currentCustomer.tax_agents" :key="person.id" type="success" style="margin-right: 5px;">
              {{ person.name }}
            </el-tag>
          </div>
          <span v-else>-</span>
        </el-descriptions-item>

        <!-- 对公账户 -->
        <el-descriptions-item label="对公账户" :span="2">
          <div v-if="currentCustomer.bank_accounts && currentCustomer.bank_accounts.length > 0">
            <el-card v-for="(account, idx) in currentCustomer.bank_accounts" :key="idx" style="margin-bottom: 10px;">
              <div style="display: flex; justify-content: space-between;">
                <div>
                  <div><strong>{{ account.bank_name }}</strong></div>
                  <div style="color: #666; font-size: 12px;">
                    账号: {{ account.account_number }} | 类型: {{ account.account_type }}
                  </div>
                  <div v-if="account.contact_phone" style="color: #666; font-size: 12px;">
                    电话: {{ account.contact_phone }}
                  </div>
                </div>
              </div>
            </el-card>
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
import { DocumentCopy } from '@element-plus/icons-vue'
import { getCustomers, createCustomer, updateCustomer, deleteCustomer, getCreditRatings, getAccountTypes, type Customer, type CreditRating, type AccountType, type BankAccount } from '@/api/customers'
import { getPeople, createPerson, updatePerson } from '@/api/people'
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

// 对公账户表单接口
interface BankAccountForm {
  id?: number
  bank_name: string
  account_number: string
  bank_code: string
  contact_phone: string
  account_type: AccountType
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
const taxAgentOptions = ref<Person[]>([])
const taxAgentIds = ref<number[]>([])

// 枚举选项
const creditRatingOptions = ref<CreditRating[]>([])
const accountTypeOptions = ref<AccountType[]>([])

// 对公账户表单列表
const bankAccountsForm = ref<BankAccountForm[]>([])

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
  taxpayer_type: undefined,
  // 新增字段 v0.4.0
  tax_agent_ids: '',
  credit_rating: undefined,
  social_security_number: '',
  yukuai_ban_password: '',
  business_scope: ''
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

// 为特定投资人索引创建搜索函数
const getSearchInvestorsForItem = (_index: number) => {
  // 为每个索引创建独立的搜索函数
  const debouncedSearch = debounce(async (qs: string, cb: (results: any[]) => void) => {
    if (!qs) {
      cb([])
      return
    }
    try {
      const res = await getPeople({
        keyword: qs
      })
      cb(res.items)
    } catch (error) {
      cb([])
    }
  }, 300)

  return (queryString: string, cb: (results: any[]) => void) => {
    debouncedSearch(queryString, cb)
  }
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

const handleAdd = async () => {
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
    taxpayer_type: undefined,
    // 新增字段 v0.4.0
    tax_agent_ids: '',
    credit_rating: undefined,
    social_security_number: '',
    yukuai_ban_password: '',
    business_scope: ''
  })
  servicePersonIds.value = []
  taxAgentIds.value = []
  bankAccountsForm.value = []
  Object.assign(representativeForm, {
    id: undefined,
    name: '',
    phone: '',
    id_card: '',
    password: ''
  })
  investorsForm.value = []
  // 预先加载服务人员和办税人列表
  await Promise.all([loadServicePersons(), loadTaxAgents()])
  dialogVisible.value = true
}

const handleEdit = async (row: Customer) => {
  isEdit.value = true
  currentCustomer.value = row
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

  // 办税人 v0.4.0
  if (row.tax_agents) {
    taxAgentIds.value = row.tax_agents.map((p: Person) => p.id)
  } else {
    taxAgentIds.value = []
  }

  // 对公账户 v0.4.0
  bankAccountsForm.value = []
  if (row.bank_accounts && row.bank_accounts.length > 0) {
    bankAccountsForm.value = row.bank_accounts.map((acc: BankAccount) => ({
      id: acc.id,
      bank_name: acc.bank_name,
      account_number: acc.account_number,
      bank_code: acc.bank_code || '',
      contact_phone: acc.contact_phone || '',
      account_type: acc.account_type
    }))
  }

  // 预先加载服务人员和办税人列表
  await Promise.all([loadServicePersons(), loadTaxAgents()])
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

    // 深度比较函数
    const deepEqual = (obj1: any, obj2: any): boolean => {
      return JSON.stringify(obj1) === JSON.stringify(obj2)
    }

    // 处理法定代表人
    let representativeId = form.representative_id
    let representativeChanged = false
    let needsRepresentativeUpdate = false

    if (representativeForm.name) {
      // 构建新的法定代表人对象
      const newRepresentative = {
        name: representativeForm.name,
        phone: representativeForm.phone || '',
        id_card: representativeForm.id_card || '',
        password: representativeForm.password || ''
      }

      // 获取原始法定代表人对象
      const originalRepresentative = currentCustomer.value?.representative

      if (originalRepresentative) {
        if (representativeForm.id === originalRepresentative.id) {
          // 同一个人，检查字段是否改变
          const originalData = {
            name: originalRepresentative.name,
            phone: originalRepresentative.phone || '',
            id_card: originalRepresentative.id_card || '',
            password: originalRepresentative.password || ''
          }
          if (!deepEqual(newRepresentative, originalData)) {
            representativeChanged = true
            needsRepresentativeUpdate = true
          } else {
            representativeId = originalRepresentative.id
          }
        } else {
          // 不同的人，需要创建/关联
          needsRepresentativeUpdate = true
        }
      } else {
        // 原来没有法定代表人，现在有，需要创建
        needsRepresentativeUpdate = true
      }

      // 只有需要改变时才更新
      if (needsRepresentativeUpdate) {
        if (representativeForm.id) {
          await updatePerson(representativeForm.id, newRepresentative)
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
      } else {
        // 没有改变，使用原来的ID
        if (originalRepresentative) {
          representativeId = originalRepresentative.id
        }
      }
    }

    // 处理投资人
    const investors = []
    let investorsChanged = false
    const originalInvestorInfos = currentCustomer.value?.investors ? JSON.parse(currentCustomer.value.investors) : []

    // 如果填写了法定代表人但投资人列表为空，自动将法定代表人添加为投资人（100%）
    if (representativeId && investorsForm.value.length === 0) {
      const newInvInfo = { person_id: representativeId, share_ratio: 100, investment_records: [] }
      if (!deepEqual(newInvInfo, originalInvestorInfos[0])) {
        investorsChanged = true
      }
      investors.push(newInvInfo)
    }

    for (const investor of investorsForm.value) {
      if (investor.name) {
        let personId = investor.id
        let needsInvestorUpdate = false

        // 构建新的投资人对象
        const newInvestorData = {
          name: investor.name,
          phone: investor.phone || '',
          id_card: investor.id_card || ''
        }

        // 检查投资人信息是否改变
        if (investor.id) {
          const originalInv = currentCustomer.value?.investor_list?.find((p: Person) => p.id === investor.id)
          if (originalInv) {
            const originalData = {
              name: originalInv.name,
              phone: originalInv.phone || '',
              id_card: originalInv.id_card || ''
            }
            if (!deepEqual(newInvestorData, originalData)) {
              needsInvestorUpdate = true
            } else {
              personId = investor.id
            }
          } else {
            needsInvestorUpdate = true
          }
        } else {
          needsInvestorUpdate = true
        }

        // 只有需要改变时才更新
        if (needsInvestorUpdate) {
          investorsChanged = true
          if (investor.id) {
            await updatePerson(investor.id, newInvestorData)
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
        } else {
          // 没有改变，使用原来的ID
          if (investor.id) {
            personId = investor.id
          }
        }

        // 检查投资信息是否改变
        const invInfo = { person_id: personId, share_ratio: investor.share_ratio || 0, investment_records: investor.investment_records || [] }
        const originalInvInfo = originalInvestorInfos.find((i: any) => i.person_id === personId)
        if (!deepEqual(invInfo, originalInvInfo)) {
          investorsChanged = true
        }

        investors.push(invInfo)
      }
    }

    const submitData: any = {
      ...form,
      representative_id: representativeId,
      investors: JSON.stringify(investors),
      service_person_ids: servicePersonIds.value.join(','),
      // 新增字段 v0.4.0
      tax_agent_ids: taxAgentIds.value.join(','),
      bank_accounts: bankAccountsForm.value
    }

    if (isEdit.value) {
      // 检查客户信息是否改变
      let customerChanged = false

      // 检查服务人员是否改变
      const originalServicePersonIds = currentCustomer.value?.service_persons?.map((p: Person) => p.id).sort() || []
      const newServicePersonIds = [...servicePersonIds.value].sort()
      if (JSON.stringify(originalServicePersonIds) !== JSON.stringify(newServicePersonIds)) {
        customerChanged = true
      }

      // 检查其他字段是否改变
      const fieldsToCheck = ['name', 'phone', 'address', 'tax_number', 'type', 'registered_capital',
                               'license_registration_date', 'tax_registration_date', 'tax_office',
                               'tax_administrator', 'tax_administrator_phone', 'taxpayer_type',
                               // 新增字段 v0.4.0
                               'tax_agent_ids', 'credit_rating', 'social_security_number',
                               'yukuai_ban_password', 'business_scope']
      for (const field of fieldsToCheck) {
        const origVal = (currentCustomer.value as any)[field]
        const newVal = (submitData as any)[field]
        if (origVal !== newVal) {
          // 对于 undefined 和 空值的特殊处理
          if ((origVal === undefined || origVal === null || origVal === '') &&
              (newVal === undefined || newVal === null || newVal === '')) {
            continue
          }
          customerChanged = true
          break
        }
      }

      // 如果有任何改变，才调用更新接口
      if (customerChanged || representativeChanged || investorsChanged) {
        await updateCustomer(form.id!, submitData)
        ElMessage.success('更新成功')
      } else {
        ElMessage.info('没有修改任何内容')
      }
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
  currentCustomer.value = null
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

// 加载枚举选项
const loadEnumOptions = async () => {
  try {
    const [ratings, accountTypes] = await Promise.all([
      getCreditRatings(),
      getAccountTypes()
    ])
    creditRatingOptions.value = ratings
    accountTypeOptions.value = accountTypes
  } catch (error) {
    console.error('加载枚举选项失败:', error)
  }
}

// 加载办税人选项
const loadTaxAgents = async () => {
  try {
    const res = await getPeople()
    taxAgentOptions.value = res.items
  } catch (error) {
    console.error('加载办税人失败:', error)
  }
}

// 添加对公账户
const handleAddBankAccount = () => {
  bankAccountsForm.value.push({
    bank_name: '',
    account_number: '',
    bank_code: '',
    contact_phone: '',
    account_type: '基本户'
  })
}

// 删除对公账户
const handleRemoveBankAccount = (index: number) => {
  bankAccountsForm.value.splice(index, 1)
}

// 获取信用等级标签类型
const getCreditRatingType = (rating: CreditRating) => {
  const typeMap: Record<CreditRating, any> = {
    'A': 'success',
    'B': 'primary',
    'C': 'warning',
    'D': 'danger',
    'M': 'info'
  }
  return typeMap[rating] || ''
}

onMounted(async () => {
  await loadEnumOptions()
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

.copy-cell {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.copy-cell:hover {
  background-color: #f0f0f0;
}

.copy-cell:hover .copy-icon {
  opacity: 1;
}

.copy-icon {
  margin-left: 6px;
  font-size: 14px;
  color: #409eff;
  opacity: 0;
  transition: opacity 0.2s;
}
</style>
