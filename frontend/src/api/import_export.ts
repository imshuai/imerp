import { http } from './index'

// 导入策略
export type ImportStrategy = 'skip' | 'update' | 'create_new'

// 导入结果
export interface ImportResult {
  total: number
  success: number
  failed: number
  errors: ImportError[]
}

// 导入错误
export interface ImportError {
  row: number
  column: string
  message: string
}

// 下载模板
export function downloadTemplate(type: 'people' | 'customers'): Promise<Blob> {
  return http.get<Blob>(`/templates/${type}`, {
    responseType: 'blob'
  })
}

// 导入人员
export function importPeople(file: File, strategy: ImportStrategy): Promise<ImportResult> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('strategy', strategy)
  return http.post<ImportResult>('/import/people', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 导入客户
export function importCustomers(file: File, strategy: ImportStrategy): Promise<ImportResult> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('strategy', strategy)
  return http.post<ImportResult>('/import/customers', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 导出人员
export function exportPeople(): Promise<Blob> {
  return http.get<Blob>('/export/people', {
    responseType: 'blob'
  })
}

// 导出客户
export function exportCustomers(): Promise<Blob> {
  return http.get<Blob>('/export/customers', {
    responseType: 'blob'
  })
}

// 下载文件工具函数
export function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
