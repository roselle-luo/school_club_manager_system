<template>
  <view class="container">
    <view class="seg">
      <text :class="['seg-item', status==='approved' ? 'active' : '']" @tap="switchStatus('approved')">已加入</text>
      <text :class="['seg-item', status==='pending' ? 'active' : '']" @tap="switchStatus('pending')">申请中</text>
    </view>
    <view v-for="it in items" :key="it.clubId" class="row" @tap="onCardTap(it.clubId)">
      <view class="info">
        <text class="title">{{ it.clubName }}</text>
        <text class="meta">社长：{{ it.leaderName || '未知' }}</text>
      </view>
      <text v-if="status==='approved'" class="enter" @tap="goClub(it.clubId)">点击进入</text>
    </view>
    <view v-if="!loading && items.length===0" class="empty">暂无社团</view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
import { go } from '../../utils/router.js'
export default {
  data() { 
    return { 
      status: 'approved',
      items: [],
      leadersMap: {},
      loading: false
    } 
  },
  onShow() { this.fetch() },
  methods: {
    switchStatus(s) {
      if (this.status === s) return
      this.status = s
      this.fetch()
    },
    async fetch() {
      this.loading = true
      try {
        const data = await request({ url: '/student/memberships/my', method: 'GET', data: { page: 1, pageSize: 200, status: this.status } })
        const list = data.list || []
        const items = list.map(it => {
          const club = it.club || {}
          return {
            clubId: it.club_id || club.id,
            clubName: club.name || '',
            leaderName: this.leadersMap[it.club_id || club.id] || ''
          }
        }).filter(it => !!it.clubId)
        this.items = items
        const ids = Array.from(new Set(items.map(i => i.clubId))).filter(Boolean)
        await this.fetchLeaders(ids)
      } catch(e) {
        this.items = []
      }
      this.loading = false
    },
    onCardTap(id) {
      if (this.status !== 'approved') return
      this.goClub(id)
    },
    goClub(id) { if (!id) return; go('clubHome', { id }) },
    async fetchLeaders(ids) {
      const pending = ids.filter(id => !(id in this.leadersMap))
      if (!pending.length) return
      try {
        const results = await Promise.all(pending.map(id => request({ url: `/public/clubs/${id}`, method: 'GET' })))
        results.forEach((data, idx) => {
          const id = pending[idx]
          const leaders = data.leaders || []
          let leaderName = ''
          for (let i = 0; i < leaders.length; i++) {
            const m = leaders[i] || {}
            if ((m.role || '').toLowerCase() === 'leader') {
              leaderName = (m.user && m.user.name) || ''
              break
            }
          }
          if (!leaderName && leaders.length > 0) {
            leaderName = (leaders[0].user && leaders[0].user.name) || ''
          }
          this.leadersMap[id] = leaderName
        })
        this.items = this.items.map(it => ({ ...it, leaderName: this.leadersMap[it.clubId] || it.leaderName }))
      } catch(e) {}
    }
  }
}
</script>

<style>
.container { padding:12px; background:#f7f8fa; min-height:100vh }
.seg { display:flex; gap:12px; align-items:center; margin-bottom:12px; background:#eef0ff; border:1px solid #e6e7fb; border-radius:8px; padding:6px }
.seg-item { flex:1; text-align:center; padding:8px 0; border-radius:6px; color:#666 }
.seg-item.active { background:#7e78ff; color:#fff; font-weight:600 }
.row { padding:12px; background:#fff; border:1px solid #e4e6ff; border-radius:12px; box-shadow:0 6px 12px rgba(0,0,0,0.04); display:flex; align-items:center; justify-content:space-between; margin-bottom:12px }
.row:last-child { margin-bottom:0 }
.info { display:flex; flex-direction:column }
.title { font-weight:600 }
.meta { display:block; color:#666; margin-top:4px }
.enter { color:#7e78ff; font-weight:700 }
.empty { text-align:center; color:#888; padding:12px }
</style>
