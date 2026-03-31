const AUTHENTIK_URL = import.meta.env.VITE_AUTHENTIK_URL
const CLIENT_ID = import.meta.env.VITE_CLIENT_ID
const REDIRECT_URI = import.meta.env.VITE_REDIRECT_URI
const POST_LOGOUT_REDIRECT_URI = `${window.location.origin}${import.meta.env.BASE_URL}`

const AUTHORIZE_URL = `${AUTHENTIK_URL}/application/o/authorize/`
const TOKEN_URL = `${AUTHENTIK_URL}/application/o/token/`
const END_SESSION_URL = `${AUTHENTIK_URL}/application/o/cinema-app/end-session/`

const STORAGE = {
  ACCESS_TOKEN: 'cinema_access_token',
  ID_TOKEN: 'cinema_id_token',
  REFRESH_TOKEN: 'cinema_refresh_token',
  TOKEN_EXPIRY: 'cinema_token_expiry',
  PKCE_VERIFIER: 'cinema_pkce_verifier',
  OAUTH_STATE: 'cinema_oauth_state',
}

function base64urlEncode(buffer) {
  return btoa(String.fromCharCode(...new Uint8Array(buffer)))
    .replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
}

async function generatePKCE() {
  const verifier = base64urlEncode(crypto.getRandomValues(new Uint8Array(32)))
  const digest = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(verifier))
  const challenge = base64urlEncode(digest)
  return { verifier, challenge }
}

export async function redirectToLogin() {
  const { verifier, challenge } = await generatePKCE()
  const state = base64urlEncode(crypto.getRandomValues(new Uint8Array(16)))

  sessionStorage.setItem(STORAGE.PKCE_VERIFIER, verifier)
  sessionStorage.setItem(STORAGE.OAUTH_STATE, state)

  const params = new URLSearchParams({
    response_type: 'code',
    client_id: CLIENT_ID,
    redirect_uri: REDIRECT_URI,
    scope: 'openid profile email',
    state,
    code_challenge: challenge,
    code_challenge_method: 'S256',
  })
  window.location.href = `${AUTHORIZE_URL}?${params}`
}

export async function handleCallback() {
  const params = new URLSearchParams(window.location.search)
  const code = params.get('code')
  const state = params.get('state')

  if (!code) throw new Error('No code in callback')
  if (state !== sessionStorage.getItem(STORAGE.OAUTH_STATE)) throw new Error('State mismatch')

  const verifier = sessionStorage.getItem(STORAGE.PKCE_VERIFIER)
  sessionStorage.removeItem(STORAGE.PKCE_VERIFIER)
  sessionStorage.removeItem(STORAGE.OAUTH_STATE)

  const body = new URLSearchParams({
    grant_type: 'authorization_code',
    client_id: CLIENT_ID,
    redirect_uri: REDIRECT_URI,
    code,
    code_verifier: verifier,
  })

  const res = await fetch(TOKEN_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body,
  })
  if (!res.ok) throw new Error('Token exchange failed')

  const tokens = await res.json()
  storeTokens(tokens)
}

function storeTokens(tokens) {
  localStorage.setItem(STORAGE.ACCESS_TOKEN, tokens.access_token)
  localStorage.setItem(STORAGE.ID_TOKEN, tokens.id_token)
  if (tokens.refresh_token) localStorage.setItem(STORAGE.REFRESH_TOKEN, tokens.refresh_token)
  const expiry = Date.now() + tokens.expires_in * 1000
  localStorage.setItem(STORAGE.TOKEN_EXPIRY, String(expiry))
}

export function getAccessToken() {
  return localStorage.getItem(STORAGE.ACCESS_TOKEN)
}

export function getUserInfo() {
  const idToken = localStorage.getItem(STORAGE.ID_TOKEN)
  if (!idToken) return null
  try {
    const payload = idToken.split('.')[1]
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
  } catch {
    return null
  }
}

export function isAdmin() {
  const info = getUserInfo()
  return Array.isArray(info?.groups) && info.groups.includes('cinema-admins')
}

export function isLoggedIn() {
  return !!getAccessToken()
}

let refreshTimer = null

export function silentRefresh() {
  const expiry = Number(localStorage.getItem(STORAGE.TOKEN_EXPIRY))
  if (!expiry) return

  const msUntilRefresh = expiry - Date.now() - 60_000
  if (refreshTimer) clearTimeout(refreshTimer)

  refreshTimer = setTimeout(async () => {
    const refreshToken = localStorage.getItem(STORAGE.REFRESH_TOKEN)
    if (!refreshToken) { redirectToLogin(); return }

    const body = new URLSearchParams({
      grant_type: 'refresh_token',
      client_id: CLIENT_ID,
      refresh_token: refreshToken,
    })

    try {
      const res = await fetch(TOKEN_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body,
      })
      if (!res.ok) throw new Error('Refresh failed')
      const tokens = await res.json()
      storeTokens(tokens)
      silentRefresh()
    } catch {
      clearTokens()
      redirectToLogin()
    }
  }, Math.max(msUntilRefresh, 0))
}

export function clearTokens() {
  Object.values(STORAGE).forEach(k => localStorage.removeItem(k))
}

export function logout() {
  const idToken = localStorage.getItem(STORAGE.ID_TOKEN)
  clearTokens()
  const params = new URLSearchParams({
    id_token_hint: idToken || '',
    post_logout_redirect_uri: POST_LOGOUT_REDIRECT_URI,
  })
  window.location.href = `${END_SESSION_URL}?${params}`
}
