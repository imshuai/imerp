<template>
  <div class="approvals-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>待审批列表</span>
          <el-button type="primary" size="small" @click="loadApprovals" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table :data="approvals" v-loading="loading" stripe>
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
            <el-tag v-else type="info">{{ row.action_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="100" />
        <el-table-column prop="resource_id" label="资源ID" width="80" />
        <el-table-column label="变更内容" width="300">
          <template #default="{ row }">
            <div class="change-content">
              <div v-if="row.old_value" class="old-value">
                <span class="label">原值:</span>
                <pre>{{ formatJSON(row.old_value) }}</pre>
              </div>
              <div v-if="row.new_value" class="new-value">
                <span class="label">新值:</span>
                <pre>{{ formatJSON(row.new_value) }}</pre>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="提交时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="handleApprove(row)">
              通过
            </el-button>
            <el-button type="danger" size="small" @click="handleReject(row)">
              拒绝
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && approvals.length === 0" description="暂无待审批记录" />
    </el-card>

    <!-- 拒绝原因对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="拒绝原因" width="500px">
      <el-form :model="rejectForm" label-width="80px">
        <el-form-item label="拒绝原因">
          <el-input
            v-model="rejectForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入拒绝原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmReject" :loading="loading">
          确认拒绝
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getPendingApprovals, approveOperation, rejectOperation } from '@/api/admin'
import type { AuditLog } from '@/api/admin'

const loading = ref(false)
const approvals = ref<AuditLog[]>([])
const rejectDialogVisible = ref(false)
const rejectForm = reactive({
  log_id: 0,
  reason: ''
})

// 加载待审批列表
const loadApprovals = async () => {
  loading.value = true
  try {
    approvals.value = await getPendingApprovals()
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

// 审批通过
const handleApprove = async (row: AuditLog) => {
  try {
    await ElMessageBox.confirm('确认通过此操作？', '提示', {
      type: 'warning'
    })

    loading.value = true
    await approveOperation({ log_id: row.id })
    ElMessage.success('审批通过')
    loadApprovals()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '操作失败')
    }
  } finally {
    loading.value = false
  }
}

// 审批拒绝
const handleReject = (row: AuditLog) => {
  rejectForm.log_id = row.id
  rejectForm.reason = ''
  rejectDialogVisible.value = true
}

// 确认拒绝
const confirmReject = async () => {
  if (!rejectForm.reason.trim()) {
    ElMessage.warning('请输入拒绝原因')
    return
  }

  loading.value = true
  try {
    await rejectOperation({
      log_id: rejectForm.log_id,
      reason: rejectForm.reason
    })
    ElMessage.success('已拒绝')
    rejectDialogVisible.value = false
    loadApprovals()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
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
  return new Date(str).toLocaleString('zh-CN')
}

onMounted(() => {
  loadApprovals()
})
</script>

<style scoped>
.approvals-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.change-content {
  font-size: 12px;
}

.change-content .old-value {
  color: #909399;
  margin-bottom: 8px;
}

.change-content .new-value {
  color: #67C23A;
}

.change-content .label {
  font-weight: bold;
  margin-right: 8px;
}

.change-content pre {
  margin: 4px 0 0 0;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
