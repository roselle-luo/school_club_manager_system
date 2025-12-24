import request from '@/utils/request'

// Public APIs
export function getClubs(params) {
  return request({
    url: '/public/clubs',
    method: 'get',
    params
  })
}

export function getClubDetail(id) {
  return request({
    url: `/public/clubs/${id}`,
    method: 'get'
  })
}

export function getCategories() {
  return request({
    url: '/public/categories',
    method: 'get'
  })
}

// Student APIs
export function applyClub(id) {
  return request({
    url: `/student/clubs/${id}/apply`,
    method: 'post'
  })
}

export function getMyMemberships(params) {
  return request({
    url: `/student/memberships/my`,
    method: 'get',
    params
  })
}

export function quitClub(id) {
  return request({
    url: `/student/clubs/${id}/exit`,
    method: 'post'
  })
}

// Leader APIs
export function getManagedClubs(userId, params) {
  return request({
    url: `/leader/users/${userId}/clubs`,
    method: 'get',
    params
  })
}

export function getClubLeaders(clubId) {
  return request({
    url: `/leader/clubs/${clubId}/users`,
    method: 'get'
  })
}

export function getClubMembers(clubId, params) {
  return request({
    url: `/leader/clubs/${clubId}/members/users`,
    method: 'get',
    params
  })
}

export function getPendingMemberships(clubId, params) {
  return request({
    url: `/leader/clubs/${clubId}/memberships`,
    method: 'get',
    params
  })
}

export function approveMembership(clubId, membershipId) {
  return request({
    url: `/leader/clubs/${clubId}/memberships/${membershipId}/approve`,
    method: 'post'
  })
}

export function rejectMembership(clubId, membershipId) {
  return request({
    url: `/leader/clubs/${clubId}/memberships/${membershipId}/reject`,
    method: 'post'
  })
}

export function getClubAttendance(clubId, params) {
  return request({
    url: `/leader/clubs/${clubId}/attendance`,
    method: 'get',
    params
  })
}

export function getClubLogs(clubId, params) {
  return request({
    url: `/leader/clubs/${clubId}/logs`,
    method: 'get',
    params
  })
}

export function getOperationLogs(clubId, params) {
  return request({
    url: `/leader/clubs/${clubId}/logs`,
    method: 'get',
    params
  })
}

export function updateMemberRole(clubId, userId, data) {
  return request({
    url: `/leader/clubs/${clubId}/members/${userId}/role`,
    method: 'post',
    data
  })
}

export function kickMember(clubId, userId) {
  return request({
    url: `/leader/clubs/${clubId}/members/${userId}`,
    method: 'delete'
  })
}

export function getClubAnnouncements(clubId, params) {
  return request({
    url: `/leader/clubs/${clubId}/announcements`,
    method: 'get',
    params
  })
}

export function createAnnouncement(clubId, data) {
  return request({
    url: `/leader/clubs/${clubId}/announcements`,
    method: 'post',
    data
  })
}

export function updateAnnouncement(clubId, id, data) {
  return request({
    url: `/leader/clubs/${clubId}/announcements/${id}`,
    method: 'put',
    data
  })
}

export function deleteAnnouncement(clubId, id) {
  return request({
    url: `/leader/clubs/${clubId}/announcements/${id}`,
    method: 'delete'
  })
}

// Member APIs
export function getMyAttendance(params) {
  return request({
    url: '/member/attendance/my',
    method: 'get',
    params
  })
}

export function registerActivity(activityId) {
  return request({
    url: `/member/activities/${activityId}/register`,
    method: 'post'
  })
}

export function getRegisterStatus(activityId) {
  return request({
    url: `/member/activities/${activityId}/register`,
    method: 'get'
  })
}

export function cancelRegisterActivity(activityId) {
  return request({
    url: `/member/activities/${activityId}/register`,
    method: 'delete'
  })
}

// Admin APIs
export function dissolveClub(clubId) {
  return request({
    url: `/admin/clubs/${clubId}`,
    method: 'delete'
  })
}
