import { http } from './index'
import type { Person } from './people'
import type { Agreement } from './agreements'

// 客户类型
export type CustomerType = '有限公司' | '个人独资企业' | '合伙企业' | '个体工商户'

// 信用等级
export type CreditRating = 'A' | 'B' | 'C' | 'D' | 'M'

// 纳税人类型
export type TaxpayerType = '一般纳税人' | '小规模纳税人'

// 投资人信息
export interface Investor {
  person_id: number
  share_ratio: number
}

// 出资记录
export interface InvestmentRecord {
  date: string
  amount: number
}

// 客户投资人关联
export interface CustomerInvestor {
  id: number
  customer_id: number
  person_id: number
  share_ratio: number
  investment_records?: InvestmentRecord[]
  person?: Person
}

// 账户类型
export type AccountType = '基本户' | '一般户' | '临时户'

// 对公账户
export interface BankAccount {
  id: number
  customer_id: number
  bank_name: string
  account_number: string
  bank_code?: string
  contact_phone?: string
  account_type: AccountType
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
  taxpayer_type?: TaxpayerType
  license_registration_date?: string  // 执照登记日
  tax_registration_date?: string      // 税务登记日
  tax_office?: string                 // 税务所
  tax_administrator?: string          // 税务管理员
  tax_administrator_phone?: string    // 税务管理员联系电话
  // 新增字段 v0.4.0
  tax_agent_ids?: string              // 办税人ID，逗号分隔
  credit_rating?: CreditRating        // 纳税人信用等级
  social_security_number?: string     // 社保号
  yukuai_ban_password?: string        // 渝快办密码
  business_scope?: string             // 经营范围
  created_at: string
  updated_at: string
  // 关联数据
  representative?: Person
  investor_list?: Person[]
  investor_relations?: CustomerInvestor[]  // 投资人关联列表
  service_persons?: Person[]
  tax_agents?: Person[]               // 办税人列表
  agreements_list?: Agreement[]
  bank_accounts?: BankAccount[]       // 对公账户列表
}

// 客户列表响应
export interface CustomerListResponse {
  total: number
  items: Customer[]
}

// 获取客户类型选项
export function getCustomerTypes(): Promise<CustomerType[]> {
  return http.get<CustomerType[]>('/customers/types')
}

// 获取信用等级选项
export function getCreditRatings(): Promise<CreditRating[]> {
  return http.get<CreditRating[]>('/customers/credit-ratings')
}

// 获取账户类型选项
export function getAccountTypes(): Promise<AccountType[]> {
  return http.get<AccountType[]>('/bank-accounts/types')
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
