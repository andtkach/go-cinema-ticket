export async function api(method, path, body) {
  const opts = { method, headers: { 'Content-Type': 'application/json' } }
  if (body) opts.body = JSON.stringify(body)
  const r = await fetch(path, opts)
  if (r.status === 204) return null
  const data = await r.json()
  if (!r.ok) throw new Error(data.error || 'request failed')
  return data
}
