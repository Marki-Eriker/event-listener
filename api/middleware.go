package api

import (
	"context"
	"fmt"
	"github.com/marki-eriker/event-listener/entity/user"
	"github.com/marki-eriker/event-listener/pkg/request"
	"github.com/marki-eriker/event-listener/pkg/token"
	"net/http"
	"strings"
)

func authMiddleware(tokenManager token.Generator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if headerParts[0] != "Bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			t, err := tokenManager.Parse(headerParts[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, ok := t.Claims.(*token.Claims)
			if !ok || !t.Valid {
				fmt.Println(ok)
				fmt.Println(claims)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, request.UserRole, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func withRoleMiddleware(roles map[user.Role]struct{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, err := request.GetUserRole(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			_, ok := roles[role]
			if !ok {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
