<template>
  <div class="audit-logs-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>审计日志</span>
          <div class="filters">
            <el-button v-if="isSuperAdmin" type="danger" size="small" @click="showClearDialog" :loading="loading">
              <el-icon><Delete /></el-icon>
              清理日志
            </el-button>
            <el-button type="primary" size="small" @click="loadLogs" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="logs" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user" label="操作人" width="150">
          <template #default="{ row }">
            {{ getUserName(row) }}
          </template>
        </el-table-column>
        <el-table-column prop="action_type" label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.action_type === 'create'" type="success">新建</el-tag>
            <el-tag v-else-if="row.action_type === 'update'" type="warning">修改</el-tag>
            <el-tag v-else-if="row.action_type === 'delete'" type="danger">删除</el-tag>
            <el-tag v-else type="info">{{ row.action_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100">
          <template #default="{ row }">
            {{ getResourceTypeName(row.resource_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="resource_name" label="资源名称" width="200" show-overflow-tooltip />
        <el-table-column label="变更内容" width="250">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="showDetail(row)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column v-if="isSuperAdmin" label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > 0"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        style="margin-top: 20px; text-align: right"
      />
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="变更详情" width="700px">
      <div v-if="currentLog">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="操作人">
            {{ getUserName(currentLog) }}
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">{{ getActionTypeName(currentLog.action_type) }}</el-descriptions-item>
          <el-descriptions-item label="资源类型">{{ getResourceTypeName(currentLog.resource_type) }}</el-descriptions-item>
          <el-descriptions-item label="资源名称">{{ currentLog.resource_name }}</el-descriptions-item>
          <el-descriptions-item label="操作时间" :span="2">{{ formatDateTime(currentLog.created_at) }}</el-descriptions-item>
        </el-descriptions>

        <div v-if="currentLog.action_type === 'update' && changes.length > 0" style="margin-top: 20px">
          <h4>变更内容</h4>
          <el-table :data="changes" stripe size="small">
            <el-table-column prop="field" label="字段" width="200">
              <template #default="{ row }">
                {{ getFieldLabel(row.field) }}
              </template>
            </el-table-column>
            <el-table-column prop="oldValue" label="原值">
              <template #default="{ row }">
                <span class="old-value">{{ formatValue(row.oldValue) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="newValue" label="新值">
              <template #default="{ row }">
                <span class="new-value">{{ formatValue(row.newValue) }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-else-if="currentLog.action_type === 'create'" style="margin-top: 20px">
          <h4>新建内容</h4>
          <div class="change-content new-value">
            <pre>{{ formatJSON(currentLog.new_value) }}</pre>
          </div>
        </div>

        <div v-else-if="currentLog.action_type === 'delete'" style="margin-top: 20px">
          <h4>删除内容</h4>
          <div class="change-content">
            <pre>{{ formatJSON(currentLog.old_value) }}</pre>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 清理对话框 -->
    <el-dialog v-model="clearDialogVisible" title="清理审计日志" width="500px">
      <el-form :model="clearForm" label-width="100px">
        <el-form-item label="清理方式">
          <el-radio-group v-model="clearMode">
            <el-radio value="all">清理全部</el-radio>
            <el-radio value="range">按时间段清理</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="clearMode === 'range'" label="开始日期">
          <el-date-picker
            v-model="clearForm.start_date"
            type="date"
            placeholder="选择开始日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item v-if="clearMode === 'range'" label="结束日期">
          <el-date-picker
            v-model="clearForm.end_date"
            type="date"
            placeholder="选择结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="clearDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="handleClear" :loading="clearLoading">确认清理</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Delete } from '@element-plus/icons-vue'
import { getAuditLogs, deleteAuditLog, clearAuditLogs } from '@/api/admin'
import type { AuditLog } from '@/api/admin'
import { getStoredUser } from '@/api/auth'

interface Change {
  field: string
  oldValue: any
  newValue: any
}

const loading = ref(false)
const clearLoading = ref(false)
const logs = ref<AuditLog[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const detailDialogVisible = ref(false)
const clearDialogVisible = ref(false)
const currentLog = ref<AuditLog | null>(null)
const clearMode = ref<'all' | 'range'>('all')
const clearForm = ref({
  start_date: '',
  end_date: ''
})

// 判断是否是超级管理员
const isSuperAdmin = computed(() => {
  const user = getStoredUser()
  return user?.role === 'super_admin'
})

// 计算变更内容
const changes = computed<Change[]>(() => {
  if (!currentLog.value || currentLog.value.action_type !== 'update') {
    return []
  }

  try {
    const oldValue = JSON.parse(currentLog.value.old_value || '{}')
    const newValue = JSON.parse(currentLog.value.new_value || '{}')

    const result: Change[] = []
    const allKeys = new Set([...Object.keys(oldValue), ...Object.keys(newValue)])

    // 跳过某些不需要显示的字段
    const skipFields = ['id', 'created_at', 'updated_at', 'DeletedAt']

    allKeys.forEach(key => {
      if (skipFields.includes(key)) return

      const oldVal = oldValue[key]
      const newVal = newValue[key]

      if (JSON.stringify(oldVal) !== JSON.stringify(newVal)) {
        result.push({
          field: key,
          oldValue: oldVal,
          newValue: newVal
        })
      }
    })

    return result
  } catch {
    return []
  }
})

// 加载审计日志
const loadLogs = async () => {
  loading.value = true
  try {
    const result = await getAuditLogs({
      offset: (currentPage.value - 1) * pageSize.value,
      limit: pageSize.value
    })
    logs.value = result.items
    total.value = result.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

// 显示详情
const showDetail = (row: AuditLog) => {
  currentLog.value = row
  detailDialogVisible.value = true
}

// 删除单条日志
const handleDelete = async (row: AuditLog) => {
  try {
    await ElMessageBox.confirm('确认删除这条审计日志？', '提示', {
      type: 'warning'
    })

    await deleteAuditLog(row.id)
    ElMessage.success('删除成功')
    loadLogs()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 显示清理对话框
const showClearDialog = () => {
  clearDialogVisible.value = true
  clearMode.value = 'all'
  clearForm.value = { start_date: '', end_date: '' }
}

// 清理日志
const handleClear = async () => {
  if (clearMode.value === 'range') {
    if (!clearForm.value.start_date || !clearForm.value.end_date) {
      ElMessage.warning('请选择开始和结束日期')
      return
    }
  }

  try {
    await ElMessageBox.confirm(
      clearMode.value === 'all'
        ? '确认清理全部审计日志？此操作不可恢复！'
        : `确认清理 ${clearForm.value.start_date} 至 ${clearForm.value.end_date} 的审计日志？此操作不可恢复！`,
      '警告',
      { type: 'error' }
    )

    clearLoading.value = true
    const data = clearMode.value === 'all'
      ? {}
      : { start_date: clearForm.value.start_date, end_date: clearForm.value.end_date }

    const result = await clearAuditLogs(data)
    ElMessage.success(`清理成功，共删除 ${result.count} 条记录`)
    clearDialogVisible.value = false
    loadLogs()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '清理失败')
    }
  } finally {
    clearLoading.value = false
  }
}

// 分页
const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadLogs()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadLogs()
}

// 获取操作人名称
const getUserName = (row: AuditLog): string => {
  if (row.user) {
    if (row.user.role === 'service_person' && row.user.person) {
      return row.user.person.name
    }
    return row.user.username
  }
  return '-'
}

// 获取操作类型名称
const getActionTypeName = (type: string) => {
  const map: Record<string, string> = {
    create: '新建',
    update: '修改',
    delete: '删除'
  }
  return map[type] || type
}

// 获取资源类型名称
const getResourceTypeName = (type: string) => {
  const map: Record<string, string> = {
    customer: '客户',
    task: '任务',
    agreement: '协议',
    payment: '收款',
    person: '人员'
  }
  return map[type] || type
}

// 获取字段标签
const getFieldLabel = (field: string) => {
  const map: Record<string, string> = {
    title: '任务标题',
    description: '任务描述',
    status: '任务状态',
    priority: '优先级',
    due_date: '截止日期',
    completed_at: '完成时间',

    name: '名称',
    tax_number: '税号',
    phone: '电话',
    address: '地址',
    representative_id: '法定代表人',
    investors: '投资人',
    service_person_ids: '服务人员',
    agreement_ids: '关联协议',

    agreement_number: '协议编号',
    agreement_type: '协议类型',
    amount: '金额',
    start_date: '开始日期',
    end_date: '结束日期',

    payment_date: '收款日期',
    payment_method: '收款方式',
    notes: '备注',

    is_service_person: '是否服务人员',
    id_card: '身份证号'
  }
  return map[field] || field
}

// 格式化值
const formatValue = (val: any): string => {
  if (val === null || val === undefined) return '-'
  if (typeof val === 'boolean') return val ? '是' : '否'
  if (typeof val === 'object') return JSON.stringify(val)
  return String(val)
}

// 格式化 JSON
const formatJSON = (str: string) => {
  try {
    const obj = JSON.parse(str)
    return JSON.stringify(obj, null, 2)
  } catch {
    return str
  }
}

// 格式化日期时间
const formatDateTime = (str: string) => {
  if (!str) return '-'
  return new Date(str).toLocaleString('zh-CN')
}

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.audit-logs-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filters {
  display: flex;
  align-items: center;
  gap: 10px;
}

.change-content {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 10px;
}

.change-content.new-value {
  background: #f0f9ff;
}

.change-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
}

.old-value {
  color: #f56c6c;
  text-decoration: line-through;
}

.new-value {
  color: #67c23a;
  font-weight: 500;
}
</style>
