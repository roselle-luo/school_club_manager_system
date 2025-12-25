<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-left">
        <div class="login-info">
          <h2>高校社团管理系统</h2>
          <p>高效便捷的社团管理平台</p>
        </div>
      </div>
      <div class="login-right">
        <div class="login-form-box">
          <div class="title">欢迎登录</div>
          <el-form :model="loginForm" :rules="rules" ref="loginFormRef" size="large" class="login-form">
            <el-form-item prop="account">
              <el-input v-model="loginForm.account" placeholder="用户名/学号" :prefix-icon="User" />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="密码"
                :prefix-icon="Lock"
                show-password
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            <el-form-item prop="captcha">
              <div class="captcha-box">
                <el-input 
                  v-model="loginForm.captcha" 
                  placeholder="验证码" 
                  :prefix-icon="Key"
                  @keyup.enter="handleLogin"
                  class="captcha-input"
                />
                <canvas 
                  ref="captchaCanvas"
                  class="captcha-img" 
                  @click="generateCaptcha" 
                  width="120"
                  height="40"
                  title="点击刷新"
                ></canvas>
              </div>
            </el-form-item>
            <div class="form-options">
               <el-checkbox v-model="loginForm.remember">记住密码</el-checkbox>
               <a class="register-link" @click="showRegisterDialog">点此申请注册社团 -></a>
            </div>
            <el-form-item>
              <el-button type="primary" class="login-btn" :loading="loading" @click="handleLogin">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
    <el-dialog v-model="registerDialogVisible" title="社团注册" width="500px">
      <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" label-width="80px">
        <div class="logo-upload-container">
          <el-upload
            class="avatar-uploader"
            action="/api/v1/public/upload/image"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload"
            name="file"
          >
            <img v-if="registerForm.logo" :src="registerForm.logo" class="avatar" />
            <div v-else class="upload-placeholder">
              <el-icon class="avatar-uploader-icon"><Plus /></el-icon>
              <div class="upload-text">上传logo</div>
            </div>
          </el-upload>
        </div>
        <el-form-item label="社团名称" prop="name">
          <el-input v-model="registerForm.name" />
        </el-form-item>
        <el-form-item label="社团分类" prop="category_id">
          <el-select v-model="registerForm.category_id" placeholder="请选择">
            <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="registerForm.contact" />
        </el-form-item>
        <el-form-item label="社团简介" prop="intro">
          <el-input v-model="registerForm.intro" type="textarea" />
        </el-form-item>
        <el-divider>申请人验证</el-divider>
        <el-form-item label="账号" prop="account">
          <el-input v-model="registerForm.account" placeholder="请输入您的学号/工号" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入登录密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="registerDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleRegister">提交</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { User, Lock, Key, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getCategories, registerClub } from '@/api/public'

const userStore = useUserStore()
const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)
const captchaCanvas = ref(null)
const generatedCode = ref('')

const registerDialogVisible = ref(false)
const registerFormRef = ref(null)
const categories = ref([])
const registerForm = reactive({
  name: '',
  logo: '',
  intro: '',
  contact: '',
  category_id: '',
  account: '',
  password: ''
})
const registerRules = {
  name: [{ required: true, message: '请输入社团名称', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择社团分类', trigger: 'change' }],
  account: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const showRegisterDialog = async () => {
  registerDialogVisible.value = true
  try {
    const res = await getCategories()
    categories.value = res.list || []
  } catch (e) {
    console.error(e)
  }
}

const handleAvatarSuccess = (response) => {
  if (response.code === 0) {
    registerForm.logo = response.data.url
  } else {
    ElMessage.error(response.msg || '上传失败')
  }
}

const beforeAvatarUpload = (rawFile) => {
  const isImg = rawFile.type === 'image/jpeg' || rawFile.type === 'image/png'
  const isLt2M = rawFile.size / 1024 / 1024 < 2

  if (!isImg) {
    ElMessage.error('Avatar picture must be JPG/PNG format!')
  }
  if (!isLt2M) {
    ElMessage.error('Avatar picture size can not exceed 2MB!')
  }
  return isImg && isLt2M
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await registerClub(registerForm)
        ElMessage.success('申请提交成功')
        registerDialogVisible.value = false
        // reset
        registerForm.name = ''
        registerForm.logo = ''
        registerForm.intro = ''
        registerForm.contact = ''
        registerForm.category_id = ''
        registerForm.account = ''
        registerForm.password = ''
      } catch (e) {
        // handled by request interceptor
      }
    }
  })
}

const loginForm = reactive({
  account: '',
  password: '',
  captcha: '',
  remember: false
})

const rules = {
  account: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
}

const generateCaptcha = () => {
  const canvas = captchaCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  const width = canvas.width
  const height = canvas.height
  
  // Clear canvas
  ctx.clearRect(0, 0, width, height)
  
  // Draw background (flower pattern simulation)
  // Random circles/petals
  for(let i = 0; i < 7; i++) {
    ctx.beginPath()
    ctx.arc(
      Math.random() * width, 
      Math.random() * height, 
      Math.random() * 20 + 10, 
      0, 
      2 * Math.PI
    )
    ctx.fillStyle = `rgba(${Math.random()*255}, ${Math.random()*255}, ${Math.random()*255}, 0.2)`
    ctx.fill()
  }
  
  // Draw random lines (interference)
  for(let i = 0; i < 5; i++) {
    ctx.beginPath()
    ctx.moveTo(Math.random() * width, Math.random() * height)
    ctx.lineTo(Math.random() * width, Math.random() * height)
    ctx.strokeStyle = `rgba(${Math.random()*255}, ${Math.random()*255}, ${Math.random()*255}, 0.5)`
    ctx.stroke()
  }

  // Generate code
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  let code = ''
  for(let i = 0; i < 4; i++) {
    code += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  generatedCode.value = code

  // Draw code
  ctx.font = 'bold 24px Arial'
  ctx.textBaseline = 'middle'
  for(let i = 0; i < 4; i++) {
    const x = (width / 4) * i + 10
    const y = height / 2
    const deg = Math.random() * 30 * Math.PI / 180 // rotation
    const sign = Math.random() > 0.5 ? 1 : -1
    
    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(deg * sign)
    ctx.fillStyle = `rgb(${Math.random()*150}, ${Math.random()*150}, ${Math.random()*150})` // darker color for text
    ctx.fillText(code[i], 0, 0)
    ctx.restore()
  }
}

onMounted(() => {
  generateCaptcha()
})

const handleLogin = async () => {
  if (!loginFormRef.value) return
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      // Verify captcha locally
      if (loginForm.captcha.toUpperCase() !== generatedCode.value) {
        ElMessage.error('验证码错误')
        generateCaptcha()
        loginForm.captcha = ''
        return
      }

      loading.value = true
      try {
        const success = await userStore.login({
          account: loginForm.account,
          password: loginForm.password
        })
        if (success) {
          ElMessage.success('登录成功')
          router.push('/')
        } else {
          generateCaptcha()
        }
      } catch (error) {
        generateCaptcha()
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f0f2f5;
  background-image: url('https://gw.alipayobjects.com/zos/rmsportal/TVYTbAXWheQpRcWDaDMu.svg'); 
  background-repeat: no-repeat;
  background-position: center 110px;
  background-size: 100%;
}

.login-content {
  display: flex;
  width: 900px;
  height: 500px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.login-left {
  flex: 1;
  background: linear-gradient(135deg, #304156 0%, #1f2d3d 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  color: white;
  flex-direction: column;
}

.login-info h2 {
  font-size: 28px;
  margin-bottom: 20px;
}

.login-info p {
  font-size: 16px;
  opacity: 0.8;
}

.login-right {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.login-form-box {
  width: 320px;
}

.title {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.login-btn {
  width: 100%;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.captcha-box {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.captcha-input {
  flex: 1;
}

.captcha-img {
  height: 40px;
  cursor: pointer;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.register-link {
  color: #409eff;
  cursor: pointer;
  font-size: 14px;
}

.logo-upload-container {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}
.avatar-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
  background-color: #f5f7fa;
}
.avatar-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}
.upload-placeholder {
  width: 120px;
  height: 120px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #8c939d;
  background-color: #f5f7fa;
}
.avatar-uploader-icon {
  font-size: 28px;
  margin-bottom: 8px;
}
.upload-text {
  font-size: 12px;
}
.avatar {
  width: 120px;
  height: 120px;
  display: block;
  border-radius: 8px;
}
</style>
