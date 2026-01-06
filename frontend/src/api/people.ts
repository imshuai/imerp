import { http } from './index'

// 人员信息
export interface Person {
  id: number
  is_service_person?: boolean
  name: string
  phone: string
  id_card: string
  password?: string
  representative_customer_ids?: string
  investor_customer_ids?: string
  service_customer_ids?: string
  created_at: string
  updated_at: string
}

// 人员列表响应
export interface PeopleListResponse {
  total: number
  items: Person[]
}

// 获取人员列表
export function getPeople(params?: {
  is_service_person?: boolean
  keyword?: string
}): Promise<PeopleListResponse> {
  return http.get<PeopleListResponse>('/people', { params })
}

// 创建人员
export function createPerson(data: Partial<Person>): Promise<Person> {
  return http.post<Person>('/people', data)
}

// 获取人员详情
export function getPerson(id: number): Promise<Person> {
  return http.get<Person>(`/people/${id}`)
}

// 更新人员
export function updatePerson(id: number, data: Partial<Person>): Promise<Person> {
  return http.put<Person>(`/people/${id}`, data)
}

// 删除人员
export function deletePerson(id: number): Promise<void> {
  return http.delete<void>(`/people/${id}`)
}

// 获取人员关联的企业
export function getPersonCustomers(id: number): Promise<any> {
  return http.get<any>(`/people/${id}/customers`)
}
