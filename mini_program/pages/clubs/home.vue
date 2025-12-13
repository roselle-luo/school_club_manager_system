<template>
  <view class="page">
    <view class="header"><text class="htitle">社团主页</text></view>
    <view class="hero">
      <image :src="club.logo || defaultCover" class="cover" mode="aspectFill" />
      <view class="overlay">
        <text class="title">{{ club.name || '社团' }}</text>
        <text class="subtitle">{{ club.category || '' }}</text>
        <view class="stats">
          <view class="stat"><text class="num">{{ counts.members }}</text><text class="lab">成员</text></view>
          <view class="stat"><text class="num">{{ counts.activities }}</text><text class="lab">活动</text></view>
          <view class="stat"><text class="num">{{ counts.history }}</text><text class="lab">历史成员</text></view>
        </view>
      </view>
    </view>
    <view class="actions">
      <view class="action" @tap="goAttendance"><view class="icon">考</view><text class="aname">考勤</text></view>
      <view class="action" @tap="goActivities"><view class="icon">活</view><text class="aname">活动</text></view>
      <view class="action" @tap="goAchievements"><view class="icon">成</view><text class="aname">成果</text></view>
      <view class="action" @tap="goMembers"><view class="icon">员</view><text class="aname">成员</text></view>
      <view class="action" @tap="goIntro"><view class="icon">介</view><text class="aname">介绍</text></view>
    </view>
    <view class="section">
      <view class="sectitlebar">
        <view class="titlebtns">
          <view class="sectitle">考勤</view>
          <button
            size="mini"
            :class="attMap.club==='signin' ? 'signout' : 'signin'"
            :loading="opLoadingId===id"
            :disabled="opLoadingId===id"
            @tap="toggleSign"
          >{{ attMap.club==='signin' ? '签退' : '签到' }}</button>
        </view>
        <text class="dur" v-if="attMap.club==='signin'">已签到：{{ durationHours }}h</text>
      </view>
    </view>
    <view class="section">
      <view class="sectitle">近期活动</view>
      <view v-for="a in previewActs" :key="a.id" class="row" @tap="openActivity(a)">
        <view class="info">
          <text class="rtitle">{{ a.subject }}</text>
          <text class="rmeta">{{ a.time }} · {{ a.place }}</text>
        </view>
      </view>
      <view v-if="!loadingActs && previewActs.length===0" class="empty">暂无活动</view>
    </view>
    <view class="bottom">
      <button class="exit" @tap="exitClub" :disabled="exiting" :loading="exiting">退出社团</button>
    </view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() { return { id: 0, club: {}, exiting: false, previewActs: [], loadingActs: false, defaultCover: '/static/logo.png', attMap: {}, opLoadingId: 0, openSigninAt: '', nowTs: Date.now(), ticker: null } },
  computed: {
    counts() {
      const members = this.club.member_count || 0
      const history = this.club.history_member_count || 0
      const activities = this.previewActs.length
      return { members, history, activities }
    },
    durationHours() {
      if (!this.openSigninAt) return '0.00'
      const start = new Date(this.openSigninAt).getTime()
      const diffMs = Math.max(0, this.nowTs - start)
      const hours = diffMs / 3600000
      return Number(hours).toFixed(2)
    }
  },
  onLoad(opts) { this.id = Number(opts.id||0); this.fetchClub(); this.fetchMyAttendance(); this.fetchPreviewActivities() },
  onUnload() { this.stopTicker() },
  methods: {
    async fetchClub() {
      if (!this.id) return
      try {
        const data = await request({ url: `/public/clubs/${this.id}`, method: 'GET' })
        this.club = { id: data.id, name: data.name, category: data.category, logo: data.logo, member_count: data.member_count || 0, history_member_count: data.history_member_count || 0, intro: data.intro || '' }
      } catch(e) { this.club = { id: this.id } }
    },
    async fetchMyAttendance() {
      if (!this.id) return
      try {
        const data = await request({ url: '/member/attendance/my', method: 'GET', data: { clubId: this.id, page: 1, pageSize: 200 } })
        const list = data.list || []
        const open = list.find(it => it.signin_at && !it.signout_at)
        const hasOpen = !!open
        this.attMap = { club: hasOpen ? 'signin' : 'signout' }
        this.openSigninAt = hasOpen ? open.signin_at : ''
        if (hasOpen) this.startTicker(); else this.stopTicker()
      } catch(e) { this.attMap = {} }
    },
    async fetchPreviewActivities() {
      if (!this.id) return
      this.loadingActs = true
      try {
        const data = await request({ url: '/public/activities', method: 'GET', data: { clubId: this.id, page: 1, pageSize: 5 } })
        this.previewActs = data.list || []
      } catch(e) { this.previewActs = [] }
      this.loadingActs = false
    },
    goAnnouncements() { uni.navigateTo({ url: `/pages/announcements/list?clubId=${this.id}` }) },
    goActivities() { uni.navigateTo({ url: `/pages/announcements/list?clubId=${this.id}&active=activity` }) },
    goAttendance() { uni.navigateTo({ url: `/pages/clubs/attendance?id=${this.id}` }) },
    goAchievements() { uni.navigateTo({ url: `/pages/clubs/achievements?id=${this.id}` }) },
    goMembers() { uni.navigateTo({ url: `/pages/clubs/members?id=${this.id}` }) },
    goIntro() { uni.navigateTo({ url: `/pages/clubs/detail?id=${this.id}` }) },
    openActivity(a) {
      if (!a || !a.id) return
      const qs = [
        ['id', a.id],
        ['clubId', this.id],
        ['subject', a.subject || ''],
        ['time', a.time || ''],
        ['place', a.place || ''],
        ['clubName', this.club.name || '']
      ].map(([k,v]) => `${k}=${encodeURIComponent(v||'')}`).join('&')
      uni.navigateTo({ url: `/pages/activities/detail?${qs}` })
    },
    async toggleSign() {
      if (!this.id || this.opLoadingId===this.id) return
      const signedIn = this.attMap.club === 'signin'
      const type = signedIn ? 'signout' : 'signin'
      this.opLoadingId = this.id
      try {
        await request({ url: `/member/clubs/${this.id}/${type}`, method: 'POST' })
        uni.showToast({ title: signedIn ? '已签退' : '已签到' })
        this.attMap.club = signedIn ? 'signout' : 'signin'
        if (this.attMap.club === 'signin') { this.openSigninAt = new Date().toISOString(); this.startTicker() } else { this.openSigninAt = ''; this.stopTicker() }
      } catch(e) { uni.showToast({ title: e.msg || '操作失败', icon: 'none' }) }
      this.opLoadingId = 0
    },
    startTicker() {
      if (this.ticker) return
      this.ticker = setInterval(() => { this.nowTs = Date.now() }, 60000)
    },
    stopTicker() {
      if (this.ticker) { clearInterval(this.ticker); this.ticker = null }
    },
    async exitClub() {
      if (!this.id || this.exiting) return
      uni.showModal({
        title: '提示',
        content: '确认退出该社团？',
        confirmText: '退出',
        confirmColor: '#e34a4a',
        success: async (res) => {
          if (res.confirm) {
            this.exiting = true
            try {
              await request({ url: `/student/clubs/${this.id}/exit`, method: 'POST' })
              uni.showToast({ title: '已退出' })
              setTimeout(()=>{ uni.navigateBack({ delta: 1 }) }, 300)
            } catch(e) { uni.showToast({ title: e.msg || '退出失败', icon: 'none' }) }
            this.exiting = false
          }
        }
      })
    }
  }
}
</script>

<style>
.page { min-height:100vh; background:#f7f8fa }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.hero { position:relative; margin:12px; }
.cover { width:100%; height:160px; border-radius:12px; background:#ddd }
.overlay { position:absolute; left:12px; right:12px; bottom:12px; display:flex; flex-direction:column; gap:6px; padding:8px 10px; border-radius:10px; background:transparent; }
.title { color:#fff; font-size:18px; font-weight:700 }
.subtitle { color:#eee; font-size:12px }
.stats { display:flex; gap:12px }
.stat { display:flex; flex-direction:column; align-items:center }
.num { color:#ffffff; font-weight:700; text-shadow:0 1px 2px rgba(0,0,0,0.35) }
.lab { font-size:12px; color:#f0f0f0; text-shadow:0 1px 2px rgba(0,0,0,0.25) }
.actions { margin:8px 12px; background:#fff; border-radius:12px; display:flex; justify-content:space-around; padding:10px 8px; box-shadow:0 6px 12px rgba(0,0,0,0.04) }
.action { display:flex; flex-direction:column; align-items:center; gap:6px }
.icon { width:42px; height:42px; border-radius:50%; background:#e9e7ff; color:#7e78ff; display:flex; align-items:center; justify-content:center; font-weight:700 }
.aname { font-size:12px; color:#666 }
.section { margin:12px; background:#fff; border-radius:12px; padding:12px; box-shadow:0 6px 12px rgba(0,0,0,0.04) }
.sectitle { font-weight:600; margin-bottom:8px }
.sectitlebar { display:flex; align-items:center; justify-content:space-between }
.titlebtns { display:flex; align-items:center; gap:10px }
.sectitlebar .sectitle { margin-bottom:0 }
.row { padding:8px 0; border-top:1px solid #f0f0f0; display:flex; align-items:center; justify-content:space-between; gap:8px }
.row:first-of-type { border-top:none }
.info { display:flex; flex-direction:column }
.rtitle { font-weight:600 }
.rmeta { display:block; color:#666; margin-top:4px }
.signin { background:#7e78ff; color:#fff }
.signout { background:#e34a4a; color:#fff }
.dur { color:#666 }
.bottom { padding:0 12px 20px 12px }
.exit { width:100%; background:linear-gradient(90deg,#e34a4a 0%, #d84242 100%); color:#fff; border-radius:24px }
</style>
