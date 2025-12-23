<template>
  <div class="app-container">
    <div class="page-header">
      <el-page-header @back="goBack">
        <template #content>
          <span class="text-large font-600 mr-3">活动详情</span>
        </template>
      </el-page-header>
    </div>

    <el-card v-loading="loading" shadow="never" class="detail-card">
      <div v-if="activity" class="activity-content">
        <div class="act-header">
          <h1 class="act-title">{{ activity.subject }}</h1>
          <div class="act-meta">
            <el-tag>{{ activity.club?.name }}</el-tag>
            <span class="meta-item"><el-icon><Clock /></el-icon> {{ activity.time }}</span>
            <span class="meta-item"><el-icon><Location /></el-icon> {{ activity.place }}</span>
          </div>
        </div>

        <el-divider />

        <div class="act-desc">
          <p>{{ activity.description || '暂无详细描述' }}</p>
        </div>

        <div class="act-actions">
           <div v-if="registerStatus && registerStatus.registered">
              <el-button type="success" plain disabled>已报名</el-button>
           </div>
           <el-button v-else type="primary" @click="handleRegister" :loading="regLoading">报名参加</el-button>
           
           <div class="signin-actions" v-if="registerStatus && registerStatus.registered">
             <el-divider direction="vertical" />
             <el-button type="primary" @click="handleSignIn" :disabled="hasSignedIn">签到</el-button>
             <el-button type="warning" @click="handleSignOut" :disabled="!hasSignedIn || hasSignedOut">签退</el-button>
           </div>
        </div>
      </div>
      <el-empty v-else description="活动不存在" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getActivityDetail } from '@/api/public'
import { registerActivity, getRegisterStatus, signInActivity, signOutActivity } from '@/api/club' // Need to add these to api/club.js
import { Clock, Location } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const id = route.params.id

const loading = ref(false)
const activity = ref(null)
const registerStatus = ref(null)
const regLoading = ref(false)
const hasSignedIn = ref(false) // Simplified state logic
const hasSignedOut = ref(false)

const getDetail = async () => {
  loading.value = true
  try {
    const res = await getActivityDetail(id)
    activity.value = res
    checkStatus()
  } finally {
    loading.value = false
  }
}

const checkStatus = async () => {
  try {
    const res = await getRegisterStatus(id)
    registerStatus.value = res // Assuming backend returns null if not registered or struct
  } catch (e) {
    // Not registered or error
    registerStatus.value = null
  }
}

const handleRegister = () => {
  ElMessageBox.confirm('确定报名参加该活动吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    regLoading.value = true
    try {
      await registerActivity(id)
      ElMessage.success('报名成功')
      checkStatus()
    } finally {
      regLoading.value = false
    }
  })
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  getDetail()
})
</script>

<style scoped>
.app-container {
  padding: 20px;
}
.page-header {
  margin-bottom: 20px;
}
.detail-card {
  min-height: 400px;
}
.act-header {
  text-align: center;
}
.act-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 15px;
}
.act-meta {
  display: flex;
  justify-content: center;
  gap: 20px;
  color: #666;
  font-size: 14px;
}
.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
}
.act-desc {
  padding: 20px 0;
  line-height: 1.8;
  color: #606266;
  font-size: 16px;
}
.act-actions {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 30px;
  gap: 15px;
}
</style>
