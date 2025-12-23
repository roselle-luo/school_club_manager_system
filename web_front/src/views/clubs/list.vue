<template>
  <div class="app-container">
    <el-card class="filter-container" shadow="never">
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="社团名称">
          <el-input v-model="queryParams.keyword" placeholder="请输入社团名称" clearable @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <el-table v-loading="loading" :data="list" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column label="社团LOGO" width="100" align="center">
          <template #default="scope">
            <el-avatar :size="50" :src="scope.row.logo || defaultLogo" shape="square" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="社团名称" min-width="150" />
        <el-table-column prop="category.name" label="社团分类" width="120" align="center" />
        <el-table-column prop="role" label="我的角色" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getRoleType(scope.row.role)">{{ getRoleLabel(scope.row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="intro" label="简介" min-width="200" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.intro || '暂无简介' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="scope">
            <el-button 
              type="danger" 
              link 
              :disabled="!isSchoolAdmin"
              @click="handleDissolve(scope.row)"
            >
              解散社团
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.size"
          :page-sizes="[10, 20, 30, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { getManagedClubs, dissolveClub } from '@/api/club'
import { Search, Refresh } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox, ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const list = ref([])
const total = ref(0)
const defaultLogo = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const isSchoolAdmin = computed(() => userStore.userInfo?.role === 'admin')

const queryParams = reactive({
  page: 1,
  size: 10,
  keyword: ''
})

const currentClub = ref(null)

const getRoleLabel = (role) => {
  const map = {
    'leader': '社长',
    'advisor': '管理员',
    'member': '成员',
    'admin': '学校管理员'
  }
  return map[role] || role
}

const getRoleType = (role) => {
  if (role === 'leader') return 'danger'
  if (role === 'advisor') return 'warning'
  if (role === 'admin') return 'success'
  return 'info'
}

const getList = async () => {
  if (!userStore.userInfo?.id) return
  loading.value = true
  try {
    const res = await getManagedClubs(userStore.userInfo.id, queryParams)
    list.value = res.list || []
    total.value = res.pagination.total
  } finally {
    loading.value = false
  }
}

const handleQuery = () => {
  queryParams.page = 1
  getList()
}

const resetQuery = () => {
  queryParams.keyword = ''
  handleQuery()
}

const handleSizeChange = (val) => {
  queryParams.size = val
  getList()
}

const handleCurrentChange = (val) => {
  queryParams.page = val
  getList()
}

const handleDissolve = (club) => {
  ElMessageBox.confirm(
    `确定要解散社团"${club.name}"吗？此操作不可恢复！`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(async () => {
      try {
        await dissolveClub(club.id)
        ElMessage.success('社团已解散')
        getList()
      } catch (error) {
        console.error(error)
      }
    })
    .catch(() => {})
}

onMounted(() => {
  if (userStore.userInfo?.id) {
    getList()
  } else {
    // Wait for user info if not ready
    const unwatch = userStore.$subscribe((mutation, state) => {
      if (state.userInfo?.id) {
        getList()
        unwatch()
      }
    })
  }
})
</script>

<style scoped>
.app-container {
  padding: 20px;
}
.filter-container {
  margin-bottom: 20px;
}
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
.detail-header {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}
.detail-logo {
  width: 100px;
  height: 100px;
  border-radius: 8px;
  object-fit: cover;
}
.detail-info p {
  margin: 8px 0;
  color: #606266;
}
</style>
