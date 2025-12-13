<template>
  <view class="page">
    <view class="header"><text class="htitle">编辑资料</text></view>
    <view class="panel">
      <view class="avatarBox" @tap="chooseAvatar">
        <image :src="avatarUrl || defaultAvatar" class="avatarImg" mode="aspectFill" />
        <text class="avatarTip">点击头像修改</text>
      </view>
      <view class="ptitle">基本信息</view>
      <view class="item">
        <text class="key">姓名</text>
        <input class="val ipt" v-model="form.name" placeholder="请输入姓名" />
      </view>
      <view class="item">
        <text class="key">手机号</text>
        <input class="val ipt" v-model="form.phone" placeholder="请输入手机号" />
      </view>
      <view class="item">
        <text class="key">性别</text>
        <view class="val gender">
          <button class="gbtn" :class="form.gender==='male'?'on':''" @tap="form.gender='male'">男</button>
          <button class="gbtn" :class="form.gender==='female'?'on':''" @tap="form.gender='female'">女</button>
          <button class="gbtn" :class="form.gender==='secret'?'on':''" @tap="form.gender='secret'">保密</button>
        </view>
      </view>
      <view class="item">
        <text class="key">学院</text>
        <input class="val ipt" v-model="form.college" placeholder="请输入学院" />
      </view>
      <view class="item">
        <text class="key">学号</text>
        <input class="val ipt" v-model="form.student_no" placeholder="请输入学号" />
      </view>
    </view>
    <view class="btnrow">
      <button class="save" @tap="submit">保存</button>
      <button class="cancel" @tap="goBack">取消</button>
    </view>
  </view>
</template>

<script>
import { request, getToken } from '../../utils/request.js'
export default {
  data() { 
    return { 
      form: { name: '', phone: '', gender: 'secret', college: '', student_no: '' },
      avatar: '',
      avatarUrl: '',
      defaultAvatar: 'https://static-aliyun.oss-cn-hangzhou.aliyuncs.com/placeholder/avatar.png'
    } 
  },
  async onShow() {
    try {
      const p = await request({ url: '/student/me', method: 'GET' })
      this.form = { 
        name: p.name || '', 
        phone: p.phone || '', 
        gender: (p.gender==='male'||p.gender==='female'||p.gender==='secret') ? p.gender : 'secret', 
        college: p.college || '', 
        student_no: p.student_no || '' 
      }
      const localAvatar = uni.getStorageSync('AVATAR') || ''
      this.avatar = p.avatar || localAvatar || ''
      this.avatarUrl = this.fullAvatar(this.avatar)
    } catch(e) {}
  },
  methods: {
    async submit() {
      try {
        const payload = Object.assign({}, this.form, this.avatar ? { avatar: this.avatar } : {})
        await request({ url: '/student/me', method: 'PUT', data: payload })
        uni.showToast({ title: '已保存' })
        setTimeout(()=>{ uni.navigateBack({ delta: 1 }) }, 300)
      } catch(e) { uni.showToast({ title: e.msg || '保存失败', icon: 'none' }) }
    },
    goBack() { uni.navigateBack({ delta: 1 }) },
    async chooseAvatar() {
      try {
        const r = await new Promise((resolve, reject) => {
          uni.chooseImage({ count: 1, sizeType: ['compressed'], success: resolve, fail: reject })
        })
        const filePath = (r.tempFilePaths && r.tempFilePaths[0]) || (r.tempFiles && r.tempFiles[0]?.path) || ''
        if (!filePath) return
        await this.uploadAvatar(filePath)
      } catch(e) {}
    },
    fullAvatar(relOrAbs) {
      if (!relOrAbs) return ''
      if (/^https?:\/\//.test(relOrAbs)) return relOrAbs
      const apiBase = uni.getStorageSync('BASE_URL') || 'http://localhost:9000/api/v1'
      const origin = apiBase.replace(/\/api\/v1$/, '')
      return origin + relOrAbs
    },
    async uploadAvatar(filePath) {
      const apiBase = uni.getStorageSync('BASE_URL') || 'http://localhost:9000/api/v1'
      const url = apiBase.replace(/\/$/, '') + '/upload/image'
      const token = getToken()
      try {
        const res = await new Promise((resolve, reject) => {
          uni.uploadFile({
            url,
            filePath,
            name: 'file',
            header: token ? { 'Authorization': 'Bearer ' + token } : {},
            success: resolve,
            fail: reject
          })
        })
        let body = {}
        try { body = JSON.parse(res.data || '{}') } catch(e) {}
        if (body && (body.status === true) && body.data && body.data.url) {
          this.avatar = body.data.url
          this.avatarUrl = this.fullAvatar(this.avatar)
          uni.setStorageSync('AVATAR', this.avatar)
          uni.showToast({ title: '头像已更新' })
        } else {
          const msg = (body && body.msg) ? body.msg : '上传失败'
          uni.showToast({ title: msg, icon: 'none' })
        }
      } catch(e) {
        uni.showToast({ title: '上传失败', icon: 'none' })
      }
    }
  }
}
</script>

<style>
.page { padding:0 12px 20px 12px; background:#f5f5f5; min-height:100vh }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; justify-content:space-between; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.back { background:transparent; color:#fff; border:1px solid rgba(255,255,255,0.6); border-radius:16px; padding:6px 10px }
.panel { margin-top:12px; background:#fff; border-radius:12px; padding:10px }
.avatarBox { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:12px 0 }
.avatarImg { width:96px; height:96px; border-radius:50%; background:#eee }
.avatarTip { color:#999; font-size:12px; margin-top:6px }
.ptitle { font-weight:600; margin-bottom:6px }
.item { display:flex; justify-content:space-between; align-items:center; padding:10px 6px; border-top:1px solid #f0f0f0 }
.item:first-of-type { border-top:none }
.key { color:#666; min-width:72px }
.val { color:#333; flex:1; text-align:right }
.ipt { border:1px solid #eee; border-radius:8px; padding:6px 8px; text-align:left; background:#fafafa }
.gender { display:flex; justify-content:flex-end; gap:8px }
.gbtn { padding:6px 12px; border:1px solid #e0e0ff; border-radius:16px; background:#fff; color:#7e78ff }
.on { background:#e9e7ff }
.btnrow { display:flex; gap:10px; margin-top:16px }
.save, .cancel { flex:1; border-radius:24px; padding:10px 0; color:#fff }
.save { background:#7e78ff }
.cancel { background:#999 }
</style>
