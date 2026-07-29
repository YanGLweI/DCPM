<template>
  <div class="login-page">
    <!-- Background with subtle gradient -->
    <div class="login-bg">
      <div class="bg-blob bg-blob-1"></div>
      <div class="bg-blob bg-blob-2"></div>
      <div class="bg-blob bg-blob-3"></div>
    </div>

    <!-- Login Card -->
    <div class="login-content">
      <div class="login-card">
        <!-- Header -->
        <div class="login-header">
          <div class="logo-icon">
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
              <rect width="48" height="48" rx="12" fill="#007AFF"/>
              <path d="M24 14C19.58 14 16 17.58 16 22C16 26.42 19.58 30 24 30C28.42 30 32 26.42 32 22C32 17.58 28.42 14 24 14ZM24 27C21.24 27 19 24.76 19 22C19 19.24 21.24 17 24 17C26.76 17 29 19.24 29 22C29 24.76 26.76 27 24 27ZM24 32C19.33 32 10 34.34 10 39V41H38V39C38 34.34 28.67 32 24 32Z" fill="white"/>
            </svg>
          </div>
          <h1 class="login-title">域账号自助管理平台</h1>
          <p class="login-subtitle">请使用域账号登录系统</p>
        </div>

        <!-- Form -->
        <el-form
          ref="formRef"
          :model="loginForm"
          :rules="rules"
          label-width="0"
          class="login-form"
          @submit.prevent="handleLogin"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="域账号"
              :prefix-icon="User"
              size="large"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              :prefix-icon="Lock"
              size="large"
              show-password
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              class="login-button"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>

        <!-- Footer -->
        <div class="login-footer">
          <p class="footer-text">安全连接 · 企业级加密保护</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入域账号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      const result = await userStore.login({
        username: loginForm.username,
        password: loginForm.password
      })
      
      if (result === 'ok') {
        ElMessage.success('登录成功')
        router.push('/account')
      } else if (result === 'expired') {
        ElMessage.warning('密码已过期，请修改密码')
        router.push('/change-password?expired=true')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '登录失败，请检查账号密码')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* Animated Background Blobs */
.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  z-index: 0;
}

.bg-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.6;
  animation: float 20s infinite ease-in-out;
  /* Dynamic sizing based on viewport */
  min-width: 40vmax;
  min-height: 40vmax;
}

.bg-blob-1 {
  width: 60vmax;
  height: 60vmax;
  background: rgba(0, 122, 255, 0.4);
  top: -20vmax;
  left: -20vmax;
  animation-delay: 0s;
}

.bg-blob-2 {
  width: 50vmax;
  height: 50vmax;
  background: rgba(118, 75, 162, 0.4);
  bottom: -15vmax;
  right: -15vmax;
  animation-delay: -7s;
}

.bg-blob-3 {
  width: 45vmax;
  height: 45vmax;
  background: rgba(52, 199, 89, 0.3);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -14s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(3vmax, -3vmax) scale(1.05);
  }
  66% {
    transform: translate(-2vmax, 2vmax) scale(0.95);
  }
}

/* Content */
.login-content {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  padding: var(--space-lg);
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: var(--radius-2xl);
  padding: var(--space-2xl) var(--space-xl);
  box-shadow: 
    0 20px 60px rgba(0, 0, 0, 0.15),
    0 8px 20px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
}

/* Header */
.login-header {
  text-align: center;
  margin-bottom: var(--space-xl);
}

.logo-icon {
  display: inline-flex;
  margin-bottom: var(--space-md);
  animation: fadeInDown 0.6s ease-out;
}

.login-title {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 var(--space-sm);
  letter-spacing: -0.5px;
  animation: fadeInDown 0.6s ease-out 0.1s both;
}

.login-subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
  animation: fadeInDown 0.6s ease-out 0.2s both;
}

/* Form */
.login-form {
  animation: fadeInUp 0.6s ease-out 0.3s both;
}

.login-form :deep(.el-form-item) {
  margin-bottom: var(--space-lg);
}

.login-form :deep(.el-input__wrapper) {
  padding: 12px 16px;
  background: var(--apple-gray-50);
  border-radius: var(--radius-md);
  box-shadow: none !important;
  border: 1px solid transparent;
  transition: all var(--transition-fast);
}

.login-form :deep(.el-input__wrapper:hover) {
  background: var(--apple-white);
  border-color: var(--apple-gray-200);
}

.login-form :deep(.el-input__wrapper.is-focus) {
  background: var(--apple-white);
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light) !important;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: var(--font-size-base);
  font-weight: 600;
  border-radius: var(--radius-md) !important;
  background: var(--apple-blue) !important;
  transition: all var(--transition-fast) !important;
}

.login-button:hover {
  background: var(--apple-blue-hover) !important;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.35);
}

.login-button:active {
  transform: translateY(0);
}

/* Footer */
.login-footer {
  text-align: center;
  margin-top: var(--space-lg);
  padding-top: var(--space-lg);
  border-top: 1px solid var(--color-border);
  animation: fadeInUp 0.6s ease-out 0.4s both;
}

.footer-text {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.footer-text::before {
  content: '';
  display: inline-block;
  width: 6px;
  height: 6px;
  background: var(--apple-green);
  border-radius: 50%;
}

/* Animations */
@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Responsive */
@media (max-width: 480px) {
  .login-content {
    padding: var(--space-md);
  }
  
  .login-card {
    padding: var(--space-xl) var(--space-lg);
    border-radius: var(--radius-xl);
  }
  
  .login-title {
    font-size: var(--font-size-xl);
  }
}
</style>
