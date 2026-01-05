import { http } from './index'
import type { Person } from './people'
import type { Agreement } from './agreements'

// 客户类型
export type CustomerType = '有限公司' | '个人独资企业' | '合伙企业' | '个体工商户'

// 投资人信息
export interface Investor {
  person_id: number
  share_ratio: number
}

// 客户信息
export interface Customer {
  id: number
  name: string
  phone: string
  address: string
  tax_number: string
  type: CustomerType
  representative_id?: number
  investor_ids?: string
  investors?: string  // JSON格式的投资人信息
  service_person_ids?: string
  agreement_ids?: string
  registered_capital?: number
  created_at: string
  updated_at: string
  // 关联数据
  representative?: Person
  investor_list?: Person[]
  service_persons?: Person[]
  agreements_list?: Agreement[]
}

// 客户列表响应
export interface CustomerListResponse {
  total: number
  items: Customer[]
}

// 获取客户列表
export function getCustomers(params?: {
  keyword?: string
  representative?: string
  investor?: string
  service_person?: string
}): Promise<CustomerListResponse> {
  return http.get<CustomerListResponse>('/customers', { params })
}

// 创建客户
export function createCustomer(data: Partial<Customer>): Promise<Customer> {
  return http.post<Customer>('/customers', data)
}

// 获取客户详情
export function getCustomer(id: number): Promise<Customer> {
  return http.get<Customer>(`/customers/${id}`)
}

// 更新客户
export function updateCustomer(id: number, data: Partial<Customer>): Promise<Customer> {
  return http.put<Customer>(`/customers/${id}`, data)
}

// 删除客户
export function deleteCustomer(id: number): Promise<void> {
  return http.delete<void>(`/customers/${id}`)
}

// 获取客户的任务列表
export function getCustomerTasks(id: number): Promise<any> {
  return http.get<any>(`/customers/${id}/tasks`)
}

// 获取客户的收款记录
export function getCustomerPayments(id: number): Promise<any> {
  return http.get<any>(`/customers/${id}/payments`)
}
