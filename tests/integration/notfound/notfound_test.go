package notfound

import (
	"murky_api/internal/app"
	"murky_api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent-path", nil)
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestNotFound(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent-path", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}
