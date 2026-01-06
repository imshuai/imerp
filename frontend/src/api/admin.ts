import { http } from './index'

// 审计日志
export interface AuditLog {
  id: number
  user_id: number
  user_type: string
  action_type: string
  resource_type: string
  resource_id?: number
  resource_name: string
  old_value: string
  new_value: string
  created_at: string
  user?: {
    id: number
    username: string
    role: string
    person?: {
      name: string
    }
  }
}

// 获取审计日志
export function getAuditLogs(params?: {
  offset?: number
  limit?: number
}): Promise<{ total: number; items: AuditLog[] }> {
  return http.get<{ total: number; items: AuditLog[] }>('/audit-logs', { params })
}

// 删除审计日志（仅超级管理员）
export function deleteAuditLog(id: number): Promise<{ message: string }> {
  return http.delete<{ message: string }>(`/admin/audit-logs/${id}`)
}

// 批量清理审计日志（仅超级管理员）
export function clearAuditLogs(data: {
  start_date?: string
  end_date?: string
}): Promise<{ message: string; count: number }> {
  return http.post<{ message: string; count: number }>('/admin/audit-logs/clear', data)
}
