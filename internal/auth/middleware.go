package auth

import (
	"context"
	"net/http"

	"github.com/Shelex/split-specs/internal/users"
	"github.com/Shelex/split-specs/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			user, err := jwt.ParseToken(token)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// create user and check if user exists in db
			if !user.Exist() {
				next.ServeHTTP(w, r)
				return
			}
			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
