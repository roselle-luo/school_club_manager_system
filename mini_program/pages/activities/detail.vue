<template>
  <view class="page">
    <view class="header"><text class="htitle">活动详情</text></view>
  <view class="card">
      <text class="title">{{ subject || '活动' }}</text>
      <view class="row">
        <text class="key">活动时间</text>
        <text class="val">{{ time || '-' }}</text>
      </view>
      <view class="row">
        <text class="key">开始时间</text>
        <text class="val">{{ formatDate(startAt) || '-' }}</text>
      </view>
      <view class="row">
        <text class="key">结束时间</text>
        <text class="val">{{ formatDate(endAt) || '-' }}</text>
      </view>
      <view class="row">
        <text class="key">活动地点</text>
        <text class="val">{{ place || '-' }}</text>
      </view>
      <view class="row">
        <text class="key">组织者</text>
        <text class="val">{{ clubName || '-' }}</text>
      </view>
      <view class="row">
        <text class="key">参与人数</text>
        <text class="val">{{ maxParticipantsDisplay }}</text>
      </view>
      <view class="row">
        <text class="key">发布时间</text>
        <text class="val">{{ formatDate(publishAt) || '-' }}</text>
      </view>
      <view class="content-block" v-if="content">
        <text class="ctitle">活动详情</text>
        <text class="cbody">{{ content }}</text>
      </view>
      <view class="ops">
        <button class="apply" @tap="toggleRegister">{{ applyText }}</button>
      </view>
    </view>
  </view>
  </template>
  
  <script>
  import { request } from '../../utils/request.js'
  export default {
    data() {
      return { id: 0, clubId: 0, subject: '', time: '', place: '', clubName: '', startAt: '', endAt: '', maxParticipants: 0, publishAt: '', content: '', registered: false, checking: false }
    },
    computed: {
      maxParticipantsDisplay() {
        const n = Number(this.maxParticipants || 0)
        return n > 0 ? String(n) : '不限'
      },
      applyText() {
        return this.registered ? '取消报名' : '我要报名'
      }
    },
    async onLoad(opts) {
      this.id = Number(opts.id || 0)
      this.clubId = Number(opts.clubId || 0)
      this.subject = opts.subject || ''
      this.time = opts.time || ''
      this.place = opts.place || ''
      this.clubName = opts.clubName || ''
      await this.fetchById()
    },
    methods: {
      async fetchById() {
        if (!this.id) return
        try {
          const it = await request({ url: `/public/activities/${this.id}`, method: 'GET' })
          this.subject = it.subject || this.subject
          this.time = it.time || this.time
          this.place = it.place || this.place
          this.clubId = Number(it.club_id || this.clubId || 0)
          this.clubName = (it.club && it.club.name) || this.clubName
          this.content = it.content || ''
          this.startAt = it.start_at || it.startTime || ''
          this.endAt = it.end_at || it.endTime || ''
          this.maxParticipants = Number(it.max_participants != null ? it.max_participants : (it.maxParticipants != null ? it.maxParticipants : 0))
          this.publishAt = it.publish_at || it.publishTime || ''
          await this.fetchRegisterStatus()
        } catch(e) {
          // fallback: query by list
          try {
            const data = await request({ url: '/public/activities', method: 'GET', data: { clubId: this.clubId || '', page: 1, pageSize: 200 } })
            const list = data.list || []
            const it2 = list.find(x => Number(x.id) === this.id) || {}
            this.subject = this.subject || it2.subject || ''
            this.time = this.time || it2.time || ''
            this.place = this.place || it2.place || ''
            this.startAt = this.startAt || it2.start_at || it2.startTime || ''
            this.endAt = this.endAt || it2.end_at || it2.endTime || ''
            this.maxParticipants = this.maxParticipants || Number(it2.max_participants != null ? it2.max_participants : (it2.maxParticipants != null ? it2.maxParticipants : 0))
            this.publishAt = this.publishAt || it2.publish_at || it2.publishTime || ''
          } catch(e2) {}
        }
      },
      async fetchRegisterStatus() {
        if (!this.id) return
        this.checking = true
        try {
          const res = await request({ url: `/member/activities/${this.id}/register`, method: 'GET' })
          this.registered = !!(res && res.registered)
        } catch(e) { this.registered = false }
        this.checking = false
      },
      formatDate(v) {
        if (!v) return ''
        try {
          const d = typeof v === 'string' ? new Date(v) : v
          if (isNaN(d.getTime())) return String(v)
          const pad = (n)=> (n<10?('0'+n):n)
          return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
        } catch(_) { return String(v) }
      },
      async toggleRegister() {
        if (!this.id) return
        if (!this.registered) {
          try {
            await request({ url: `/member/activities/${this.id}/register`, method: 'POST' })
            this.registered = true
            uni.showToast({ title: '报名成功' })
          } catch(e) {
            uni.showToast({ title: e.msg || '报名失败', icon: 'none' })
          }
        } else {
          uni.showModal({
            title: '提示',
            content: '确认取消该活动的报名？',
            confirmText: '取消报名',
            confirmColor: '#e34a4a',
            success: async (res) => {
              if (res.confirm) {
                try {
                  await request({ url: `/member/activities/${this.id}/register`, method: 'DELETE' })
                  this.registered = false
                  uni.showToast({ title: '已取消报名' })
                } catch(e) {
                  uni.showToast({ title: e.msg || '取消失败', icon: 'none' })
                }
              }
            }
          })
        }
      }
    }
  }
  </script>
  
  <style>
  .page { padding:0 12px 20px 12px; background:#f5f5f5; min-height:100vh }
  .header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; justify-content:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
  .htitle { color:#fff; font-weight:600 }
  .card { margin-top:12px; background:#fff; border-radius:12px; padding:12px }
  .title { font-size:16px; font-weight:600; margin-bottom:8px }
  .row { display:flex; justify-content:space-between; align-items:center; padding:10px 0; border-top:1px solid #f0f0f0 }
  .row:first-of-type { border-top:none }
  .key { color:#666 }
  .val { color:#333 }
  .content-block { margin-top:10px; padding-top:10px; border-top:1px solid #f0f0f0 }
  .ctitle { font-weight:600; margin-bottom:6px; display:block }
  .cbody { color:#333; line-height:1.6 }
  .ops { margin-top:16px; display:flex; justify-content:center }
  .apply { width:80%; border-radius:24px; background:#7e78ff; color:#fff }
  </style>
