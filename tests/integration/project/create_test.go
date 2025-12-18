package project

import (
	"bytes"
	"encoding/json"
	"murky_api/internal/app"
	"murky_api/internal/model"
	"murky_api/internal/project"
	"murky_api/internal/routing"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/projects", nil)
	rr := httptest.NewRecorder()

	c := app.NewTestContainer(t)
	defer c.Close()
	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCreateInvalidContentType(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/projects", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestCreateInvalidJson(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader("{]"))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateValidationError(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader("{}"))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp routing.ValidationErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []string{
		"Must not be blank.",
		"Has to be at least 2 characters long.",
		"Must only contain letters, numbers, dashes and spaces.",
	}, resp.FieldErrors["title"])
	require.Equal(t, []string{
		"Must not be blank.",
		"Has to be at least 2 characters long.",
		"Must only contain letters, numbers and dashes.",
	}, resp.FieldErrors["slug"])
}

func TestCreateNotUniqueSlugError(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	body := project.CreateRequest{Title: "Test Project", Slug: "project-1"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp routing.ValidationErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, []string{"A project with this slug already exists. Slug must be unique."}, resp.FieldErrors["slug"])
}

func TestCreateSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	token, err := c.JwtService.CreateAccessToken(model.User{Id: 1, Username: "user"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	body := project.CreateRequest{Title: "Test Project", Slug: "test-project"}
	bodyJson, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(bodyJson))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.NewMux(c).ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp project.CreateResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, body.Title, resp.Title)
	require.Equal(t, body.Slug, resp.Slug)
}
