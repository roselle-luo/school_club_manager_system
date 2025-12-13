<template>
  <view class="page">
    <view class="nav">
      <text class="nav-title">社团详情</text>
    </view>
    <view class="strip">
      <text class="strip-title">{{ club.name || '未命名社团' }}</text>
    </view>
    <view class="banner">
      <image :src="club.logo" class="cover" mode="aspectFill" />
      <view class="overlay">
        <text class="btitle">{{ club.name || '未命名社团' }}</text>
        <view class="chips">
          <text class="chip">{{ club.category || '未知类别' }}</text>
          <text class="chip">{{ club.college || '未知学院' }}</text>
        </view>
      </view>
    </view>
    <view class="stats">
      <view class="stat">
        <text class="num">{{ counts.members }}</text>
        <text class="lab">成员</text>
      </view>
      <view class="stat">
        <text class="num">{{ counts.activities }}</text>
        <text class="lab">活动</text>
      </view>
      <view class="stat">
        <text class="num">{{ counts.history }}</text>
        <text class="lab">历史成员</text>
      </view>
    </view>
    <view class="card">
      <view class="ctitle">社团简介</view>
      <view class="ccontent">{{ club.intro || '暂无简介' }}</view>
    </view>
    <view class="card">
      <view class="ctitle">基本信息</view>
      <view class="row"><text class="key">社长</text><text class="val">{{ leaderName }}</text></view>
      <view class="row"><text class="key">联系电话</text><text class="val">{{ leaderPhoneText }}</text></view>
      <view class="row"><text class="key">创建时间</text><text class="val">{{ createdAtText }}</text></view>
    </view>
    <view class="bottom">
      <button class="join" :disabled="statusLoading || joinDisabled" @tap="apply">{{ statusLoading ? '加载中...' : joinButtonText }}</button>
    </view>
  </view>
  </template>
  
  <script>
  import { request, getToken } from '../../utils/request.js'
  export default {
    data() { return { id: 0, club: {}, leaders: [], activities: [], membershipStatus: '', statusLoading: false } },
    computed: {
      counts() {
        const members = this.club.member_count || this.club.members_count || 0
        const acts = Array.isArray(this.activities) ? this.activities.length : (this.club.activity_count || 0)
        const history = this.club.history_member_count || 0
        return { members, activities: acts, history }
      },
      joinDisabled() {
        return this.mergedStatus === 'pending' || this.mergedStatus === 'approved'
      },
      joinButtonText() {
        if (this.mergedStatus === 'pending') return '申请中'
        if (this.mergedStatus === 'approved') return '已加入'
        return '加入社团'
      },
      mergedStatus() {
        return this.membershipStatus || ''
      },
      leaderName() {
        const l1 = (this.leaders || []).find(x => (x.role || '').toLowerCase().includes('leader'))
        if (l1 && l1.user) return l1.user.name || ''
        if (this.leaders && this.leaders[0] && this.leaders[0].user) return this.leaders[0].user.name || ''
        return this.club.leader_name || '未知'
      },
      leaderPhoneText() {
        const l1 = (this.leaders || []).find(x => (x.role || '').toLowerCase().includes('leader'))
        if (l1 && l1.user) return l1.user.phone || '未知'
        if (this.leaders && this.leaders[0] && this.leaders[0].user) return this.leaders[0].user.phone || '未知'
        return this.club.leader_phone || '未知'
      },
      createdAtText() {
        const v = this.club.created_at || this.club.create_time || this.club.createdAt || ''
        if (!v) return '未知'
        try {
          const d = new Date(v)
          const y = d.getFullYear()
          const m = d.getMonth() + 1
          const day = d.getDate()
          return `${y}年${m}月${day}日`
        } catch(e) { return '未知' }
      }
    },
    onLoad(opts) { this.id = Number(opts.id||0); this.fetch(); this.fetchMembershipStatus() },
    methods: {
      goBack() { uni.navigateBack({ delta: 1 }) },
      async fetch() {
        if (!this.id) return
        const data = await request({ url: `/public/clubs/${this.id}`, method: 'GET' })
        this.club = {
          id: data.id,
          name: data.name,
          logo: data.logo,
          intro: data.intro,
          category: data.category,
          college: data.college || '',
          leader_name: data.leader_name || data.leaderName || '',
          leader_phone: data.leader_phone || data.leaderPhone || '',
          created_at: data.created_at || data.createdAt || '',
          member_count: data.member_count || data.members_count || 0,
          history_member_count: data.history_member_count || 0
        }
        this.leaders = data.leaders || []
        this.activities = data.activities || []
      },
      async fetchMembershipStatus() {
        if (!this.id) return
        const token = getToken()
        if (!token) { this.membershipStatus = ''; this.statusLoading = false; await this.$nextTick(); return }
        this.statusLoading = true
        try {
          const data = await request({ url: '/student/memberships/my', method: 'GET', data: { page: 1, pageSize: 200 } })
          const list = data.list || []
          const found = list.find(it => Number(it.club_id || (it.club && it.club.id)) === this.id)
          const st = found ? (found.status || '') : ''
          this.membershipStatus = st === 'approved' ? 'approved' : (st === 'pending' ? 'pending' : '')
        } catch(e) { this.membershipStatus = '' }
        this.statusLoading = false
        await this.$nextTick()
      },
      async apply() {
        if (!this.id) return
        if (this.joinDisabled) return
        try {
          await request({ url: `/student/clubs/${this.id}/apply`, method: 'POST' })
          uni.showToast({ title: '已申请' })
          this.membershipStatus = 'pending'
          await this.$nextTick()
        } catch(e) { uni.showToast({ title: e.msg || '申请失败', icon: 'none' }) }
      }
    }
  }
  </script>
  
  <style>
  .page { padding:56px 12px calc(env(safe-area-inset-bottom) + 88px) 12px; background:#f5f5f5; min-height:100vh }
  .nav { position:fixed; left:0; right:0; top:0; height:48px; background:#7e78ff; display:flex; align-items:center; padding:0 12px; z-index:10 }
  .back { width:32px; height:32px; line-height:32px; text-align:center; color:#fff; font-size:18px; }
  .nav-title { flex:1; text-align:center; color:#fff; font-weight:600; }
  .strip { height:32px; background:#bdbdbd; border-radius:8px; display:flex; align-items:center; justify-content:center; color:#fff; margin-top:8px }
  .strip-title { font-size:14px }
  .banner { position:relative; margin-top:0; }
  .cover { width:100%; height:180px; background:#eee; border-radius:8px }
  .overlay { position:absolute; left:12px; bottom:12px; right:12px; display:flex; flex-direction:column; align-items:flex-start; padding:8px 10px; border-radius:8px; background:linear-gradient(180deg, rgba(0,0,0,0.00), rgba(0,0,0,0.35)); }
  .btitle { color:#fff; font-size:18px; font-weight:700; text-shadow:0 2px 4px rgba(0,0,0,0.35) }
  .chips { margin-top:8px; display:flex; gap:8px }
  .chip { padding:6px 10px; border-radius:16px; background:rgba(126,120,255,0.95); color:#fff; font-size:12px }
  .stats { margin-top:12px; background:#fff; border-radius:12px; display:flex; overflow:hidden; box-shadow:0 6px 12px rgba(0,0,0,0.04) }
  .stat { flex:1; padding:14px 0; display:flex; flex-direction:column; align-items:center; }
  .stat + .stat { border-left:1px solid #f0f0f0 }
  .num { font-size:22px; font-weight:700; color:#7e78ff; line-height:26px }
  .lab { font-size:12px; color:#666; margin-top:4px }
  .card { margin-top:12px; background:#fff; border-radius:12px; padding:12px }
  .ctitle { font-weight:600; margin-bottom:8px }
  .ccontent { color:#555; line-height:20px }
  .row { display:flex; justify-content:space-between; align-items:center; padding:10px 0; border-top:1px solid #f0f0f0 }
  .row:first-of-type { border-top:none }
  .key { color:#666 }
  .val { color:#333 }
  .bottom { position:fixed; left:0; right:0; bottom:0; padding:10px 12px calc(env(safe-area-inset-bottom) + 10px); background:transparent }
  .join { width:100%; background:linear-gradient(90deg,#7e78ff 0%, #6a64ff 100%); color:#fff; border-radius:24px }
  .join[disabled] { opacity:0.6 }
  </style>
