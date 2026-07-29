<template>
  <div class="account-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-content">
        <div class="header-left">
          <div class="header-logo">
            <svg width="32" height="32" viewBox="0 0 48 48" fill="none">
              <rect width="48" height="48" rx="12" fill="#007AFF"/>
              <path d="M24 14C19.58 14 16 17.58 16 22C16 26.42 19.58 30 24 30C28.42 30 32 26.42 32 22C32 17.58 28.42 14 24 14ZM24 27C21.24 27 19 24.76 19 22C19 19.24 21.24 17 24 17C26.76 17 29 19.24 29 22C29 24.76 26.76 27 24 27ZM24 32C19.33 32 10 34.34 10 39V41H38V39C38 34.34 28.67 32 24 32Z" fill="white"/>
            </svg>
          </div>
          <h1 class="header-title">域账号自助管理平台</h1>
        </div>
        <div class="header-right">
          <span class="user-badge">
            <svg width="16" height="16" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"/>
            </svg>
            {{ userStore.username }}
          </span>
          <el-button type="danger" plain class="btn-logout" @click="handleLogout">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 4px">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            退出
          </el-button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="main-content">
      <!-- Account Info Card -->
      <section class="info-section">
        <div class="card info-card">
          <div class="card-header">
            <div class="card-title">
              <svg width="20" height="20" viewBox="0 0 20 20" fill="var(--apple-blue)">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
              </svg>
              <h2>账号信息</h2>
            </div>
          </div>
          
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">账号</span>
              <span class="info-value">{{ userStore.username }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">密码状态</span>
              <span class="info-value">
                <span class="status-badge" :class="passwordStatusClass">
                  {{ passwordStatusText }}
                </span>
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">密码过期时间</span>
              <span class="info-value">
                <template v-if="userStore.passwordNeverExpires">
                  <span class="status-badge status-info">永不过期</span>
                </template>
                <template v-else>
                  {{ userStore.passwordExpiresAt || '加载中...' }}
                </template>
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">剩余天数</span>
              <span class="info-value">
                <template v-if="userStore.passwordNeverExpires">
                  <span class="status-badge status-info">-</span>
                </template>
                <template v-else>
                  <span :class="{ 'text-warning': userStore.isExpiringSoon, 'text-danger': userStore.isExpired }">
                    {{ userStore.daysRemaining }} 天
                  </span>
                </template>
              </span>
            </div>
          </div>
          
          <!-- Status Alerts -->
          <div v-if="userStore.passwordNeverExpires" class="status-alert status-alert-info">
            <svg width="18" height="18" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
            </svg>
            <span>当前账号已设置密码永不过期</span>
          </div>
          
          <div v-else-if="userStore.isExpiringSoon && !userStore.isExpired" class="status-alert status-alert-warning">
            <svg width="18" height="18" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
            </svg>
            <span>您的密码将在 <strong>{{ userStore.daysRemaining }}</strong> 天后过期，请尽快修改</span>
          </div>
          
          <div v-else-if="userStore.isExpired" class="status-alert status-alert-error">
            <svg width="18" height="18" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
            </svg>
            <span>您的密码已过期，请立即修改</span>
          </div>
        </div>
      </section>

      <!-- Change Password Card -->
      <section class="password-section">
        <div class="card password-card">
          <div class="card-header">
            <div class="card-title">
              <svg width="20" height="20" viewBox="0 0 20 20" fill="var(--apple-blue)">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd"/>
              </svg>
              <h2>修改密码</h2>
            </div>
          </div>
          
          <el-form
            ref="formRef"
            :model="passwordForm"
            :rules="rules"
            label-position="top"
            class="password-form"
          >
            <el-form-item label="旧密码" prop="oldPassword">
              <el-input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入当前密码"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="请输入新密码"
                show-password
                @input="checkPasswordStrength"
              />
              <!-- Strength Indicator -->
              <div v-if="passwordForm.newPassword" class="strength-indicator">
                <div class="strength-bars">
                  <div class="strength-bar" :class="{ active: strengthLevel >= 1 }"></div>
                  <div class="strength-bar" :class="{ active: strengthLevel >= 2 }"></div>
                  <div class="strength-bar" :class="{ active: strengthLevel >= 3 }"></div>
                  <div class="strength-bar" :class="{ active: strengthLevel >= 4 }"></div>
                </div>
                <span class="strength-text" :style="{ color: strengthColor }">{{ strengthText }}</span>
              </div>
              <div class="password-tips">
                <svg width="14" height="14" viewBox="0 0 20 20" fill="currentColor" class="tip-icon">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
                </svg>
                <span>至少14位，包含大写字母、小写字母、数字、特殊字符中的至少3类</span>
              </div>
            </el-form-item>
            
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
                show-password
              />
            </el-form-item>
            
            <div class="form-actions">
              <el-button
                type="primary"
                :loading="loading"
                class="btn-primary"
                @click="handleChangePassword"
              >
                确认修改
              </el-button>
              <el-button class="btn-secondary" @click="resetForm">
                重置
              </el-button>
            </div>
          </el-form>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '../stores/user'
import { changePassword } from '../api'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// Password strength
const strengthLevel = computed(() => {
  const pwd = passwordForm.newPassword
  if (!pwd) return 0
  let strength = 0
  if (pwd.length >= 14) strength++
  if (/[A-Z]/.test(pwd)) strength++
  if (/[a-z]/.test(pwd)) strength++
  if (/[0-9]/.test(pwd)) strength++
  if (/[!@#$%^&*()_+\-=\[\]{}|;':",./<>?]/.test(pwd)) strength++
  if (strength <= 2) return 1
  if (strength === 3) return 2
  if (strength === 4) return 3
  return 4
})

const strengthPercentage = ref(0)
const strengthColor = ref('#909399')
const strengthText = ref('')

const checkPasswordStrength = () => {
  const pwd = passwordForm.newPassword
  if (!pwd) {
    strengthPercentage.value = 0
    return
  }
  
  let strength = 0
  if (pwd.length >= 14) strength++
  if (/[A-Z]/.test(pwd)) strength++
  if (/[a-z]/.test(pwd)) strength++
  if (/[0-9]/.test(pwd)) strength++
  if (/[!@#$%^&*()_+\-=\[\]{}|;':",./<>?]/.test(pwd)) strength++
  
  switch (strength) {
    case 0:
    case 1:
    case 2:
      strengthPercentage.value = 20
      strengthColor.value = '#FF3B30'
      strengthText.value = '弱'
      break
    case 3:
      strengthPercentage.value = 60
      strengthColor.value = '#FF9500'
      strengthText.value = '中'
      break
    case 4:
      strengthPercentage.value = 80
      strengthColor.value = '#007AFF'
      strengthText.value = '强'
      break
    case 5:
      strengthPercentage.value = 100
      strengthColor.value = '#34C759'
      strengthText.value = '很强'
      break
  }
}

// Password status
const passwordStatusClass = computed(() => {
  if (userStore.passwordNeverExpires) return 'status-info'
  if (userStore.isExpired) return 'status-danger'
  if (userStore.isExpiringSoon) return 'status-warning'
  return 'status-success'
})

const passwordStatusText = computed(() => {
  if (userStore.passwordNeverExpires) return '永不过期'
  if (userStore.isExpired) return '已过期'
  if (userStore.isExpiringSoon) return '即将过期'
  return '正常'
})

// Validation rules
const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 14, message: '密码长度不能少于14位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

onMounted(async () => {
  await userStore.fetchUserInfo()
})

const handleChangePassword = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      await changePassword({
        username: userStore.username,
        old_password: passwordForm.oldPassword,
        new_password: passwordForm.newPassword
      })
      
      ElMessage.success('密码修改成功')
      
      await ElMessageBox.confirm(
        '密码修改成功，是否立即退出重新登录？',
        '提示',
        {
          confirmButtonText: '退出登录',
          cancelButtonText: '继续操作',
          type: 'success'
        }
      )
      
      handleLogout()
    } catch (error: any) {
      if (error === 'cancel') return
      const message = error.response?.data?.message || error.message || '密码修改失败'
      ElMessage.error(message)
    } finally {
      loading.value = false
    }
  })
}

const resetForm = () => {
  formRef.value?.resetFields()
  strengthPercentage.value = 0
}

const handleLogout = () => {
  ElMessageBox.confirm(
    '确定要退出登录吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    userStore.logout()
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped>
.account-page {
  min-height: 100vh;
  background: var(--color-bg-page);
}

/* Header */
.page-header {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 960px;
  margin: 0 auto;
  padding: var(--space-md) var(--space-xl);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.header-logo {
  flex-shrink: 0;
}

.header-title {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.3px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.user-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--apple-blue-light);
  color: var(--apple-blue);
  border-radius: var(--radius-full);
  font-size: var(--font-size-sm);
  font-weight: 500;
}

.btn-logout {
  border-radius: var(--radius-md) !important;
  font-weight: 500 !important;
}

/* Main Content */
.main-content {
  max-width: 960px;
  margin: 0 auto;
  padding: var(--space-xl);
  display: flex;
  flex-direction: column;
  gap: var(--space-xl);
}

/* Cards */
.card {
  background: var(--apple-white);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
  transition: transform var(--transition-normal), box-shadow var(--transition-normal);
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.card-header {
  padding: var(--space-lg) var(--space-xl);
  border-bottom: 1px solid var(--color-border);
}

.card-title {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.card-title h2 {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

/* Info Grid */
.info-grid {
  padding: var(--space-xl);
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-lg);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  font-weight: 500;
}

.info-value {
  font-size: var(--font-size-base);
  color: var(--color-text-primary);
  font-weight: 500;
}

/* Status Badges */
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: 600;
}

.status-success {
  background: rgba(52, 199, 89, 0.15);
  color: #248A3D;
}

.status-warning {
  background: rgba(255, 149, 0, 0.15);
  color: #C93400;
}

.status-danger {
  background: rgba(255, 59, 48, 0.15);
  color: #BF2C00;
}

.status-info {
  background: rgba(0, 122, 255, 0.1);
  color: var(--apple-blue);
}

/* Status Alerts */
.status-alert {
  margin: 0 var(--space-xl) var(--space-xl);
  padding: var(--space-md);
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  font-size: var(--font-size-sm);
}

.status-alert-info {
  background: rgba(0, 122, 255, 0.08);
  color: var(--apple-blue);
  border: 1px solid rgba(0, 122, 255, 0.15);
}

.status-alert-warning {
  background: rgba(255, 149, 0, 0.08);
  color: #C93400;
  border: 1px solid rgba(255, 149, 0, 0.15);
}

.status-alert-error {
  background: rgba(255, 59, 48, 0.08);
  color: #BF2C00;
  border: 1px solid rgba(255, 59, 48, 0.15);
}

/* Text Colors */
.text-warning {
  color: #C93400;
  font-weight: 600;
}

.text-danger {
  color: #BF2C00;
  font-weight: 600;
}

/* Password Form */
.password-form {
  padding: var(--space-xl);
}

.password-form :deep(.el-form-item) {
  margin-bottom: var(--space-lg);
}

.password-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--color-text-primary);
  padding-bottom: 6px;
}

.password-form :deep(.el-input__wrapper) {
  padding: 10px 14px;
  background: var(--apple-gray-50);
  border-radius: var(--radius-md);
  box-shadow: none !important;
  border: 1px solid transparent;
  transition: all var(--transition-fast);
}

.password-form :deep(.el-input__wrapper:hover) {
  background: var(--apple-white);
  border-color: var(--apple-gray-200);
}

.password-form :deep(.el-input__wrapper.is-focus) {
  background: var(--apple-white);
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light) !important;
}

/* Strength Indicator */
.strength-indicator {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  margin-top: var(--space-sm);
}

.strength-bars {
  display: flex;
  gap: 4px;
  flex: 1;
}

.strength-bar {
  height: 4px;
  flex: 1;
  background: var(--apple-gray-100);
  border-radius: 2px;
  transition: background var(--transition-fast);
}

.strength-bar.active {
  background: currentColor;
}

.strength-text {
  font-size: var(--font-size-xs);
  font-weight: 600;
  min-width: 24px;
}

/* Password Tips */
.password-tips {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  margin-top: var(--space-sm);
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  line-height: 1.4;
}

.tip-icon {
  flex-shrink: 0;
  margin-top: 1px;
  color: var(--apple-blue);
}

/* Form Actions */
.form-actions {
  display: flex;
  gap: var(--space-md);
  margin-top: var(--space-xl);
}

.btn-primary {
  height: 44px;
  padding: 0 var(--space-xl);
  font-weight: 600 !important;
  border-radius: var(--radius-md) !important;
}

.btn-secondary {
  height: 44px;
  border-radius: var(--radius-md) !important;
}

/* Responsive */
@media (max-width: 768px) {
  .main-content {
    padding: var(--space-lg);
  }
  
  .info-grid {
    grid-template-columns: 1fr;
    gap: var(--space-md);
    padding: var(--space-lg);
  }
  
  .header-content {
    padding: var(--space-md);
  }
  
  .header-title {
    font-size: var(--font-size-base);
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-primary,
  .btn-secondary {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .user-badge span {
    display: none;
  }
  
  .card {
    border-radius: var(--radius-lg);
  }
}
</style>
