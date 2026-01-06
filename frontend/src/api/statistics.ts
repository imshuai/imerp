import { http } from './index'

// 首页概览统计
export interface OverviewStats {
  customer_count: number
  pending_task_count: number
  active_agreement_count: number
  monthly_payment: number
  yearly_payment: number
}

// 任务统计
export interface TaskStats {
  pending: number
  in_progress: number
  completed: number
}

// 收款统计
export interface PaymentStats {
  total_amount: number
  count: number
}

// 获取首页概览统计
export function getOverview(): Promise<OverviewStats> {
  return http.get<OverviewStats>('/statistics/overview')
}

// 获取任务统计
export function getTaskStats(): Promise<TaskStats> {
  return http.get<TaskStats>('/statistics/tasks')
}

// 获取收款统计
export function getPaymentStats(params?: {
  start_date?: string
  end_date?: string
}): Promise<PaymentStats> {
  return http.get<PaymentStats>('/statistics/payments', { params })
}
