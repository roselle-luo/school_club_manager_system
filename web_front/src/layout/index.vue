<template>
  <div class="common-layout">
    <el-container class="layout-container">
      <el-aside width="220px" class="aside">
        <div class="logo">
          <img src="@/assets/logo.svg" alt="Logo" class="logo-img" />
          <span class="logo-text">高校社团管理系统</span>
        </div>
        <el-menu active-text-color="#409EFF" background-color="#304156" class="el-menu-vertical" text-color="#bfcbd9"
          :default-active="$route.path" router>
          <el-menu-item index="/dashboard">
            <el-icon>
              <House />
            </el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-sub-menu index="1">
            <template #title>
              <el-icon>
                <Key />
              </el-icon>
              <span>权限管理</span>
            </template>
            <el-menu-item index="/permission/users">用户管理</el-menu-item>
          </el-sub-menu>
          <el-sub-menu index="2">
            <template #title>
              <el-icon>
                <School />
              </el-icon>
              <span>社团管理</span>
            </template>
            <el-menu-item index="/clubs/list">社团列表</el-menu-item>
            <el-menu-item index="/clubs/audit">社团审核</el-menu-item>
          </el-sub-menu>
          <el-menu-item index="/admin/attendance">
            <el-icon>
              <Timer />
            </el-icon>
            <span>考勤管理</span>
          </el-menu-item>
          <el-menu-item index="/announcement">
            <el-icon>
              <Bell />
            </el-icon>
            <span>公告管理</span>
          </el-menu-item>
          <el-menu-item index="/logs">
            <el-icon>
              <Bell />
            </el-icon>
            <span>管理日志</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header class="header">
          <div class="header-left">
            <el-icon class="collapse-btn">
              <Expand />
            </el-icon>
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ $route.meta.title }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <div class="header-right">
            <el-dropdown>
              <span class="el-dropdown-link">
                {{ userStore.userInfo?.name || 'Admin' }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="goToProfile">个人中心</el-dropdown-item>
                  <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { House, Key, School, User, Setting, Expand, ArrowDown, Timer, Bell } from '@element-plus/icons-vue'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const goToProfile = () => {
  router.push('/profile')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  background-color: #304156;
  color: white;
  display: flex;
  flex-direction: column;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b3649;
  color: white;
  font-weight: bold;
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 10px;
}

.el-menu-vertical {
  border-right: none;
}

.header {
  background-color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e6e6e6;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
}
</style>
