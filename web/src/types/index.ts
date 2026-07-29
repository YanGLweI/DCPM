// API 响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

// 登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应数据
export interface LoginResponseData {
  status: 'ok' | 'expired' | 'error'
  message: string
  token?: string
  username?: string
  password_expires_at?: string
  days_remaining?: number
  password_never_expires?: boolean
}

// 修改密码请求
export interface ChangePasswordRequest {
  username: string
  old_password: string
  new_password: string
}

// 用户信息
export interface UserInfo {
  username: string
  password_expires_at: string
  days_remaining: number
  password_never_expires: boolean
}
