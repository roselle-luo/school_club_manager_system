import request from '@/utils/request'

export function getMyProfile() {
  return request({
    url: '/student/me',
    method: 'get'
  })
}

export function updateMyProfile(data) {
  return request({
    url: '/student/me',
    method: 'put',
    data
  })
}

export function changePassword(data) {
  return request({
    url: '/student/password',
    method: 'put',
    data
  })
}
