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
        <el-form-item label="搜索">
          <el-input 
            v-model="queryParams.keyword" 
            placeholder="请输入姓名或学号" 
            clearable 
            @keyup.enter="handleQuery" 
            style="width: 200px"
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
        <el-table-column prop="user.name" label="申请人姓名" width="120" />
        <el-table-column prop="user.student_no" label="学号" width="120" />
        <el-table-column prop="user.college" label="学院" width="150" show-overflow-tooltip />
        <el-table-column prop="user.phone" label="联系电话" width="120" />
        <el-table-column prop="created_at" label="申请时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag type="warning">待审核</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="scope">
            <el-button
              type="success"
              link
              size="small"
              @click="handleApprove(scope.row)"
            >
              同意
            </el-button>
            <el-button
              type="danger"
              link
              size="small"
              @click="handleReject(scope.row)"
            >
              拒绝
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
import { getManagedClubs, getClubs, getPendingMemberships, approveMembership, rejectMembership } from '@/api/club'
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
  keyword: ''
})

const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const getClubList = async () => {
  if (!userStore.userInfo) return
  
  try {
    let res
    if (userStore.userInfo.role.code === 'admin') {
      res = await getClubs({ pageSize: 1000 })
      clubList.value = res.list || []
    } else {
      res = await getManagedClubs(userStore.userInfo.id, { pageSize: 1000 })
      // Backend returns array directly for getManagedClubs
      clubList.value = Array.isArray(res) ? res : (res.list || [])
    }
    
    // 默认选中第一个社团
    if (clubList.value.length > 0) {
      queryParams.club_id = clubList.value[0].id
    }
    
    getList()
  } catch (error) {
    console.error('Failed to fetch clubs:', error)
    getList()
  }
}

const getList = async () => {
  if (!queryParams.club_id) {
    list.value = []
    total.value = 0
    return
  }

  loading.value = true
  try {
    const res = await getPendingMemberships(queryParams.club_id, {
      page: queryParams.page,
      pageSize: queryParams.size,
      keyword: queryParams.keyword
    })
    list.value = res.list || []
    total.value = res.pagination ? res.pagination.total : 0
  } catch (error) {
    console.error('Failed to fetch pending memberships:', error)
    list.value = []
  } finally {
    loading.value = false
  }
}

const handleQuery = () => {
  queryParams.page = 1
  getList()
}

const resetQuery = () => {
  if (clubList.value.length > 0) {
    queryParams.club_id = clubList.value[0].id
  } else {
    queryParams.club_id = ''
  }
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

const handleApprove = (row) => {
  ElMessageBox.confirm(
    `确定要同意 ${row.user.name} 加入社团吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'success',
    }
  ).then(async () => {
    try {
      await approveMembership(queryParams.club_id, row.id)
      ElMessage.success('操作成功')
      getList()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleReject = (row) => {
  ElMessageBox.confirm(
    `确定要拒绝 ${row.user.name} 加入社团吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await rejectMembership(queryParams.club_id, row.id)
      ElMessage.success('操作成功')
      getList()
    } catch (error) {
      console.error(error)
    }
  })
}

onMounted(() => {
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
