import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/public/login',
    method: 'post',
    data
  })
}

export function getUserInfo() {
  return request({
    url: '/student/me', // Updated to match likely endpoint based on router.go
    method: 'get'
  })
}
