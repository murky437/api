package routing

import (
	"context"
	"database/sql"
	"murky_api/internal/config"
	mycontext "murky_api/internal/context"
	"murky_api/internal/jwt"
	"murky_api/internal/model"
	"net/http"
	"slices"
	"strings"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	for i := 0; i < len(m); i++ {
		h = m[i](h)
	}
	return h
}

func RequireJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			WriteContentTypeNotJsonResponse(w)
			return
		}
		next(w, r)
	}
}

func RequireAuth(db *sql.DB, jwtService jwt.Service, guestDb *sql.DB) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") || len(auth) <= 7 {
				WriteUnauthorizedResponse(w)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims, err := jwtService.ParseAccessToken(tokenStr)
			if err != nil {
				WriteUnauthorizedResponse(w)
				return
			}

			currentDb := db
			if claims.IsGuest {
				currentDb = guestDb
			}

			var exists bool
			exists, err = model.UserExists(currentDb, claims.Username)
			if err != nil || !exists {
				WriteUnauthorizedResponse(w)
				return
			}

			ctx := context.WithValue(r.Context(), mycontext.AccessTokenKey, claims)
			ctx = context.WithValue(ctx, mycontext.DbKey, currentDb)
			next(w, r.WithContext(ctx))
		}
	}
}

func CorsMiddleware(conf *config.Config) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if slices.Contains(conf.AllowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	}
}
