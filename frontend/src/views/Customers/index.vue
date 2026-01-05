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
        <el-table-column prop="registered_capital" label="注册资本" width="140">
          <template #default="{ row }">
            {{ row.registered_capital ? row.registered_capital.toLocaleString() + ' 元' : '-' }}
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
      width="800px"
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
          <el-input-number v-model="form.registered_capital" :min="0" style="width: 280px" />
          <span style="margin-left: 10px">元</span>
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

    <!-- 人员详情弹窗 -->
    <el-dialog v-model="personDialogVisible" title="人员详情" width="500px" @close="handlePersonDialogClose">
      <el-form :model="personForm" label-width="100px">
        <el-form-item label="姓名">
          <el-input v-model="personForm.name" />
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
import { getPeople, createPerson, updatePerson } from '@/api/people'
import type { Customer } from '@/api/customers'
import type { Person } from '@/api/people'
import { smartCopy, debounce } from '@/utils/clipboard'

// 投资人表单接口
interface InvestorForm {
  id?: number
  name: string
  phone: string
  id_card: string
  share_ratio: number
}

const loading = ref(false)
const tableData = ref<Customer[]>([])
const dialogVisible = ref(false)
const personDialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

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

// 人员详情表单
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

  // 获取或创建该索引的防抖函数
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
    share_ratio: 0
  })
}

// 删除投资人
const handleRemoveInvestor = (index: number) => {
  investorsForm.value.splice(index, 1)
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
        // 查找对应的人员信息
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
          share_ratio: info.share_ratio || 0
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
    // 用于记录身份证号到人员ID的映射，处理重复身份证的情况
    const idCardToPersonId = new Map<string, number>()

    // 处理法定代表人 - 创建或更新
    let representativeId = form.representative_id
    if (representativeForm.name) {
      if (representativeForm.id) {
        // 更新已有人员
        await updatePerson(representativeForm.id, {
          name: representativeForm.name,
          phone: representativeForm.phone,
          id_card: representativeForm.id_card,
          password: representativeForm.password
        })
        representativeId = representativeForm.id
        // 记录身份证号映射
        if (representativeForm.id_card) {
          idCardToPersonId.set(representativeForm.id_card, representativeForm.id)
        }
      } else {
        // 创建新人员前，先检查是否已存在相同身份证号
        if (representativeForm.id_card) {
          const existingRes = await getPeople({ keyword: representativeForm.id_card })
          const existing = existingRes.items.find((p: Person) => p.id_card === representativeForm.id_card)
          if (existing) {
            // 找到已存在的人员，更新其信息
            await updatePerson(existing.id, {
              name: representativeForm.name,
              phone: representativeForm.phone,
              password: representativeForm.password
            })
            representativeId = existing.id
            idCardToPersonId.set(representativeForm.id_card, existing.id)
          } else {
            // 没有找到，创建新人员
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
          // 没有身份证号，直接创建
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

    // 处理投资人 - 创建或更新
    const investors = []
    for (const investor of investorsForm.value) {
      if (investor.name) {
        let personId = investor.id

        if (investor.id) {
          // 更新已有人员
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
          // 创建新人员前，先检查是否已存在相同身份证号
          if (investor.id_card && idCardToPersonId.has(investor.id_card)) {
            // 已存在该身份证号（可能是刚创建的法定代表人），复用
            personId = idCardToPersonId.get(investor.id_card)!
          } else if (investor.id_card) {
            // 检查数据库中是否已存在
            const existingRes = await getPeople({ keyword: investor.id_card })
            const existing = existingRes.items.find((p: Person) => p.id_card === investor.id_card)
            if (existing) {
              // 找到已存在的人员，更新其信息
              await updatePerson(existing.id, {
                name: investor.name,
                phone: investor.phone
              })
              personId = existing.id
              idCardToPersonId.set(investor.id_card, existing.id)
            } else {
              // 没有找到，创建新人员
              const newPerson = await createPerson({
                name: investor.name,
                phone: investor.phone,
                id_card: investor.id_card
              })
              personId = newPerson.id
              idCardToPersonId.set(investor.id_card!, newPerson.id)
            }
          } else {
            // 没有身份证号，直接创建
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
          share_ratio: investor.share_ratio || 0
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

// 查看人员详情
const handleViewPerson = (person: Person) => {
  Object.assign(personForm, person)
  personDialogVisible.value = true
}

// 保存人员信息
const handleSavePerson = async () => {
  if (personForm.id) {
    await updatePerson(personForm.id, personForm)
    ElMessage.success('人员信息已保存')
    personDialogVisible.value = false
    loadData()
  }
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
