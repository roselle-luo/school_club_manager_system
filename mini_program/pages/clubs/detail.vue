<template>
  <view class="wrap">
    <image :src="club.logo" class="cover" mode="aspectFill" />
    <view class="title">{{ club.name }}</view>
    <view class="category">{{ club.category }}</view>
    <view class="intro">{{ club.intro }}</view>
    <view class="sub">负责人</view>
    <view class="leaders">
      <view v-for="m in leaders" :key="m.id" class="leader">
        <text>{{ m.user.name }}</text>
        <text class="phone">{{ m.user.phone }}</text>
        <text class="role">{{ m.role }}</text>
      </view>
    </view>
    <view class="sub">近期活动</view>
    <view class="acts">
      <view v-for="a in activities" :key="a.id" class="act">
        <text>{{ a.subject }}</text>
        <text class="time">{{ a.time }}</text>
        <text class="place">{{ a.place }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import { request } from '../../utils/request.js'
export default {
  data() { return { id: 0, club: {}, leaders: [], activities: [] } },
  onLoad(opts) { this.id = Number(opts.id||0); this.fetch() },
  methods: {
    async fetch() {
      if (!this.id) return
      const data = await request({ url: `/public/clubs/${this.id}`, method: 'GET' })
      this.club = { id: data.id, name: data.name, logo: data.logo, intro: data.intro, category: data.category }
      this.leaders = data.leaders || []
      this.activities = data.activities || []
    }
  }
}
</script>

<style>
.wrap { padding:12px }
.cover { width:100%; height:160px; background:#eee; border-radius:8px }
.title { font-size:18px; font-weight:700; margin-top:8px }
.category { color:#888; margin-top:2px }
.intro { color:#555; margin-top:6px }
.sub { margin-top:10px; font-weight:600 }
.leader, .act { display:flex; gap:8px; padding:6px 0; border-bottom:1px solid #f0f0f0 }
.phone { color:#666 }
.role { color:#007AFF }
.time, .place { color:#666 }
</style>
