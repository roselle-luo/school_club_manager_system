import request from '@/utils/request'

export function getAllAttendance(params) {
  return request({
    url: '/admin/attendance',
    method: 'get',
    params
  })
}

export function forceSignOut(id) {
  return request({
    url: `/leader/attendance/${id}/signout`,
    method: 'post'
  })
}

export function deleteAttendance(id) {
  return request({
    url: `/leader/attendance/${id}`,
    method: 'delete'
  })
}
