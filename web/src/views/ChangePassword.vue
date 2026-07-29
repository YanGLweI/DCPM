<template>
  <div class="change-password-container">
    <el-card class="change-password-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <h2>修改密码</h2>
          <el-tag v-if="isExpired" type="danger" size="large">密码已过期</el-tag>
        </div>
      </template>
      
      <el-alert
        v-if="isExpired"
        title="您的密码已过期"
        description="请修改密码后才能继续使用系统"
        type="error"
        show-icon
        :closable="false"
        style="margin-bottom: 20px"
      />
      
      <el-form
        ref="formRef"
        :model="passwordForm"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="账号" prop="username">
          <el-input
            v-model="passwordForm.username"
            placeholder="请输入域账号"
            disabled
          />
        </el-form-item>
        
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            placeholder="请输入旧密码"
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
          <div class="password-strength" v-if="passwordForm.newPassword">
            <span>密码强度：</span>
            <el-progress
              :percentage="strengthPercentage"
              :color="strengthColor"
              :stroke-width="10"
              style="width: 200px; display: inline-block; vertical-align: middle"
            />
            <span :style="{ color: strengthColor }">{{ strengthText }}</span>
          </div>
          <div class="password-tips">
            <el-text type="info" size="small">
              密码要求：至少14位，包含大写字母、小写字母、数字、特殊字符中的至少3类
            </el-text>
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
        
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            @click="handleChangePassword"
          >
            确认修改
          </el-button>
          <el-button @click="handleBackToLogin">
            返回登录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { changePassword } from '../api'

const router = useRouter()
const route = useRoute()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isExpired = ref(false)

const passwordForm = reactive({
  username: '',
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码强度
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
      strengthColor.value = '#f56c6c'
      strengthText.value = '弱'
      break
    case 3:
      strengthPercentage.value = 60
      strengthColor.value = '#e6a23c'
      strengthText.value = '中'
      break
    case 4:
      strengthPercentage.value = 80
      strengthColor.value = '#409eff'
      strengthText.value = '强'
      break
    case 5:
      strengthPercentage.value = 100
      strengthColor.value = '#67c23a'
      strengthText.value = '很强'
      break
  }
}

// 表单验证规则
const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  username: [
    { required: true, message: '请输入域账号', trigger: 'blur' }
  ],
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

onMounted(() => {
  // 检查是否是从密码过期跳转来的
  if (route.query.expired === 'true') {
    isExpired.value = true
  }
  
  // 从 sessionStorage 获取临时用户名
  const tempUsername = sessionStorage.getItem('temp_username')
  if (tempUsername) {
    passwordForm.username = tempUsername
  }
})

const handleChangePassword = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      await changePassword({
        username: passwordForm.username,
        old_password: passwordForm.oldPassword,
        new_password: passwordForm.newPassword
      })
      
      ElMessage.success('密码修改成功，请重新登录')
      
      // 清除临时数据
      sessionStorage.removeItem('temp_username')
      
      // 跳转登录页
      router.push('/login')
    } catch (error: any) {
      const message = error.response?.data?.message || error.message || '密码修改失败'
      ElMessage.error(message)
    } finally {
      loading.value = false
    }
  })
}

const handleBackToLogin = () => {
  sessionStorage.removeItem('temp_username')
  router.push('/login')
}
</script>

<style scoped>
.change-password-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.change-password-card {
  width: 500px;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
  color: #303133;
  font-size: 24px;
}

.password-strength {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.password-tips {
  margin-top: 4px;
}
</style>
