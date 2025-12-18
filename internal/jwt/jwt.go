package jwt

import (
	"fmt"
	"murky_api/internal/config"
	"murky_api/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
}

func NewService(conf *config.Config) Service {
	return Service{
		accessTokenSecret:  []byte(conf.AccessTokenSecret),
		refreshTokenSecret: []byte(conf.RefreshTokenSecret),
	}
}

type AccessTokenClaims struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	IsGuest  bool   `json:"is_guest"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (service Service) CreateAccessToken(user model.User, expiresAt time.Time) (string, error) {
	claims := AccessTokenClaims{
		user.Id,
		user.Username,
		user.IsGuest,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(service.accessTokenSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (service Service) CreateRefreshToken(username string, expiresAt time.Time) (string, error) {
	claims := RefreshTokenClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(service.refreshTokenSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (service Service) ParseAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return service.accessTokenSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (service Service) ParseRefreshToken(tokenStr string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return service.refreshTokenSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
