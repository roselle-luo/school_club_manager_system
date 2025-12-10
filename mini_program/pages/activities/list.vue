<template>
  <view class="container">
    <view class="filters">
      <input class="search" placeholder="搜索主题/地点" v-model="keyword" @confirm="doSearch" />
    </view>
    <view v-for="a in list" :key="a.id" class="row">
      <text class="title">{{ a.subject }}</text>
      <text class="meta">{{ a.time }} · {{ a.place }}</text>
    </view>
    <view class="pager">
      <button size="mini" @tap="prev" :disabled="page<=1">上一页</button>
      <text class="pg">{{ page }}/{{ totalPages }}</text>
      <button size="mini" @tap="next" :disabled="page>=totalPages">下一页</button>
    </view>
    <view v-if="!loading && list.length===0" class="empty">暂无活动</view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() { return { list: [], page: 1, pageSize: 10, total: 0, keyword: '', loading: false } },
  computed: { totalPages() { return Math.max(1, Math.ceil(this.total / this.pageSize)) } },
  onShow() { this.fetch() },
  methods: {
    async fetch() {
      this.loading = true
      const data = await request({ url: `/public/activities`, method: 'GET', data: { page: this.page, pageSize: this.pageSize, keyword: this.keyword } })
      this.list = data.list || []
      const p = data.pagination || { total: 0 }
      this.total = p.total || 0
      this.loading = false
    },
    prev() { if (this.page>1) { this.page--; this.fetch() } },
    next() { if (this.page<this.totalPages) { this.page++; this.fetch() } },
    doSearch() { this.page = 1; this.fetch() }
  }
}
</script>

<style>
.container { padding:12px }
.filters { display:flex; gap:8px; align-items:center; margin-bottom:8px }
.search { flex:1; border:1px solid #ddd; border-radius:8px; padding:8px }
.row { padding:8px; border-bottom:1px solid #f0f0f0 }
.title { font-weight:600 }
.meta { display:block; color:#666; margin-top:4px }
.pager { display:flex; justify-content:center; align-items:center; gap:12px; padding:12px }
.pg { color:#333 }
.empty { text-align:center; color:#888; padding:12px }
</style>
