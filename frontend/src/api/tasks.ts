import { http } from './index'

// 任务状态
export type TaskStatus = '待处理' | '进行中' | '已完成'

// 任务信息
export interface Task {
  id: number
  customer_id: number
  title: string
  description: string
  status: TaskStatus
  due_date: string
  completed_at?: string
  created_at: string
  updated_at: string
  customer?: {
    id: number
    name: string
    tax_number: string
  }
}

// 任务列表响应
export interface TaskListResponse {
  total: number
  items: Task[]
}

// 获取任务列表
export function getTasks(params?: {
  keyword?: string
  status?: TaskStatus
  customer_id?: number
}): Promise<TaskListResponse> {
  return http.get<TaskListResponse>('/tasks', { params })
}

// 创建任务
export function createTask(data: Partial<Task>): Promise<Task> {
  return http.post<Task>('/tasks', data)
}

// 获取任务详情
export function getTask(id: number): Promise<Task> {
  return http.get<Task>(`/tasks/${id}`)
}

// 更新任务
export function updateTask(id: number, data: Partial<Task>): Promise<Task> {
  return http.put<Task>(`/tasks/${id}`, data)
}

// 删除任务
export function deleteTask(id: number): Promise<void> {
  return http.delete<void>(`/tasks/${id}`)
}
