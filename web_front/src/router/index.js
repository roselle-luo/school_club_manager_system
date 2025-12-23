import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layout/index.vue'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: { title: 'Login' }
    },
    {
      path: '/',
      component: Layout,
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/index.vue'),
          meta: { title: '首页', icon: 'House' }
        },
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('@/views/profile/index.vue'),
          meta: { title: '个人中心', icon: 'User' }
        },
        {
          path: 'admin/attendance',
          name: 'AdminAttendance',
          component: () => import('@/views/admin/attendance/index.vue'),
          meta: { title: '考勤管理', icon: 'Timer' }
        },
        {
          path: 'permission/users',
          name: 'UserPermission',
          component: () => import('@/views/permission/users/index.vue'),
          meta: { title: '用户管理' }
        },
        {
          path: 'clubs/list',
          name: 'ClubList',
          component: () => import('@/views/clubs/list.vue'),
          meta: { title: '社团列表' }
        },
        {
          path: 'clubs/audit',
          name: 'ClubAudit',
          component: () => import('@/views/dashboard/index.vue'), // Placeholder
          meta: { title: '社团审核' }
        },
        {
          path: 'activities/:id',
          name: 'ActivityDetail',
          component: () => import('@/views/activities/detail.vue'),
          meta: { title: '活动详情', hidden: true }
        },
        {
          path: 'announcement',
          name: 'AnnouncementManage',
          component: () => import('@/views/announcement/index.vue'),
          meta: { title: '公告管理' }
        }
      ]
    },
    // Add more routes here later based on sidebar items in screenshots
    // e.g., Club Management, System Management, etc.
  ]
})

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const token = userStore.token

  if (token) {
    if (to.path === '/login') {
      next('/')
    } else {
      if (!userStore.userInfo) {
        try {
          await userStore.getUserInfo()
          next()
        } catch (error) {
          userStore.logout()
          next('/login')
        }
      } else {
        next()
      }
    }
  } else {
    if (to.path === '/login') {
      next()
    } else {
      next('/login')
    }
  }
})

export default router
