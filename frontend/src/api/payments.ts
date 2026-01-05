import { http } from './index'

// 收款方式
export type PaymentMethod = '转账' | '现金' | '支票' | '其他'

// 收款记录
export interface Payment {
  id: number
  customer_id: number
  agreement_id?: number
  amount: number
  payment_date: string
  payment_method: PaymentMethod
  period: string
  remark: string
  created_at: string
  updated_at: string
  customer?: {
    id: number
    name: string
    tax_number: string
  }
  agreement?: {
    id: number
    agreement_number: string
    fee_type: string
    amount: number
  }
}

// 收款列表响应
export interface PaymentListResponse {
  total: number
  items: Payment[]
}

// 获取收款列表
export function getPayments(params?: {
  customer_id?: number
  start_date?: string
  end_date?: string
}): Promise<PaymentListResponse> {
  return http.get<PaymentListResponse>('/payments', { params })
}

// 创建收款记录
export function createPayment(data: Partial<Payment>): Promise<Payment> {
  return http.post<Payment>('/payments', data)
}

// 获取收款详情
export function getPayment(id: number): Promise<Payment> {
  return http.get<Payment>(`/payments/${id}`)
}

// 更新收款记录
export function updatePayment(id: number, data: Partial<Payment>): Promise<Payment> {
  return http.put<Payment>(`/payments/${id}`, data)
}

// 删除收款记录
export function deletePayment(id: number): Promise<void> {
  return http.delete<void>(`/payments/${id}`)
}
