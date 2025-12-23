<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="box-card">
          <div class="user-profile">
            <div class="box-center">
              <el-upload
                class="avatar-uploader"
                action=""
                :http-request="handleUpload"
                :show-file-list="false"
                :before-upload="beforeAvatarUpload"
              >
                <img v-if="userStore.userInfo?.avatar" :src="userStore.userInfo.avatar" class="user-avatar" />
                <div v-else class="user-avatar-placeholder">
                  <span>点击上传头像</span>
                </div>
              </el-upload>
            </div>
            <div class="box-center">
              <div class="user-name text-center">{{ userStore.userInfo?.name || 'User' }}</div>
              <div class="user-role text-center text-muted">{{ userStore.userInfo?.role || 'Student' }}</div>
            </div>
          </div>

          <div class="user-bio">
            <div class="user-education user-bio-section">
              <div class="user-bio-section-header">
                <span><el-icon><User /></el-icon> 用户信息</span>
              </div>
              <div class="user-bio-section-body">
                <div class="text-muted">
                  <div class="info-item">
                    <span class="label">账户:</span>
                    <span>{{ userStore.userInfo?.account }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">手机:</span>
                    <span>{{ userStore.userInfo?.phone || '暂无' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">学号:</span>
                    <span>{{ userStore.userInfo?.student_no || '暂无' }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">学院:</span>
                    <span>{{ userStore.userInfo?.college || '暂无' }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <el-tabs v-model="activeTab">
            <el-tab-pane label="基本资料" name="info">
              <el-form :model="infoForm" :rules="infoRules" ref="infoFormRef" label-width="100px">
                <el-form-item label="姓名" prop="name">
                  <el-input v-model="infoForm.name" />
                </el-form-item>
                <el-form-item label="性别" prop="gender">
                  <el-select v-model="infoForm.gender" placeholder="请选择">
                    <el-option label="男" value="male" />
                    <el-option label="女" value="female" />
                    <el-option label="保密" value="secret" />
                  </el-select>
                </el-form-item>
                <el-form-item label="学院" prop="college">
                  <el-input v-model="infoForm.college" />
                </el-form-item>
                <el-form-item label="学号" prop="student_no">
                  <el-input v-model="infoForm.student_no" />
                </el-form-item>
                <el-form-item label="手机" prop="phone">
                  <el-input v-model="infoForm.phone" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="handleUpdateInfo" :loading="loading">保存</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <el-tab-pane label="修改密码" name="password">
              <el-form :model="pwdForm" :rules="pwdRules" ref="pwdFormRef" label-width="100px">
                <el-form-item label="旧密码" prop="old_password">
                  <el-input v-model="pwdForm.old_password" type="password" show-password />
                </el-form-item>
                <el-form-item label="新密码" prop="new_password">
                  <el-input v-model="pwdForm.new_password" type="password" show-password />
                </el-form-item>
                <el-form-item label="确认密码" prop="confirm_password">
                  <el-input v-model="pwdForm.confirm_password" type="password" show-password />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="handleChangePassword" :loading="loading">保存</el-button>
                  <el-button @click="resetPwdForm">重置</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { updateMyProfile, changePassword } from '@/api/student'
import { uploadImage } from '@/api/common'
import { ElMessage } from 'element-plus'
import { User } from '@element-plus/icons-vue'

const userStore = useUserStore()
const activeTab = ref('info')
const loading = ref(false)

const infoForm = reactive({
  name: '',
  gender: '',
  college: '',
  student_no: '',
  phone: ''
})

const infoRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }]
}

const pwdForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validatePass2 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== pwdForm.new_password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const pwdRules = {
  old_password: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  new_password: [{ required: true, message: '请输入新密码', trigger: 'blur' }],
  confirm_password: [{ validator: validatePass2, trigger: 'blur' }]
}

const infoFormRef = ref(null)
const pwdFormRef = ref(null)

onMounted(async () => {
  if (!userStore.userInfo) {
    await userStore.getUserInfo()
  }
  Object.assign(infoForm, {
    name: userStore.userInfo.name,
    gender: userStore.userInfo.gender,
    college: userStore.userInfo.college,
    student_no: userStore.userInfo.student_no,
    phone: userStore.userInfo.phone
  })
})

const handleUpdateInfo = async () => {
  if (!infoFormRef.value) return
  await infoFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await updateMyProfile(infoForm)
        ElMessage.success('保存成功')
        await userStore.getUserInfo()
      } finally {
        loading.value = false
      }
    }
  })
}

const handleChangePassword = async () => {
  if (!pwdFormRef.value) return
  await pwdFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await changePassword({
          old_password: pwdForm.old_password,
          new_password: pwdForm.new_password
        })
        ElMessage.success('修改成功，请重新登录')
        userStore.logout()
        // redirect to login? handled by store or router guard mostly, but logout clears token
        window.location.reload()
      } catch (error) {
        // handled by request interceptor
      } finally {
        loading.value = false
      }
    }
  })
}

const resetPwdForm = () => {
  if (pwdFormRef.value) pwdFormRef.value.resetFields()
}

const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('上传头像图片只能是 JPG/PNG 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('上传头像图片大小不能超过 2MB!')
  }
  return isJPG && isLt2M
}

const handleUpload = async (options) => {
  const formData = new FormData()
  formData.append('file', options.file)
  try {
    const res = await uploadImage(formData)
    await updateMyProfile({ avatar: res.url })
    ElMessage.success('头像上传成功')
    await userStore.getUserInfo()
  } catch (error) {
    // Error handled by interceptor
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}
.box-center {
  margin: 0 auto;
  display: table;
}
.text-center {
  text-align: center;
}
.text-muted {
  color: #777;
}
.user-profile {
  .user-name {
    font-weight: bold;
  }
  .box-center {
    padding-top: 10px;
  }
  .user-role {
    padding-top: 10px;
    font-weight: 400;
    font-size: 14px;
  }
  .box-social {
    padding-top: 30px;
    .el-table {
      border-top: 1px solid #dfe6ec;
    }
  }
  .user-bio {
    margin-top: 20px;
    color: #606266;
    span {
      padding-left: 4px;
    }
    .user-bio-section {
      font-size: 14px;
      padding: 15px 0;
      .user-bio-section-header {
        border-bottom: 1px solid #dfe6ec;
        padding-bottom: 10px;
        margin-bottom: 10px;
        font-weight: bold;
      }
    }
  }
}

.user-avatar {
  cursor: pointer;
  width: 100px;
  height: 100px;
  border-radius: 50%;
}
.user-avatar-placeholder {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  cursor: pointer;
}
.info-item {
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
}
.label {
  font-weight: bold;
}
</style>
