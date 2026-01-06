<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6" v-for="stat in stats" :key="stat.title">
        <el-card class="stat-card" @click="handleStatClick(stat)">
          <div class="stat-content">
            <div class="stat-icon" :style="{ backgroundColor: stat.color }">
              <el-icon :size="32" :color="stat.iconColor">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-title">{{ stat.title }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>任务统计</span>
            </div>
          </template>
          <div v-if="taskStats">
            <el-row :gutter="20">
              <el-col :span="8">
                <div class="task-stat">
                  <div class="task-stat-value">{{ taskStats.pending }}</div>
                  <div class="task-stat-label">待处理</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="task-stat">
                  <div class="task-stat-value">{{ taskStats.in_progress }}</div>
                  <div class="task-stat-label">进行中</div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="task-stat">
                  <div class="task-stat-value">{{ taskStats.completed }}</div>
                  <div class="task-stat-label">已完成</div>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>收款统计</span>
            </div>
          </template>
          <div v-if="paymentStats">
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="payment-stat">
                  <div class="payment-stat-value">¥{{ paymentStats.total_amount.toLocaleString() }}</div>
                  <div class="payment-stat-label">总金额</div>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="payment-stat">
                  <div class="payment-stat-value">{{ paymentStats.count }}</div>
                  <div class="payment-stat-label">笔数</div>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getOverview, getTaskStats, getPaymentStats } from '@/api/statistics'
import type { TaskStats, PaymentStats } from '@/api/statistics'

const router = useRouter()

const stats = ref([
  { title: '客户总数', value: 0, icon: 'OfficeBuilding', color: '#ecf5ff', iconColor: '#409EFF', route: '/customers' },
  { title: '待办任务', value: 0, icon: 'List', color: '#fef0f0', iconColor: '#F56C6C', route: '/tasks' },
  { title: '有效协议', value: 0, icon: 'Document', color: '#f0f9ff', iconColor: '#67C23A', route: '/agreements' },
  { title: '本月收款', value: '¥0', icon: 'Money', color: '#fdf6ec', iconColor: '#E6A23C', route: '/payments' }
])

const taskStats = ref<TaskStats | null>(null)
const paymentStats = ref<PaymentStats | null>(null)

const loadData = async () => {
  try {
    const overview = await getOverview()
    stats.value[0].value = overview.customer_count
    stats.value[1].value = overview.pending_task_count
    stats.value[2].value = overview.active_agreement_count
    stats.value[3].value = '¥' + overview.monthly_payment.toLocaleString()

    taskStats.value = await getTaskStats()
    paymentStats.value = await getPaymentStats()
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

const handleStatClick = (stat: any) => {
  if (stat.route) {
    router.push(stat.route)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.stat-title {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.task-stat,
.payment-stat {
  text-align: center;
  padding: 20px 0;
}

.task-stat-value,
.payment-stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #409EFF;
}

.task-stat-label,
.payment-stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}
</style>
