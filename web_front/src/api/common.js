import request from '@/utils/request'

export function uploadImage(data) {
  return request({
    url: '/upload/image',
    method: 'post',
    headers: { 'Content-Type': 'multipart/form-data' },
    data
  })
}
