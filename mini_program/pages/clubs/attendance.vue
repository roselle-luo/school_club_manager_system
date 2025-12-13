<template>
  <view class="page">
    <view class="header"><text class="htitle">考勤</text></view>
    <view class="list">
      <view v-for="it in list" :key="it.id" class="row">
        <view class="col">
          <text class="title">社团打卡</text>
          <text class="meta" v-if="it.signin_at">{{ signMeta(it) }}</text>
          <text class="meta" v-else>{{ dateText(it.created_at) }}</text>
        </view>
        <button v-if="it.signin_at && !it.signout_at" size="mini" @tap="signout(it)" :loading="opLoadingId===it.id" :disabled="opLoadingId===it.id">签退</button>
      </view>
      <view v-if="!loading && list.length===0" class="empty">暂无记录</view>
    </view>
    <view class="pager">
      <button size="mini" @tap="prev" :disabled="page<=1">上一页</button>
      <text class="pg">{{ page }}/{{ totalPages }}</text>
      <button size="mini" @tap="next" :disabled="page>=totalPages">下一页</button>
    </view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() { return { id: 0, list: [], loading: false, page: 1, pageSize: 10, total: 0, opLoadingId: 0 } },
  computed: {
    totalPages() { return Math.max(1, Math.ceil(this.total / this.pageSize)) },
    totalHours() {
      const hours = (this.list || []).reduce((acc, it) => acc + ((it.signout_at && (it.duration_hours || 0) > 0) ? (it.duration_hours || 0) : 0), 0)
      if (hours > 0) return Number(hours).toFixed(2)
      const minsFallback = (this.list || []).reduce((acc, it) => acc + ((it.signout_at && it.duration_minutes > 0) ? it.duration_minutes : 0), 0)
      return (minsFallback / 60).toFixed(2)
    }
  },
  onLoad(opts) { this.id = Number(opts.id||0); this.fetch() },
  methods: {
    async fetch() {
      if (!this.id) return
      this.loading = true
      try {
        const data = await request({ url: '/member/attendance/my', method: 'GET', data: { clubId: this.id, page: this.page, pageSize: this.pageSize } })
        const p = data.pagination || { total: 0 }
        this.total = p.total || 0
        this.list = data.list || []
      } catch(e) { this.list = [] }
      this.loading = false
    },
    async signout(it) {
      if (!it || this.opLoadingId===it.id) return
      this.opLoadingId = it.id
      try {
        const res = await request({ url: `/member/clubs/${this.id}/signout`, method: 'POST' })
        if (!res) {
          this.list = this.list.filter(row => row.id !== it.id)
          uni.showToast({ title: '不足1分钟，记录已取消', icon: 'none' })
        } else {
          const idx = this.list.findIndex(row => row.id === it.id)
          if (idx >= 0) this.list.splice(idx, 1, { ...it, ...res })
          uni.showToast({ title: '已签退' })
        }
      } catch(e) {
        uni.showToast({ title: e?.msg || '签退失败', icon: 'none' })
      }
      this.opLoadingId = 0
    },
    prev() { if (this.page>1) { this.page--; this.fetch() } },
    next() { if (this.page<this.totalPages) { this.page++; this.fetch() } },
    typeText(t) { return t==='signin' ? '签到' : (t==='signout' ? '签退' : (t||'')) },
    dateText(v) { if (!v) return ''; try { const d = new Date(v); return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')} ${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}` } catch(e) { return '' } },
    signMeta(it) {
      const inText = this.dateText(it.signin_at)
      if (!it.signout_at) return `签到：${inText} · 正在签到中`
      const outText = this.dateText(it.signout_at)
      const hours = (it.duration_hours != null) ? Number(it.duration_hours).toFixed(2) : ((it.duration_minutes || 0) / 60).toFixed(2)
      const durText = `${hours}h`
      return `签到：${inText} · 签退：${outText} · 本次时长：${durText}`
    }
  }
}
</script>

<style>
.page { min-height:100vh; background:#f7f8fa }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.sum { color:#fff; margin-left:12px; font-size: 24rpx }
.list { padding:12px }
.row { padding:8px; border-bottom:1px solid #f0f0f0; display:flex; align-items:center; justify-content:space-between; gap:8px }
.col { flex:1 }
.title { font-weight:600 }
.meta { display:block; color:#666; margin-top:4px }
.empty { text-align:center; color:#888; padding:12px }
.pager { display:flex; align-items:center; justify-content:center; gap:12px; padding:12px }
.pg { color:#666 }
</style>
