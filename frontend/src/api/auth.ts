import { http, getToken, setToken as _setToken, removeToken as _removeToken } from './index'

// Re-export functions for use in other modules
export const setToken = _setToken
export const removeToken = _removeToken

// 登录请求
export interface LoginRequest {
  username?: string
  password?: string
  person_id?: number
}

// 登录响应
export interface LoginResponse {
  token: string
  user_id: number
  username: string
  role: string
  must_change_password: boolean
}

// 用户信息
export interface UserInfo {
  id: number
  username?: string
  name?: string
  role: string
  must_change_password?: boolean
  is_manager?: boolean
  is_service_person?: boolean
  person?: any
}

// 修改密码请求
export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

// 登录
export function login(data: LoginRequest): Promise<LoginResponse> {
  return http.post<LoginResponse>('/auth/login', data)
}

// 获取当前用户信息
export function getCurrentUser(): Promise<UserInfo> {
  return http.get<UserInfo>('/user/me')
}

// 修改密码
export function changePassword(data: ChangePasswordRequest): Promise<{ message: string }> {
  return http.post<{ message: string }>('/user/change-password', data)
}

// 登出
export function logout() {
  removeToken()
  window.location.href = '/login'
}

// 检查是否已登录
export function isLoggedIn(): boolean {
  return !!getToken()
}

// 获取存储的用户信息
export function getStoredUser(): UserInfo | null {
  const userStr = localStorage.getItem('erp_user')
  return userStr ? JSON.parse(userStr) : null
}

// 存储用户信息
export function setStoredUser(user: UserInfo): void {
  localStorage.setItem('erp_user', JSON.stringify(user))
}

// 移除用户信息
export function removeStoredUser(): void {
  localStorage.removeItem('erp_user')
}
