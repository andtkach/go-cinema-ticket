# Authentication Integration Plan

IDP: Authentik (self-hosted, already running via docker-compose)
Auth flow: OAuth2 Authorization Code + PKCE
Token validation: JWT (local, stateless — no introspection per request)

---

## Decisions

- Booking routes require authentication; seat listing is public
- Nginx proxies Authentik — client never calls Authentik directly
- Silent token refresh via refresh token rotation
- Two user groups in Authentik: **Admins** and **Clients**
  - Clients: can browse movies, view seats, book tickets
  - Admins: future management capabilities (manage movies, seats, users)

---

## Phase 1 — Authentik Configuration (manual, in Authentik admin UI)

1. Create two groups:
   - `cinema-admins`
   - `cinema-clients`

2. Create an OAuth2/OIDC Provider:
   - Name: `cinema-provider`
   - Client type: Public (no client secret — PKCE only)
   - Signing algorithm: RS256
   - Scopes: `openid`, `profile`, `email`
   - Add a custom scope or use the `groups` claim to expose group membership in the token
   - Redirect URIs:
     - `https://localhost:17091/app/callback` (production via nginx)
     - `http://localhost:17070/callback` (local dev via Vite)

3. Create an Application:
   - Name: `Cinema App`
   - Slug: `cinema-app`
   - Provider: `cinema-provider`
   - Bind a policy that allows only users in `cinema-admins` or `cinema-clients`

4. Note down:
   - Client ID
   - Authentik base URL (e.g. `https://localhost:17091/idp`)
   - Discovery URL: `<base>/application/o/cinema-app/.well-known/openid-configuration`

---

## Phase 2 — Nginx: Proxy Authentik

File: `config/nginx/nginx.conf`

1. Add upstream block pointing to `server-idp:9000`
2. Add location block `/idp/` that proxies to Authentik:
   - Strip or rewrite the `/idp` prefix appropriately
   - Forward `Host`, `X-Real-IP`, `X-Forwarded-For`, `X-Forwarded-Proto` headers
3. Client will reach Authentik at `https://localhost:17091/idp/application/o/...`
   — same origin as the app, no CORS issues

---

## Phase 3 — Go Server: Auth Middleware

### 3.1 Dependencies
```
go get github.com/coreos/go-oidc/v3/oidc
go get golang.org/x/oauth2
```

### 3.2 New package: `server/internal/auth/`

**`middleware.go`**
- `NewMiddleware(ctx, issuerURL, clientID)` — initializes `oidc.Provider` from
  Authentik discovery URL on startup; creates `oidc.IDTokenVerifier`
- `RequireAuth(next http.Handler) http.Handler` — middleware func:
  - Extracts `Authorization: Bearer <token>` from request header
  - Verifies JWT locally (signature via cached JWKS, `exp`, `aud`, `iss`)
  - Extracts claims: `sub` (user ID), `groups` (group membership)
  - Injects `UserID` and `Groups` into request context via typed key
  - Returns `401` if token missing or invalid
- `RequireGroup(group string, next http.Handler) http.Handler` — optional middleware
  for future admin-only routes:
  - Reads groups from context
  - Returns `403` if user is not in the required group

**`context.go`**
- Typed context keys to avoid collisions
- Helper getters: `UserIDFromContext(ctx)`, `GroupsFromContext(ctx)`

### 3.3 Changes to `server/cmd/main.go`
- Read env vars: `AUTHENTIK_ISSUER_URL`, `AUTHENTIK_CLIENT_ID`
- Initialize auth middleware on startup
- Apply `RequireAuth` to booking routes only:
  ```
  POST /movies/{movieID}/seats/{seatID}/hold   ← protected
  PUT  /sessions/{sessionID}/confirm           ← protected
  DELETE /sessions/{sessionID}                 ← protected
  GET  /movies                                 ← public
  GET  /movies/{movieID}/seats                 ← public
  GET  /                                       ← public (static files)
  ```
- Remove any logic that accepts client-supplied `user_id`

### 3.4 Changes to `server/internal/booking/handler.go`
- Replace client-supplied `user_id` from request body with `UserIDFromContext(ctx)`
- Booking is now tied to the authenticated identity from the token

### 3.5 Environment variables (add to `.env` / server config)
```
AUTHENTIK_ISSUER_URL=https://localhost:17091/idp/application/o/cinema-app
AUTHENTIK_CLIENT_ID=<client-id-from-authentik>
```

---

## Phase 4 — React Client: PKCE Auth Flow

### 4.1 New file: `client/src/auth.js`

Constants (from `VITE_*` env vars):
- `VITE_AUTHENTIK_URL` — e.g. `https://localhost:17091/idp`
- `VITE_CLIENT_ID`
- `VITE_REDIRECT_URI` — e.g. `https://localhost:17091/app/callback`

Functions:
- `redirectToLogin()` — generate PKCE verifier + challenge (via `crypto.subtle`),
  store verifier in `sessionStorage`, redirect to Authentik authorize endpoint
- `handleCallback()` — read `code` from URL params, exchange for tokens using
  stored verifier, save tokens, redirect to `/`
- `getAccessToken()` — return stored access token
- `getUserInfo()` — decode ID token payload (base64 decode, no library needed),
  return `{ sub, email, name, groups }`
- `isAdmin()` — check if groups claim includes `cinema-admins`
- `silentRefresh()` — use stored refresh token to call token endpoint, update
  stored access token; schedule next refresh before expiry
- `logout()` — clear tokens from storage, redirect to Authentik logout endpoint

Token storage: `localStorage`
- `cinema_access_token`
- `cinema_id_token`
- `cinema_refresh_token`
- `cinema_token_expiry` (unix timestamp for refresh scheduling)

### 4.2 Changes to `client/src/App.jsx`
- On mount:
  1. If current path is `/callback` → call `handleCallback()` then redirect
  2. If no access token in storage → call `redirectToLogin()`
  3. Otherwise → call `silentRefresh()` to schedule background refresh
- Pass `getUserInfo()` result as prop instead of the random `userID`
- Show user name/email in the header

### 4.3 Changes to `client/src/api.js`
- Add `Authorization: Bearer <token>` header to all requests using `getAccessToken()`
- On `401` response → call `redirectToLogin()`

### 4.4 Changes to `client/src/components/Header.jsx`
- Display authenticated user name/email instead of random ID
- Add logout button calling `logout()`

### 4.5 New file: `client/.env`
```
VITE_AUTHENTIK_URL=https://localhost:17091/idp
VITE_CLIENT_ID=<client-id-from-authentik>
VITE_REDIRECT_URI=https://localhost:17091/app/callback
```

### 4.6 New file: `client/.env.development`
```
VITE_AUTHENTIK_URL=https://localhost:17091/idp
VITE_CLIENT_ID=<client-id-from-authentik>
VITE_REDIRECT_URI=http://localhost:17070/callback
```

---

## Phase 5 — Silent Token Refresh

Logic in `auth.js` `silentRefresh()`:
1. Read `cinema_token_expiry` from storage
2. Calculate time until expiry minus 60-second buffer
3. `setTimeout(() => { ... }, timeUntilRefresh)`:
   - POST to Authentik token endpoint with `grant_type=refresh_token`
   - Update `cinema_access_token` and `cinema_token_expiry` in storage
   - Reschedule next refresh recursively
4. If refresh fails (refresh token expired) → call `redirectToLogin()`
5. Call `silentRefresh()` once on app mount after verifying token exists

---

## Phase 6 — Groups and Future Admin Routes

### Backend
- `RequireGroup("cinema-admins", handler)` middleware is ready from Phase 3
  but not applied to any routes yet (no admin endpoints exist)
- Future admin routes will be added under `/admin/...` and wrapped with this middleware

### Frontend
- `isAdmin()` utility is ready in `auth.js`
- Future admin UI sections will gate rendering behind `isAdmin()` check

---

## File Change Summary

| File | Change |
|------|--------|
| `config/nginx/nginx.conf` | Add `/idp/` proxy location to `server-idp:9000` |
| `server/go.mod` | Add `go-oidc`, `oauth2` dependencies |
| `server/internal/auth/middleware.go` | New — OIDC middleware |
| `server/internal/auth/context.go` | New — context key helpers |
| `server/cmd/main.go` | Init auth, apply middleware to booking routes, remove user_id from body |
| `server/internal/booking/handler.go` | Read user_id from context instead of request body |
| `client/src/auth.js` | New — PKCE flow, token storage, silent refresh |
| `client/src/App.jsx` | Auth check on mount, callback handling, pass real user info |
| `client/src/api.js` | Add Bearer token header, handle 401 |
| `client/src/components/Header.jsx` | Show real user info, add logout button |
| `client/.env` | New — Vite env vars for production |
| `client/.env.development` | New — Vite env vars for dev |
| `.env` | Add `AUTHENTIK_ISSUER_URL`, `AUTHENTIK_CLIENT_ID` |

---

## Open Questions (resolve before implementation)

- [ ] What port/domain will Authentik be accessible on in production?
      (Currently `server-idp` listens on 9000/9443 internally)
## Resolved Decisions

| Question | Answer |
|----------|--------|
| Authentik port/domain | Nginx proxies at `https://localhost:17091/idp/` |
| Groups claim location | Access token — configure a custom property mapping in Authentik that adds `groups` to the access token scope |
| Booking store | Keep Redis; `postgres-app` will be used for future features |
