<template>
  <div class="audit-logs-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>审计日志</span>
          <div class="filters">
            <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 150px; margin-right: 10px" @change="loadLogs">
              <el-option label="全部" value="" />
              <el-option label="待审批" value="pending" />
              <el-option label="已通过" value="approved" />
              <el-option label="已拒绝" value="rejected" />
            </el-select>
            <el-button type="primary" size="small" @click="loadLogs" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="logs" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_type" label="操作人类型" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.user_type === 'super_admin'" type="danger">超级用户</el-tag>
            <el-tag v-else-if="row.user_type === 'manager'" type="warning">管理员</el-tag>
            <el-tag v-else type="info">服务人员</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action_type" label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.action_type === 'create'" type="success">新建</el-tag>
            <el-tag v-else-if="row.action_type === 'update'" type="warning">修改</el-tag>
            <el-tag v-else-if="row.action_type === 'delete'" type="danger">删除</el-tag>
            <el-tag v-else-if="row.action_type === 'approve'" type="success">通过</el-tag>
            <el-tag v-else-if="row.action_type === 'reject'" type="danger">拒绝</el-tag>
            <el-tag v-else type="info">{{ row.action_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100" />
        <el-table-column prop="resource_id" label="资源ID" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'pending'" type="warning">待审批</el-tag>
            <el-tag v-else-if="row.status === 'approved'" type="success">已通过</el-tag>
            <el-tag v-else-if="row.status === 'rejected'" type="danger">已拒绝</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="变更内容" width="250">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="showDetail(row)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="备注/原因" width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="操作时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
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
          <el-descriptions-item label="操作人类型">
            <el-tag v-if="currentLog.user_type === 'super_admin'" type="danger">超级用户</el-tag>
            <el-tag v-else-if="currentLog.user_type === 'manager'" type="warning">管理员</el-tag>
            <el-tag v-else type="info">服务人员</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">{{ currentLog.action_type }}</el-descriptions-item>
          <el-descriptions-item label="资源类型">{{ currentLog.resource_type }}</el-descriptions-item>
          <el-descriptions-item label="资源ID">{{ currentLog.resource_id }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag v-if="currentLog.status === 'pending'" type="warning">待审批</el-tag>
            <el-tag v-else-if="currentLog.status === 'approved'" type="success">已通过</el-tag>
            <el-tag v-else-if="currentLog.status === 'rejected'" type="danger">已拒绝</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作时间">{{ formatDateTime(currentLog.created_at) }}</el-descriptions-item>
        </el-descriptions>

        <div v-if="currentLog.old_value || currentLog.new_value" style="margin-top: 20px">
          <h4>变更内容</h4>
          <div v-if="currentLog.old_value" class="change-content">
            <h5>原值：</h5>
            <pre>{{ formatJSON(currentLog.old_value) }}</pre>
          </div>
          <div v-if="currentLog.new_value" class="change-content new-value">
            <h5>新值：</h5>
            <pre>{{ formatJSON(currentLog.new_value) }}</pre>
          </div>
        </div>

        <div v-if="currentLog.reason" style="margin-top: 20px">
          <h4>备注/原因：</h4>
          <p>{{ currentLog.reason }}</p>
        </div>

        <div v-if="currentLog.approved_at" style="margin-top: 20px">
          <h4>审批信息：</h4>
          <p>审批人ID：{{ currentLog.approved_by }}</p>
          <p>审批时间：{{ formatDateTime(currentLog.approved_at) }}</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getAuditLogs } from '@/api/admin'
import type { AuditLog } from '@/api/admin'

const loading = ref(false)
const logs = ref<AuditLog[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const statusFilter = ref('')
const detailDialogVisible = ref(false)
const currentLog = ref<AuditLog | null>(null)

// 加载审计日志
const loadLogs = async () => {
  loading.value = true
  try {
    const result = await getAuditLogs({
      status: statusFilter.value,
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

// 分页
const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadLogs()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadLogs()
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

.change-content h5 {
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: bold;
}

.change-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
}
</style>
