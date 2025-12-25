import request from '@/utils/request'

export function getAnnouncements(params) {
  return request({
    url: '/public/announcements',
    method: 'get',
    params
  })
}

export function getActivities(params) {
  return request({
    url: '/public/activities',
    method: 'get',
    params
  })
}

export function getActivityDetail(id) {
  return request({
    url: `/public/activities/${id}`,
    method: 'get'
  })
}

export function getCategories() {
  return request({
    url: '/public/categories',
    method: 'get'
  })
}

export function registerClub(data) {
  return request({
    url: '/public/clubs/register',
    method: 'post',
    data
  })
}

