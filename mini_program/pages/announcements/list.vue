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
    <view v-for="item in list" :key="item.id" class="row">
      <text class="title">{{ titleOf(item) }}</text>
      <text class="meta">{{ metaOf(item) }}</text>
      <text v-if="active==='announcement'" class="tag">公开</text>
    </view>
    <view class="pager">
      <button size="mini" @tap="prev" :disabled="page<=1">上一页</button>
      <text class="pg">{{ page }}/{{ totalPages }}</text>
      <button size="mini" @tap="next" :disabled="page>=totalPages">下一页</button>
    </view>
    <view v-if="!loading && list.length===0" class="empty">暂无数据</view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() { return { list: [], page: 1, pageSize: 10, total: 0, keyword: '', loading: false, active: 'announcement' } },
  computed: { totalPages() { return Math.max(1, Math.ceil(this.total / this.pageSize)) } },
  onShow() { this.fetch() },
  methods: {
    async fetch() {
      this.loading = true
      const url = this.active==='announcement' ? '/public/announcements' : '/public/activities'
      const data = await request({ url, method: 'GET', data: { page: this.page, pageSize: this.pageSize, keyword: this.keyword } })
      this.list = data.list || []
      const p = data.pagination || { total: 0 }
      this.total = p.total || 0
      this.loading = false
    },
    prev() { if (this.page>1) { this.page--; this.fetch() } },
    next() { if (this.page<this.totalPages) { this.page++; this.fetch() } },
    doSearch() { this.page = 1; this.fetch() },
    setTab(t) { this.active = t; this.page = 1; this.fetch() },
    titleOf(item) { return this.active==='announcement' ? item.title : item.subject },
    metaOf(item) { return this.active==='announcement' ? (item.club ? item.club.name : '') : `${item.time} · ${item.place}` }
  }
}
</script>

<style>
.container { padding:0 12px 12px 12px }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.filters { display:flex; gap:8px; align-items:center; margin-bottom:8px }
.search { flex:1; border:1px solid #ddd; border-radius:20px; padding:8px 12px; background:#f6f6f6 }
.seg { display:flex; gap:8px }
.segbtn { padding:6px 10px; border:1px solid #e0e0ff; border-radius:16px; background:#fff; color:#7e78ff }
.on { background:#e9e7ff }
.row { padding:8px; border-bottom:1px solid #f0f0f0 }
.title { font-weight:600 }
.meta { display:block; color:#666; margin-top:4px }
.tag { position:relative; float:right; color:#2ecc71 }
.pager { display:flex; justify-content:center; align-items:center; gap:12px; padding:12px }
.pg { color:#333 }
.empty { text-align:center; color:#888; padding:12px }
</style>
