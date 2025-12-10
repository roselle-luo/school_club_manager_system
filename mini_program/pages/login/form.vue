<template>
  <view class="page">
    <view class="header"><text class="htitle">登录</text></view>
    <view class="card">
      <view class="avatar"></view>
      <view class="brand">快乐社团平台</view>
      <input class="ipt" placeholder="用户名" v-model="form.account" />
      <input class="ipt" placeholder="密码" password v-model="form.password" />
      <button class="primary" @tap="submit">登录</button>
      <view class="tip">还没有账号？<text class="action" @tap="goRegister">立即注册</text></view>
    </view>
  </view>
</template>

<script>
import { request, setToken } from '../../utils/request.js'
export default {
  data() { return { form: { account: '', password: '' } } },
  methods: {
    async submit() {
      try {
        const data = await request({ url: '/public/login', method: 'POST', data: this.form })
        if (data && data.token) {
          setToken(data.token)
          uni.showToast({ title: '登录成功' })
          uni.switchTab({ url: '/pages/mine/memberships' })
        } else {
          uni.showToast({ title: '登录失败', icon: 'none' })
        }
      } catch(e) { uni.showToast({ title: e.msg || '登录失败', icon: 'none' }) }
    }
    ,
    goRegister() { uni.redirectTo({ url: '/pages/register/form' }) }
  }
}
</script>

<style>
.page { min-height:100vh; display:flex; flex-direction:column; background:#f7f8fa; padding:0 16px }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.card { width:92%; max-width:360px; margin:auto 0; background:#fff; border-radius:14px; padding:20px; box-shadow:0 8px 20px rgba(0,0,0,0.06); display:flex; flex-direction:column; align-items:center }
.avatar { width:72px; height:72px; border-radius:50%; background:#7e78ff; opacity:0.9 }
.brand { margin:10px 0 14px 0; color:#333; font-weight:600 }
.ipt { width:100%; border:1px solid #e6e6e6; border-radius:12px; padding:12px 16px; font-size:14px; background:#fafafa; margin:10px 0 }
.primary { width:100%; margin-top:10px; background:#7e78ff; color:#fff; border:none; border-radius:12px; padding:3px; font-weight:600 }
.tip { margin-top:12px; color:#666 }
.action { color:#7e78ff }
</style>
