<template>
  <view class="container">
    <view class="header"><text class="htitle">首页</text></view>
    <view class="filters">
      <input class="search" placeholder="搜索社团" v-model="keyword" @confirm="doSearch" />
      <scroll-view class="chips" scroll-x>
        <view class="chipwrap">
          <button v-for="c in chipList" :key="c.key" class="chip" :class="activeChipClass(c)" @tap="setCat(c)">{{ c.name }}</button>
        </view>
      </scroll-view>
    </view>
    <view class="grid">
      <view v-for="item in list" :key="item.id" class="card">
        <image :src="item.logo" mode="aspectFill" class="logo" />
        <view class="info">
          <text class="name">{{ item.name }}</text>
          <text class="intro">{{ item.intro }}</text>
          <view class="meta">
            <text class="count">{{ (item.activities||[]).length }}个活动</text>
            <button class="act" v-if="btnState(item.id)==='join'" @tap="join(item.id)">加入</button>
            <button class="act danger" v-else-if="btnState(item.id)==='exit'" @tap="exitClub(item.id)">退出</button>
            <button class="act disabled" v-else>审批中</button>
          </view>
        </view>
      </view>
    </view>
    <view v-if="loading" class="loading">加载中...</view>
    <view v-if="!loading && list.length===0" class="empty">暂无社团</view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() {
    return { list: [], page: 1, pageSize: 10, total: 0, keyword: '', chipKeyword: '', selectedChipKey: 'cat-0', categories: [], categoryId: 0, loading: false, hasMore: true, chipList: [ {id:0,key:'cat-0',name:'全部', mock:false} ], memberships: {} }
  },
  computed: {
    totalPages() { return Math.max(1, Math.ceil(this.total / this.pageSize)) },
    categoryLabel() { return this.categoryId ? (this.categories.find(c=>c.id===this.categoryId)?.name||'全部类别') : '全部类别' }
  },
  async onShow() { if (!this.categories.length) await this.fetchCategories(); await this.fetchMemberships(); this.resetAndFetch() },
  onReachBottom() { this.loadMore() },
  methods: {
    async fetchCategories() {
      try {
        const data = await request({ url: '/public/categories', method: 'GET', data: { page: 1, pageSize: 100 } })
        this.categories = data.list || []
        const more = (this.categories||[]).map(c=>({id:c.id,key:'cat-'+c.id,name:c.name, mock:false}))
        this.chipList = [{id:0,key:'cat-0',name:'全部', mock:false}, ...more]
        this.appendMockCategories()
      } catch(e) {}
    },
    appendMockCategories() {
      const presets = ['文学类','音乐类','志愿服务','社会实践','电竞类','摄影类','创业类']
      // 去重后随机取3个
      const existNames = new Set(this.chipList.map(i=>i.name))
      const candidates = presets.filter(n=>!existNames.has(n))
      for (let i=0; i<Math.min(3, candidates.length); i++) {
        const idx = Math.floor(Math.random() * candidates.length)
        const name = candidates.splice(idx,1)[0]
        this.chipList.push({ id: 0, key: 'mock-'+name, name, mock: true })
      }
    },
    async fetchMemberships() {
      try {
        const data = await request({ url: '/student/memberships/my', method: 'GET', data: { page: 1, pageSize: 200 } })
        const ms = data.list || []
        const m = {}
        ms.forEach(it=>{ m[it.club_id || (it.club && it.club.id) ] = it.status })
        this.memberships = m
      } catch(e) { this.memberships = {} }
    },
    async fetch(append=false) {
      this.loading = true
      const kw = this.keyword || this.chipKeyword
      const data = await request({ url: `/public/clubs`, method: 'GET', data: { page: this.page, pageSize: this.pageSize, keyword: kw, categoryId: this.categoryId || '' } })
      const items = data.list || []
      const p = data.pagination || { total: 0 }
      this.total = p.total || 0
      this.list = append ? (this.list.concat(items)) : items
      this.hasMore = this.list.length < this.total
      this.loading = false
    },
    loadMore() { if (this.loading || !this.hasMore) return; this.page++; this.fetch(true) },
    resetAndFetch() { this.page = 1; this.hasMore = true; this.list = []; this.fetch(false) },
    goDetail(id) { uni.navigateTo({ url: `/pages/clubs/detail?id=${id}` }) },
    doSearch() { this.chipKeyword = ''; this.resetAndFetch() },
    onCatChange(e) { const idx = Number(e.detail.value||0); const item = this.categories[idx]; this.categoryId = item ? item.id : 0; this.resetAndFetch() },
    setCat(c) {
      this.selectedChipKey = c.key
      if (c.mock) { this.chipKeyword = c.name; this.categoryId = 0 }
      else { this.categoryId = c.id || 0; this.chipKeyword = '' }
      this.resetAndFetch()
    },
    activeChipClass(c) {
      return c.key === this.selectedChipKey ? 'on' : ''
    },
    btnState(cid) {
      const st = this.memberships[cid]
      if (!st) return 'join'
      if (st === 'approved') return 'exit'
      return 'pending'
    },
    async join(cid) {
      try { await request({ url: `/student/clubs/${cid}/apply`, method: 'POST' }); uni.showToast({ title: '已申请' }); await this.fetchMemberships() } catch(e) {}
    },
    async exitClub(cid) {
      try { await request({ url: `/student/clubs/${cid}/exit`, method: 'POST' }); uni.showToast({ title: '已退出' }); await this.fetchMemberships() } catch(e) {}
    }
  }
}
</script>

<style>
.container { padding: 0 12px 12px 12px; background:#f5f5f5; min-height:100vh }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.filters { display:flex; flex-direction:column; gap:8px; align-items:stretch; margin:12px 0 8px 0; background:#fff; border-radius:12px; padding:12px }
.search { flex:1; border:1px solid #ddd; border-radius:20px; padding:8px 12px; background:#f6f6f6 }
.chips { margin-top:8px; white-space:nowrap }
.chipwrap { display:flex; gap:6px; flex-wrap:nowrap; padding:4px 2px }
.chip { padding:6px 12px; border:1px solid #e0e0ff; border-radius:18px; background:#fff; color:#7e78ff; font-size:12px; flex-shrink:0 }
.on { background:#e9e7ff }
.grid { display:flex; flex-wrap:wrap; margin:-6px }
.card { width: calc(50% - 12px); margin:6px; box-sizing:border-box; background:#fff; border-radius: 10px; padding: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); display:flex; flex-direction:column }
.logo { width: 100%; height: 90px; border-radius: 6px; background: #eee }
.info { margin-top: 6px; flex:1; display:flex; flex-direction:column }
.name { font-weight: 600; white-space:nowrap; overflow:hidden; text-overflow:ellipsis }
.intro { color: #666; margin-top: 4px; overflow:hidden; display:-webkit-box; -webkit-line-clamp:2; -webkit-box-orient:vertical; line-height:18px; min-height:36px; font-size:12px }
.meta { display:flex; align-items:center; justify-content:space-between; margin-top:auto }
.count { color:#666 }
.act { padding:1px 10px; border:1px solid #7e78ff; color:#7e78ff; border-radius:16px; background:#fff; font-size:12px }
.danger { border-color:#e34a4a; color:#e34a4a }
.disabled { border-color:#ddd; color:#999 }
.loading { text-align:center; color:#666; padding:10px }
.empty { text-align:center; color:#888; padding:12px; display:flex; align-items:center; justify-content:center; min-height:200px }
</style>
