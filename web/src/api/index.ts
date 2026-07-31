import axios from 'axios'
import type { ApiResponse, LoginRequest, LoginResponseData, ChangePasswordRequest, UserInfo } from '../types'

// 创建 axios 实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 添加 JWT Token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理 401
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // 登录接口的 401 是账号密码错误，不需要跳转，直接抛出
      const isLoginRequest = error.config?.url?.includes('/auth/login')
      if (!isLoginRequest) {
        // Token 过期或无效，清除本地存储并跳转登录页
        localStorage.removeItem('token')
        localStorage.removeItem('username')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

// 登录
export async function login(data: LoginRequest): Promise<ApiResponse<LoginResponseData>> {
  const response = await api.post('/auth/login', data)
  return response.data
}

// 修改密码
export async function changePassword(data: ChangePasswordRequest): Promise<ApiResponse> {
  const response = await api.post('/password/change', data)
  return response.data
}

// 获取用户信息
export async function getUserInfo(): Promise<ApiResponse<UserInfo>> {
  const response = await api.get('/user/info')
  return response.data
}

// 健康检查
export async function healthCheck(): Promise<ApiResponse> {
  const response = await api.get('/health')
  return response.data
}

export default api
