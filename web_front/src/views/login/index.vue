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
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { User, Lock, Key } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)
const captchaCanvas = ref(null)
const generatedCode = ref('')

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
</style>
