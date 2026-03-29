import { getAccessToken, redirectToLogin } from './auth'

const BASE = import.meta.env.BASE_URL.replace(/\/$/, '')

export async function api(method, path, body) {
  const headers = { 'Content-Type': 'application/json' }
  const token = getAccessToken()
  if (token) headers['Authorization'] = `Bearer ${token}`

  const opts = { method, headers }
  if (body) opts.body = JSON.stringify(body)

  const r = await fetch(BASE + path, opts)
  if (r.status === 401) { redirectToLogin(); return }
  if (r.status === 204) return null
  const data = await r.json()
  if (!r.ok) throw new Error(data.error || 'request failed')
  return data
}
