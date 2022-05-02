package server

import (
	"context"
	"net/http"

	"gitlab.com/1f320/pronouns/backend/server/auth"
)

// maybeAuth is a globally-used middleware.
func (s *Server) maybeAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := s.Auth.Claims(token)
		if err != nil {
			// if we get here, a token was supplied but it's invalid--return an error
		}

		ctx := context.WithValue(r.Context(), ctxKeyClaims, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// MustAuth makes a valid token required
func MustAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, ok := ClaimsFromContext(r.Context())
		if !ok {

		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// ClaimsFromContext returns the auth.Claims in the context, if any.
func ClaimsFromContext(ctx context.Context) (auth.Claims, bool) {
	v := ctx.Value(ctxKeyClaims)
	if v == nil {
		return auth.Claims{}, false
	}

	claims, ok := v.(auth.Claims)
	return claims, ok
}
