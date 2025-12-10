<template>
  <view class="container">
    <view class="header"><text class="htitle">社团</text></view>
    <input class="search" placeholder="搜索社团名称或简介……" v-model="keyword" />
    <view class="tabs">
      <button :class="['tab', activeStatus===''?'on':'']" @tap="setStatus('')">全部</button>
      <button :class="['tab', activeStatus==='pending'?'on':'']" @tap="setStatus('pending')">审批中</button>
      <button :class="['tab', activeStatus==='approved'?'on':'']" @tap="setStatus('approved')">已通过</button>
      <button :class="['tab', activeStatus==='rejected'?'on':'']" @tap="setStatus('rejected')">已拒绝</button>
      <button :class="['tab', activeStatus==='quit'?'on':'']" @tap="setStatus('quit')">已退出</button>
    </view>
    <view class="row" v-for="m in filteredList" :key="m.id">
      <text class="title">{{ m.club.name }}</text>
      <text class="meta">{{ m.role }} · {{ m.status }}</text>
    </view>
    <view class="pager">
      <button size="mini" @tap="prev" :disabled="page<=1">上一页</button>
      <text class="pg">{{ page }}/{{ totalPages }}</text>
      <button size="mini" @tap="next" :disabled="page>=totalPages">下一页</button>
    </view>
    <view v-if="!loading && list.length===0" class="empty">暂无记录</view>
    <view class="login">
      <button size="mini" @tap="goLogin">登录</button>
      <button size="mini" @tap="goClubs">去申请加入社团</button>
    </view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() { return { list: [], page: 1, pageSize: 10, total: 0, activeStatus: '', loading: false, keyword: '' } },
  computed: { 
    totalPages() { return Math.max(1, Math.ceil(this.total / this.pageSize)) },
    filteredList() { const kw=(this.keyword||'').trim(); if(!kw) return this.list; return this.list.filter(it=> (it.club?.name||'').includes(kw) || (it.club?.intro||'').includes(kw)) }
  },
  onShow() { this.fetch() },
  methods: {
    async fetch() {
      this.loading = true
      try {
        const data = await request({ url: `/student/memberships/my`, method: 'GET', data: { page: this.page, pageSize: this.pageSize, status: this.activeStatus } })
        this.list = data.list || []
        const p = data.pagination || { total: 0 }
        this.total = p.total || 0
      } catch(e) { this.list = []; this.total = 0 }
      this.loading = false
    },
    prev() { if (this.page>1) { this.page--; this.fetch() } },
    next() { if (this.page<this.totalPages) { this.page++; this.fetch() } },
    setStatus(st) { this.activeStatus = st; this.page = 1; this.fetch() },
    goLogin() { uni.navigateTo({ url: '/pages/login/form' }) },
    goClubs() { uni.switchTab({ url: '/pages/clubs/list' }) }
  }
}
</script>

<style>
.container { padding:0 12px 12px 12px }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.search { width:100%; border:1px solid #e6e6e6; border-radius:20px; padding:10px 12px; font-size:14px; background:#fafafa; margin:8px 0 }
.tabs { display:flex; flex-wrap:wrap; gap:8px; margin-bottom:8px }
.tab { padding:6px 10px; border:1px solid #ddd; border-radius:8px; background:#fff; color:#333 }
.on { border-color:#2e5cff; color:#2e5cff }
.row { padding:8px; border-bottom:1px solid #f0f0f0 }
.title { font-weight:600 }
.meta { display:block; color:#666; margin-top:4px }
.pager { display:flex; justify-content:center; align-items:center; gap:12px; padding:12px }
.pg { color:#333 }
.empty { text-align:center; color:#888; padding:12px }
.login { display:flex; justify-content:center; padding:12px }
</style>
