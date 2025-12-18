package auth

import (
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/auth"
	"murky_api/internal/constants"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRefreshAccessTokenSuccess(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh-access-token", nil)
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

	require.Equal(t, http.StatusOK, rr.Code)

	var resp auth.RefreshAccessTokenResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.NotEmpty(t, resp.AccessToken)
}

// TODO: more tests
