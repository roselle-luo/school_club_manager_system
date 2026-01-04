<template>
  <div class="app-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>入团申请审核</span>
          <el-select v-if="clubs.length > 1" v-model="currentClubId" placeholder="选择社团" @change="handleClubChange">
            <el-option
              v-for="item in clubs"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
          <span v-else-if="clubs.length === 1" class="current-club">{{ clubs[0].name }}</span>
        </div>
      </template>

      <el-empty v-if="!currentClubId && clubs.length === 0" description="您没有管理的社团" />

      <el-table v-else v-loading="loading" :data="list" style="width: 100%">
        <el-table-column prop="user.name" label="姓名" width="120" />
        <el-table-column prop="user.student_no" label="学号" width="150" />
        <el-table-column prop="user.college" label="学院" width="150" />
        <el-table-column prop="created_at" label="申请时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="success" size="small" @click="handleAudit(scope.row, 'approved')">通过</el-button>
            <el-button type="danger" size="small" @click="handleAudit(scope.row, 'rejected')">拒绝</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getManagedClubs, getPendingMemberships, approveMembership, rejectMembership } from '@/api/club'
import { ElMessage, ElMessageBox } from 'element-plus'

const userStore = useUserStore()
const loading = ref(false)
const clubs = ref([])
const currentClubId = ref(null)
const list = ref([])

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString()
}

const fetchManagedClubs = async () => {
  if (!userStore.userInfo?.id) return
  try {
    const res = await getManagedClubs(userStore.userInfo.id)
    clubs.value = res || []
    if (clubs.value.length > 0) {
        currentClubId.value = clubs.value[0].id
        getList()
    }
  } catch (error) {
    console.error(error)
  }
}

const getList = async () => {
  if (!currentClubId.value) return
  loading.value = true
  try {
    const res = await getPendingMemberships(currentClubId.value)
    // 后端返回结构为 { list: [], pagination: {} }
    list.value = res?.list || []
  } catch (error) {
    console.error(error)
    list.value = []
  } finally {
    loading.value = false
  }
}

const handleClubChange = () => {
  getList()
}

const handleAudit = (row, status) => {
  const actionText = status === 'approved' ? '通过' : '拒绝'
  ElMessageBox.confirm(`确定要${actionText} "${row.user.name}" 的入团申请吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: status === 'approved' ? 'success' : 'warning'
  }).then(async () => {
    try {
      if (status === 'approved') {
          await approveMembership(currentClubId.value, row.id)
      } else {
          await rejectMembership(currentClubId.value, row.id)
      }
      
      ElMessage.success('操作成功')
      getList()
    } catch (error) {
      // Error handled by interceptor
    }
  }).catch(() => {})
}

onMounted(() => {
  fetchManagedClubs()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.current-club {
    font-weight: bold;
    color: #409EFF;
}
</style>
