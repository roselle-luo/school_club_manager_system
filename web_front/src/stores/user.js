import { defineStore } from 'pinia'
import { login, getUserInfo } from '@/api/auth'
import { getManagedClubs } from '@/api/club'
import { ElMessage } from 'element-plus'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null
  }),
  getters: {
    isLoggedIn: (state) => !!state.token
  },
  actions: {
    async login(loginForm) {
      try {
        const res = await login(loginForm)
        // res should be { token: "..." }
        this.token = res.token
        localStorage.setItem('token', res.token)
        const userInfo = await this.getUserInfo()
        
        // Check Permissions: Only School Admin or Club Managers allowed
        let hasPermission = false
        
        // 1. Check if School Admin
        if (userInfo.role === 'admin') {
          hasPermission = true
        } else {
          // 2. Check if manages any clubs
          try {
            const clubsRes = await getManagedClubs(userInfo.id)
            if (clubsRes.list && clubsRes.list.length > 0) {
              hasPermission = true
            }
          } catch (err) {
            console.error('Check managed clubs failed', err)
          }
        }

        if (!hasPermission) {
          ElMessage.error('您没有权限登录管理端，仅限社团管理员或学校管理员访问')
          this.logout()
          return false
        }

        return true
      } catch (error) {
        this.logout()
        return false
      }
    },
    async getUserInfo() {
      try {
        const res = await getUserInfo()
        this.userInfo = res
        return res
      } catch (error) {
        console.error('Get user info failed', error)
      }
    },
    logout() {
      this.token = ''
      this.userInfo = null
      localStorage.removeItem('token')
    }
  }
})
