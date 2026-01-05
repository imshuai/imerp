import { http } from './index'

// 服务人员类型（固定为服务人员）
export type ServicePersonType = '服务人员'

// 服务人员信息（包含客户数量）
export interface ServicePersonnel {
  id: number
  type: ServicePersonType
  name: string
  phone: string
  password: string
  id_card: string
  service_customer_ids?: string
  customer_count: number
  created_at: string
  updated_at: string
}

// 服务人员列表响应
export interface ServicePersonnelListResponse {
  total: number
  items: ServicePersonnel[]
}

// 获取服务人员列表
export function getServicePersonnel(params?: {
  keyword?: string
}): Promise<ServicePersonnelListResponse> {
  return http.get<ServicePersonnelListResponse>('/service-personnel', { params })
}

// 创建服务人员
export function createServicePersonnel(data: Partial<ServicePersonnel>): Promise<ServicePersonnel> {
  return http.post<ServicePersonnel>('/service-personnel', data)
}

// 获取人员详情（复用原有API）
export function getServicePerson(id: number): Promise<ServicePersonnel> {
  return http.get<ServicePersonnel>(`/service-personnel/${id}`)
}

// 更新服务人员（复用原有API）
export function updateServicePersonnel(id: number, data: Partial<ServicePersonnel>): Promise<ServicePersonnel> {
  return http.put<ServicePersonnel>(`/service-personnel/${id}`, data)
}

// 删除服务人员（复用原有API）
export function deleteServicePersonnel(id: number): Promise<void> {
  return http.delete<void>(`/service-personnel/${id}`)
}
