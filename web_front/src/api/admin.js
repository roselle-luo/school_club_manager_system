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

export function getPendingClubs(params) {
  return request({
    url: '/admin/clubs/audit',
    method: 'get',
    params
  })
}

export function auditClub(id, status) {
  return request({
    url: `/admin/clubs/${id}/audit`,
    method: 'post',
    data: { status }
  })
}
