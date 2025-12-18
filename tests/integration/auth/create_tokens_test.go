package auth

import (
	"bytes"
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/auth"
	"murky_api/internal/routing"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateTokensInvalidContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/create-tokens", nil)
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestCreateTokensInvalidJson(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/create-tokens", strings.NewReader("{]"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateTokensValidationError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/auth/create-tokens", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp routing.ValidationErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []string{"Must not be blank."}, resp.FieldErrors["username"])
	require.Equal(t, []string{"Must not be blank."}, resp.FieldErrors["password"])
}

func TestCreateTokensUserNotFound(t *testing.T) {
	body := auth.CreateTokensRequest{Username: "invaliduser", Password: "pass"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/create-tokens", bytes.NewReader(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp routing.ValidationErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []string{"Invalid credentials"}, resp.GeneralErrors)
}

func TestCreateTokensInvalidPassword(t *testing.T) {

	body := auth.CreateTokensRequest{Username: "user", Password: "pass1"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/create-tokens", bytes.NewReader(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp routing.ValidationErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []string{"Invalid credentials"}, resp.GeneralErrors)
}

func TestCreateTokensSuccess(t *testing.T) {
	body := auth.CreateTokensRequest{Username: "user", Password: "pass"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/create-tokens", bytes.NewReader(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp auth.CreateTokensResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.NotEmpty(t, resp.AccessToken)

	cookie, err := http.ParseSetCookie(rr.Header().Get("Set-Cookie"))
	require.NoError(t, err)

	require.Equal(t, "refresh_token", cookie.Name)
	require.Equal(t, "/", cookie.Path)
	require.Equal(t, true, cookie.HttpOnly)
	require.Equal(t, false, cookie.Secure) // TODO: change this to test for true
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	require.Greater(t, cookie.Expires, time.Now().UTC().Add(time.Hour*24*7).Add(-time.Minute))

}
