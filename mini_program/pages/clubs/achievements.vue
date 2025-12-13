<template>
  <view class="page">
    <view class="header"><text class="htitle">成果展示</text></view>
    <view class="list">
      <view class="empty">暂未开放</view>
    </view>
  </view>
</template>

<script>
export default {
  data() { return { id: 0 } },
  onLoad(opts) { this.id = Number(opts.id||0) }
}
</script>

<style>
.page { min-height:100vh; background:#f7f8fa }
.header { height: 88rpx; background: #7e78ff; display:flex; align-items:center; padding:0 12px; border-bottom-left-radius:12px; border-bottom-right-radius:12px }
.htitle { color:#fff; font-weight:600 }
.list { padding:12px }
.empty { text-align:center; color:#888; padding:24px }
</style>
