<template>
  <div class="app-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>社团创建审批</span>
          <el-radio-group v-model="filterStatus" size="small" @change="handleFilterChange">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="pending">待审核</el-radio-button>
            <el-radio-button label="approved">已通过</el-radio-button>
            <el-radio-button label="rejected">已拒绝</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <el-table v-loading="loading" :data="list" style="width: 100%">
        <el-table-column label="Logo" width="80">
          <template #default="scope">
            <el-image
              v-if="scope.row.logo"
              style="width: 50px; height: 50px"
              :src="scope.row.logo"
              :preview-src-list="[scope.row.logo]"
              fit="cover"
            />
            <span v-else>无</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="社团名称" width="150" />
        <el-table-column prop="category.name" label="分类" width="100" />
        <el-table-column prop="intro" label="简介" show-overflow-tooltip />
        <el-table-column prop="contact" label="联系方式" width="150" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ getStatusLabel(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="scope">
            <div v-if="scope.row.status === 'pending'">
              <el-button type="success" size="small" @click="handleAudit(scope.row, 'approved')">通过</el-button>
              <el-button type="danger" size="small" @click="handleAudit(scope.row, 'rejected')">拒绝</el-button>
            </div>
            <span v-else class="text-gray-500 text-sm">已完成</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPendingClubs, auditClub } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const list = ref([])
const filterStatus = ref('pending') // Default to pending

const getList = async () => {
  loading.value = true
  try {
    const params = {}
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    const res = await getPendingClubs(params)
    list.value = res || []
  } finally {
    loading.value = false
  }
}

const handleFilterChange = () => {
  getList()
}

const getStatusType = (status) => {
  const map = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

const getStatusLabel = (status) => {
  const map = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return map[status] || status
}

const handleAudit = (row, status) => {
  const actionText = status === 'approved' ? '通过' : '拒绝'
  ElMessageBox.confirm(`确定要${actionText}社团 "${row.name}" 的创建申请吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: status === 'approved' ? 'success' : 'warning'
  }).then(async () => {
    try {
      await auditClub(row.id, status)
      ElMessage.success('操作成功')
      getList()
    } catch (e) {
      // handled by interceptor
    }
  }).catch(() => {})
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
