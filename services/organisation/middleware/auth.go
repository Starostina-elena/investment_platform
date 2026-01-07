package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Starostina-elena/investment_platform/services/organisation/auth"
)

type ctxKey string

const ctxUserKey ctxKey = "user"

type UserClaims struct {
	UserID int
	Admin  bool
	Banned bool
}

func FromContext(ctx context.Context) *UserClaims {
	if v := ctx.Value(ctxUserKey); v != nil {
		if uc, ok := v.(*UserClaims); ok {
			return uc
		}
	}
	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if authz == "" {
			http.Error(w, "missing authorization", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}
		token := parts[1]
		claims, err := auth.ParseAndVerify(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		uc := &UserClaims{UserID: claims.UserID, Admin: claims.Admin, Banned: claims.Banned}
		ctx := context.WithValue(r.Context(), ctxUserKey, uc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}