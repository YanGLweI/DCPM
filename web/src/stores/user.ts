import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, getUserInfo as getUserInfoApi } from '../api'
import type { LoginRequest } from '../types'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const username = ref(localStorage.getItem('username') || '')
  const passwordExpiresAt = ref('')
  const daysRemaining = ref(0)
  const passwordNeverExpires = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const isExpired = computed(() => !passwordNeverExpires.value && daysRemaining.value <= 0)
  const isExpiringSoon = computed(() => !passwordNeverExpires.value && daysRemaining.value > 0 && daysRemaining.value <= 7)

  // 登录
  async function login(loginData: LoginRequest) {
    try {
      const response = await loginApi(loginData)
      
      if (response.data?.status === 'ok') {
        // 登录成功，保存信息
        token.value = response.data.token || ''
        username.value = response.data.username || loginData.username
        passwordExpiresAt.value = response.data.password_expires_at || ''
        daysRemaining.value = response.data.days_remaining || 0
        passwordNeverExpires.value = response.data.password_never_expires || false
        
        localStorage.setItem('token', token.value)
        localStorage.setItem('username', username.value)
        
        return 'ok'
      } else if (response.data?.status === 'expired') {
        // 密码过期
        sessionStorage.setItem('temp_username', loginData.username)
        return 'expired'
      } else {
        throw new Error(response.data?.message || '登录失败')
      }
    } catch (error: any) {
      // axios 401 错误：从响应中提取后端返回的错误信息
      const message = error.response?.data?.message || error.message || '登录失败，请检查账号密码'
      throw new Error(message)
    }
  }

  // 获取用户信息
  async function fetchUserInfo() {
    try {
      const response = await getUserInfoApi()
      if (response.data) {
        passwordExpiresAt.value = response.data.password_expires_at
        daysRemaining.value = response.data.days_remaining
        passwordNeverExpires.value = response.data.password_never_expires || false
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }

  // 登出
  function logout() {
    token.value = ''
    username.value = ''
    passwordExpiresAt.value = ''
    daysRemaining.value = 0
    passwordNeverExpires.value = false
    localStorage.removeItem('token')
    localStorage.removeItem('username')
  }

  return {
    token,
    username,
    passwordExpiresAt,
    daysRemaining,
    passwordNeverExpires,
    isLoggedIn,
    isExpired,
    isExpiringSoon,
    login,
    fetchUserInfo,
    logout
  }
})
