package server

import (
	"context"
	"net/http"
	"strings"

	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server/auth"
	"github.com/go-chi/render"
)

// maybeAuth is a globally-used middleware.
func (s *Server) maybeAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := s.Auth.Claims(token)
		if err != nil {
			render.Status(r, errCodeStatuses[ErrInvalidToken])
			render.JSON(w, r, APIError{
				Code:    ErrInvalidToken,
				Message: errCodeMessages[ErrInvalidToken],
			})
			return
		}

		// "valid" here refers to existence and expiry date, not whether the token is known
		valid, err := s.DB.TokenValid(r.Context(), claims.UserID, claims.TokenID)
		if err != nil {
			log.Errorf("validating token for user %v: %v", claims.UserID, err)
			render.Status(r, errCodeStatuses[ErrInternalServerError])
			render.JSON(w, r, APIError{
				Code:    ErrInternalServerError,
				Message: errCodeMessages[ErrInternalServerError],
			})
			return
		}

		if !valid {
			render.Status(r, errCodeStatuses[ErrInvalidToken])
			render.JSON(w, r, APIError{
				Code:    ErrInvalidToken,
				Message: errCodeMessages[ErrInvalidToken],
			})
			return
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
			render.Status(r, errCodeStatuses[ErrForbidden])
			render.JSON(w, r, APIError{
				Code:    ErrForbidden,
				Message: errCodeMessages[ErrForbidden],
			})
			return
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
