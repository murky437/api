package auth

import (
	"murky_api/internal/app"
	"murky_api/internal/constants"
	"murky_api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDeleteRefreshTokenSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	expiresAt := time.Now().UTC().Add(time.Hour)
	refreshToken, err := c.JwtService.CreateRefreshToken("user", expiresAt)
	require.NoError(t, err)
	_, err = c.Db.Exec(
		`INSERT INTO refresh_token (user_id, jwt, expires_at) VALUES (?, ?, ?)`,
		1,
		refreshToken,
		expiresAt.Format(constants.SqliteDateTimeFormat),
	)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/delete-refresh-token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
	cookie, err := http.ParseSetCookie(rr.Header().Get("Set-Cookie"))
	require.NoError(t, err)

	require.Equal(t, "refresh_token", cookie.Name)
	require.Equal(t, "", cookie.Value)
	require.Equal(t, "/", cookie.Path)
	require.Equal(t, true, cookie.HttpOnly)
	require.Equal(t, -1, cookie.MaxAge)

	token, err := model.FindRefreshTokenByJwt(c.Db, refreshToken)
	require.NoError(t, err)
	require.Nil(t, token)
}

// TODO: more tests
