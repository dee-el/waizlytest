package commonhttpmiddleware

import (
	"context"
	"net/http"
	"strings"

	"waizlytest/common/contextkey"
	commonerr "waizlytest/common/errors"
	commonhttpenc "waizlytest/common/http/encoder"
	authservice "waizlytest/services/auth"
)

type AuthMiddleware struct {
	authz authservice.Authz
}

func NewAuthMiddleware(authz authservice.Authz) *AuthMiddleware {
	return &AuthMiddleware{
		authz: authz,
	}
}

func (md *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			e := commonerr.ErrorForbidden
			commonhttpenc.JSONErrorEncoder(ctx, w, e)
			return
		}

		values := strings.Fields(authHeader)
		if len(values) != 2 || values[0] != "Bearer" {
			e := commonerr.ErrorForbidden
			commonhttpenc.JSONErrorEncoder(ctx, w, e)
			return
		}

		token := values[1]
		id, err := md.authz.ValidateToken(token)
		if err != nil {
			commonhttpenc.JSONErrorEncoder(ctx, w, err)
			return
		}

		ctx = context.WithValue(ctx, contextkey.UserIDContextKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
