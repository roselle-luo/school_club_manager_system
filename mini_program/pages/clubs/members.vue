<template>
  <view class="page">
    <view class="header"><text class="htitle">成员列表</text></view>
    <view class="list">
      <view v-for="u in list" :key="u.id" class="row">
        <text class="title">{{ u.name }}</text>
        <text class="meta">{{ u.student_no || '' }} · {{ u.phone || '' }}</text>
      </view>
      <view v-if="!loading && list.length===0" class="empty">暂无数据或无权限</view>
    </view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() { return { id: 0, list: [], loading: false } },
  onLoad(opts) { this.id = Number(opts.id||0); this.fetch() },
  methods: {
    async fetch() {
      if (!this.id) return
      this.loading = true
      try {
        const data = await request({ url: `/leader/clubs/${this.id}/users`, method: 'GET', data: { page: 1, pageSize: 100 } })
        this.list = data.list || []
      } catch(e) { this.list = [] }
      this.loading = false
    }
  }
}
</script>

<style>
.page { min-height:100vh; background:#f7f8fa }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.list { padding:12px }
.row { padding:8px; border-bottom:1px solid #f0f0f0 }
.title { font-weight:600 }
.meta { display:block; color:#666; margin-top:4px }
.empty { text-align:center; color:#888; padding:24px }
</style>
