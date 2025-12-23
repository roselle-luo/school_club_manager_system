<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card class="box-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span><el-icon><Bell /></el-icon> 最新公告</span>
            </div>
          </template>
          <div v-loading="loadingAnn">
             <div v-if="announcements.length > 0" class="list-wrapper">
               <div v-for="item in announcements" :key="item.id" class="list-item">
                 <span class="item-title text-ellipsis">{{ item.title }}</span>
                 <span class="item-club text-muted">{{ item.club?.name }}</span>
                 <span class="item-time text-muted">{{ formatDate(item.created_at) }}</span>
               </div>
             </div>
             <el-empty v-else description="暂无公告" />
          </div>
        </el-card>

        <el-card class="box-card" shadow="hover" style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span><el-icon><Calendar /></el-icon> 近期活动</span>
            </div>
          </template>
           <div v-loading="loadingAct">
             <div v-if="activities.length > 0" class="list-wrapper">
               <div v-for="item in activities" :key="item.id" class="list-item" @click="openActivity(item.id)" style="cursor: pointer;">
                 <span class="item-title text-ellipsis">{{ item.subject }}</span>
                 <span class="item-club text-muted">{{ item.club?.name }}</span>
                 <span class="item-time text-muted">{{ item.time }}</span>
               </div>
             </div>
             <el-empty v-else description="暂无活动" />
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card class="box-card user-card" shadow="hover">
          <div class="user-info">
            <el-avatar :size="80" :src="userStore.userInfo?.avatar" />
            <h3>{{ userStore.userInfo?.name }}</h3>
            <p>{{ userStore.userInfo?.college }} - {{ userStore.userInfo?.student_no }}</p>
          </div>
          <el-divider />
          <div class="user-stats">
             <div class="stat-item">
               <div class="stat-val">{{ userStore.userInfo?.club_count || 0 }}</div>
               <div class="stat-label">已加入社团</div>
             </div>
             <div class="stat-item">
               <div class="stat-val">{{ userStore.userInfo?.activity_count || 0 }}</div>
               <div class="stat-label">参与活动</div>
             </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAnnouncements, getActivities } from '@/api/public'
import { useUserStore } from '@/stores/user'
import { Bell, Calendar } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore()
const loadingAnn = ref(false)
const announcements = ref([])
const loadingAct = ref(false)
const activities = ref([])

const getAnnList = async () => {
  loadingAnn.value = true
  try {
    const res = await getAnnouncements({ pageSize: 5 })
    announcements.value = res.list
  } finally {
    loadingAnn.value = false
  }
}

const getActList = async () => {
  loadingAct.value = true
  try {
    const res = await getActivities({ pageSize: 5 })
    activities.value = res.list
  } finally {
    loadingAct.value = false
  }
}

const openActivity = (id) => {
  router.push(`/activities/${id}`)
}

const formatDate = (str) => {
  return dayjs(str).format('MM-DD HH:mm')
}

onMounted(async () => {
  if (!userStore.userInfo) {
    await userStore.getUserInfo()
  }
  getAnnList()
  getActList()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}
.box-card {
  margin-bottom: 20px;
}
.card-header {
  display: flex;
  align-items: center;
  font-weight: bold;
}
.card-header .el-icon {
  margin-right: 8px;
}
.list-item {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}
.list-item:last-child {
  border-bottom: none;
}
.item-title {
  flex: 1;
  font-size: 14px;
  color: #333;
}
.item-club {
  width: 120px;
  text-align: right;
  margin-right: 15px;
}
.item-time {
  width: 100px;
  text-align: right;
  font-size: 12px;
}
.text-muted {
  color: #999;
}
.text-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.user-card {
  text-align: center;
}
.user-info h3 {
  margin: 10px 0 5px;
}
.user-info p {
  color: #666;
  font-size: 13px;
}
.user-stats {
  display: flex;
  justify-content: space-around;
}
.stat-val {
  font-size: 20px;
  font-weight: bold;
  color: #409EFF;
}
.stat-label {
  font-size: 12px;
  color: #666;
  margin-top: 5px;
}
</style>
