<template>
  <div class="account-container">
    <el-row :gutter="20">
      <!-- 用户信息卡片 -->
      <el-col :span="24">
        <el-card class="info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <h3>账号信息</h3>
              <el-button type="danger" @click="handleLogout">退出登录</el-button>
            </div>
          </template>
          
          <el-descriptions :column="2" border>
            <el-descriptions-item label="账号">
              {{ userStore.username }}
            </el-descriptions-item>
            <el-descriptions-item label="密码状态">
              <el-tag :type="passwordStatusType">{{ passwordStatusText }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="密码过期时间">
              <template v-if="userStore.passwordNeverExpires">
                <el-tag type="info">永不过期</el-tag>
              </template>
              <template v-else>
                {{ userStore.passwordExpiresAt || '加载中...' }}
              </template>
            </el-descriptions-item>
            <el-descriptions-item label="剩余天数">
              <template v-if="userStore.passwordNeverExpires">
                <el-tag type="info">-</el-tag>
              </template>
              <template v-else>
                <span :class="{ 'text-warning': userStore.isExpiringSoon, 'text-danger': userStore.isExpired }">
                  {{ userStore.daysRemaining }} 天
                </span>
              </template>
            </el-descriptions-item>
          </el-descriptions>
          
          <el-alert
            v-if="userStore.passwordNeverExpires"
            title="当前账号已设置密码永不过期"
            type="info"
            show-icon
            :closable="false"
            style="margin-top: 16px"
          />
          
          <el-alert
            v-else-if="userStore.isExpiringSoon && !userStore.isExpired"
            :title="`您的密码将在 ${userStore.daysRemaining} 天后过期，请尽快修改`"
            type="warning"
            show-icon
            :closable="false"
            style="margin-top: 16px"
          />
          
          <el-alert
            v-else-if="userStore.isExpired"
            title="您的密码已过期，请立即修改"
            type="error"
            show-icon
            :closable="false"
            style="margin-top: 16px"
          />
        </el-card>
      </el-col>
      
      <!-- 修改密码卡片 -->
      <el-col :span="24" style="margin-top: 20px">
        <el-card class="password-card" shadow="hover">
          <template #header>
            <h3>修改密码</h3>
          </template>
          
          <el-form
            ref="formRef"
            :model="passwordForm"
            :rules="rules"
            label-width="100px"
            style="max-width: 500px"
          >
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
              <el-button @click="resetForm">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
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

// 密码状态
const passwordStatusType = computed(() => {
  if (userStore.passwordNeverExpires) return 'info'
  if (userStore.isExpired) return 'danger'
  if (userStore.isExpiringSoon) return 'warning'
  return 'success'
})

const passwordStatusText = computed(() => {
  if (userStore.passwordNeverExpires) return '永不过期'
  if (userStore.isExpired) return '已过期'
  if (userStore.isExpiringSoon) return '即将过期'
  return '正常'
})

// 表单验证规则
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
  // 获取最新用户信息
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
      
      // 确认是否继续
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
.account-container {
  min-height: 100vh;
  padding: 20px;
  background: #f5f7fa;
}

.info-card, .password-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  color: #303133;
}

h3 {
  margin: 0;
  color: #303133;
}

.text-warning {
  color: #e6a23c;
  font-weight: bold;
}

.text-danger {
  color: #f56c6c;
  font-weight: bold;
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
