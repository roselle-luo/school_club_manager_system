<template>
  <view class="container">
    <view class="header"><text class="htitle">公告与活动</text></view>
    <view class="filters">
      <input class="search" placeholder="搜索公告或活动" v-model="keyword" @confirm="doSearch" />
      <view class="seg">
        <button :class="['segbtn', active==='announcement'?'on':'']" @tap="setTab('announcement')">公告</button>
        <button :class="['segbtn', active==='activity'?'on':'']" @tap="setTab('activity')">活动</button>
      </view>
    </view>
    <view v-for="item in list" :key="item.id" class="card" @tap="openItem(item)">
      <view class="card-head">
        <text class="title">{{ titleOf(item) }}</text>
        <text v-if="active==='announcement'" class="tag">公开</text>
      </view>
      <text class="meta">{{ meta1Of(item) }}</text>
      <text class="meta" v-if="meta2Of(item)">{{ meta2Of(item) }}</text>
      <text class="meta" v-if="active==='activity' && item.club">{{ item.club.name }}</text>
    </view>
    <view v-if="!loading && list.length===0" class="empty">暂无数据</view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() { return { list: [], page: 1, pageSize: 10, total: 0, keyword: '', loading: false, active: 'announcement', clubId: 0 } },
  computed: { totalPages() { return Math.max(1, Math.ceil(this.total / this.pageSize)) } },
  onLoad(opts) { if (opts && opts.clubId) this.clubId = Number(opts.clubId||0); if (opts && opts.active) this.active = opts.active; },
  onShow() { this.fetch() },
  methods: {
    async fetch() {
      this.loading = true
      const url = this.active==='announcement' ? '/public/announcements' : '/public/activities'
      const data = await request({ url, method: 'GET', data: { page: this.page, pageSize: this.pageSize, keyword: this.keyword, clubId: this.clubId || '' } })
      this.list = data.list || []
      const p = data.pagination || { total: 0 }
      this.total = p.total || 0
      this.loading = false
    },
    doSearch() { this.page = 1; this.fetch() },
    setTab(t) { this.active = t; this.page = 1; this.fetch() },
    titleOf(item) { return this.active==='announcement' ? item.title : item.subject },
    meta1Of(item) {
      if (this.active==='announcement') {
        const d = item.date || item.publishDate || (item.createdAt ? String(item.createdAt).slice(0,10) : '')
        return d
      } else {
        if (item.time) return item.time
        const st = item.startTime || ''
        const et = item.endTime || ''
        if (st && et) return `${st} - ${et}`
        return st || et || ''
      }
    },
    meta2Of(item) {
      if (this.active==='announcement') {
        return item.club ? item.club.name : ''
      } else {
        return item.place || ''
      }
    },
    openItem(item) {
      if (this.active !== 'activity') return
      const id = item.id
      const clubId = item.club_id || (item.club && item.club.id) || this.clubId || ''
      const qs = [
        ['id', id],
        ['clubId', clubId],
        ['subject', item.subject || ''],
        ['time', item.time || ''],
        ['place', item.place || ''],
        ['clubName', (item.club && item.club.name) || '']
      ].map(([k,v]) => `${k}=${encodeURIComponent(v||'')}`).join('&')
      uni.navigateTo({ url: `/pages/activities/detail?${qs}` })
    }
  }
}
</script>

<style>
.container { padding:0 12px 12px 12px; background:#f4f4f4; min-height:100vh }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.filters { display:flex; flex-direction:column; gap:10px; margin-top:10px }
.search { width:100%; border:1px solid #ddd; border-radius:20px; padding:10px 14px; background:#f6f6f6 }
.seg { display:flex; gap:10px }
.segbtn { flex:1; border:1px solid #ddd; border-radius:20px; background:#fff; color:#999; font-size:14px }
.on { background:#7e78ff; color:#fff; border-color:#7e78ff }
.card { background:#fff; border:1px solid #eaeaea; border-radius:12px; padding:10px 12px; margin:10px 0; }
.card-head { display:flex; align-items:center; justify-content:space-between; }
.title { font-weight:600 }
.meta { display:block; color:#666; margin-top:6px }
.tag { padding:2px 8px; background:#2ecc71; color:#fff; border-radius:12px; font-size:12px }
.empty { text-align:center; color:#888; padding:12px }
</style>
