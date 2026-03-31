package auth

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type Middleware struct {
	verifier *oidc.IDTokenVerifier
}

// NewMiddleware initialises the OIDC middleware by fetching the provider
// discovery document from issuerURL.  A TLS-insecure HTTP client is used so
// that the local self-signed nginx certificate is accepted.
func NewMiddleware(ctx context.Context, issuerURL, clientID string) (*Middleware, error) {
	insecure := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	ctx = oidc.ClientContext(ctx, insecure)

	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	return &Middleware{verifier: verifier}, nil
}

// RequireAuth rejects requests without a valid Bearer token with 401.
// On success it stores the subject (user ID) and groups in the request context.
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, err := m.verifier.Verify(r.Context(), rawToken)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var claims struct {
			Groups            []string `json:"groups"`
			PreferredUsername string   `json:"preferred_username"`
			Email             string   `json:"email"`
		}
		_ = idToken.Claims(&claims)

		username := claims.PreferredUsername
		if username == "" {
			username = claims.Email
		}
		if username == "" {
			username = idToken.Subject
		}

		ctx := WithUserID(r.Context(), idToken.Subject)
		ctx = WithUsername(ctx, username)
		ctx = WithGroups(ctx, claims.Groups)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireGroup returns 403 if the authenticated user is not in the given group.
func (m *Middleware) RequireGroup(group string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, g := range GroupsFromContext(r.Context()) {
			if g == group {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
	})
}
