import { http } from './index'

// 收费类型
export type FeeType = '月度' | '季度' | '年度'

// 协议状态
export type AgreementStatus = '有效' | '已过期' | '已取消'

// 协议信息
export interface Agreement {
  id: number
  customer_id: number
  agreement_number: string
  start_date: string
  end_date: string
  fee_type: FeeType
  amount: number
  status: AgreementStatus
  created_at: string
  updated_at: string
  customer?: {
    id: number
    name: string
    tax_number: string
  }
  payments?: Payment[]
}

// 协议列表响应
export interface AgreementListResponse {
  total: number
  items: Agreement[]
}

// 获取协议列表
export function getAgreements(params?: {
  keyword?: string
  status?: AgreementStatus
  customer_id?: number
}): Promise<AgreementListResponse> {
  return http.get<AgreementListResponse>('/agreements', { params })
}

// 创建协议
export function createAgreement(data: Partial<Agreement>): Promise<Agreement> {
  return http.post<Agreement>('/agreements', data)
}

// 获取协议详情
export function getAgreement(id: number): Promise<Agreement> {
  return http.get<Agreement>(`/agreements/${id}`)
}

// 更新协议
export function updateAgreement(id: number, data: Partial<Agreement>): Promise<Agreement> {
  return http.put<Agreement>(`/agreements/${id}`, data)
}

// 删除协议
export function deleteAgreement(id: number): Promise<void> {
  return http.delete<void>(`/agreements/${id}`)
}
