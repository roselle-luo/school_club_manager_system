<template>
  <div class="app-container">
    <el-card class="filter-container" shadow="never">
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="社团名称">
          <el-select
            v-model="queryParams.club_id"
            placeholder="请选择社团"
            clearable
            filterable
            style="width: 200px"
            @change="handleQuery"
          >
            <el-option
              v-for="item in clubList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="成员姓名">
          <el-input v-model="queryParams.user_name" placeholder="请输入成员姓名" clearable @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item label="学号">
          <el-input v-model="queryParams.student_no" placeholder="请输入学号" clearable @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item label="考勤日期">
          <el-date-picker
            v-model="queryParams.date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            clearable
            @change="handleQuery"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
          <el-button icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-container" shadow="never">
      <el-table
        v-loading="loading"
        :data="list"
        style="width: 100%"
        border
      >
        <el-table-column type="index" label="序号" width="50" align="center" />
        <el-table-column prop="club.name" label="社团名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="user.name" label="成员姓名" width="100" />
        <el-table-column prop="user.student_no" label="成员学号" width="120" />
        <el-table-column prop="user.college" label="成员学院" width="120" show-overflow-tooltip />
        <el-table-column prop="type" label="考勤类型" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.type === 'activity' ? 'success' : 'primary'">
              {{ scope.row.type === 'activity' ? '活动考勤' : '值班考勤' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="考勤时间" width="320">
          <template #default="scope">
             <div>签到: {{ formatDate(scope.row.signin_at) }}</div>
             <div v-if="scope.row.signout_at">签退: {{ formatDate(scope.row.signout_at) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="duration_hours" label="时长(小时)" width="100" align="center" />
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="scope">
            <el-button
              v-if="!scope.row.signout_at"
              type="warning"
              link
              size="small"
              @click="handleForceSignOut(scope.row)"
            >
              强制签退
            </el-button>
            <el-button
              v-else
              type="danger"
              link
              size="small"
              @click="handleDelete(scope.row)"
            >
              删除
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
import { ref, reactive, onMounted } from 'vue'
import { getAllAttendance, forceSignOut, deleteAttendance } from '@/api/admin'
import { getManagedClubs, getClubs } from '@/api/club'
import { useUserStore } from '@/stores/user'
import { Search, Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'

const userStore = useUserStore()
const loading = ref(false)
const list = ref([])
const total = ref(0)
const clubList = ref([])

const queryParams = reactive({
  page: 1,
  size: 10,
  club_id: '',
  user_name: '',
  student_no: '',
  date: ''
})

const getClubList = async () => {
  if (!userStore.userInfo) return
  
  try {
    let res
    if (userStore.userInfo.role.code === 'admin') {
      // 管理员获取所有社团
      res = await getClubs({ pageSize: 1000 })
    } else {
      // 负责人获取管理的社团
      res = await getManagedClubs(userStore.userInfo.id, { pageSize: 1000 })
    }
    
    clubList.value = res.list || []
    
    // 默认选中第一个社团
    if (clubList.value.length > 0) {
      queryParams.club_id = clubList.value[0].id
    }
    
    // 获取到社团列表后再获取考勤列表，确保有默认筛选条件（如果有）
    getList()
  } catch (error) {
    console.error('Failed to fetch clubs:', error)
    // 即使获取社团失败，也尝试获取考勤列表
    getList()
  }
}

const getList = async () => {
  loading.value = true
  try {
    const res = await getAllAttendance(queryParams)
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const handleQuery = () => {
  queryParams.page = 1
  getList()
}

const resetQuery = () => {
  // 不重置社团选择，因为是必选的上下文（或者根据需求重置为默认第一个？）
  // 如果要重置为默认第一个：
  if (clubList.value.length > 0) {
    queryParams.club_id = clubList.value[0].id
  } else {
    queryParams.club_id = ''
  }
  queryParams.user_name = ''
  queryParams.student_no = ''
  queryParams.date = ''
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

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return dayjs(dateStr).format('YYYY-MM-DD HH:mm:ss')
}

const handleForceSignOut = (row) => {
  ElMessageBox.confirm(
    '确认要强制签退该记录吗？系统将以当前时间作为签退时间。',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await forceSignOut(row.id)
      ElMessage.success('操作成功')
      getList()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleDelete = (row) => {
  ElMessageBox.confirm(
    '确认要删除该考勤记录吗？此操作不可恢复！',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteAttendance(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error(error)
    }
  })
}

onMounted(() => {
  // 先获取社团列表，获取成功后会自动调用getList
  getClubList()
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
</style>
