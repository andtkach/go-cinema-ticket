import { getAccessToken, redirectToLogin } from './auth'

export const BASE = import.meta.env.BASE_URL.replace(/\/$/, '')

export async function api(method, path, body) {
  const headers = { 'Content-Type': 'application/json' }
  const token = getAccessToken()
  if (token) headers['Authorization'] = `Bearer ${token}`

  const opts = { method, headers }
  if (body) opts.body = JSON.stringify(body)

  const r = await fetch(BASE + path, opts)
  if (r.status === 401) { redirectToLogin(); return }
  if (r.status === 204) return null
  let data
  try { data = await r.json() } catch { data = {} }
  if (!r.ok) throw new Error(data.error || `${r.status} ${r.statusText}`)
  return data
}

export async function fetchPublic(path) {
  try {
    const r = await fetch(BASE + path)
    if (!r.ok) return null
    return r.json()
  } catch {
    return null
  }
}
