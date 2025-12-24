<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>用户权限管理</span>
        </div>
      </template>
      
      <div class="filter-container">
        <el-select v-model="currentClubId" placeholder="选择社团" @change="handleClubChange" style="width: 200px; margin-right: 10px;">
          <el-option v-for="club in managedClubs" :key="club.id" :label="club.name" :value="club.id" />
        </el-select>
        <el-input v-model="queryParams.keyword" placeholder="搜索姓名/学号" style="width: 200px; margin-right: 10px;" clearable @keyup.enter="handleQuery" />
        <el-button type="primary" @click="handleQuery">查询</el-button>
      </div>

      <el-table v-loading="loading" :data="list" border style="width: 100%; margin-top: 20px;">
        <el-table-column prop="user.account" label="用户名" width="120" />
        <el-table-column prop="user.name" label="姓名" width="120" />
        <el-table-column prop="user.student_no" label="学号" width="120" />
        <el-table-column prop="user.college" label="学院" min-width="150" show-overflow-tooltip />
        <el-table-column prop="user.phone" label="手机号" width="120" />
        <el-table-column prop="role" label="角色" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getRoleType(scope.row.role)">{{ getRoleLabel(scope.row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="scope">
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="handleEditRole(scope.row)"
              :disabled="!canEdit(scope.row)"
            >
              修改权限
            </el-button>
            <el-button 
              type="danger" 
              link 
              size="small" 
              @click="handleKick(scope.row)"
              :disabled="!canKick(scope.row)"
            >
              删除成员
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.size"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="getList"
        />
      </div>
    </el-card>

    <!-- 权限修改弹窗 -->
    <el-dialog v-model="dialogVisible" title="修改成员权限" width="400px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="当前角色">
          <el-tag :type="getRoleType(currentRow?.role)">{{ getRoleLabel(currentRow?.role) }}</el-tag>
        </el-form-item>
        <el-form-item label="新角色">
          <el-select v-model="form.role" placeholder="请选择角色">
            <el-option 
              v-for="role in availableRoles" 
              :key="role.value" 
              :label="role.label" 
              :value="role.value"
              :disabled="role.disabled"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitRoleChange" :loading="submitLoading">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { getManagedClubs, getClubMembers, updateMemberRole, kickMember } from '@/api/club'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const userStore = useUserStore()
const managedClubs = ref([])
const currentClubId = ref(null)
const loading = ref(false)
const list = ref([])
const total = ref(0)
const currentUserRoleInClub = ref('') // 当前用户在当前选中社团的角色

const queryParams = reactive({
  page: 1,
  size: 10,
  keyword: '',
  status: 'approved' // Only show approved members
})

const dialogVisible = ref(false)
const submitLoading = ref(false)
const currentRow = ref(null)
const form = reactive({
  role: ''
})

const roles = [
  { label: '社长', value: 'leader', level: 3 },
  { label: '管理员', value: 'advisor', level: 2 },
  { label: '普通成员', value: 'member', level: 1 }
]

const getRoleLevel = (role) => {
  const found = roles.find(r => r.value === role)
  return found ? found.level : 0
}

const getRoleLabel = (role) => {
  const found = roles.find(r => r.value === role)
  return found ? found.label : role
}

const getRoleType = (role) => {
  if (role === 'leader') return 'danger'
  if (role === 'advisor') return 'warning'
  return 'info'
}

// Check if current user can edit target user's role
const canEdit = (row) => {
  // School Admin can edit everyone
  if (userStore.userInfo?.role === 'admin') return true
  
  // Club Leader/Admin can only edit users with lower level
  const myLevel = getRoleLevel(currentUserRoleInClub.value)
  const targetLevel = getRoleLevel(row.role)
  
  return myLevel > targetLevel
}

// Check if current user can kick target user
const canKick = (row) => {
  // Leader cannot be kicked by anyone (including admin)
  if (row.role === 'leader') return false

  // School Admin can kick anyone (except leader, checked above)
  if (userStore.userInfo?.role === 'admin') return true
  
  // Club Leader/Admin can only kick users with lower level
  const myLevel = getRoleLevel(currentUserRoleInClub.value)
  const targetLevel = getRoleLevel(row.role)
  
  return myLevel > targetLevel
}

// Available roles to assign: only roles <= my level
const availableRoles = computed(() => {
  const isAdmin = userStore.userInfo?.role === 'admin'
  const myLevel = getRoleLevel(currentUserRoleInClub.value)
  
  return roles.map(r => ({
    ...r,
    // Disable if:
    // 1. Not admin AND role level > my level (cannot assign higher role)
    // 2. Not admin AND role is leader (only admin can assign leader)
    disabled: (!isAdmin && r.level > myLevel) || (!isAdmin && r.value === 'leader')
  }))
})

const fetchManagedClubs = async () => {
  if (!userStore.userInfo?.id) return
  try {
    const res = await getManagedClubs(userStore.userInfo.id, { pageSize: 1000 })
    // Backend returns array directly for getManagedClubs
    managedClubs.value = Array.isArray(res) ? res : (res.list || [])
    if (managedClubs.value.length > 0) {
      currentClubId.value = managedClubs.value[0].id
      updateCurrentRole()
      getList()
    }
  } catch (error) {
    console.error(error)
  }
}

const updateCurrentRole = () => {
  const club = managedClubs.value.find(c => c.id === currentClubId.value)
  if (club) {
    currentUserRoleInClub.value = club.role || 'member'
  }
}

const getList = async () => {
  if (!currentClubId.value) return
  loading.value = true
  try {
    const res = await getClubMembers(currentClubId.value, queryParams)
    list.value = res.list || []
    total.value = res.pagination.total
    
    // Sort by role level desc
    list.value.sort((a, b) => getRoleLevel(b.role) - getRoleLevel(a.role))
  } finally {
    loading.value = false
  }
}

const handleClubChange = () => {
  queryParams.page = 1
  updateCurrentRole()
  getList()
}

const handleQuery = () => {
  queryParams.page = 1
  getList()
}

const handleEditRole = (row) => {
  currentRow.value = row
  form.role = row.role
  dialogVisible.value = true
}

const handleKick = (row) => {
  ElMessageBox.confirm(
    `确定要将成员 "${row.user.name}" 踢出社团吗？此操作不可恢复！`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(async () => {
      try {
        await kickMember(currentClubId.value, row.user.id)
        ElMessage.success('成员已踢出')
        getList()
      } catch (error) {
        console.error(error)
      }
    })
    .catch(() => {})
}

const submitRoleChange = async () => {
  if (!currentRow.value) return
  submitLoading.value = true
  try {
    await updateMemberRole(currentClubId.value, currentRow.value.user.id, { role: form.role })
    ElMessage.success('权限修改成功')
    dialogVisible.value = false
    getList()
  } catch (error) {
    console.error(error)
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  if (!userStore.userInfo) {
    userStore.getUserInfo().then(() => fetchManagedClubs())
  } else {
    fetchManagedClubs()
  }
})
</script>

<style scoped>
.app-container {
  padding: 20px;
}
.filter-container {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
