package context

import (
	"database/sql"
	"errors"
	"murky_api/internal/jwt"
	"net/http"
)

type ctxKey int

const (
	AccessTokenKey ctxKey = iota
	DbKey
)

type AuthUser struct {
	Id       int
	Username string
}

func GetAccessTokenClaims(r *http.Request) *jwt.AccessTokenClaims {
	if claims, ok := r.Context().Value(AccessTokenKey).(*jwt.AccessTokenClaims); ok {
		return claims
	}
	return nil
}

func GetCurrentUser(r *http.Request) AuthUser {
	claims := GetAccessTokenClaims(r)
	if claims == nil {
		return AuthUser{}
	}
	return AuthUser{
		Id:       claims.Id,
		Username: claims.Username,
	}
}

func GetDb(r *http.Request) (*sql.DB, error) {
	db, ok := r.Context().Value(DbKey).(*sql.DB)
	if !ok {
		return nil, errors.New("db not found in context")
	}
	return db, nil
}
