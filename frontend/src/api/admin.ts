import { http } from './index'

// 审批日志
export interface AuditLog {
  id: number
  user_id: number
  user_type: string
  action_type: string
  resource_type: string
  resource_id?: number
  old_value: string
  new_value: string
  status: 'pending' | 'approved' | 'rejected'
  approved_by?: number
  approved_at?: string
  reason?: string
  created_at: string
  updated_at: string
  user?: any
}

// 管理员用户
export interface AdminUser {
  id: number
  username: string
  role: string
  person_id?: number
  must_change_password: boolean
  last_password_change?: string
  person?: any
}

// 设置管理员请求
export interface SetManagerRequest {
  person_id: number
  is_manager: boolean
}

// 审批请求
export interface ApprovalRequest {
  log_id: number
  reason?: string
}

// 获取待审批列表
export function getPendingApprovals(): Promise<AuditLog[]> {
  return http.get<AuditLog[]>('/admin/approvals/pending')
}

// 审批通过
export function approveOperation(data: ApprovalRequest): Promise<{ message: string }> {
  return http.post<{ message: string }>('/admin/approvals/approve', data)
}

// 审批拒绝
export function rejectOperation(data: ApprovalRequest): Promise<{ message: string }> {
  return http.post<{ message: string }>('/admin/approvals/reject', data)
}

// 获取审计日志
export function getAuditLogs(params?: {
  status?: string
  offset?: number
  limit?: number
}): Promise<{ total: number; items: AuditLog[] }> {
  return http.get<{ total: number; items: AuditLog[] }>('/admin/audit-logs', { params })
}

// 获取管理员列表
export function getAdminUsers(): Promise<AdminUser[]> {
  return http.get<AdminUser[]>('/admin/users')
}

// 获取服务人员列表
export function getServicePeople(): Promise<any[]> {
  return http.get<any[]>('/admin/service-people')
}

// 创建管理员
export function createAdminUser(data: {
  username: string
  password: string
  person_id: number
}): Promise<AdminUser> {
  return http.post<AdminUser>('/admin/users', data)
}

// 删除管理员
export function deleteAdminUser(id: number): Promise<{ message: string }> {
  return http.delete<{ message: string }>(`/admin/users/${id}`)
}

// 设置/取消管理员
export function setManager(data: SetManagerRequest): Promise<{ message: string }> {
  return http.post<{ message: string }>('/admin/set-manager', data)
}
