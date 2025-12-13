import { getToken } from './request.js'

const ROUTES = {
  login: { path: '/pages/login/form' },
  register: { path: '/pages/register/form' },
  clubsList: { path: '/pages/clubs/list', tab: true },
  clubsDetail: { path: '/pages/clubs/detail' },
  clubHome: { path: '/pages/clubs/home' },
  activitiesList: { path: '/pages/activities/list', tab: true },
  announcementsList: { path: '/pages/announcements/list', tab: true },
  mineHome: { path: '/pages/mine/memberships', tab: true, auth: true },
  mineEdit: { path: '/pages/mine/edit', auth: true }
}

const TAB_PAGES = new Set([
  '/pages/clubs/list',
  '/pages/activities/list',
  '/pages/announcements/list',
  '/pages/mine/memberships'
])

function buildUrl(path, params = {}) {
  const qs = Object.keys(params).length
    ? '?' + Object.keys(params).map(k => `${encodeURIComponent(k)}=${encodeURIComponent(params[k])}`).join('&')
    : ''
  return path + qs
}

function isTabPath(path) {
  return TAB_PAGES.has(path)
}

let navLock = false
function release() { navLock = false }

export function go(nameOrPath, params = {}, options = {}) {
  if (navLock) return
  navLock = true
  const route = typeof nameOrPath === 'string' && ROUTES[nameOrPath] ? ROUTES[nameOrPath] : { path: nameOrPath }
  const path = route.path
  const url = buildUrl(path, params)
  const needAuth = !!route.auth
  const token = getToken()
  const replace = !!options.replace
  const relaunch = !!options.relaunch
  if (needAuth && !token) {
    uni.reLaunch({ url: ROUTES.login.path, complete: () => { setTimeout(release, 300) } })
    return
  }
  const isTab = isTabPath(path)
  try {
    if (relaunch) {
      uni.reLaunch({ url, complete: () => { setTimeout(release, 300) } })
    } else if (isTab) {
      uni.switchTab({ url, complete: () => { setTimeout(release, 300) } })
    } else if (replace) {
      uni.redirectTo({ url, complete: () => { setTimeout(release, 300) } })
    } else {
      uni.navigateTo({ url, complete: () => { setTimeout(release, 300) } })
    }
  } catch (e) {
    uni.reLaunch({ url, complete: () => { setTimeout(release, 300) } })
  }
}

export function back(delta = 1) {
  if (navLock) return
  navLock = true
  try {
    uni.navigateBack({ delta, complete: () => { setTimeout(release, 300) } })
  } catch (e) {
    setTimeout(release, 300)
  }
}

export function relaunchTo(nameOrPath, params = {}) {
  return go(nameOrPath, params, { relaunch: true })
}

export function switchTo(nameOrPath) {
  const route = typeof nameOrPath === 'string' && ROUTES[nameOrPath] ? ROUTES[nameOrPath] : { path: nameOrPath }
  return go(route.path, {}, {})
}

export function urlOf(nameOrPath, params = {}) {
  const route = typeof nameOrPath === 'string' && ROUTES[nameOrPath] ? ROUTES[nameOrPath] : { path: nameOrPath }
  return buildUrl(route.path, params)
}

export { ROUTES }
