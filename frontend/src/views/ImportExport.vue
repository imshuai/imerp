<template>
  <div class="import-export-page">
    <el-row :gutter="20">
      <!-- 人员导入导出 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>人员导入导出</span>
            </div>
          </template>

          <div class="section">
            <h4>下载模板</h4>
            <el-button @click="downloadPeopleTemplate">
              <el-icon><Download /></el-icon>
              下载人员导入模板
            </el-button>
          </div>

          <el-divider />

          <div class="section">
            <h4>导入人员</h4>
            <el-upload
              :auto-upload="false"
              :on-change="handlePeopleImport"
              :show-file-list="false"
              accept=".xlsx,.xls"
            >
              <el-button type="primary">
                <el-icon><Upload /></el-icon>
                选择Excel文件
              </el-button>
            </el-upload>
          </div>

          <el-divider />

          <div class="section">
            <h4>导出人员</h4>
            <el-button type="success" @click="exportPeople" :loading="exportingPeople">
              <el-icon><Upload /></el-icon>
              导出人员数据
            </el-button>
          </div>
        </el-card>
      </el-col>

      <!-- 客户导入导出 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>客户导入导出</span>
            </div>
          </template>

          <div class="section">
            <h4>下载模板</h4>
            <el-button @click="downloadCustomersTemplate">
              <el-icon><Download /></el-icon>
              下载客户导入模板
            </el-button>
          </div>

          <el-divider />

          <div class="section">
            <h4>导入客户</h4>
            <el-upload
              :auto-upload="false"
              :on-change="handleCustomersImport"
              :show-file-list="false"
              accept=".xlsx,.xls"
            >
              <el-button type="primary">
                <el-icon><Upload /></el-icon>
                选择Excel文件
              </el-button>
            </el-upload>
          </div>

          <el-divider />

          <div class="section">
            <h4>导出客户</h4>
            <el-button type="success" @click="exportCustomers" :loading="exportingCustomers">
              <el-icon><Upload /></el-icon>
              导出客户数据
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 导入冲突策略选择对话框 -->
    <el-dialog v-model="strategyDialogVisible" title="选择导入策略" width="400px">
      <el-alert
        title="当数据冲突时，请选择处理方式："
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      />
      <el-radio-group v-model="selectedStrategy">
        <el-radio label="skip">跳过已存在的记录</el-radio>
        <el-radio label="update">更新已存在的记录</el-radio>
        <el-radio label="create_new">修改标识后创建新记录</el-radio>
      </el-radio-group>
      <template #footer>
        <el-button @click="strategyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmImport" :loading="importing">确定</el-button>
      </template>
    </el-dialog>

    <!-- 导入结果对话框 -->
    <el-dialog v-model="resultDialogVisible" title="导入结果" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="总行数">{{ importResult.total }}</el-descriptions-item>
        <el-descriptions-item label="成功">{{ importResult.success }}</el-descriptions-item>
        <el-descriptions-item label="失败">{{ importResult.failed }}</el-descriptions-item>
      </el-descriptions>

      <el-alert
        v-if="importResult.errors.length > 0"
        title="错误详情"
        type="error"
        :closable="false"
        style="margin-top: 20px"
      >
        <el-table :data="importResult.errors" max-height="300" size="small">
          <el-table-column prop="row" label="行号" width="80" />
          <el-table-column prop="column" label="列名" width="120" />
          <el-table-column prop="message" label="错误信息" />
        </el-table>
      </el-alert>

      <template #footer>
        <el-button type="primary" @click="resultDialogVisible = false">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  downloadTemplate,
  importPeople,
  importCustomers,
  exportPeople,
  exportCustomers,
  downloadBlob,
  type ImportResult,
  type ImportStrategy
} from '@/api/import_export'

const strategyDialogVisible = ref(false)
const resultDialogVisible = ref(false)
const importing = ref(false)
const exportingPeople = ref(false)
const exportingCustomers = ref(false)

const selectedStrategy = ref<ImportStrategy>('skip')
const currentImportType = ref<'people' | 'customers'>('people')
const currentFile = ref<File | null>(null)

const importResult = ref<ImportResult>({
  total: 0,
  success: 0,
  failed: 0,
  errors: []
})

const downloadPeopleTemplate = async () => {
  try {
    const blob = await downloadTemplate('people')
    downloadBlob(blob, '人员导入模板.xlsx')
    ElMessage.success('模板下载成功')
  } catch (error) {
    console.error('下载模板失败:', error)
  }
}

const downloadCustomersTemplate = async () => {
  try {
    const blob = await downloadTemplate('customers')
    downloadBlob(blob, '客户导入模板.xlsx')
    ElMessage.success('模板下载成功')
  } catch (error) {
    console.error('下载模板失败:', error)
  }
}

const handlePeopleImport = (file: File) => {
  currentFile.value = file
  currentImportType.value = 'people'
  strategyDialogVisible.value = true
}

const handleCustomersImport = (file: File) => {
  currentFile.value = file
  currentImportType.value = 'customers'
  strategyDialogVisible.value = true
}

const confirmImport = async () => {
  if (!currentFile.value) return

  importing.value = true
  try {
    let result: ImportResult
    if (currentImportType.value === 'people') {
      result = await importPeople(currentFile.value, selectedStrategy.value)
    } else {
      result = await importCustomers(currentFile.value, selectedStrategy.value)
    }

    importResult.value = result
    strategyDialogVisible.value = false
    resultDialogVisible.value = true

    if (result.failed === 0) {
      ElMessage.success(`导入成功！共 ${result.success} 条记录`)
    } else {
      ElMessage.warning(`导入完成，成功 ${result.success} 条，失败 ${result.failed} 条`)
    }
  } catch (error) {
    console.error('导入失败:', error)
  } finally {
    importing.value = false
  }
}

const exportPeople = async () => {
  exportingPeople.value = true
  try {
    const blob = await exportPeople()
    downloadBlob(blob, `人员导出_${new Date().getTime()}.xlsx`)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  } finally {
    exportingPeople.value = false
  }
}

const exportCustomers = async () => {
  exportingCustomers.value = true
  try {
    const blob = await exportCustomers()
    downloadBlob(blob, `客户导出_${new Date().getTime()}.xlsx`)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  } finally {
    exportingCustomers.value = false
  }
}
</script>

<style scoped>
.import-export-page {
  padding: 0;
}

.section {
  margin-bottom: 20px;
}

.section h4 {
  margin: 0 0 10px 0;
  color: #606266;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
