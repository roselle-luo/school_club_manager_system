<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>公告管理</span>
        </div>
      </template>

      <div class="filter-container">
        <el-select v-model="currentClubId" placeholder="选择社团" @change="handleClubChange" style="width: 200px; margin-right: 10px;">
          <el-option v-for="club in managedClubs" :key="club.id" :label="club.name" :value="club.id" />
        </el-select>
        <el-button type="primary" @click="handleAdd" :disabled="!currentClubId">发布公告</el-button>
      </div>

      <el-table v-loading="loading" :data="list" border style="width: 100%; margin-top: 20px;">
        <el-table-column prop="title" label="标题" min-width="150" show-overflow-tooltip />
        <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="scope" label="范围" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.scope === 'public' ? 'success' : 'info'">
              {{ scope.row.scope === 'public' ? '公开' : '内部' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">修改</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
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

    <!-- 公告编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="form.id ? '修改公告' : '发布公告'" width="500px">
      <el-form :model="form" label-width="80px" :rules="rules" ref="formRef">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input type="textarea" v-model="form.content" :rows="4" placeholder="请输入内容" />
        </el-form-item>
        <el-form-item label="范围" prop="scope">
          <el-radio-group v-model="form.scope">
            <el-radio label="public">公开</el-radio>
            <el-radio label="internal">内部</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getManagedClubs, getClubAnnouncements, createAnnouncement, updateAnnouncement, deleteAnnouncement } from '@/api/club'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const userStore = useUserStore()
const managedClubs = ref([])
const currentClubId = ref(null)
const list = ref([])
const total = ref(0)
const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref(null)

const queryParams = reactive({
  page: 1,
  size: 10
})

const form = reactive({
  id: null,
  title: '',
  content: '',
  scope: 'public'
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  scope: [{ required: true, message: '请选择范围', trigger: 'change' }]
}

const getClubs = async () => {
  if (!userStore.userInfo?.id) return
  const res = await getManagedClubs(userStore.userInfo.id, { pageSize: 100 })
  // Backend returns array directly for getManagedClubs
  managedClubs.value = Array.isArray(res) ? res : (res.list || [])
  if (managedClubs.value.length > 0) {
    currentClubId.value = managedClubs.value[0].id
    getList()
  }
}

const handleClubChange = () => {
  queryParams.page = 1
  getList()
}

const getList = async () => {
  if (!currentClubId.value) return
  loading.value = true
  try {
    const res = await getClubAnnouncements(currentClubId.value, queryParams)
    list.value = res.list || []
    total.value = res.pagination.total
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  form.id = null
  form.title = ''
  form.content = ''
  form.scope = 'public'
  dialogVisible.value = true
}

const handleEdit = (row) => {
  form.id = row.id
  form.title = row.title
  form.content = row.content
  form.scope = row.scope
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该公告吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteAnnouncement(currentClubId.value, row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (form.id) {
          await updateAnnouncement(currentClubId.value, form.id, form)
          ElMessage.success('修改成功')
        } else {
          await createAnnouncement(currentClubId.value, form)
          ElMessage.success('发布成功')
        }
        dialogVisible.value = false
        getList()
      } catch (error) {
        console.error(error)
      } finally {
        submitting.value = false
      }
    }
  })
}

onMounted(() => {
  getClubs()
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
