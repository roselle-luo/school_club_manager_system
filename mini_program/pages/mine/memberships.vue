<template>
  <view class="page">
    <view class="header"><text class="htitle">我的</text></view>
    <view class="hero">
      <image :src="avatarUrl || defaultAvatar" class="avatar" mode="aspectFill" />
      <text class="uname">{{ profile.name || profile.account || '未登录' }}</text>
      <view class="stats">
        <view class="stat">
          <text class="num">{{ profile.club_count||0 }}</text>
          <text class="label">社团</text>
        </view>
        <view class="stat">
          <text class="num">{{ profile.activity_count||0 }}</text>
          <text class="label">活动</text>
        </view>
      </view>
    </view>
    <view class="panel">
      <view class="ptitle">基本信息</view>
      <view class="item">
        <text class="key">姓名</text><text class="val">{{ profile.name || '未设置' }}</text>
      </view>
      <view class="item">
        <text class="key">邮箱</text><text class="val">{{ profile.email || '未设置' }}</text>
      </view>
      <view class="item">
        <text class="key">手机号</text><text class="val">{{ profile.phone || '未设置' }}</text>
      </view>
      <view class="item">
        <text class="key">性别</text><text class="val">{{ genderText(profile.gender) }}</text>
      </view>
      <view class="item">
        <text class="key">学院</text><text class="val">{{ profile.college || '未设置' }}</text>
      </view>
    </view>
    <view class="btnrow">
      <button class="edit" @tap="editProfile">修改个人信息</button>
      <button class="logout" @tap="logout">退出登录</button>
    </view>
  </view>
</template>

<script>
import { request, clearToken } from '../../utils/request.js'
import { go, relaunchTo } from '../../utils/router.js'
export default {
  data() { return { profile: {}, navLock: false, avatarUrl: '', defaultAvatar: 'https://static-aliyun.oss-cn-hangzhou.aliyuncs.com/placeholder/avatar.png' } },
  onShow() { this.fetch(); this.loadAvatar() },
  methods: {
    async fetch() {
      try { this.profile = await request({ url: '/student/me', method: 'GET' }) || {} } catch(e) { this.profile = {} }
    },
    loadAvatar() {
      const local = uni.getStorageSync('AVATAR') || ''
      const fromProfile = this.profile && this.profile.avatar ? this.profile.avatar : ''
      const relOrAbs = fromProfile || local
      if (!relOrAbs) { this.avatarUrl = '' ; return }
      this.avatarUrl = this.fullAvatar(relOrAbs)
    },
    fullAvatar(relOrAbs) {
      if (!relOrAbs) return ''
      if (/^https?:\/\//.test(relOrAbs)) return relOrAbs
      const apiBase = uni.getStorageSync('BASE_URL') || 'http://localhost:9000/api/v1'
      const origin = apiBase.replace(/\/api\/v1$/, '')
      return origin + relOrAbs
    },
    genderText(g) { if(!g) return '未设置'; if (g==='male' || g==='M' || g==='男') return '男'; if (g==='female' || g==='F' || g==='女') return '女'; return '保密' },
    editProfile() { if (this.navLock) return; this.navLock = true; go('mineEdit') ; setTimeout(()=>{ this.navLock=false }, 320) },
    logout() {
      uni.showModal({
        title: '提示',
        content: '确认退出登录？',
        confirmText: '退出',
        confirmColor: '#e34a4a',
        success: (res) => {
          if (res.confirm) {
            clearToken()
            uni.showToast({ title: '已退出' })
            relaunchTo('login')
          }
        }
      })
    }
  }
}
</script>

<style>
.page { padding:0 12px 20px 12px; background:#f5f5f5; min-height:100vh }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.hero { margin-top:12px; background:#7e78ff; border-radius:12px; padding:18px; display:flex; flex-direction:column; align-items:center; color:#fff }
.avatar { width:72px; height:72px; border-radius:50%; border:2px solid #fff; background:transparent }
.uname { margin-top:10px }
.stats { margin-top:14px; width:100%; display:flex; justify-content:space-around }
.stat { display:flex; flex-direction:column; align-items:center }
.num { font-size:18px; font-weight:600 }
.label { font-size:12px; opacity:0.9 }
.panel { margin-top:12px; background:#fff; border-radius:12px; padding:10px }
.ptitle { font-weight:600; margin-bottom:6px }
.item { display:flex; justify-content:space-between; align-items:center; padding:10px 6px; border-top:1px solid #f0f0f0 }
.item:first-of-type { border-top:none }
.key { color:#666 }
.val { color:#333 }
.btnrow { display:flex; gap:10px; margin-top:16px; width:100%; margin-left:0; margin-right:0 }
.edit { flex:1; background:#7e78ff; color:#fff; border-radius:24px }
.logout { flex:1; background:#e34a4a; color:#fff; border-radius:24px }
</style>
