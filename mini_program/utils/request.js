const BASE_URL = uni.getStorageSync('BASE_URL') || 'http://localhost:9000/api/v1'

export function setToken(t) { uni.setStorageSync('TOKEN', t) }
export function getToken() { return uni.getStorageSync('TOKEN') }
export function clearToken() { uni.removeStorageSync('TOKEN') }

function handleUnauthorized() {
  clearToken()
  const pages = getCurrentPages()
  const last = pages[pages.length - 1]
  const path = last && last.route ? '/' + last.route : ''
  if (path !== '/pages/login/form') {
    uni.reLaunch({ url: '/pages/login/form' })
  }
}

export function request({ url, method = 'GET', data = {}, header = {} }) {
  const token = getToken()
  const h = Object.assign({}, header)
  if (token) h['Authorization'] = 'Bearer ' + token
  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE_URL + url,
      method,
      data,
      header: h,
      success: (res) => {
        const body = res.data || {}
        if (res.statusCode === 401 || body.code === 401) {
          handleUnauthorized()
          reject(body)
          return
        }
        if (body.status) {
          resolve(body.data)
        } else {
          reject(body)
        }
      },
      fail: (err) => reject(err)
    })
  })
}

export function setBaseUrl(v) { uni.setStorageSync('BASE_URL', v) }
